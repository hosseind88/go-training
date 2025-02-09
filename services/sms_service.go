package services

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client *twilio.RestClient
}

func NewSMSService() *SMSService {
	client := twilio.NewRestClient()

	return &SMSService{
		client: client,
	}
}

func (s *SMSService) SendVerificationCode(to, code string, template string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(fmt.Sprintf(template, code))

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}

	return nil
}
