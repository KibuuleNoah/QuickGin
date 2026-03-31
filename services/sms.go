package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TwilioConfig holds Twilio credentials loaded from environment variables.
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string // E.164, e.g. +12025551234
}

// SMSService sends SMS messages via Twilio's REST API.
type SMSService struct {
	cfg    TwilioConfig
	client *http.Client
}

func NewSMSService(cfg TwilioConfig) *SMSService {
	return &SMSService{cfg: cfg, client: &http.Client{}}
}

// SendOTP sends a 6-digit OTP to the given E.164 phone number.
func (s *SMSService) SendOTP(to, otp string) error {
	msgBody := fmt.Sprintf("Your verification code is %s. It expires in 5 minutes.", otp)

	apiURL := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		s.cfg.AccountSID,
	)

	data := url.Values{}
	data.Set("To", to)
	data.Set("From", s.cfg.FromNumber)
	data.Set("Body", msgBody)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("build twilio request: %w", err)
	}
	req.SetBasicAuth(s.cfg.AccountSID, s.cfg.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var twilioErr struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}
		json.Unmarshal(body, &twilioErr)
		return fmt.Errorf("twilio error %d: %s", twilioErr.Code, twilioErr.Message)
	}
	return nil
}
