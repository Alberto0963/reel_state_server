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
	"time"

	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/nacl/auth"
)

func LoginHandler(c *gin.Context) {
	// secretKey := []byte("your-secret-key") // Replace with your own secret key
	auth.Login(c, models.Pool) // Pass the DB connection and secret key to the Login function
}

func LoginWithGoogleHandler(c *gin.Context) {
	// secretKey := []byte("your-secret-key") // Replace with your own secret key
	auth.HandleGoogleLogin(c, models.Pool) // Pass the DB connection and secret key to the Login function
}

func RegisterHandler(c *gin.Context) {
	auth.Register(c)
}

func RegisterHandlerWithGogle(c *gin.Context) {
	auth.HandleGoogleRegister(c)
}

func UpdatePasswordHandler(c *gin.Context) {
	auth.ResetPassword(c)
}

func UpdateUsernameHandler(c *gin.Context) {
	auth.UpdateUserName(c)
}

func UpdatePhoneNumberHandler(c *gin.Context) {
	auth.UpdatePhoneNumber(c)
}


func CurrentUserHandler(c *gin.Context) {
	auth.CurrentUser(c)
}

func UserByIdHandler(c *gin.Context) {

	actualUserID, _ := token.ExtractTokenID(c)

	id := c.Query("id")

	profileId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	result, err := models.GetUserByID(uint(profileId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	followers, err := models.Getfollowers(int(profileId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	imFollower, err := models.Imfollower(int(profileId), int(actualUserID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"data":       result,
		"followers":  followers,
		"actualUser": actualUserID,
		"imfollower": imFollower,
	})
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

	u := models.UserUpdate{}

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

func GetPublicMemberShips(c *gin.Context) {

	// userID, _ := token.ExtractTokenID(c)
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

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": cat})
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

func Setlike(c *gin.Context) {

	actualUserID, _ := token.ExtractTokenID(c)
	// var u *models.Likes
	id := c.Query("id")
	profileId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	u, err := models.GetProfilefollowers(int(profileId), int(actualUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user"})
		return
	}

	if u.Id == 0 {
		u = models.Likes{}
		u.Id_profile = uint(profileId)
		u.Id_user = actualUserID
		models.SaveLikeProfile(&u)

	} else {
		models.UpdateLikeProfile(&u)
	}

	// if(u == nil){
	// 	u = *models.Likes{}

	// 	u.Id_profile = uint(profileId)
	// 	u.Id_user = actualUserID
	// 	models.SaveLikeProfile(u)
	// }

	// _, err = u.UpdateUser()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// token, err := token.GenerateToken(u.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "update success", "d": u})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func CreateSubscription(c *gin.Context) {

	var sub models.Createsubscription
	actualUserID, _ := token.ExtractTokenID(c)

	// Bind JSON to struct
	if err := c.ShouldBind(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub.IdUser = int(actualUserID)

	_, err := sub.CreateSubscription()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully"})
}

func CancelSubscription(c *gin.Context) {
	var sub models.CancelSubscription

	// Bind JSON to struct
	if err := c.ShouldBind(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"
	now := time.Now()
	formattedDate := now.Format("2006-01-02 15:04:05.999999999 -0700 MST")

	// Parse the date string
	parsedTime, err := time.Parse(layout, formattedDate)
	if err != nil {
		// fmt.Println("Error parsing date:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Error Parse Time": err.Error()})

		return
	}

	err = models.CancelPaypalSubscription(sub.PaypalSubscriptionId, sub.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error Paypal Sub": err.Error()})
		return
	}

	err = models.CancelSubscriptionFunction(sub.PaypalSubscriptionId, parsedTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error Calcel sub": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription cancelled successfully"})
}

func GetUserSubscription(c *gin.Context) {

	// var sub models.Createsubscription
	actualUserID, _ := token.ExtractTokenID(c)

	// Bind JSON to struct
	// if err := c.ShouldBind(&sub); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// sub.IdUser = int(actualUserID)

	subs, err := models.GetSubscription(int(actualUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"data": subs})
}
