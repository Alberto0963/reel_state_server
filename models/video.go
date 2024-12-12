package models

import (
	// "encoding/json"
	"math/rand"
	"reelState/utils/location"
	"regexp"
	"strings"

	// "strings"
	"time"

	"gorm.io/gorm"
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
	Type             int        `gorm:"size:255;not null;" json:"type"`
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

func SetAvailable(idvideo int, typev int) (Video, error) {
	var err error
	dbConn := Pool
	var vid Video
	err = dbConn.Where("id = ?", idvideo).Find(&vid).Error
	if err != nil {
		return vid, err
	}

	vid.Type = typev

	err = dbConn.Save(vid).Error
	if err != nil {
		return Video{}, err
	}
	return vid, nil
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

func GetVideo(id int, id_user int) (FeedVideo, error) {
	var err error
	dbConn := Pool
	var vid FeedVideo
	// var ads []FeedVideo
	// var user User
	// typeUser := 0

	// err = dbConn.Unscoped().Find(&vid).Error

	// err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
	// if err != nil {
	// 	return vid, err
	// }

	// var videos []Video
	result := dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("videos.id = ?", id).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error

	if result != nil {
		// http.Error(w, "Database error", http.StatusInternalServerError)
		return vid, err
	}

	return vid, nil

}

func FetchAllVideos(id_user int, sale_type int, typeV int, page int, idvideo *int) ([]FeedVideo, error) {
	var (
		err      error
		vid      []FeedVideo
		ads      []FeedVideo
		typeUser int
		pageSize = 10
		offset   = (page - 1) * pageSize
		dbConn   = Pool
	)

	// Obtener el usuario y determinar su tipo
	user, err := GetUserByID(uint(id_user))
	if err == nil {
		typeUser = user.IdMembership
	}

	// Obtener los videos normales
	err = dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("sale_type_id = ? AND type = ?", sale_type, typeV).
		Limit(pageSize).
		Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error

	if err != nil {
		return vid, err
	}

	// Si se proporciona un idvideo, buscarlo y colocarlo en la posición 0
	if idvideo != nil {
		var specialVideo FeedVideo
		err = dbConn.Table("videos").
			Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
			Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
			Where("videos.id = ?", *idvideo).
			Preload("SaleType").
			Preload("SaleCategory").
			Preload("User").
			Unscoped().
			First(&specialVideo).Error

		if err == nil {
			vid = append([]FeedVideo{specialVideo}, vid...) // Coloca el video especial en la primera posición
		}
	}

	// Si se proporcionó un idvideo, asegúrate de que no esté duplicado
	if idvideo != nil && len(vid) > 1 && vid[1].Id == *idvideo {
		vid = vid[1:] // Elimina el video duplicado en la lista (mantén el que está en la posición 0)
	}

	// Mezclar los videos aleatoriamente, excepto el primero si es el video especial
	if len(vid) > 1 {
		rand.NewSource(time.Now().UnixNano())
		rand.Shuffle(len(vid)-1, func(i, j int) { vid[i+1], vid[j+1] = vid[j+1], vid[i+1] })
	}

	// Procesar los videos con anuncios
	if (typeUser == 0 || typeUser == 100000 || typeUser == 100004) && len(vid) > 0 {
		pageSize = 2
		offset = (page - 1) * pageSize
		err = dbConn.Table("videos").
			Select("videos.*").
			Where("type = ?", 3).
			Limit(pageSize).
			Offset(offset).
			Unscoped().
			Find(&ads).Error

		if err == nil && len(ads) > 0 {
			vid = append(vid, ads[0])
		}
	}

	return vid, nil
}

// func FetchAllVideos(id_user int, sale_type int, typeV int, page int) ([]FeedVideo, error) {
// 	var err error
// 	dbConn := Pool
// 	var vid []FeedVideo
// 	var ads []FeedVideo
// 	var user User
// 	typeUser := 0

// 	user, err = GetUserByID(uint(id_user))
// 	if err == nil {
// 		typeUser = user.IdMembership

// 	}

// 	pageSize := 10
// 	// Calculate the offset based on the page number and page size
// 	offset := (page - 1) * pageSize
// 	// err = dbConn.Unscoped().Find(&vid).Error

// 	// err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
// 	// if err != nil {
// 	// 	return vid, err
// 	// }

// 	// var videos []Video
// 	result := dbConn.Table("videos").
// 		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 		Where("sale_type_id = ? && type = ?", sale_type, typeV).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().
// 		Find(&vid).Error

// 	if result != nil {
// 		// http.Error(w, "Database error", http.StatusInternalServerError)
// 		return vid, err
// 	}

// 	// responseJSON, err := json.Marshal(videos)
// 	// if err != nil {

// 	// 	return vid,err
// 	// }

// 	rand.NewSource(time.Now().UnixNano())

// 	rand.Shuffle(len(vid), func(i, j int) { vid[i], vid[j] = vid[j], vid[i] })

// 	// Process videos with ads
// 	if (typeUser == 0 || typeUser == 100000 || typeUser == 100004) && len(vid) > 0 {
// 		// var tempVideos []FeedVideo
// 		adsIndex := 0
// 		pageSize = 2
// 		offset := (page - 1) * pageSize

// 		err = dbConn.Table("videos").
// 			Select("videos.*").

// 			// Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 			Where("type = ?", 3).Limit(pageSize).Offset(offset).Unscoped().
// 			Find(&ads).Error

// 		if len(ads) != 0 {
// 			vid = append(vid, ads[adsIndex])
// 		}
// 		// for i, video := range vid {
// 		// 	// Add the new element to the end of the list
// 		// 	tempVideos = append(tempVideos, video)
// 		// 	if (i+1)%5 == 0 && video.Type != 2 {
// 		// 		// Skip this video or handle as needed
// 		// 		tempVideos = append(tempVideos, ads[adsIndex])
// 		// 		adsIndex++
// 		// 		// continue
// 		// 	}
// 		// }
// 		// vid = tempVideos
// 	}

// 	return vid, nil

// }

// func SearchVideos(search string, page int, id_user int) ([]FeedVideo, error) {
// 	var err error
// 	dbConn := Pool
// 	var vid []FeedVideo
// 	var videos []FeedVideo

// 	pageSize := 12

// 	// Calculate the offset based on the page number and page size
// 	offset := (page - 1) * pageSize
// 	locations := strings.Fields(search)

// 	result := dbConn.Model(&Video{}).
// 		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 		Where("description like ? && type = 1", "%"+search+"%").
// 		Limit(pageSize).
// 		Offset(offset).
// 		Preload("SaleType").
// 		Preload("SaleCategory").
// 		Preload("User").
// 		Unscoped().
// 		Find(&vid).Error
// 	if result != nil {
// 		// http.Error(w, "Database error", http.StatusInternalServerError)
// 		return videos, err
// 	}
// 	videos = append(videos, vid...)

// 	price := findPrices(search)
// 	if price != nil {

// 		result = dbConn.Model(&Video{}).
// 			Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 			Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 			Where("price IN (?) && type = 1", price).
// 			Limit(pageSize).
// 			Offset(offset).
// 			Preload("SaleType").
// 			Preload("SaleCategory").
// 			Preload("User").
// 			Unscoped().
// 			Find(&vid).Error
// 		if result != nil {
// 			// http.Error(w, "Database error", http.StatusInternalServerError)
// 			return videos, err
// 		}
// 		videos = append(videos, vid...)
// 	}

// 	// Dynamically build the WHERE clause to use LIKE for each keyword
// 	// Dynamically build the query with LIKE conditions for each keyword
// 	for _, keyword := range locations {
// 		likePattern := "%" + keyword + "%"
// 		dbConn = dbConn.Or(" type = 1 && location LIKE ?", likePattern)
// 	}

// 	result = dbConn.Model(&Video{}).
// 		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 		// Where("?", whereClause).
// 		Limit(pageSize).
// 		Offset(offset).
// 		Preload("SaleType").
// 		Preload("SaleCategory").
// 		Preload("User").
// 		Unscoped().
// 		Find(&vid).Error
// 	if result != nil {
// 		// http.Error(w, "Database error", http.StatusInternalServerError)
// 		return videos, err
// 	}

// 	videos = append(videos, vid...)

// 	rand.NewSource(time.Now().UnixNano())

// 	rand.Shuffle(len(videos), func(i, j int) { videos[i], videos[j] = videos[j], videos[i] })
// 	return videos, nil

// }
func SearchVideos(search string, page int, id_user int) ([]FeedVideo, error) {
	var (
		dbConn    = Pool
		videos    []FeedVideo
		pageSize  = 12
		offset    = (page - 1) * pageSize
		locations = strings.Fields(search)
		err       error
	)

	// Helper function to fetch videos based on a condition
	fetchVideos := func(condition string, args ...interface{}) ([]FeedVideo, error) {
		var vid []FeedVideo
		result := dbConn.Model(&Video{}).
			Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
			Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
			Where(condition, args...).
			Limit(pageSize).
			Offset(offset).
			Preload("SaleType").
			Preload("SaleCategory").
			Preload("User").
			Unscoped().
			Find(&vid)
		return vid, result.Error
	}

	// Search by description
	vid, err := fetchVideos("description LIKE ? AND type = 1", "%"+search+"%")
	if err != nil {
		return videos, err
	}
	videos = append(videos, vid...)

	// // Search by price if applicable
	// price := findPrices(search)
	// if price != nil {
	// 	vid, err = fetchVideos("price IN (?) AND type = 1", price)
	// 	if err != nil {
	// 		return videos, err
	// 	}
	// 	videos = append(videos, vid...)
	// }
	// Search by price if applicable
	pricePatterns := findPricePatterns(search)
	if len(pricePatterns) > 0 {
		for _, pattern := range pricePatterns {
			vid, err = fetchVideos("price LIKE ? AND type = 1", "%"+pattern+"%")
			if err != nil {
				return videos, err
			}
			videos = append(videos, vid...)
		}
	}

	// Search by location keywords
	for _, keyword := range locations {
		likePattern := "%" + keyword + "%"
		vid, err = fetchVideos("location LIKE ? AND type = 1", likePattern)
		if err != nil {
			return videos, err
		}
		videos = append(videos, vid...)
	}

	// Remove duplicates (optional) and shuffle the results
	videos = removeDuplicates(videos)
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(videos), func(i, j int) { videos[i], videos[j] = videos[j], videos[i] })

	return videos, nil
}

// Helper function to remove duplicate videos
func removeDuplicates(videos []FeedVideo) []FeedVideo {
	seen := make(map[int]bool)
	var uniqueVideos []FeedVideo
	for _, video := range videos {
		if !seen[video.Id] {
			seen[video.Id] = true
			uniqueVideos = append(uniqueVideos, video)
		}
	}
	return uniqueVideos
}

// func FetchAllCategoryVideos(id_user int, sale_type int, typeV int, categoryId, page int) ([]FeedVideo, error) {
// 	var err error
// 	dbConn := Pool
// 	var vid []FeedVideo
// 	var ads []FeedVideo

// 	pageSize := 12
// 	var user User
// 	typeUser := 0
// 	// Calculate the offset based on the page number and page size
// 	offset := (page - 1) * pageSize

// 	user, err = GetUserByID(uint(id_user))
// 	if err == nil {
// 		typeUser = user.IdMembership

// 	}
// 	// err = dbConn.Unscoped().Find(&vid).Error

// 	// err = dbConn.Model(&Video{}).Where("sale_type_id = ? && is_vip = ?", sale_type, isvip).Limit(pageSize).Offset(offset).Preload("SaleType").Preload("SaleCategory").Preload("User").Unscoped().Find(&vid).Error
// 	// if err != nil {
// 	// 	return vid, err
// 	// }

// 	// var videos []Video
// 	result := dbConn.Table("videos").
// 		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 		Where("sale_type_id = ? && type = ? && sale_category_id = ? ", sale_type, typeV, categoryId).
// 		// Where("sale_category_id = ? && is_vip = ?", sale_type, isvip).
// 		Limit(pageSize).
// 		Offset(offset).
// 		Preload("SaleType").
// 		Preload("SaleCategory").
// 		Preload("User").
// 		Unscoped().
// 		Find(&vid).Error

// 	if result != nil {
// 		// http.Error(w, "Database error", http.StatusInternalServerError)
// 		return vid, err
// 	}

// 	// responseJSON, err := json.Marshal(videos)
// 	// if err != nil {

// 	// 	return vid,err
// 	// }

// 	rand.NewSource(time.Now().UnixNano())

// 	rand.Shuffle(len(vid), func(i, j int) { vid[i], vid[j] = vid[j], vid[i] })

// 	// Process videos with ads
// 	if (typeUser == 0 || typeUser == 1 || typeUser == 8 || typeUser == 2 || typeUser == 3) && len(vid) > 0 {
// 		// var tempVideos []FeedVideo
// 		adsIndex := 0
// 		pageSize = 2
// 		offset := (page - 1) * pageSize

// 		err = dbConn.Table("videos").
// 			Select("videos.*").

// 			// Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
// 			Where("type = ?", 3).Limit(pageSize).Offset(offset).Unscoped().
// 			Find(&ads).Error

// 		if len(ads) != 0 {
// 			vid = append(vid, ads[adsIndex])
// 		}
// 		// for i, video := range vid {
// 		// 	// Add the new element to the end of the list
// 		// 	tempVideos = append(tempVideos, video)
// 		// 	if (i+1)%5 == 0 && video.Type != 2 {
// 		// 		// Skip this video or handle as needed
// 		// 		tempVideos = append(tempVideos, ads[adsIndex])
// 		// 		adsIndex++
// 		// 		// continue
// 		// 	}
// 		// }
// 		// vid = tempVideos
// 	}
// 	return vid, nil

// }

func FetchAllCategoryVideos(id_user int, sale_type int, typeV int, categoryId, page int, idvideo *int) ([]FeedVideo, error) {
	var (
		err      error
		vid      []FeedVideo
		ads      []FeedVideo
		typeUser int
		pageSize = 12
		offset   = (page - 1) * pageSize
		dbConn   = Pool
	)

	// Obtener el usuario y determinar su tipo
	user, err := GetUserByID(uint(id_user))
	if err == nil {
		typeUser = user.IdMembership
	}

	// Obtener los videos normales
	err = dbConn.Table("videos").
		Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
		Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
		Where("sale_type_id = ? AND type = ? AND sale_category_id = ?", sale_type, typeV, categoryId).
		Limit(pageSize).
		Offset(offset).
		Preload("SaleType").
		Preload("SaleCategory").
		Preload("User").
		Unscoped().
		Find(&vid).Error

	if err != nil {
		return vid, err
	}

	// Si se proporciona un idvideo, buscarlo y colocarlo en la posición 0
	if idvideo != nil {
		var specialVideo FeedVideo
		err = dbConn.Table("videos").
			Select("videos.*, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
			Joins("LEFT JOIN users_videos_favorites ON videos.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", id_user).
			Where("videos.id = ?", *idvideo).
			Preload("SaleType").
			Preload("SaleCategory").
			Preload("User").
			Unscoped().
			First(&specialVideo).Error

		if err == nil {
			vid = append([]FeedVideo{specialVideo}, vid...) // Coloca el video especial en la primera posición
		}
	}

	// Si se proporcionó un idvideo, asegúrate de que no esté duplicado
	if idvideo != nil && len(vid) > 1 && vid[1].Id == *idvideo {
		vid = vid[1:] // Elimina el video que ya está en la posición 0
	}

	// Mezclar los videos aleatoriamente, excepto el primero si es el video especial
	if len(vid) > 1 {
		rand.NewSource(time.Now().UnixNano())
		rand.Shuffle(len(vid)-1, func(i, j int) { vid[i+1], vid[j+1] = vid[j+1], vid[i+1] })
	}

	// Insertar anuncios si el usuario es del tipo adecuado
	if (typeUser == 0 || typeUser == 1 || typeUser == 2 || typeUser == 3 || typeUser == 8) && len(vid) > 0 {
		pageSize = 2
		offset = (page - 1) * pageSize
		err = dbConn.Table("videos").
			Select("videos.*").
			Where("type = ?", 3).
			Limit(pageSize).
			Offset(offset).
			Unscoped().
			Find(&ads).Error

		if err == nil && len(ads) > 0 {
			vid = append(vid, ads[0])
		}
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
		Where("type = ?", 1).
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

func findPrices(text string) []string {
	// This regex pattern is quite basic and might need to be refined depending on your needs
	// It's designed to match:
	// - Optional currency symbols ($, €, etc.) or currency codes (USD, EUR, etc.) at the start
	// - Numbers, which can include thousands separators (,) and decimal points (.)
	// - Optional currency codes (USD, EUR, etc.) at the end
	var pattern = `(?i)(\$\s*|€\s*|£\s*|¥\s*|usd\s*|eur\s*|mxn\s*|jpy\s*)?\d{1,8}(,\d{3})*(\.\d{1,2})?(\s*(usd|eur|gbp|mxn))?`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	return matches
}

func findPricePatterns(search string) []string {
	// Implementa la lógica para extraer patrones de precios de la cadena de búsqueda.
	// Por ejemplo, podrías buscar todos los números en la cadena.
	var patterns []string
	priceRegex := regexp.MustCompile(`\d+`)
	prices := priceRegex.FindAllString(search, -1)
	for _, price := range prices {
		patterns = append(patterns, price)
	}
	return patterns
}

// func FetchAllCategoryVideosWithFilters(userID, saleID, typeVideo, category, page int, idVideo *int, userLat, userLon float64) ([]FeedVideo, error) {
// 	var videos []FeedVideo
// 	db := Pool

// 	// Define el rango máximo de distancia (en kilómetros)
// 	const maxDistance = 50.0
// 	const limitPerPriority = 10 // Límite por prioridad (ajustable)

// 	// // Subconsulta para calcular la distancia

// 	// Subconsulta para calcular la distancia y filtrar por rango
// 	subquery := db.Table("videos").
// 		Select("videos.*, (6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance", userLat, userLon, userLat).
// 		Having("distance <= ?", maxDistance)

// 	// Consulta principal
// 	query := db.Table("(?) AS v", subquery).
// 		Select("v.*, users.medal_type, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
// 		Joins("INNER JOIN users ON v.id_user = users.id").
// 		Joins("LEFT JOIN users_videos_favorites ON v.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", userID).
// 		Where("videos.sale_type_id = ? AND videos.type = ? AND videos.sale_category_id = ?", saleID, typeVideo, category)

// 		// Scan(&results).Error

// 	if idVideo != nil {
// 		query = query.Or("videos.id = ?", *idVideo)
// 	}

// 	// Iterar por las prioridades (1, 2, 3)
// 	for medalType := 1; medalType <= 3; medalType++ {
// 		var group []FeedVideo

// 		// Subconsulta para obtener videos de esta prioridad con aleatoriedad
// 		err := query.Where("users.medal_type = ?", medalType).
// 			Order("RAND()").         // Orden aleatorio
// 			Limit(limitPerPriority). // Límite por grupo
// 			Find(&group).Error
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Agregar el grupo al resultado final
// 		videos = append(videos, group...)
// 	}

// 	// Paginación sobre el resultado final
// 	start := (page - 1) * 20
// 	end := start + 20
// 	if start >= len(videos) {
// 		return []FeedVideo{}, nil // Página vacía
// 	}
// 	if end > len(videos) {
// 		end = len(videos)
// 	}
// 	return videos[start:end], nil
// }

func FetchAllCategoryVideosWithFilters(userID, saleID int, typeVideo, page int, idVideo *int, userLat, userLon float64) ([]FeedVideo, error) {
	// var videos []FeedVideo
	var (
		// err      error
		videos []FeedVideo
		// ads      []FeedVideo
		// typeUser int
		pageSize = 12
		offset   = (page - 1) * pageSize
	)

	// Obtener el usuario y determinar su tipo
	// user, err := GetUserByID(uint(userID))
	// if err == nil {
	// 	typeUser = user.IdMembership
	// }

	// Define el rango máximo de distancia (en kilómetros)
	const maxDistance = 50.0
	const limitPerPriority = 10 // Límite por prioridad (ajustable)
	db := Pool
	// Función auxiliar para construir la consulta principal
	buildQuery := func(ignoreDistance bool) *gorm.DB {
		subquery := db.Table("videos").
			Select("videos.*, (6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance", userLat, userLon, userLat).
			Limit(pageSize).
			Offset(offset)

		if !ignoreDistance {
			subquery = subquery.Having("distance <= ?", maxDistance)
		}

		query := db.Table("(?) AS v", subquery).
			Select("v.*, users.medal_type, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
			Joins("INNER JOIN users ON v.id_user = users.id").
			Joins("LEFT JOIN users_videos_favorites ON v.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", userID).
			Where("v.sale_type_id = ? AND v.type = ?", saleID, typeVideo).
			Limit(pageSize).
			Offset(offset).
			Preload("SaleType").
			Preload("SaleCategory").
			Preload("User")
		if idVideo != nil {
			query = query.Or("v.id = ?", *idVideo)
		}

		return query
	}

	// Intentar primero con el filtro de distancia
	query := buildQuery(false)

	// Iterar por las prioridades (1, 2, 3)
	for medalType := 1; medalType <= 3; medalType++ {
		var group []FeedVideo

		err := query.Where("users.medal_type = ?", medalType).
			Order("RAND()").
			Limit(limitPerPriority).
			Find(&group).Error
		if err != nil {
			return nil, err
		}

		videos = append(videos, group...)
	}

	// Si no se encontraron videos, intentar sin el filtro de distancia
	if len(videos) == 0 {
		query = buildQuery(true)

		for medalType := 1; medalType <= 3; medalType++ {
			var group []FeedVideo

			err := query.Where("users.medal_type = ?", medalType).
				Order("RAND()").
				Limit(limitPerPriority).
				Find(&group).Error
			if err != nil {
				return nil, err
			}

			videos = append(videos, group...)
		}
	}

	// Paginación sobre el resultado final
	// start := (page - 1) * 20
	// end := start + 20
	// if start >= len(videos) {
	// 	return []FeedVideo{}, nil // Página vacía
	// }
	// if end > len(videos) {
	// 	end = len(videos)
	// }
	return videos, nil
}

func FetchAllVideosWithFilters(userID, saleID, typeVideo, category int, page int, idVideo *int, userLat, userLon float64) ([]FeedVideo, error) {
	// var videos []FeedVideo
	var (
		// err      error
		videos []FeedVideo
		// ads      []FeedVideo
		// typeUser int
		pageSize = 12
		offset   = (page - 1) * pageSize
	)

	// Obtener el usuario y determinar su tipo
	// user, err := GetUserByID(uint(userID))
	// if err == nil {
	// 	typeUser = user.IdMembership
	// }

	// Define el rango máximo de distancia (en kilómetros)
	const maxDistance = 50.0
	const limitPerPriority = 10 // Límite por prioridad (ajustable)
	db := Pool
	// Función auxiliar para construir la consulta principal
	buildQuery := func(ignoreDistance bool) *gorm.DB {
		subquery := db.Table("videos").
			Select("videos.*, (6371 * acos(cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + sin(radians(?)) * sin(radians(latitude)))) AS distance", userLat, userLon, userLat)
			// Limit(pageSize).
			// Offset(offset)

		if !ignoreDistance {
			subquery = subquery.Having("distance <= ?", maxDistance)
		}

		query := db.Table("(?) AS v", subquery).
			Select("v.*, users.medal_type, IF(users_videos_favorites.id IS NULL, 0, 1) AS is_favorite").
			Joins("INNER JOIN users ON v.id_user = users.id").
			Joins("LEFT JOIN users_videos_favorites ON v.id = users_videos_favorites.id_video AND users_videos_favorites.id_user = ?", userID).
			Where("v.sale_type_id = ? AND v.type = ? AND v.sale_category_id = ?", saleID, typeVideo, category).
			Limit(pageSize).
			Offset(offset).
			Preload("SaleType").
			Preload("SaleCategory").
			Preload("User")
		if idVideo != nil {
			query = query.Or("v.id = ?", *idVideo)
		}

		return query
	}

	// Intentar primero con el filtro de distancia
	query := buildQuery(false)
	// Iterar por las prioridades (1, 2, 3)
	for medalType := 1; medalType <= 3; medalType++ {
		var group []FeedVideo

		err := query.Where("users.medal_type = ?", medalType).
			Order("RAND()").
			Limit(limitPerPriority).
			Find(&group).Error
		if err != nil {
			return nil, err
		}

		videos = append(videos, group...)
	}

	// Si no se encontraron videos, intentar sin el filtro de distancia
	if len(videos) == 0 {
		query = buildQuery(true)

		for medalType := 1; medalType <= 3; medalType++ {
			var group []FeedVideo

			err := query.Where("users.medal_type = ?", medalType).
				Order("RAND()").
				Limit(limitPerPriority).
				Find(&group).Error
			if err != nil {
				return nil, err
			}

			videos = append(videos, group...)
		}
	}

	return videos, nil
}
