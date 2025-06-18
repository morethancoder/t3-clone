package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QueryParam string

const (
	QueryParamExpand QueryParam = "expand"
	QueryParamFields QueryParam = "fields"
)

type AuthRecord struct {
	CollectionID    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	ID              string `json:"id"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Verified        bool   `json:"verified"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
}

type OAuth2AuthRequest struct {
	Provider     string                 `json:"provider"`
	Code         string                 `json:"code"`
	CodeVerifier string                 `json:"codeVerifier"`
	RedirectURL  string                 `json:"redirectURL"`
	CreateData   map[string]interface{} `json:"createData,omitempty"`
}

type OAuth2AuthResponse struct {
	Token  string `json:"token"`
	Record AuthRecord `json:"record"`
	Meta struct {
		ID           string                 `json:"id"`
		Name         string                 `json:"name"`
		Username     string                 `json:"username"`
		Email        string                 `json:"email"`
		AvatarURL    string                 `json:"avatarURL"`
		AccessToken  string                 `json:"accessToken"`
		RefreshToken string                 `json:"refreshToken"`
		Expiry       string                 `json:"expiry"`
		IsNew        bool                   `json:"isNew"`
		RawUser      map[string]interface{} `json:"rawUser"`
	} `json:"meta"`
}

func (db *db) AuthWithOAuth2(request OAuth2AuthRequest, queryParams map[QueryParam]string) (*OAuth2AuthResponse, error) {
	authEndpoint := "/api/collections/users/auth-with-oauth2"

	queryString := ""
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(string(key), value)
		}
		queryString = "?" + query.Encode()
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", db.Url+authEndpoint+queryString, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	var authResponse OAuth2AuthResponse
	if err := json.Unmarshal(respBody, &authResponse); err != nil {
		return nil, err
	}

	return &authResponse, nil
}


// AuthRefreshResponse represents the response structure for the auth-refresh endpoint
type AuthRefreshResponse struct {
	Token  string `json:"token"`
	Record AuthRecord `json:"record"`
}

func (d *db) AuthRefresh(token string) (*AuthRefreshResponse, error) {
	if token == "" {
		return nil, errors.New("[ERROR] Authorization token is required")
	}

	url := fmt.Sprintf("%s/api/collections/users/auth-refresh", d.Url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[ERROR] request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResponse AuthRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, fmt.Errorf("[ERROR] failed to parse response: %w", err)
	}

	return &authResponse, nil
}
