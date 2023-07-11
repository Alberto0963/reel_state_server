package models

import "errors"

// "math/rand"

// "time"

// "gorm.io/gorm"


type Category struct {
	// gorm.Model
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID uint `gorm:"not null;unique" json:"id"`
	Category string `gorm:"size:13;not null;unique" json:"category"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	
}

func (Category) TableName() string {
    return "sales_categories"
}
// func GenerateRandomName() string {
// 	rand.Seed(time.Now().UnixNano())
// 	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	length := 10
// 	name := make([]byte, length)
// 	for i := 0; i < length; i++ {
// 		name[i] = chars[rand.Intn(len(chars))]
// 	}
// 	return  "reel_state." + string( name)
// }

// func (v *Video) SaveVideo() (*Video, error) {
// 	var err error
// 	dbConn := Pool

// 	err = dbConn.Create(&v).Error
// 	if err != nil {
// 		return &Video{}, err
// 	}
// 	return v, nil

// }


func GetCategory() ([]Category, error) {

	var cat []Category
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Find(&cat).Error; err != nil {
		return cat, errors.New("Categories not found!")
	}

	// u.PrepareGive()

	return cat, nil

}