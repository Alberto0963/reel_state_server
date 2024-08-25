package models

import (
	"errors"

	// "gorm.io/gorm"
)

// "math/rand"

// "time"

// "gorm.io/gorm"


type Sponsors struct {
	// gorm.Model
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID uint `gorm:"not null;unique" json:"id"`
	Promotion string `gorm:"size:13;not null;unique" json:"promotion"`
	DaysPromotion     string `gorm:"size:255;not null;unique" json:"days_promotion"`
	Price  string `gorm:"size:255;not null;unique" json:"price"`
	CancelationDays  string `gorm:"size:255;not null;unique" json:"cancelation_days"`
	RegionCode  string `gorm:"size:255;not null;unique" json:"region_code"`
	PaypalId  string `gorm:"size:255;not null;unique" json:"paypal_id"`

}

func (Sponsors) TableName() string {
    return "sponsors"
}



func GetSponsors(region_code string) ([]Sponsors, error) {

	var cat []Sponsors
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Where("region_code = ?",region_code).Find(&cat).Error; err != nil {
		return cat, errors.New("Sponsors not found! ")
	}

	// u.PrepareGive()

	return cat, nil

}