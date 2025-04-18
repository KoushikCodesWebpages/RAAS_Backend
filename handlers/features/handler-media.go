
package features

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"RAAS/config"
)

// MediaUploadHandler handles media uploads to Azure Blob Storage
type MediaUploadHandler struct {
	blobServiceClient *azblob.ServiceURL
}

// NewMediaUploadHandler creates a new MediaUploadHandler with the provided Azure Blob service client
func NewMediaUploadHandler(blobServiceClient *azblob.ServiceURL) *MediaUploadHandler {
	return &MediaUploadHandler{
		blobServiceClient: blobServiceClient,
	}
}

// UploadMedia uploads a file to the specified Azure Blob container
func (h *MediaUploadHandler) UploadMedia(c *gin.Context, containerName string) (string, error) {
	// Get the file from the form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get the container URL dynamically
	containerURL := h.blobServiceClient.NewContainerURL(containerName)
	blobURL := containerURL.NewBlockBlobURL(header.Filename)

	// Upload the file to Azure Blob Storage
	_, err = azblob.UploadStreamToBlockBlob(context.Background(), file, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	// Construct the file URL for the uploaded file
	fileURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", config.Cfg.AzureStorageAccount, containerName, header.Filename)
	return fileURL, nil
}

// ValidateFileType checks if the uploaded file has a valid extension
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

// HandleUpload handles the media file upload
func (h *MediaUploadHandler) HandleUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Validate file type
	if !h.ValidateFileType(header) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Determine which container to use based on some criteria
	// Example: We can set the container name based on the request
	containerName := c.DefaultQuery("container", "default-container") // Pass this dynamically, or define for different routes

	// Upload the file to the appropriate container
	fileURL, err := h.UploadMedia(c, containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Return the file URL in the response
	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL})
}

// GetBlobServiceClient creates a new Azure Blob service client using credentials
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
