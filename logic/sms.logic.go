package logic

import (
	"errors"
	"github.com/nyaruka/phonenumbers"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"os"
)

func SendPhoneVerification(phoneNumber string) error {
	client := twilio.NewRestClient()

	if client == nil {
		return errors.New("failed to create twilio client")
	}
	channel := "sms"

	params := &openapi.CreateVerificationParams{
		To:      &phoneNumber,
		Channel: &channel,
	}

	_, err := client.VerifyV2.CreateVerification(os.Getenv("TWILIO_VERIFY_SERVICE_SID"), params)
	return err
}

func CheckPhoneVerification(code, phoneNumber string) error {
	client := twilio.NewRestClient()

	if client == nil {
		return errors.New("failed to create twilio client")
	}

	params := &openapi.CreateVerificationCheckParams{
		Code: &code,
		To:   &phoneNumber,
	}

	_, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("TWILIO_VERIFY_SERVICE_SID"), params)
	return err
}

func FormatPhoneNumber(phoneNumber string) (string, error) {
	number, err := phonenumbers.Parse(phoneNumber, "US")
	if err != nil {
		return "", err
	}
	to := phonenumbers.Format(number, phonenumbers.E164)
	return to, nil
}
