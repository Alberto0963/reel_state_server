package models

import (
	"errors"
)

type Membership struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID             uint   `gorm:"not null;unique" json:"id"`
	Membership     string `gorm:"size:13;not null;unique" json:"membership"`
	Description    string `gorm:"size:255;not null;unique" json:"description"`
	Price          string `gorm:"size:100;not null;" json:"price"`
	ProductCode    string `gorm:"size:255;not null;" json:"product_code"`
	CurrencyCode   string `gorm:"size:255;" json:"currency_code"`
	MembershipCode string `gorm:"size:255;not null;" json:"memberships_code"`
}

func GetMemberShips(code string, page int) ([]Membership, error) {

	// var m Membership
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	// result
	pageSize := 12
	var member []Membership

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize

	if err := dbConn.Model(&Membership{}).Where("currency_code = ?", code).Limit(pageSize).Offset(offset).Find(&member).Error; err != nil {
		return member, errors.New("membership not found!")
	}

	return member, nil

}
