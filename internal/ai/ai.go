package airequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vivalabelousov2025/go-worker/internal/dto"
)

func AiRequest(prompt string) (dto.AiResponse, error) {
	apiKey := "sk-179cda0b06d741dfad6969c3282b25fe"

	requestBody := dto.AiRequest{
		Model: "deepseek-chat",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return dto.AiResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return dto.AiResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.AiResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dto.AiResponse{}, err
	}

	var r dto.AiResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return dto.AiResponse{}, err
	}

	return r, nil
}
