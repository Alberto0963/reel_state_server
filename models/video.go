package models

import (
	"math/rand"

	"time"

	// "gorm.io/gorm"
)


type Video struct {
	// gorm.Model
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID uint `gorm:"not null;unique" json:"id"`
	Video_url string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location     string `gorm:"size:100;not null;" json:"location"`
	Area string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`
	Price string `gorm:"size:255;not null;" json:"price"`
	Id_user string `gorm:"size:255;not null;" json:"id_user"`
	Sale_type_id string `gorm:"size:255;not null;" json:"sale_type_id"`
	Sale_category_id string `gorm:"size:255;not null;" json:"sale_category_id"`
}

func GenerateRandomName() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	name := make([]byte, length)
	for i := 0; i < length; i++ {
		name[i] = chars[rand.Intn(len(chars))]
	}
	return  "reel_state." + string( name)
}

func (v *Video) SaveVideo() (*Video, error) {
	var err error
	dbConn := Pool

	err = dbConn.Create(&v).Error
	if err != nil {
		return &Video{}, err
	}
	return v, nil

}