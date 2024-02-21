package models

import (
	"errors"

	// "gorm.io/gorm"
)

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