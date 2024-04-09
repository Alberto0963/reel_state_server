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

func UpdatePasswordHandler(c *gin.Context) {
	auth.ResetPassword(c)
}

func UpdateUsernameHandler(c *gin.Context) {
	auth.UpdateUserName(c)
}

func CurrentUserHandler(c *gin.Context) {
	auth.CurrentUser(c)
}

func UserByIdHandler(c *gin.Context) {

	id := c.Query("id")

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := models.GetUserByID(uint(userID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": result})
}

func SendVerificationCode(c *gin.Context) {
	auth.SendVerificationCode(c)
}

func ValidateVerificationCode(c *gin.Context) {
	auth.ValidateVerificationCode(c)
}

func SearchProfile(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	username := c.Query("username")
	// username, err := strconv.ParseUint(p, 10, 64)

	// page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))

	cat, err := models.SearchProfile(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetMyVideos(c *gin.Context) {

	userID, _ := token.ExtractTokenID(c)
	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	typeV := c.Query("type")
	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Type"})
		return
	}
	// page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))

	cat, err := models.GetMyVideos(int(userID), int(page), int(typeVideo))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetUserVideos(c *gin.Context) {

	id := c.Query("id")

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	typeV := c.Query("type")
	typeVideo, err := strconv.ParseUint(typeV, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}
	// userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	cat, err := models.GetMyVideos(int(userID), int(page), int(typeVideo))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
}

func GetUserFavoritesVideos(c *gin.Context) {

	userID, _ := token.ExtractTokenID(c)

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	cat, err := models.GetMyFavoritesVideos(int(userID), int(page))
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
	profileImagePath := filepath.Join("public/profile_images", imageFileName)

	userID, _ := token.ExtractTokenID(c)

	u := models.User{}

	u.ProfileImage = profileImagePath
	u.ID = userID

	_, err = u.UpdateProfileImageUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(profileImage, url+profileImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image is updated"})
}

func UpdateCoverImageUserName(c *gin.Context) {

	coverImage, err := c.FormFile("cover_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
		return
	}
	// Generate a random file name for the profile image
	imageFileName := models.GenerateRandomName() + filepath.Ext(coverImage.Filename)
	url := os.Getenv("MY_URL")
	coverImagePath := filepath.Join("public/profile_images", imageFileName)

	userID, _ := token.ExtractTokenID(c)

	u := models.User{}

	u.Cover_image = coverImagePath
	u.ID = userID

	_, err = u.UpdateCoverImageUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(coverImage, url+coverImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image is updated"})
}

func GetMemberShips(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := c.Query("page")
	page, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
		return
	}

	currencyCode := c.Query("currencyCode")
	// currencyCode, err := strconv.ParseUint(p, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page"})
	// 	return
	// }

	cat, err := models.GetMemberShips(currencyCode, int(page))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat, "user": u})
}

func DeleteUserVideo(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := c.Query("idVideo")
	idVideo, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = models.DeleteUserVideo(int(idVideo), int(user_id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}


