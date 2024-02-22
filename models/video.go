package models

import (
	// "encoding/json"
	"math/rand"
	"reelState/utils/location"
	"time"
	// "gorm.io/gorm/utils"
	// "gorm.io/gorm"
)

type Video struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id              int    `gorm:"not null;unique" json:"id"`
	Video_url       string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location        string `gorm:"size:100;not null;" json:"location"`
	Area            string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`

	Price            string     `gorm:"size:255;not null;" json:"price"`
	Id_user          uint       `gorm:"size:255;not null;" json:"id_user"`
	User             PublicUser `gorm:"references:Id; foreignKey:Id_user"`
	Sale_type_id     int        `gorm:"size:255;not null;" json:"sale_type_id"`
	SaleType         Type       `gorm:"references:id; foreignKey:sale_type_id"`
	Sale_category_id int        `gorm:"size:255;not null;" json:"sale_category_id"`
	SaleCategory     Category   `gorm:"references:id; foreignKey:sale_category_id"`
	Image_cover      string     `gorm:"size:255;not null;" json:"image_cover"`
	Latitude         float64    `gorm:"size:255;not null;" json:"latitude"`
	Longitude        float64    `gorm:"size:255;not null;" json:"longitude"`
	Type             float64    `gorm:"size:255;not null;" json:"type"`
}

type FeedVideo struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id              int    `gorm:"not null;unique" json:"id"`
	Video_url       string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location        string `gorm:"size:100;not null;" json:"location"`
	Area            string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`

	Price            string     `gorm:"size:255;not null;" json:"price"`
	Id_user          uint       `gorm:"size:255;not null;" json:"id_user"`
	User             PublicUser `gorm:"references:Id; foreignKey:Id_user"`
	Sale_type_id     int        `gorm:"size:255;not null;" json:"sale_type_id"`
	SaleType         Type       `gorm:"references:id; foreignKey:sale_type_id"`
	Sale_category_id int        `gorm:"size:255;not null;" json:"sale_category_id"`
	SaleCategory     Category   `gorm:"references:id; foreignKey:sale_category_id"`
	Image_cover      string     `gorm:"size:255;not null;" json:"image_cover"`
	Latitude         float64    `gorm:"size:255;not null;" json:"latitude"`
	Longitude        float64    `gorm:"size:255;not null;" json:"longitude"`
	Is_favorite      string     `gorm:"size:255;not null;" json:"is_favorite"`
	Type             float64    `gorm:"size:255;not null;" json:"type"`
}

type MyVideo struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id              int    `gorm:"not null;unique" json:"id"`
	Video_url       string `gorm:"size:13;not null;unique" json:"video_url"`
	Description     string `gorm:"size:255;not null;unique" json:"description"`
	Location        string `gorm:"size:100;not null;" json:"location"`
	Area            string `gorm:"size:255;not null;" json:"area"`
	Property_number string `gorm:"size:255;not null;" json:"property_number"`

	Price            string   `gorm:"size:255;not null;" json:"price"`
	Id_user          uint     `gorm:"size:255;not null;" json:"id_user"`
	Sale_type_id     int      `gorm:"size:255;not null;" json:"sale_type_id"`
	SaleType         Type     `gorm:"references:id; foreignKey:sale_type_id"`
	Sale_category_id int      `gorm:"size:255;not null;" json:"sale_category_id"`
	SaleCategory     Category `gorm:"references:id; foreignKey:sale_category_id"`
	Image_cover      string   `gorm:"size:255;not null;" json:"image_cover"`
	Type             float64  `gorm:"size:255;not null;" json:"type"`
}

func (MyVideo) TableName() string {
	return "videos"
}

func GenerateRandomName() string {
	rand.NewSource(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10
	name := make([]byte, length)
	for i := 0; i < length; i++ {
		name[i] = chars[rand.Intn(len(chars))]
	}
	return "reel_state." + string(name)
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

func (v *Video) EditVideo() (*Video, error) {
	var err error
	dbConn := Pool
	var vid Video
	err = dbConn.Where("id_user = ? && id = ?", v.Id_user, v.Id).Find(&vid).Error
	if err != nil {
		return &vid, err
	}

	vid.Description = v.Description
	vid.Location = v.Location
	vid.Area = v.Area
	vid.Property_number = v.Property_number
	vid.Price = v.Price
	// uidserID, _ := v.ExtractTokenID(c)
	// vid.Id_user = userID
	vid.Latitude = v.Latitude
	vid.Longitude = v.Longitude

	vid.Sale_type_id = v.Sale_type_id

	vid.Sale_category_id = v.Sale_category_id

	err = dbConn.Save(&vid).Error
	if err != nil {
		return &Video{}, err
	}
	return v, nil

}

func FetchAllVideos(id_user int, sale_type int, typeV int, page int) ([]FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid []FeedVideo
	var ads []FeedVideo
	var user User
	typeUser := 0

	user,err = GetUserByID(uint(id_user))
	if(err == nil ){
		typeUser = user.IdMembership

	}

	pageSize := 10
	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize
	// err = dbConn.Unscoped().Find(&vid).Error

	// err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	// if err != nil {
	// 	return vid, err
	// }

	// var videos []Video
	result := dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("sale_type_id = ? && type = ?", sale_type, typeV).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().
		Find(&vid).Error

	if result != nil {
		// http.Error(w, "Database error", http.StatusInternalServerError)
		return vid, err
	}

	// responseJSON, err := json.Marshal(videos)
	// if err != nil {

	// 	return vid,err
	// }

	rand.NewSource(time.Now().UnixNano())

	rand.Shuffle(len(vid), func(i, j int) { vid[i], vid[j] = vid[j], vid[i] })

	// Process videos with ads
	if (typeUser == 0 || typeUser == 1 || typeUser == 8 || typeUser == 2 || typeUser == 3) && len(vid) > 0{
		// var tempVideos []FeedVideo
		adsIndex := 0
		pageSize = 2
		offset := (page - 1) * pageSize

		err = dbConn.Table("videos").
			Select("videos.*").

			// Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
			Where("type = ?", 3).Limit(pageSize).Offset(offset).Unscoped().
			Find(&ads).Error

		if len(ads) != 0{
			vid = append(vid, ads[adsIndex])
		}
		// for i, video := range vid {
		// 	// Add the new element to the end of the list
		// 	tempVideos = append(tempVideos, video)
		// 	if (i+1)%5 == 0 && video.Type != 2 {
		// 		// Skip this video or handle as needed
		// 		tempVideos = append(tempVideos, ads[adsIndex])
		// 		adsIndex++
		// 		// continue
		// 	}
		// }
		// vid = tempVideos
	}

	return vid, nil

}

func SearchVideos(search string, page int, id_user int) ([]FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid []FeedVideo
	pageSize := 12

	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize

	result := dbConn.Model(&Video{}).
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("location like ? type = 1", "%"+search+"%").
		Limit(pageSize).
		Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error

	// Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
	// Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
	// Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).
	// Limit(pageSize).Offset(offset).
	// Preload("SaleType").
	// Preload("SaleCategory").
	// Preload("User").
	// Unscoped().
	// Find(&vid).Error

	if result != nil {
		// http.Error(w, "Database error", http.StatusInternalServerError)
		return vid, err
	}

	rand.NewSource(time.Now().UnixNano())

	rand.Shuffle(len(vid), func(i, j int) { vid[i], vid[j] = vid[j], vid[i] })
	return vid, nil

}

func FetchAllCategoryVideos(id_user int, sale_type int, typeV int, categoryId, page int) ([]FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid []FeedVideo
	var ads []FeedVideo

	pageSize := 12
	var user User
	typeUser := 0
	// Calculate the offset based on the page number and page size
	offset := (page - 1) * pageSize

	user,err = GetUserByID(uint(id_user))
	if(err == nil ){
		typeUser = user.IdMembership

	}
	// err = dbConn.Unscoped().Find(&vid).Error

	// err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	// if err != nil {
	// 	return vid, err
	// }

	// var videos []Video
	result := dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("sale_type_id = ? && type = ? && sale_category_id = ? ", sale_type, typeV, categoryId).
		// Where("sale_category_id = ? && is_vip = ?", sale_type, isvip).
		Limit(pageSize).
		Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error

	if result != nil {
		// http.Error(w, "Database error", http.StatusInternalServerError)
		return vid, err
	}

	// responseJSON, err := json.Marshal(videos)
	// if err != nil {

	// 	return vid,err
	// }

	rand.NewSource(time.Now().UnixNano())

	rand.Shuffle(len(vid), func(i, j int) { vid[i], vid[j] = vid[j], vid[i] })

		// Process videos with ads
		if (typeUser == 0 || typeUser == 1 || typeUser == 8 || typeUser == 2 || typeUser == 3) && len(vid) > 0{
			// var tempVideos []FeedVideo
			adsIndex := 0
			pageSize = 2
			offset := (page - 1) * pageSize
	
			err = dbConn.Table("videos").
				Select("videos.*").
	
				// Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
				Where("type = ?", 3).Limit(pageSize).Offset(offset).Unscoped().
				Find(&ads).Error
	
			if len(ads) != 0{
				vid = append(vid, ads[adsIndex])
			}
			// for i, video := range vid {
			// 	// Add the new element to the end of the list
			// 	tempVideos = append(tempVideos, video)
			// 	if (i+1)%5 == 0 && video.Type != 2 {
			// 		// Skip this video or handle as needed
			// 		tempVideos = append(tempVideos, ads[adsIndex])
			// 		adsIndex++
			// 		// continue
			// 	}
			// }
			// vid = tempVideos
		}
	return vid, nil

}

func GetPlacesAroundLocation(centerLat, centerLon float64, maxDistance float64, id_user int) ([]Video, error) {
	var err error
	dbConn := Pool
	var vid []Video
	// pageSize := 12
	var nearbyPlaces []Video

	// Calculate the offset based on the page number and page size
	// offset := (page - 1) * pageSize

	err = dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		// Where("sale_type_id = ? && is_vip = ? && sale_category_id = ? ", sale_type, isvip,categoryId).
		// Where("sale_category_id = ? && is_vip = ?", sale_type, isvip).
		// Limit(pageSize).
		// Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error
	// .Where("sale_type_id = ? && is_vip = ?", sale_type, isvip)
	if err != nil {
		return vid, err
	}

	// // query := fmt.Sprintf("SELECT name, latitude, longitude FROM places")
	// rows, err := db.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var places []string

	for _, place := range vid {
		distance := location.HaversineDistance(centerLat, centerLon, place.Latitude, place.Longitude)
		if distance <= maxDistance {
			nearbyPlaces = append(nearbyPlaces, place)
		}
	}

	return nearbyPlaces, nil
}
