package models

import (
	"errors"
	"html"
	"os"
	// "strconv"
	"time"

	// "mime/multipart"

	// "io"
	// "mime/multipart"
	// "os"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID                       uint      `gorm:"not null;unique" json:"id"`
	Phone                    string    `gorm:"size:13;not null;unique" json:"phone"`
	Username                 string    `gorm:"size:255;not null;unique" json:"username"`
	Password                 string    `gorm:"size:100;not null;" json:"password"`
	ProfileImage             string    `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time `gorm:"size:255;" json:"expiration_membership_date"`
	IdMembership             int       `gorm:"size:255;not null;" json:"id_membership"`
	RenovationActive         int       `gorm:"size:255;not null;" json:"renovation_active"`
}

type PublicUser struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id                       uint      `gorm:"not null;unique" json:"id"`
	Phone                    string    `gorm:"size:13;not null;unique" json:"phone"`
	Username                 string    `gorm:"size:255;not null;unique" json:"username"`
	ProfileImage             string    `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time `gorm:"size:255;" json:"expiration_membership_date"`
	IdMembership             int       `gorm:"size:255;not null;" json:"id_membership"`
	RenovationActive         int       `gorm:"size:255;not null;" json:"renovation_active"`
	Videos                   []MyVideo `gorm:"references:id; foreignKey:id_user"`
}

func (PublicUser) TableName() string {
	return "users"
}

func (updatedUser *User) UpdateProfileImageUser() (User, error) {
	dbConn := Pool

	// Fetch the existing user from the database
	var user User
	if err := dbConn.First(&user, updatedUser.ID).Error; err != nil {
		return user, err
	}

	oldImage := user.ProfileImage
	// Update the user fields with the new values
	user.ProfileImage = updatedUser.ProfileImage
	// user.Email = updatedUser.Email
	// Update other user fields as needed...

	// Save the changes to the database

	if err := dbConn.Save(&user).Error; err != nil {
		return user, err
	}
	pathOldImage := os.Getenv("MY_URL")

	
	 deleteImage(pathOldImage + oldImage)
	// if err != nil {
	// 	return user, err
	// }

	return user, nil
}

func deleteImage(imagePath string) error {
	err := os.Remove(imagePath)
	if err != nil {
		// Handle the error
		return err
	}

	return nil
}

func GetUserByID(uid uint) (User, error) {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&User{}).First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func UsernameExists(username string) bool {

	var u PublicUser
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	result := dbConn.Model(&PublicUser{}).Where("username = ?", username).First(&u)

	if result.Error == nil {
		// The username already exists
		return true
	}
	// u.PrepareGive()

	return false

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() (*User, error) {

	var err error
	dbConn := Pool

	err = dbConn.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func GetMyVideos(id_user int, page int) ([]MyVideo, error) {
	var err error
	dbConn := Pool
	var vid []MyVideo
	// Get page number and page size from query parameters
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	// err = dbConn.Unscoped().Find(&vid).Error
	err = dbConn.Model(&MyVideo{}).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Unscoped().Find(&vid).Error
	if err != nil {
		return vid, err
	}
	return vid, nil

}
