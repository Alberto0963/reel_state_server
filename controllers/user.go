package controllers

import (
	// "io"
	// "mime/multipart"

	// "os"

	"net/http"
	"os"
	"path/filepath"
	"reelState/auth"
	"reelState/models"
	"reelState/utils/token"
	"strconv"

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

func SendVerificationCode(c *gin.Context) {
	auth.SendVerificationCode(c)
}

func ValidateVerificationCode(c *gin.Context) {
	auth.ValidateVerificationCode(c)
}

func GetMyVideos(c *gin.Context) {

	userID, _ := token.ExtractTokenID(c)
	page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))

	cat, err := models.GetMyVideos(int(userID), page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func ValidateUserName(c *gin.Context) {

	username := c.DefaultPostForm("UserName", "")

	if models.UsernameExists(username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username is available"})
}

func UpdateProfileImageUserName(c *gin.Context) {

	profileImage, err := c.FormFile("profile_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
		return
	}
	// Generate a random file name for the profile image
	imageFileName := models.GenerateRandomName() + filepath.Ext(profileImage.Filename)
	url := os.Getenv("MY_URL")
	profileImagePath := filepath.Join(url +"public/profile_images", imageFileName)
	
	userID, _ := token.ExtractTokenID(c)

	u := models.User{}

	u.ProfileImage = profileImagePath
	u.ID = userID

	_, err = u.UpdateProfileImageUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Save the profile image
	err = c.SaveUploadedFile(profileImage, url+ profileImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image is updated"})
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
