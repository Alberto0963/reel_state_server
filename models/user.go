package models

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"os"
	"path/filepath"

	// "path/filepath"

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
	Link                     string    `gorm:"size:255" json:"link"`
	Ventas                   int       `gorm:"size:255" json:"ventas"`
	TotalReviews             int       `gorm:"size:255" json:"total_reviews"`
	AverageRating            string    `gorm:"size:255" json:"average_rating"`
	VideoCount               int       `gorm:"size:255" json:"video_count"`
	MedalType                int       `gorm:"size:255" json:"medal_type"`
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
	Link                     string    `gorm:"size:255" json:"link"`

	MedalType                int       `gorm:"size:255" json:"medal_type"`
	DeviceToken              string    `gorm:"size:255" json:"device_token"`

	Email string `gorm:"" json:"email"`
}

type UserDB struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID                       uint           `gorm:"not null;unique" json:"id"`
	Phone                    sql.NullString `gorm:"size:13;unique" json:"phone"`
	Username                 string         `gorm:"size:255;not null;unique" json:"username"`
	Password                 string         `gorm:"size:100;not null;" json:"password"`
	ProfileImage             string         `gorm:"size:255;not null;" json:"profileImage"`
	ExpirationMembershipDate time.Time      `gorm:"size:255;" json:"expiration_membership_date"`
	IdMembership             int            `gorm:"size:255;not null;" json:"id_membership"`
	RenovationActive         int            `gorm:"size:255;not null;" json:"renovation_active"`
	// Cover_image              sql.NullString `gorm:"size:255;" json:"cover_image"`
	Description string         `gorm:"size:255" json:"description"`
	Email       sql.NullString `gorm:"" json:"email"`
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

func (UserDB) TableName() string {
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

func (updatedUser *UserUpdate) UpdateCoverImageUser() (User, error) {
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

	if err := dbConn.Where("id = ?", uid).Find(&u).
		Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

	return u, nil

}

func GetUserDeviceToken(likedUserIDs []int) ([]string, error) {

	// var u UserUpdate
	// Obtain a connection from the pool
	dbConn := Pool

	var tokens []string
	if err := dbConn.Table("users").Select("device_token").Where("id IN (?) AND device_token IS NOT NULL", likedUserIDs).Find(&tokens).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tokens"})
		return tokens, err
	}

	return tokens, nil

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

func GetUserByPhoneUpdate(phone string) (UserDB, error) {

	var u UserDB
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()

	if err := dbConn.Model(&User{}).Where("phone = ?", phone).Find(&u).Error; err != nil {
		return u, errors.New("User not found! ")
	}

	// u.PrepareGive()

	return u, nil

}

func DeleteUserByID(uid uint) error {

	var u User
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	// if err := dbConn.Model(&User{}).
	// Define the subquery as a raw SQL string
	// videosCountSubquery := "(SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v"

	// if err := dbConn.Model(&User{}).Where("id = ?", uid).Delete(&u).Error; err != nil {
	// 	return u, errors.New("User not found! ")
	// }

	// u.PrepareGive()

	// return u, nil
	// Check if user exists
	if err := dbConn.First(&u, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return err
		}
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return err
	}

	// Delete the user
	if err := dbConn.Delete(&u).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return err
	}

	return nil
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

func (u *UserDB) SaveUser() (*UserDB, error) {

	var err error
	dbConn := Pool

	err = dbConn.Create(&u).Error
	if err != nil {
		return &UserDB{}, err
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

func (u *UserDB) UpdateUser() (*UserDB, error) {

	var err error
	dbConn := Pool

	err = dbConn.Save(&u).Error
	if err != nil {
		return &UserDB{}, err
	}
	return u, nil
}

// func (user *UserUpdate) UpdatePassword(newPassword string) error {
// 	// var err error
// 	dbConn := Pool

// 	return dbConn.Model(&user).Update("password", newPassword).Error
// }

func (u *UserDB) BeforeSave(tx *gorm.DB) error {

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


// func (u *UserUpdate) BeforeSave(tx *gorm.DB) error {

// 	//turn password into hash
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = string(hashedPassword)

// 	//remove spaces in username
// 	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

// 	return nil

// }

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

// func DeleteUserVideo(id_video int, id_user int) error {
// 	var err error
// 	dbConn := Pool
// 	var vid MyVideo

// 	if err = dbConn.Where("id_user = ? && id = ?", id_user, id_video).Find(&vid).Error; err != nil {
// 		return err
// 	}
// 	pathOldImage := os.Getenv("MY_URL")

// 	deleteImage(pathOldImage + vid.Image_cover)
// 	deleteImage(pathOldImage + vid.Video_url)

// 	if err = dbConn.Where("id_user = ? && id = ?", id_user, id_video).Delete(&vid).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func DeleteUserVideo(id_video int, id_user int) error {
	dbConn := Pool
	var vid Video

	// Buscar el video
	err := dbConn.Where("id_user = ? AND id = ?", id_user, id_video).First(&vid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("video no encontrado: %w", err)
		}
		return fmt.Errorf("error al buscar el video: %w", err)
	}

	// Construir rutas de archivos a eliminar
	baseURL := os.Getenv("MY_URL")
	if baseURL == "" {
		return fmt.Errorf("la variable de entorno MY_URL no está configurada")
	}

	imagePath := filepath.Join(baseURL, vid.Image_cover)
	videoPath := filepath.Join(baseURL, vid.Video_url)

	// Eliminar archivos asociados
	if err := deleteImage(imagePath); err != nil {
		return fmt.Errorf("error al eliminar la imagen: %w", err)
	}
	if err := deleteImage(videoPath); err != nil {
		return fmt.Errorf("error al eliminar el video: %w", err)
	}

	// Eliminar el registro de la base de datos
	err = dbConn.Where("id_user = ? AND id = ?", id_user, id_video).Delete(&vid).Error
	if err != nil {
		return fmt.Errorf("error al eliminar el registro del video: %w", err)
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

func SepararNombreCompleto(nombreCompleto string) (string, string) {
	// Dividimos el string por espacios
	partes := strings.Split(nombreCompleto, " ")

	// Si hay más de un elemento, el primer elemento es el nombre y el resto son los apellidos
	if len(partes) >= 2 {
		nombre := partes[0]
		apellidos := strings.Join(partes[1:], " ") // Reunimos los apellidos en un solo string
		return nombre, apellidos
	}

	// Si no hay suficientes partes, devolvemos el string completo como nombre y los apellidos vacíos
	return nombreCompleto, ""
}

func ConvertToUser(dbUser UserDB) UserUpdate {
	var phone string
	if dbUser.Phone.Valid {
		phone = dbUser.Phone.String // Si es un valor válido
	} else {
		phone = "" // Si es NULL, lo convertimos a una cadena vacía
	}

	return UserUpdate{
		ID:    dbUser.ID,
		Phone: phone, // Asignamos el valor convertido
	}
}

func Setnull(phoneNumber string) sql.NullString {
	var phone sql.NullString

	if phoneNumber == "" {
		phone = sql.NullString{
			String: "",    // El valor aquí no importa cuando Valid es false
			Valid:  false, // Se trata como NULL en la base de datos
		}
	} else {
		phone = sql.NullString{
			String: phoneNumber, // Asigna el valor del número de teléfono
			Valid:  true,        // Indica que el valor es válido
		}
	}

	return phone
}
