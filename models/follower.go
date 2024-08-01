package models

import (

)

type Likes struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id              int    `gorm:"not null;unique" json:"id"`
	Id_user          uint       `gorm:"size:255;not null;" json:"id_user"`
	Id_profile          uint       `gorm:"size:255;not null;" json:"id_profile"`

}