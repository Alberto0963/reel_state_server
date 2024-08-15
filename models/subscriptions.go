package models

import (
	// "errors"
	// "fmt"
	// "html"
	// "os"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	NextBillingTime      time.Time  `json:"next_billing_time"` 
	StartedAt			 time.Time  `json:"started_at"` 
}

type CancelSubscription struct {
	PaypalSubscriptionId string  `json:"paypal_subscription_id"`
	Reason               string `json:"reason"`
}

type Createsubscription struct {
	ID                   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int    `gorm:"not null" json:"id_user"`
	MembershipId         int    `gorm:"not null" json:"membership_id"`
	Renewal              bool   `json:"renewal"`
	PaypalSubscriptionId string `gorm:"size:255" json:"paypal_subscription_id"`
	CurrencyCode         string `gorm:"size:255" json:"currency_code"`

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

// func CancelSubscription(subID string, date time.Time) error {
// 	var err error
// 	dbConn := Pool

// 	var subscriber Subscription

// 	if err = dbConn.Where("paypal_subscription_id = ?", subID).First(&subscriber).Error; err != nil {
// 		return err
// 	}

// 	subscriber.Renewal = false
// 	subscriber.RenewalCancelledAt = date

// 	if err = dbConn.Save(&subscriber).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// CancelSubscription cancels a PayPal subscription given its ID and an access token
func CancelPaypalSubscription(subscriptionID, reason string) error {

	accessToken := os.Getenv("Paypal_accessToken")
	url := fmt.Sprintf("%ssubscriptions/%s/cancel", os.Getenv("Paypal_apiUrl"), subscriptionID)
	// if !isSandbox {
	// 	url = fmt.Sprintf("https://api.paypal.com/v1/billing/subscriptions/%s/cancel", subscriptionID)
	// }

	// Create the request body
	body := map[string]string{
		"reason": reason,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}


func GetSubscription(userId int) ([]Subscription, error) {
	var err error
	dbConn := Pool

	var subscriber []Subscription

	if err = dbConn.Where("id_user = ? && renewal = 1", userId).Find(&subscriber).Error; err != nil {
		return []Subscription{} ,err
	}

	return subscriber, nil
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
