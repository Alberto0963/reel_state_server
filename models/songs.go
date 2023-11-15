package models

import (
	"errors"

)

type Song struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID                       int      `gorm:"not null;unique" json:"id"`
	Name                    string    `gorm:"size:13;not null;unique" json:"name"`

}


func (Song) TableName() string {
	return "songs"
}

func GetAllSongs() ([]Song, error) {

	var songs []Song
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Find(&songs).Error; err != nil {
		return songs, errors.New("songs not found!")
	}

	return songs, nil

}


