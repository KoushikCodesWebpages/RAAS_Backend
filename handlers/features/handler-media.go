package features

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"log"
	"net/url"
	"github.com/spf13/viper"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	fileURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", os.Getenv("AZURE_STORAGE_ACCOUNT"), h.containerName, header.Filename)
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


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetBlobServiceClient() *azblob.ServiceURL {
	accountName := viper.GetString("AZURE_STORAGE_ACCOUNT")
	accountKey := viper.GetString("AZURE_STORAGE_KEY")
	if accountName == "" || accountKey == "" {
		log.Fatal("Azure storage account or key is not set")
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials")
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	u, err := url.Parse(serviceURL)
	if err != nil {
		log.Fatal("Invalid service URL")
	}
	URL := azblob.NewServiceURL(*u, p)
	return &URL
}
func uploadFileToBlobStorage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid file"})
		return
	}

	blobServiceClient := GetBlobServiceClient()
	containerURL := blobServiceClient.NewContainerURL(os.Getenv("AZURE_BLOB_CONTAINER"))
	blobURL := containerURL.NewBlockBlobURL(header.Filename)

	_, err = azblob.UploadStreamToBlockBlob(context.Background(), file, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}