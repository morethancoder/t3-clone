package db

import (
	"encoding/json"
	"io"
	"net/http"
	"errors"
)

type AuthMethods struct {
	MFA      MFAConfig      `json:"mfa"`
	OAuth2   OAuth2Config   `json:"oauth2"`
	OTP      OTPConfig      `json:"otp"`
	Password PasswordConfig `json:"password"`
}

type MFAConfig struct {
	Duration int  `json:"duration"`
	Enabled  bool `json:"enabled"`
}

type OTPConfig struct {
	Duration int  `json:"duration"`
	Enabled  bool `json:"enabled"`
}

type PasswordConfig struct {
	Enabled        bool     `json:"enabled"`
	IdentityFields []string `json:"identityFields"`
}

type OAuth2Config struct {
	Enabled   bool             `json:"enabled"`
	Providers []OAuth2Provider `json:"providers"`
}

type OAuth2Provider struct {
	AuthURL             string `json:"authURL"`
	AuthUrl             string `json:"authUrl"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	CodeVerifier        string `json:"codeVerifier"`
	DisplayName         string `json:"displayName"`
	Name                string `json:"name"`
	State               string `json:"state"`
}

func (db *db) GetAuthMethods() (*AuthMethods, error) {
	authMethodsEndpoint := "/api/collections/users/auth-methods"

	req, err := http.NewRequest("GET", db.Url+authMethodsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := db.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respBody))
	}

	var authMethods AuthMethods
	if err := json.Unmarshal(respBody, &authMethods); err != nil {
		return nil, err
	}

	return &authMethods, nil
}

