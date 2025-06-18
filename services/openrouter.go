package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"morethancoder/t3-clone/utils"
	"net/http"
	"os"
)

type OpenRouterInstance struct {
	APIKey string
	Url    string
	Client *http.Client
}

var OpenRouter = NewOpenRouter()

func NewOpenRouter() *OpenRouterInstance {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		utils.Log.Fatal("OPENROUTER_API_KEY is not set")
	}
	return &OpenRouterInstance{
		APIKey: apiKey,
		Url:    "https://openrouter.ai/api/v1/chat/completions",
		Client: &http.Client{},
	}
}

type ImageURL struct {
	URL string `json:"url"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type Message struct {
	Role      string      `json:"role"`
	Content   interface{} `json:"content"` // Can be string or []Content
	Refusal   *string     `json:"refusal"`
	Reasoning *string     `json:"reasoning"`
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Logprobs           interface{} `json:"logprobs"`
	FinishReason       string      `json:"finish_reason"`
	NativeFinishReason string      `json:"native_finish_reason"`
	Index              int         `json:"index"`
	Message            Message     `json:"message"`
}

type UsageDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type Usage struct {
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokensDetails     UsageDetails            `json:"prompt_tokens_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type OpenRouterResponse struct {
	ID                string   `json:"id"`
	Provider          string   `json:"provider"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Choices           []Choice `json:"choices"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Usage             Usage    `json:"usage"`
}

func (o *OpenRouterInstance) Request(req OpenRouterRequest) (*OpenRouterResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", o.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)

	httpRes, err := o.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", httpRes.StatusCode, body)
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (o *OpenRouterInstance) RequestStream(req OpenRouterRequest, handler func(data []byte)) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", o.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)

	httpRes, err := o.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != 200 {
		body, _ := io.ReadAll(httpRes.Body)
		return fmt.Errorf("non-200 status code: %d, body: %s", httpRes.StatusCode, body)
	}

	decoder := json.NewDecoder(httpRes.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunkData, err := json.Marshal(chunk)
		if err != nil {
			return err
		}

		handler(chunkData)
	}

	return nil
}
