package models

import (
	"errors"

	"gorm.io/gorm"
	// "gorm.io/gorm"
)

// "math/rand"

// "time"

// "gorm.io/gorm"

type Sponsors struct {
	// gorm.Model
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID              uint   `gorm:"not null;unique" json:"id"`
	Promotion       string `gorm:"size:13;not null;unique" json:"promotion"`
	DaysPromotion   string `gorm:"size:255;not null;unique" json:"days_promotion"`
	Price           string `gorm:"size:255;not null;unique" json:"price"`
	CancelationDays string `gorm:"size:255;not null;unique" json:"cancelation_days"`
	RegionCode      string `gorm:"size:255;not null;unique" json:"region_code"`
	PaypalId        string `gorm:"size:255;not null;unique" json:"paypal_id"`
	OpenpayId       string `gorm:"size:255;not null;" json:"openpay_id"`
}

func (Sponsors) TableName() string {
	return "sponsors"
}

func GetSponsors(region_code string) ([]Sponsors, error) {

	var cat []Sponsors
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Where("region_code = ?", region_code).Find(&cat).Error; err != nil {
		return cat, errors.New("Sponsors not found! ")
	}

	// u.PrepareGive()

	return cat, nil

}

// CreateSponsor creates a new sponsor in the database
func CreateSponsor(sponsor Sponsors) (*Sponsors, error) {
	dbConn := Pool // Get database connection

	if err := dbConn.Create(&sponsor).Error; err != nil {
		return nil, errors.New("failed to create sponsor: " + err.Error())
	}

	return &sponsor, nil
}

// EditSponsor updates an existing sponsor by its ID
func EditSponsor(id uint, updatedData Sponsors) (*Sponsors, error) {
	dbConn := Pool
	var sponsor Sponsors

	// Find the sponsor by ID
	if err := dbConn.First(&sponsor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("sponsor not found")
		}
		return nil, errors.New("failed to retrieve sponsor: " + err.Error())
	}

	// Update the sponsor fields
	if err := dbConn.Model(&sponsor).Updates(updatedData).Error; err != nil {
		return nil, errors.New("failed to update sponsor: " + err.Error())
	}

	return &sponsor, nil
}

// DeleteSponsor removes a sponsor from the database by its ID
func DeleteSponsor(id uint) error {
	dbConn := Pool
	var sponsor Sponsors

	// Find the sponsor by ID
	if err := dbConn.First(&sponsor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("sponsor not found")
		}
		return errors.New("failed to retrieve sponsor: " + err.Error())
	}

	// Delete the sponsor
	if err := dbConn.Delete(&sponsor).Error; err != nil {
		return errors.New("failed to delete sponsor: " + err.Error())
	}

	return nil
}
