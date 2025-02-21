package SMS

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSMS(number string, code string, signature string) {
    // Get your Twilio credentials from environment variables
    accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
    authToken := os.Getenv("TWILIO_AUTH_TOKEN")
    phone := os.Getenv("TWILIO_PHONE")

    // Create a Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
    // Set the sender and receiver phone numbers

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(phone)

	params.SetBody(fmt.Sprintf("Your ReelState verification code is: %s, Don't share this code with anyone. %s", code,signature))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}


func GenerateRandomCode(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)

}

