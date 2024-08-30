package models

import (
	"errors"
	"fmt"
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
	CanUpload                bool      `gorm:"" json:"can_upload"`
	Email                    string    `gorm:"" json:"email"`
}

type UserUpdate struct {
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
	Email                    string    `gorm:"" json:"email"`
}

type PublicUser struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id                       uint      `gorm:"not null;unique" json:"id"`
	Phone                    string    `gorm:"size:13;not null;unique" json:"phone"`
	Username                 string    `gorm:"size:255;not null;unique" json:"username"`
	ProfileImage             string    `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time `gorm:"size:255;" json:"expiration_membership_date"`
	Id_Membership            int       `gorm:"size:255;not null;" json:"id_membership"`
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

func (UserUpdate) TableName() string {
	return "users"
}

func (updatedUser *UserUpdate) UpdateProfileImageUser() (UserUpdate, error) {
	dbConn := Pool

	// Fetch the existing user from the database
	var user UserUpdate
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
	// if err := dbConn.Model(&User{}).
	// Define the subquery as a raw SQL string
	// videosCountSubquery := "(SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v"

	if err := dbConn.Table("view_user_upload_status").
		// Joins("LEFT JOIN memberships ON memberships.id = users.id_membership").
		// Joins("LEFT JOIN (SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v ON v.id_user = users.id").
		Where("id = ?", uid).Find(&u).
		Error; err != nil {
		return u, errors.New("User not found! ")
	}

	u.PrepareGive()

	return u, nil

}

func GetUserByIdToUpdate(uid uint) (UserUpdate, error) {

	var u UserUpdate
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	// if err := dbConn.Model(&User{}).
	// Define the subquery as a raw SQL string
	// videosCountSubquery := "(SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v"

	if err := dbConn.Table("view_user_upload_status").
		// Joins("LEFT JOIN memberships ON memberships.id = users.id_membership").
		// Joins("LEFT JOIN (SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v ON v.id_user = users.id").
		Where("id = ?", uid).Find(&u).
		Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

	return u, nil

}

func GetUserByIDWithVideos(uid uint) (PublicUser, error) {

	var u PublicUser
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&PublicUser{}).Preload("Videos").First(&u, uid).Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

	return u, nil

}

func GetUserByPhone(phone string) (UserUpdate, error) {

	var u UserUpdate
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&User{}).Where("phone = ?", phone).Find(&u).Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

	return u, nil

}

func GetUserByEmail(email string) (User, error) {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	result := dbConn.Where("email = ?", email).First(&u)

	// if err := dbConn.Model(&User{}).Where("email = ?", email).Find(&u).Error; err != nil {
	// 	return u, errors.New("User not found! ")
	// }
	// Check the type of error returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// User not found is not a database error but a "normal" flow error
		return User{}, errors.New("User not found! ")
	} else if result.Error != nil {
		// Some other error occurred during the query execution
		return User{}, result.Error
	}

	// u.PrepareGive()

	return u, nil

}

func SearchProfile(username string) ([]PublicUser, error) {

	var u []PublicUser
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&PublicUser{}).Where("username like ?", "%"+username+"%").Find(&u).Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

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

func (u *UserUpdate) SaveUser() (*UserUpdate, error) {

	var err error
	dbConn := Pool

	err = dbConn.Create(&u).Error
	if err != nil {
		return &UserUpdate{}, err
	}
	return u, nil
}

func (u *UserUpdate) UpdateUser() (*UserUpdate, error) {

	var err error
	dbConn := Pool

	err = dbConn.Save(&u).Error
	if err != nil {
		return &UserUpdate{}, err
	}
	return u, nil
}

func (user *UserUpdate) UpdatePassword( newPassword string) error {
	// var err error
	dbConn := Pool

	return dbConn.Model(&user).Update("password", newPassword).Error
}

func (u *UserUpdate) BeforeSave(tx *gorm.DB) error {

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

func GetMyVideos(id_user int, page int, typeVideo int) ([]FeedVideo, error) {
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
			Where("videos.type = ?", typeVideo).
			Order("videos.created_at DESC").
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

	if err = dbConn.Where("id_user = ? && id = ?", id_user, id_video).Find(&vid).Error; err != nil {
		return err
	}
	pathOldImage := os.Getenv("MY_URL")

	deleteImage(pathOldImage + vid.Image_cover)
	deleteImage(pathOldImage + vid.Video_url)

	if err = dbConn.Where("id_user = ? && id = ?", id_user, id_video).Delete(&vid).Error; err != nil {
		return err
	}

	return nil
}

func ValidateUserType(id_user int) error {
	u, err := GetUserByIDWithVideos(uint(id_user))
	if err != nil {
		// Return the error as is if it's from GetUserByIDWithVideos
		return err
	}

	countVideos := len(u.Videos)

	switch {
	case countVideos >= 1 && (u.Id_Membership == 1 || u.Id_Membership == 8):
		// User has basic membership but trying to upload more than allowed
		return fmt.Errorf("no puedes publicar mas videos, límite alcanzado para tu membresía")

	case countVideos >= 20 && (u.Id_Membership == 2 || u.Id_Membership == 3):
		// User has a mid-tier membership but is trying to upload more than allowed
		return fmt.Errorf("no puedes publicar mas videos, límite alcanzado para tu membresía")

	case countVideos >= 50 && (u.Id_Membership == 4 || u.Id_Membership == 5):
		// User has a high-tier membership but is trying to upload more than allowed
		return fmt.Errorf("no puedes publicar mas videos, límite alcanzado para tu membresía")
	}

	// If none of the conditions are met, then there's no error
	return nil
}

func Getfollowers(id_profile int) (int, error) {
	var err error
	dbConn := Pool
	var foll int //[]likes

	err =
		dbConn.Table("likes").
			Select("count(*) AS followers").Where("id_profile = ?", id_profile).
			Find(&foll).Error
	// dbConn.Model(&MyVideo{}).Where("id_user = ?", id_user).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	if err != nil {
		return foll, err
	}

	return foll, nil

}

func GetProfilefollowers(id_profile int, userid int) (Likes, error) {
	var err error
	dbConn := Pool
	var foll Likes

	err =
		dbConn.Where("id_profile = ? && id_user = ?", id_profile, userid).
			Find(&foll).Error
	// dbConn.Model(&MyVideo{}).Where("id_user = ?", id_user).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	if err != nil {
		return foll, err
	}

	return foll, nil

}

func Imfollower(id_profile int, id_user int) (bool, error) {
	var err error
	dbConn := Pool
	// var foll bool//[]likes
	var count int

	err =
		dbConn.Table("likes").
			Select("count(*)").Where("id_profile = ? && id_user = ?", id_profile, id_user).
			Find(&count).Error
	// dbConn.Model(&MyVideo{}).Where("id_user = ?", id_user).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	}

	return false, nil

}

func SaveLikeProfile(like *Likes) (*Likes, error) {
	var err error
	dbConn := Pool

	dbConn.Save(&like)

	if err != nil {
		return &Likes{}, err
	}
	return like, nil
}

func UpdateLikeProfile(like *Likes) (*Likes, error) {
	var err error
	dbConn := Pool

	dbConn.Delete(&like)

	if err != nil {
		return &Likes{}, err
	}
	return like, nil
}

func FindByPhone(phone string) error {
	var err error
	dbConn := Pool
	var v UserUpdate

	err = dbConn.Where("phone = ?", phone).First(&v).Error
	return err
}

