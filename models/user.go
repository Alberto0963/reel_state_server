package models

import (
	"errors"
	"html"
	// "mime/multipart"

	// "io"
	// "mime/multipart"
	// "os"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID uint `gorm:"not null;unique" json:"id"`
	Phone string `gorm:"size:13;not null;unique" json:"phone"`
	Username     string `gorm:"size:255;not null;unique" json:"username"`
	Password     string `gorm:"size:100;not null;" json:"password"`
	ProfileImage string `gorm:"size:255;not null;" json:"profileImage"`
}

func GetUserByID(uid uint) (User, error) {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

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
