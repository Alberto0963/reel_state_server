package models

type VerificationCode struct {
    // gorm.Model `gorm:"softDelete:false"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`
	ID uint `gorm:"not null;unique" json:"id"`
	Phone string `gorm:"size:13;not null;unique" json:"phone"`
	Code     string `gorm:"size:255;not null;unique" json:"code"`
}

func (VerificationCode) TableName() string {
    return "verification_code"
}


func (v *VerificationCode) SaveVerificationCode() (*VerificationCode, error) {

	var err error
	dbConn := Pool

	err = dbConn.Create(&v).Error
	if err != nil {
		return &VerificationCode{}, err
	}
	return v, nil
}

func (c *VerificationCode) CodeIsValid() (*VerificationCode, error){
	
	var err error
	dbConn := Pool
	var verificationCode VerificationCode
	// Check if the verification code exists in the database and is not expired
	err = dbConn.Where("code = ? AND phone = ? ", c.Code, c.Phone).First(&verificationCode).Error

	if err != nil {
		return &VerificationCode{}, err
		
	}

	return c,nil
}
