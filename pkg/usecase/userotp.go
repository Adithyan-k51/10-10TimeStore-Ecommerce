package usecase

import (
	"context"
	"ecommerce/pkg/commonhelp/requests.go"
	"ecommerce/pkg/config"
	interfaces "ecommerce/pkg/usecase/interface"
	"fmt"
	"strings"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type OtpUseCase struct {
	cfg config.Config
}

func NewOtpUseCase(cfg config.Config) interfaces.OtpUseCase {
	return &OtpUseCase{
		cfg: cfg,
	}
}

func (c *OtpUseCase) SendOTP(ctx context.Context, mobno requests.OTPreq) (string, error) {
	// Validate phone number
	if mobno.Phone == "" {
		return "", fmt.Errorf("phone number is required")
	}

	// Ensure phone number starts with +
	if !strings.HasPrefix(mobno.Phone, "+") {
		mobno.Phone = "+" + mobno.Phone
	}

	// Validate Twilio configuration
	if c.cfg.ACCOUNTSID == "" || c.cfg.AUTHTOCKEN == "" || c.cfg.SERVICES_ID == "" {
		return "", fmt.Errorf("twilio configuration is incomplete")
	}

	// Initialize Twilio client
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: c.cfg.ACCOUNTSID,
		Password: c.cfg.AUTHTOCKEN,
	})

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(mobno.Phone)
	params.SetChannel("sms")

	resp, err := twilioClient.VerifyV2.CreateVerification(c.cfg.SERVICES_ID, params)
	if err != nil {
		return "", fmt.Errorf("failed to send OTP: %w", err)
	}

	if resp == nil || resp.Status == nil {
		return "", fmt.Errorf("received empty response from Twilio")
	}

	return *resp.Status, nil
}

func (c *OtpUseCase) VerifyOTP(ctx context.Context, userData requests.Otpverifier) error {
	// Validate input
	if userData.Phone == "" || userData.Pin == "" {
		return fmt.Errorf("phone number and PIN are required")
	}

	// Ensure phone number starts with +
	if !strings.HasPrefix(userData.Phone, "+") {
		userData.Phone = "+" + userData.Phone
	}

	// Validate Twilio configuration
	if c.cfg.ACCOUNTSID == "" || c.cfg.AUTHTOCKEN == "" || c.cfg.SERVICES_ID == "" {
		return fmt.Errorf("twilio configuration is incomplete")
	}

	// Initialize Twilio client
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: c.cfg.ACCOUNTSID,
		Password: c.cfg.AUTHTOCKEN,
	})

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(userData.Phone)
	params.SetCode(userData.Pin)

	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(c.cfg.SERVICES_ID, params)
	if err != nil {
		return fmt.Errorf("failed to verify OTP: %w", err)
	}

	if resp == nil || resp.Status == nil {
		return fmt.Errorf("received empty response from Twilio")
	}

	if *resp.Status != "approved" {
		return fmt.Errorf("invalid OTP code")
	}

	return nil
}
