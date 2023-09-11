package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// "fmt"

// "github.com/jinzhu/gorm"



type Favorites struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id int `gorm:"not null;unique" json:"id"`
	Id_user int `gorm:"size:13;not null;unique" json:"id_user"`
	Id_video     int `gorm:"size:255;not null;unique" json:"id_video"`
	// Location     string `gorm:"size:100;not null;" json:"location"`
	// Area string `gorm:"size:255;not null;" json:"area"`
	// Property_number string `gorm:"size:255;not null;" json:"property_number"`

	// Price string `gorm:"size:255;not null;" json:"price"`
	// Id_user uint `gorm:"size:255;not null;" json:"id_user"`
	// Sale_type_id int `gorm:"size:255;not null;" json:"sale_type_id"`
	// SaleType Type `gorm:"references:id; foreignKey:sale_type_id"`
	// Sale_category_id int `gorm:"size:255;not null;" json:"sale_category_id"`
	// SaleCategory Category `gorm:"references:id; foreignKey:sale_category_id"`
	// Image_cover string `gorm:"size:255;not null;" json:"image_cover"`

}

func (Favorites) TableName() string {
	return "users_videos_favorites"
}


func IsVideoFavorite(id_user int, id_video int) (error) {
	dbConn := Pool

	// Fetch the existing user from the database
	var fav Favorites
	if err := dbConn.Where("id_user = ? && id_video = ?", id_user,id_video).First(&fav).Error; 
	err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// No record found, which means 'fav' is null
			fmt.Println("fav is null")
			return err
		} 
		//else {
		// 	// Handle other errors
		// 	return err
		// }
		return err
	}

	// oldImage := vid.
	// Update the user fields with the new values
	// user.ProfileImage = updatedUser.ProfileImage
	// user.Email = updatedUser.Email
	// Update other user fields as needed...

	// Save the changes to the database

	// if err := dbConn.Save(&user).Error; err != nil {
	// 	return user, err
	// }
	// pathOldImage := os.Getenv("MY_URL")

	
	//  deleteImage(pathOldImage + oldImage)
	// if err != nil {
	// 	return user, err
	// }

	return  nil
}

func SetVideoFavorite(fav Favorites) (Favorites,error) {
	// dbConn = Pool

	// Fetch the existing user from the database
	var err error
	dbConn := Pool

	err = dbConn.Create(&fav).Error
	if err != nil {
		return fav, err
	}
	return fav, nil

	// oldImage := vid.
	// Update the user fields with the new values
	// user.ProfileImage = updatedUser.ProfileImage
	// user.Email = updatedUser.Email
	// Update other user fields as needed...

	// Save the changes to the database

	// if err := dbConn.Save(&user).Error; err != nil {
	// 	return user, err
	// }
	// pathOldImage := os.Getenv("MY_URL")

	
	//  deleteImage(pathOldImage + oldImage)
	// if err != nil {
	// 	return user, err
	// }

	// return  nil
}

func DeleteFavoritetByID(id_user int, id_video int) error {
	var err error
	dbConn := Pool
	var fav Favorites
	if err = dbConn.Where("id_user = ? && id_video = ?", id_user,id_video).Delete(&fav).Error; err != nil {
		return err
	}
	return nil
}