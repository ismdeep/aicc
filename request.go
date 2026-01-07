package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestData struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseData struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func Request(endpoint string, model string, key string, input string) (string, error) {
	requestURL := fmt.Sprintf("%v/chat/completions", endpoint)

	requestData := RequestData{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	requestDataRaw, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed json.Marshal, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestDataRaw))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func RequestContent(endpoint string, model string, key string, input string) (string, error) {
	raw, err := Request(endpoint, model, key, input)
	if err != nil {
		return "", err
	}

	var responseData ResponseData
	if err := json.Unmarshal([]byte(raw), &responseData); err != nil {
		return "", fmt.Errorf("failed json.Unmarshal, content: %v, err: %w", raw, err)
	}

	if len(responseData.Choices) == 0 {
		return "", fmt.Errorf("no choices in the response, content: %v", raw)
	}

	return responseData.Choices[0].Message.Content, nil
}
