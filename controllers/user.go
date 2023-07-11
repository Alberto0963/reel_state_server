package controllers

import (
	// "io"
	// "mime/multipart"
	
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

func getCategoryAndTypeHandler(c *gin.Context) {
	
}

// func HandleVideoUpload(c *gin.Context) {
// 	file, err := c.FormFile("video")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}



	// Create the destination file



// 	// Generate a random file name
// 	fileName := models.GenerateRandomName() + filepath.Ext(file.Filename)

// 	// Create the destination file
// 	destPath := filepath.Join("public/videos", fileName)

// 	err =  saveVideoFile(file, destPath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var input models.Video

// 	if err := c.ShouldBind(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}


// 	v := models.Video{}

// 	_, err = v.SaveVideo()

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully"})
// }


// func saveVideoFile(file *multipart.FileHeader, destination string) error {
// 	src, err := file.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()

// 	dst, err := os.Create(destination)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	_, err = io.Copy(dst, src)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
