package models

import ()

type Likes struct {
	// gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	Id         int  `gorm:"not null;unique" json:"id"`
	Id_user    uint `gorm:"size:255;not null;" json:"id_user"`
	Id_profile uint `gorm:"size:255;not null;" json:"id_profile"`
}

func GetlikedUsersId(userID int) ([]int, error) {

	dbConn := Pool

	var likedUserIDs []int
	if err := dbConn.Table("likes").Select("id_user").Where("id_profile = ?", userID).Find(&likedUserIDs).Error; err != nil {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener likes"})
		return likedUserIDs, err
	}
	return likedUserIDs, nil
}
