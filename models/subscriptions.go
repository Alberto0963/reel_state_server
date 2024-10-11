package models

import (
	// "errors"
	// "fmt"
	// "html"
	// "os"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	SMS "reelState/utils"
	"time"

	"gorm.io/gorm"
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
	NextBillingTime      time.Time `json:"next_billing_time"`
	StartedAt            time.Time `json:"started_at"`
}

type UpdateSubscription struct {
	ID                   int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int       `gorm:"not null" json:"id_user"`
	MembershipId         int       `gorm:"not null" json:"membership_id"`
	Renewal              bool      `json:"renewal"`
	PaypalSubscriptionId string    `gorm:"size:255" json:"paypal_subscription_id"`
	RenewalCancelledAt   time.Time `json:"renewal_cancelled_at"`
	StartedAt            time.Time `json:"started_at"`
	CustomerId           string `gorm:"size:255" json:"customer_id"`

}

type SubscriptionView struct {
	ID                   int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int       `gorm:"not null" json:"id_user"`
	MembershipId         int       `gorm:"not null" json:"membership_id"`
	Membership           string    `gorm:"not null" json:"membership"`
	Price                int       `gorm:"not null" json:"price"`
	Renewal              bool      `json:"renewal"`
	PaypalSubscriptionId string    `gorm:"size:255" json:"paypal_subscription_id"`
	RenewalCancelledAt   time.Time `json:"renewal_cancelled_at"`
	NextBillingTime      time.Time `json:"next_billing_time"`
	StartedAt            time.Time `json:"started_at"`
}

type CancelSubscription struct {
	Id string `json:"id"`
	Reason               string `json:"reason"`
}

type Createsubscription struct {
	ID                   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	IdUser               int    `gorm:"not null" json:"id_user"`
	MembershipId         int    `gorm:"not null" json:"membership_id"`
	Renewal              bool   `json:"renewal"`
	PaypalSubscriptionId string `gorm:"size:255" json:"paypal_subscription_id"`
	CurrencyCode         string `gorm:"size:255" json:"currency_code"`
	CustomerId           string `gorm:"size:255" json:"customer_id"`

	// RenewalCancelledAt   time.Time `json:"renewal_cancelled_at"`
}

func (SubscriptionView) TableName() string {
	return "UserMemberships"
}

func (UpdateSubscription) TableName() string {
	return "paypal_subscriptions"
}

func (Subscription) TableName() string {
	return "paypal_subscriptions"
}

func (Createsubscription) TableName() string {
	return "paypal_subscriptions"
}


func GetSubscriptionByID(uid uint) (UpdateSubscription, error) {

	var s UpdateSubscription
	// Obtain a connection from the pool
	dbConn := Pool
	// defer dbConn.Close()
	// if err := dbConn.Model(&User{}).
	// Define the subquery as a raw SQL string
	// videosCountSubquery := "(SELECT id_user, COUNT(*) as video_count FROM videos GROUP BY id_user) as v"

	if err := dbConn.Model(&UpdateSubscription{}).First(&s, uid).Error; err != nil {
		return s, errors.New("User not found! ")
	}

	return s, nil

}

func (sub *Createsubscription) CreateSubscription() (*Createsubscription, error) {
	var err error
	dbConn := Pool

	if err = dbConn.Create(&sub).Error; err != nil {
		return &Createsubscription{}, err
	}
	return sub, nil
}

func CancelSubscriptionFunction(subID int, date time.Time) error {
	var err error
	dbConn := Pool

	var subscriber Subscription

	if err = dbConn.Where("id = ?", subID).First(&subscriber).Error; err != nil {
		return err
	}

	subscriber.Renewal = false
	subscriber.RenewalCancelledAt = date

	if err = dbConn.Save(&subscriber).Error; err != nil {
		return err
	}

	return nil
}

func CancelSubscriptionWebhook(paypal_subscription_id string, date time.Time) error {
	var err error
	dbConn := Pool

	var subscriber Subscription

	if err = dbConn.Where("paypal_subscription_id = ?", paypal_subscription_id).First(&subscriber).Error; err != nil {
		return err
	}

	subscriber.Renewal = false
	subscriber.RenewalCancelledAt = date

	if err = dbConn.Save(&subscriber).Error; err != nil {
		return err
	}

	return nil
}

func CancelSubscriptionIfActive(userID string, reason string, cancelDate time.Time, typesub string) error {
	dbConn := Pool
	var subscriber UpdateSubscription

	// Verificar si el usuario tiene una suscripción activa con el ID proporcionado
	if err := dbConn.Where("renewal = 1 AND id_user = ?", userID).First(&subscriber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no active subscription found for user with ID %s", userID)
		}
		return fmt.Errorf("error querying subscription: %w", err)
	}

	switch typesub {
	case "openpay":
		// Cancelar la suscripción en PayPal
		if err := CancelOpenPaySubscription(subscriber.CustomerId,subscriber.PaypalSubscriptionId); err != nil {
			return fmt.Errorf("error cancelling openpay subscription: %w", err)
		}
	case "paypal":
		// Cancelar la suscripción en PayPal
		if err := CancelPaypalSubscription(subscriber.PaypalSubscriptionId, reason); err != nil {
			return fmt.Errorf("error cancelling PayPal subscription: %w", err)
		}
	}

	// Actualizar la suscripción en la base de datos
	subscriber.Renewal = false
	subscriber.RenewalCancelledAt = cancelDate
	if err := dbConn.Save(&subscriber).Error; err != nil {
		return fmt.Errorf("error updating subscription: %w", err)
	}

	return nil
}

func CancelPaypalSubscription(subscriptionID, reason string) error {

	client_id := os.Getenv("Client_id")
	secret_Key := os.Getenv("Secret_Key")

	url := fmt.Sprintf("%ssubscriptions/%s/cancel", os.Getenv("Paypal_apiUrl"), subscriptionID)
	// if !isSandbox {
	// 	url = fmt.Sprintf("https://api.paypal.com/v1/billing/subscriptions/%s/cancel", subscriptionID)
	// }

	tokenPaypal, err := SMS.FetchPayPalToken(client_id, secret_Key)
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}

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
	req.Header.Set("Authorization", "Bearer "+tokenPaypal)

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

func GetSubscription(userId int) ([]SubscriptionView, error) {
	var err error
	dbConn := Pool

	var subscriber []SubscriptionView

	if err = dbConn.Where("id_user = ? && renewal = 1", userId).Find(&subscriber).Error; err != nil {
		return []SubscriptionView{}, err
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
