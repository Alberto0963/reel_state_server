package models

import (
	"errors"
	"html"
	"os"
	// "reelState/utils/token"

	// "strconv"
	"time"

	// "mime/multipart"

	// "io"
	// "mime/multipart"
	// "os"

	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID                       uint      `gorm:"not null;unique" json:"id"`
	Phone                    string    `gorm:"size:13;not null;unique" json:"phone"`
	Username                 string    `gorm:"size:255;not null;unique" json:"username"`
	Password                 string    `gorm:"size:100;not null;" json:"password"`
	ProfileImage             string    `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time `gorm:"size:255;" json:"expiration_membership_date"`
	IdMembership             int       `gorm:"size:255;not null;" json:"id_membership"`
	RenovationActive         int       `gorm:"size:255;not null;" json:"renovation_active"`
	Cover_image              string    `gorm:"size:255;not null;" json:"cover_image"`
	Description              string    `gorm:"size:255" json:"description"`

}

type PublicUser struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID                       uint      `gorm:"not null;unique" json:"id"`
	Phone                    string    `gorm:"size:13;not null;unique" json:"phone"`
	Username                 string    `gorm:"size:255;not null;unique" json:"username"`
	ProfileImage             string    `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time `gorm:"size:255;" json:"expiration_membership_date"`
	Id_Membership             int       `gorm:"size:255;not null;" json:"id_membership"`
	RenovationActive         int       `gorm:"size:255;not null;" json:"renovation_active"`
	Videos                   []MyVideo `gorm:"references:id; foreignKey:id_user"`
	Description              string    `gorm:"size:255" json:"description"`

}

func (PublicUser) TableName() string {
	return "users"
}

func (User) TableName() string {
	return "users"
}

func (updatedUser *User) UpdateProfileImageUser() (User, error) {
	dbConn := Pool

	// Fetch the existing user from the database
	var user User
	if err := dbConn.First(&user, updatedUser.ID).Error; err != nil {
		return user, err
	}

	oldImage := user.ProfileImage
	// Update the user fields with the new values
	user.ProfileImage = updatedUser.ProfileImage
	// user.Email = updatedUser.Email
	// Update other user fields as needed...

	// Save the changes to the database

	if err := dbConn.Save(&user).Error; err != nil {
		return user, err
	}
	pathOldImage := os.Getenv("MY_URL")

	deleteImage(pathOldImage + oldImage)
	// if err != nil {
	// 	return user, err
	// }

	return user, nil
}

func (updatedUser *User) UpdateCoverImageUser() (User, error) {
	dbConn := Pool

	// Fetch the existing user from the database
	var user User
	if err := dbConn.First(&user, updatedUser.ID).Error; err != nil {
		return user, err
	}

	oldImage := user.Cover_image
	// Update the user fields with the new values
	user.Cover_image = updatedUser.Cover_image
	// user.Email = updatedUser.Email
	// Update other user fields as needed...

	// Save the changes to the database

	if err := dbConn.Save(&user).Error; err != nil {
		return user, err
	}
	pathOldImage := os.Getenv("MY_URL")

	deleteImage(pathOldImage + oldImage)
	// if err != nil {
	// 	return user, err
	// }

	return user, nil
}

func deleteImage(imagePath string) error {
	err := os.Remove(imagePath)
	if err != nil {
		// Handle the error
		return err
	}

	return nil
}

func GetUserByID(uid uint) (User, error) {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&User{}).First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}


func GetUserByIDWithVideos(uid uint) (PublicUser, error) {

	var u PublicUser
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&PublicUser{}).First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	// u.PrepareGive()

	return u, nil

}

func GetUserByPhone(phone string) (User, error) {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&User{}).Where("phone = ?", phone).Find(&u).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func UsernameExists(username string) bool {

	var u PublicUser
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	result := dbConn.Model(&PublicUser{}).Where("username = ?", username).First(&u)

	if result.Error == nil {
		// The username already exists
		return true
	}
	// u.PrepareGive()

	return false

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() (*User, error) {

	var err error
	dbConn := Pool

	err = dbConn.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser() (*User, error) {

	var err error
	dbConn := Pool

	err = dbConn.Save(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func GetMyVideos(id_user int, page int) ([]FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid []FeedVideo
	// userID, _ := token.ExtractTokenID(c)

	// Get page number and page size from query parameters
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 12

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	// err = dbConn.Unscoped().Find(&vid).Error
	err = 
	dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		// Where("sale_type_id = ? && is_vip = ? && sale_category_id = ? ", sale_type, isvip,categoryId).
		// Where("sale_category_id = ? && is_vip = ?", sale_type, isvip).
		Where("videos.id_user = ?", id_user).
		Limit(pageSize).
		Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error
	// dbConn.Model(&MyVideo{}).Where("id_user = ?", id_user).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	if err != nil {
		return vid, err
	}
	return vid, nil

}

func GetMyFavoritesVideos(id_user int, page int) ([]FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid []FeedVideo
	// userID, _ := token.ExtractTokenID(c)

	// Get page number and page size from query parameters
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 12

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	// err = dbConn.Unscoped().Find(&vid).Error
	err = dbConn.Model(&MyVideo{}).
		Joins("INNER JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video").
		Where("users_videos_favorites.id_user = ?", id_user).
		Limit(pageSize).Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error
	// var favoriteVideos []Video
	// err = db.
	// 	Table("videos").
	// 	Select("videos.*").
	// 	Joins("INNER JOIN videosfavorites ON videos.id = videosfavorites.video_id").
	// 	Where("videosfavorites.user_id = ?", userID).
	// 	Find(&favoriteVideos).Error

	// if err != nil {
	// 	log.Fatal(err)
	// }
	if err != nil {
		return vid, err
	}
	return vid, nil

}


func DeleteUserVideo(id_video int, id_user int) error {
	var err error
	dbConn := Pool
	var vid Video

	
	if err = dbConn.Where("id_user = ? && id = ?", id_user,id_video).Find(&vid).Error; err != nil {
		return err
	}
	pathOldImage := os.Getenv("MY_URL")

	deleteImage(pathOldImage + vid.Image_cover )
	deleteImage(pathOldImage + vid.Video_url )

	if err = dbConn.Where("id_user = ? && id = ?", id_user,id_video).Delete(&vid).Error; err != nil {
		return err
	}

	return nil
}

