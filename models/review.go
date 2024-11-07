package models

import (
	// "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Review representa la estructura de la tabla reviews
type Review struct {
	ID        int     `gorm:"primaryKey" json:"id"`
	Rating    float64 `gorm:"type:double(10,2)" json:"rating" binding:"required"`
	Review    string  `gorm:"type:varchar(250)" json:"review" binding:"required"`
	IDUser    int     `gorm:"not null" json:"id_user"`
	IDProfile int     `gorm:"not null" json:"id_profile" binding:"required"`
}

// TableName indica el nombre de la tabla en la base de datos
func (Review) TableName() string {
	return "reviews"
}

// CreateOrUpdate inserta una nueva reseña o actualiza una existente
func (r *Review) CreateOrUpdate() error {

	dbConn := Pool

	// Utilizamos GORM para el INSERT ON DUPLICATE KEY UPDATE
	return dbConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id_user"}, {Name: "id_profile"}}, // Unicidad en id_user y id_profile
		UpdateAll: true,                                                     // Actualizar todos los campos si ya existe
	}).Create(r).Error
}

// CreateOrUpdate inserta una nueva reseña o actualiza una existente
func GetReviewsProfile(idProfile int) ([]Review, error) {

	dbConn := Pool

	var reviews []Review
	if err := dbConn.Where("id_profile = ?", idProfile).Find(&reviews).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener reseñas"})
		return reviews,err
	}
	return reviews,nil
}

// DeleteReview elimina una reseña específica por su ID
func DeleteReview(reviewID int) error {
    dbConn := Pool

    // Eliminación de la reseña por su ID
    if err := dbConn.Delete(&Review{}, reviewID).Error; err != nil {
        return err
    }
    return nil
}

// DeleteReviewByUserAndProfile elimina una reseña específica por el ID de usuario y el perfil
func DeleteReviewByUserAndProfile(userID int, profileID int) error {
    dbConn := Pool

    // Eliminación de la reseña que coincide con el id_user y el id_profile
    if err := dbConn.Where("id_user = ? AND id_profile = ?", userID, profileID).Delete(&Review{}).Error; err != nil {
        return err
    }
    return nil
}