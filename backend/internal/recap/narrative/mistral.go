package narrative

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultMistralEndpoint = "https://api.mistral.ai/v1/chat/completions"

type MistralHTTPProvider struct {
	APIKey   string
	Model    string
	Endpoint string
	Timeout  time.Duration
	Client   *http.Client
}

func (p MistralHTTPProvider) Generate(ctx context.Context, input Input) (ProviderResult, error) {
	if p.APIKey == "" {
		return ProviderResult{}, errors.New("mistral api key is empty")
	}
	model := p.Model
	if model == "" {
		model = "mistral-small-latest"
	}
	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = defaultMistralEndpoint
	}
	timeout := p.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	requestContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	client := p.Client
	if client == nil {
		client = &http.Client{}
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return ProviderResult{}, err
	}
	requestBody := map[string]any{
		"model":       model,
		"temperature": 0,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": string(inputJSON)},
		},
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":        "recap_narrative",
				"strict":      true,
				"description": nil,
				"schema_definition": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"required":             []string{"summary_title", "summary_text"},
					"properties": map[string]any{
						"summary_title": map[string]any{"type": "string", "minLength": 1, "maxLength": 120},
						"summary_text":  map[string]any{"type": "string", "minLength": 1, "maxLength": 500},
					},
				},
			},
		},
	}
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		return ProviderResult{}, err
	}
	req, err := http.NewRequestWithContext(requestContext, http.MethodPost, endpoint, bytes.NewReader(encoded))
	if err != nil {
		return ProviderResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return ProviderResult{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return ProviderResult{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ProviderResult{}, fmt.Errorf("mistral status %d", resp.StatusCode)
	}
	var decoded struct {
		Model   string `json:"model"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return ProviderResult{}, err
	}
	if len(decoded.Choices) == 0 || decoded.Choices[0].Message.Content == "" {
		return ProviderResult{}, errors.New("empty mistral response")
	}
	actualModel := decoded.Model
	if actualModel == "" {
		actualModel = model
	}
	return ProviderResult{Content: []byte(decoded.Choices[0].Message.Content), Model: actualModel}, nil
}

const systemPrompt = `Ты формируешь только безопасный текст персонального recap на русском языке.
Верни строго JSON с полями summary_title и summary_text, без дополнительных полей.
Тон дружелюбный и нейтральный, без оценки личности.
Не придумывай факты и числа. Используй числа только из safe_facts; год можно брать из поля year.
Не добавляй рекомендации, обещания выгоды, призывы к действию, URL или сравнения с другими пользователями.`
