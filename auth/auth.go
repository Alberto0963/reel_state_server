package auth

import (
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
	"gorm.io/gorm"
	// "gorm.io/gorm"
	// "gorm.io/gorm"
)

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
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


	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(input.Password)); err != nil {
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
	c.JSON(http.StatusOK, gin.H{"token": token, "isVip": isVip,"canUpload": user.CanUpload})
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:""`
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

	u := models.User{}

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

	u, err := models.GetUserByID(userID)
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

	c.JSON(http.StatusOK, gin.H{"message": "update success",})
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
