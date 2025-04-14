package features

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"log"
	"net/url"
	
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"


	"RAAS/config"

)

type MediaUploadHandler struct {
	blobServiceClient *azblob.ServiceURL
	containerName     string
}

func NewMediaUploadHandler(blobServiceClient *azblob.ServiceURL, containerName string) *MediaUploadHandler {
	return &MediaUploadHandler{
		blobServiceClient: blobServiceClient,
		containerName:     containerName,
	}
}

func (h *MediaUploadHandler) UploadMedia(c *gin.Context) (string, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	containerURL := h.blobServiceClient.NewContainerURL(h.containerName)
	blobURL := containerURL.NewBlockBlobURL(header.Filename)
	_, err = azblob.UploadStreamToBlockBlob(context.Background(), file, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", config.Cfg.AzureStorageAccount, h.containerName, header.Filename)
	return fileURL, nil
}

func (h *MediaUploadHandler) ValidateFileType(fileHeader *multipart.FileHeader) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".docx", ".pdf"}
	fileExtension := filepath.Ext(fileHeader.Filename)
	for _, ext := range allowedExtensions {
		if fileExtension == ext {
			return true
		}
	}
	return false
}

func (h *MediaUploadHandler) HandleUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	if !h.ValidateFileType(header) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	fileURL, err := h.UploadMedia(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL})
}


func GetBlobServiceClient() *azblob.ServiceURL {
	if config.Cfg.AzureStorageAccount == "" || config.Cfg.AzureStorageKey == "" {
		log.Fatal("Azure storage account or key is not set")
	}
	credential, err := azblob.NewSharedKeyCredential(config.Cfg.AzureStorageAccount, config.Cfg.AzureStorageKey)
	if err != nil {
		log.Fatal("Invalid credentials")
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", config.Cfg.AzureStorageAccount)
	u, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatal("Invalid service URL")
	}
	URL := azblob.NewServiceURL(*u, p)
	return &URL
}
