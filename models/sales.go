package models

import (
    "gorm.io/gorm"
)

type Sale struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null" json:"user_id"`
    Sales     int       `gorm:"default:0" json:"sales"`
}

func (Sale) TableName() string {
    return "sales"
}

// Método para buscar o crear un registro de ventas para un usuario
func (s *Sale) IncrementOrCreate( userID uint) error {
	dbConn := Pool
    // Busca el registro de ventas del usuario
    if err := dbConn.Where("user_id = ?", userID).First(&s).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Si no existe, crea un nuevo registro
            s.UserID = userID
            s.Sales = 1
            if err := dbConn.Create(s).Error; err != nil {
                return err
            }
        } else {
            return err
        }
    } else {
        // Si el registro ya existe, incrementa las ventas
        s.Sales += 1
        if err := dbConn.Save(s).Error; err != nil {
            return err
        }
    }
    return nil
}
