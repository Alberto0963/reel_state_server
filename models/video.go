package models

import (
	"math/rand"
	"time"
	// "gorm.io/gorm"
)


type Video struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id int `gorm:"not null;unique" json:"id"`
	Video_url string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location     string `gorm:"size:100;not null;" json:"location"`
	Area string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`

	Price string `gorm:"size:255;not null;" json:"price"`
	Id_user uint `gorm:"size:255;not null;" json:"id_user"`
	User PublicUser `gorm:"references:Id; foreignKey:Id_user"`
	Sale_type_id int `gorm:"size:255;not null;" json:"sale_type_id"`
	SaleType Type `gorm:"references:id; foreignKey:sale_type_id"`
	Sale_category_id int `gorm:"size:255;not null;" json:"sale_category_id"`
	SaleCategory Category `gorm:"references:id; foreignKey:sale_category_id"`
	Image_cover string `gorm:"size:255;not null;" json:"image_cover"`

}

type MyVideo struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id int `gorm:"not null;unique" json:"id"`
	Video_url string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location     string `gorm:"size:100;not null;" json:"location"`
	Area string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`

	Price string `gorm:"size:255;not null;" json:"price"`
	Id_user uint `gorm:"size:255;not null;" json:"id_user"`
	Sale_type_id int `gorm:"size:255;not null;" json:"sale_type_id"`
	SaleType Type `gorm:"references:id; foreignKey:sale_type_id"`
	Sale_category_id int `gorm:"size:255;not null;" json:"sale_category_id"`
	SaleCategory Category `gorm:"references:id; foreignKey:sale_category_id"`
	Image_cover string `gorm:"size:255;not null;" json:"image_cover"`

}

func (MyVideo) TableName() string {
	return "videos"
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

func  FetchAllVideos(sale_type int, isvip int,page int) ([]Video, error) {
	var err error
	dbConn := Pool
	var vid []Video
	pageSize := 10

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	// err = dbConn.Unscoped().Find(&vid).Error
	err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Order("RAND()").Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	if err != nil {
		return vid, err
	}
	return vid, nil

}



