package utils

import (
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendPasswordResetOTPSMS sends a password reset OTP via SMS
func SendPasswordResetOTPSMS(toPhone string, otp string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	params := &api.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody("Your OTP is: " + otp)

	_, err := client.Api.CreateMessage(params)
	return err
}
