package controllers

import (
	// "io"
	// "mime/multipart"
	"net/http"
	"path/filepath"
	// "os"

	"reelState/auth"
	"reelState/models"

	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/nacl/auth"
)

func LoginHandler(c *gin.Context) {
	// secretKey := []byte("your-secret-key") // Replace with your own secret key
	auth.Login(c, models.Pool) // Pass the DB connection and secret key to the Login function
}

func RegisterHandler(c *gin.Context) {
	auth.Register(c)
}

func CurrentUserHandler(c *gin.Context) {
	auth.CurrentUser(c)
}

func HandleVideoUpload(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a random file name
	fileName := models.GenerateRandomName() + filepath.Ext(file.Filename)

	// Create the destination file
	destPath := filepath.Join("home/reelstate/go/reel_state_server/public/videos", fileName)

	err = models.SaveVideo(file, destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
}
