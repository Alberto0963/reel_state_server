package auth

import (
	// "io/ioutil"
	// "io"
	// "database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	// "os"

	// "os"

	// "path/filepath"
	"reelState/models"
	SMS "reelState/utils"
	"reelState/utils/token"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	// "gorm.io/gorm"
	// "gorm.io/gorm"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	followers, err := models.Getfollowers(int(user_id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	imFollower, err := models.Imfollower(int(user_id), int(user_id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u, "followers": followers, "imfollower": imFollower})
}

type LoginInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context, db *gorm.DB) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	if err := db.Table("view_user_upload_status").Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	// if err != nil {
	//     panic(err)
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	isVip := false
	if user.IdMembership == 6 || user.IdMembership == 7 {
		isVip = true
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "isVip": isVip, "canUpload": user.CanUpload})
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:""`
}

type LoginWithGoogle struct {
	Token string `json:"token " binding:"required"`
	// Password string `json:"password" binding:"required"`
	// Phone    string `json:"phone" binding:"required"`
	// Code     string `json:"code" binding:""`
}

type ResetPasswordInput struct {
	// Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:""`
}

type UpdateUserNameInput struct {
	Username string `json:"username" binding:"required"`
	// Password string `json:"password" binding:"required"`
	// Phone    string `json:"phone" binding:"required"`
	// Code     string `json:"code" binding:""`
}

type VerificationPhoneInput struct {
	Phone        string `json:"phone" binding:"required"`
	AppSignature string `json:"appSignature" binding:"required"`
}

type VerificationCodeInput struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func Register(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// profileImage, err := c.FormFile("profile_image")
	// if err != nil {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from the request"})
	//     return
	// }
	// Generate a random file name for the profile image
	// imageFileName := models.GenerateRandomName() + filepath.Ext(profileImage.Filename)
	// url := os.Getenv("MY_URL")

	// Create the destination path for saving the image
	// profileImagePath := filepath.Join("public/profile_images", imageFileName)

	u := models.UserUpdate{}

	u.Username = input.Username
	u.Password = input.Password
	u.Phone = input.Phone

	// u.ProfileImage = profileImagePath
	u.ExpirationMembershipDate = time.Now()
	u.IdMembership = 1
	// // Get the current date and time
	// currentTime := time.Now()

	// // Extract the date from the current time
	// // currentDate := currentTime.Format("2006-01-02")

	// date := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	// // Define a duration to add
	// duration := 24 * time.Hour

	// // Add the duration to the date
	// sum := date.Add(duration)
	vc := models.VerificationCode{}
	vc.Code = input.Code
	vc.Phone = input.Phone

	// _, err = vc.CodeIsValid()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
	// 	return
	// }

	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := token.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Save the profile image
	// err = c.SaveUploadedFile(profileImage, url+ profileImagePath)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "registration success", "token": token})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func ResetPassword(c *gin.Context) {

	var input ResetPasswordInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByPhone(input.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user"})
		return
	}
	// u.Username = input.Username
	u.Password = input.Password
	// u.Phone = input.Phone

	// u.ProfileImage = profileImagePath
	// u.ExpirationMembershipDate = time.Now()
	// u.IdMembership = 1

	vc := models.VerificationCode{}
	vc.Code = input.Code
	vc.Phone = input.Phone

	_, err = vc.CodeIsValid()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	_, err = u.UpdateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := token.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update success", "token": token})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func UpdateUserName(c *gin.Context) {

	var input UpdateUserNameInput
	userID, _ := token.ExtractTokenID(c)

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByIdToUpdate(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user"})
		return
	}
	u.Username = input.Username

	_, err = u.UpdateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// token, err := token.GenerateToken(u.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "update success"})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func SendVerificationCode(c *gin.Context) {

	var input VerificationPhoneInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := models.VerificationCode{}
	var code = SMS.GenerateRandomCode(6)

	v.Code = code
	v.Phone = input.Phone

	_, err := v.SaveVerificationCode()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SMS.SendSMS(v.Phone, code, input.AppSignature)
	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

func ValidateVerificationCode(c *gin.Context) {

	var input VerificationCodeInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v := models.VerificationCode{}
	// var code = SMS.GenerateRandomCode(6)

	v.Code = input.Code
	v.Phone = input.Phone

	_, err := v.CodeIsValid()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "code is valid"})
	// c.JSON(http.StatusOK, gin.H{"message": "validated!"})

}

var googleOauthConfig = &oauth2.Config{
	ClientID:     "863989854330-88369t3et8090geknm71tj9rjve196ti.apps.googleusercontent.com",
	ClientSecret: "com.reel_state_mx",
	Endpoint:     google.Endpoint,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
}

func HandleGoogleLogin(c *gin.Context,db *gorm.DB) {
	// ctx := r.Context()
	var input LoginWithGoogle

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenSource := googleOauthConfig.TokenSource(c, &oauth2.Token{
		AccessToken: input.Token,
	})

	newToken, err := tokenSource.Token() // Verifica y posiblemente refresca el token
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	print(newToken)

	client := oauth2.NewClient(c, tokenSource)
	userData, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if userData.StatusCode != 200 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	defer userData.Body.Close()
	userInfo, err := io.ReadAll(userData.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var userInfoModel GoogleUserInfo
	if err := json.Unmarshal(userInfo, &userInfoModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//  var localuser  models.User
	localuser := models.User{}
	if err := db.Table("view_user_upload_status").Where("email = ?", userInfoModel.Email).First(&localuser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := token.GenerateToken(localuser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	isVip := false
	if localuser.IdMembership == 6 || localuser.IdMembership == 7 {
		isVip = true
	}

	// Envía userInfo al cliente Flutter o procesa según necesites
	c.JSON(http.StatusOK, gin.H{"message": "token is valid", "GoogleUser": userInfoModel, "ReelStateUser": localuser, "token": token, "isVip": isVip,"canUpload": localuser.CanUpload})
	// return
}

func googleService(c *gin.Context, token string) (GoogleUserInfo, error) {
	// ctx := r.Context()
	// var input LoginWithGoogle
	var userGoogle GoogleUserInfo

	tokenSource := googleOauthConfig.TokenSource(c, &oauth2.Token{
		AccessToken: token,
	})

	newToken, err := tokenSource.Token() // Verifica y posiblemente refresca el token
	if err != nil {
		return userGoogle, err

	}
	print(newToken)

	client := oauth2.NewClient(c, tokenSource)
	userData, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if userData.StatusCode != 200 {
		return userGoogle, err

	}

	defer userData.Body.Close()
	userInfo, err := io.ReadAll(userData.Body)
	if err != nil {
		return userGoogle, err

	}
	// var userInfoModel GoogleUserInfo
	if err := json.Unmarshal(userInfo, &userGoogle); err != nil {
		return userGoogle, err

	}

	return userGoogle, nil

}

func HandleGoogleRegister(c *gin.Context) {

	var input LoginWithGoogle

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInfoModel, err := googleService(c, input.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	//  var localuser  models.User
	// Check if user already exists in database by email
	localuser, err := models.GetUserByEmail(userInfoModel.Email)
	log.Printf("user:",localuser)
	// log
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	if err != nil {

		// No existing user, let's register a new one
		var saveUser models.UserUpdate

		saveUser.Email = userInfoModel.Email
		saveUser.Username = userInfoModel.Name
		// saveUser.Username = userInfoModel.Name
		saveUser.ExpirationMembershipDate = time.Now()
		saveUser.IdMembership = 1
		usr, err := saveUser.SaveUser()
		print(usr)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database error:": err})
			return
		}

		token, err := token.GenerateToken(saveUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database error:": err})
			return
		}
		isVip := false
		if saveUser.IdMembership == 6 || saveUser.IdMembership == 7 {
			isVip = true
		}

		// Envía userInfo al cliente Flutter o procesa según necesites
		c.JSON(http.StatusOK, gin.H{"GoogleUser": userInfoModel, "ReelStateUser": saveUser, "token": token, "isVip": isVip})

	} else {
		// User already registered
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered"})
		return
	}

	// return
}
