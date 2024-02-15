package models

import (
	"errors"

	// "gorm.io/gorm"
)

// "math/rand"

// "time"

// "gorm.io/gorm"


type TypeReports struct {
	// gorm.Model `gorm:"softDelete:false"`
	Id uint `gorm:"not null;unique" json:"id"`
	Type_Report string `gorm:"size:13;not null;unique" json:"type_report"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	
}

func (TypeReports) TableName() string {
    return "reports_types"
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


func GetReportsTypes(regioncode string) ([]TypeReports, error) {

	var types []TypeReports
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Where("region_code = ?", regioncode).Find(&types).Error; err != nil {
		return types, errors.New("types not found! ")
	}

	

	return types, nil

}