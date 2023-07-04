package sms_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supabase/gotrue/internal/conf"
	"github.com/supabase/gotrue/internal/utilities"
)

type GenericProvider struct {
	Config *conf.GenericProviderConfiguration
}

// Creates a SmsProvider with the gateway Config
func NewGenericProvider(config conf.GenericProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &GenericProvider{
		Config: &config,
	}, nil
}

func (t *GenericProvider) SendMessage(phone string, message string, channel string) (string, error) {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message)
	default:
		return "", fmt.Errorf("channel type %q is not supported for gateway", channel)
	}
}

// Send an SMS containing the OTP with gateway
func (t *GenericProvider) SendSms(phone string, message string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"recipient": phone,
		"body":      message,
		"sender":    t.Config.Sender,
	})
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: defaultTimeout}
	r, err := http.NewRequest("POST", t.Config.Url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/json")
	if len(t.Config.BearerToken) > 0 {
		r.Header.Add("Authorization", "Bearer "+t.Config.BearerToken)
	}
	res, err := client.Do(r)
	defer utilities.SafeClose(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode/100 != 2 {
		return "", fmt.Errorf("Unexpected response while calling the SMS gateway: %v", res.StatusCode)
	}

	return "", nil
}
