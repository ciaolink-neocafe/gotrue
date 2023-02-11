package sms_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/netlify/gotrue/internal/conf"
)

type GatewayProvider struct {
	Config *conf.GatewayProviderConfiguration
}

// Creates a SmsProvider with the gateway Config
func NewGatewayProvider(config conf.GatewayProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &GatewayProvider{
		Config: &config,
	}, nil
}

// Send an SMS containing the OTP with gateway
func (t *GatewayProvider) SendSms(phone string, message string) error {
	body, err := json.Marshal(map[string]string{
		"recipient": phone,
		"body":      message,
		"sender":    t.Config.Sender,
	})
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: defaultTimeout}
	r, err := http.NewRequest("POST", t.Config.Url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	if len(t.Config.BearerToken) > 0 {
		r.Header.Add("Authorization", "Bearer "+t.Config.BearerToken)
	}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("Unexpected response while calling the SMS gateway: %v", res.StatusCode)
	}

	return nil
}
