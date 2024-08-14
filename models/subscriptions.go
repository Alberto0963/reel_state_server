package models

import (
	// "errors"
	// "fmt"
	// "html"
	// "os"
	"time"
	// "strings"
	// "golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"
)

type Subscription struct {
	ID                   int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int       `gorm:"not null" json:"id_user"`
	MembershipId         int       `gorm:"not null" json:"membership_id"`
	Renewal              bool      `json:"renewal"`
	PaypalSubscriptionId string    `gorm:"size:255" json:"paypal_subscription_id"`
	RenewalCancelledAt   time.Time `json:"renewal_cancelled_at"`
}

type Createsubscription struct {
	ID                   int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int       `gorm:"not null" json:"id_user"`
	MembershipId         int       `gorm:"not null" json:"membership_id"`
	Renewal              bool      `json:"renewal"`
	PaypalSubscriptionId string    `gorm:"size:255" json:"paypal_subscription_id"`
	// RenewalCancelledAt   time.Time `json:"renewal_cancelled_at"`
}

func (Subscription) TableName() string {
	return "paypal_subscriptions"
}

func (Createsubscription) TableName() string {
	return "paypal_subscriptions"
}

func (sub *Createsubscription) CreateSubscription() (*Createsubscription, error) {
	var err error
	dbConn := Pool

	if err = dbConn.Create(&sub).Error; err != nil {
		return &Createsubscription{}, err
	}
	return sub, nil
}

func CancelSubscriptionFunction(subID string, date time.Time) error {
	var err error
	dbConn := Pool

	var subscriber Subscription

	if err = dbConn.Where("paypal_subscription_id = ?", subID).First(&subscriber).Error; err != nil {
		return err
	}

	subscriber.Renewal = false
	subscriber.RenewalCancelledAt = date

	if err = dbConn.Save(&subscriber).Error; err != nil {
		return err
	}

	return nil
}

// func CancelUserSubscriptionFunction(userid int, curr_code string) error {
// 	var err error
// 	dbConn := Pool

// 	var subscriber UserUpdate

// 	if err = dbConn.Where("id = ?", userid).First(&subscriber).Error; err != nil {
// 		return err
// 	}

// 	if curr_code == "MXN" {
// 		subscriber.IdMembership = 100005
// 	} else {
// 		subscriber.IdMembership = 100001
// 	}	

// 	// subscriber.RenewalCancelledAt = date

// 	if err = dbConn.Save(&subscriber).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }
