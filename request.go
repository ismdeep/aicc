package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type ModelData struct {
	Created int64  `json:"created"`
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned by"`
}

type ModelsResponseData struct {
	Data []ModelData `json:"data"`
}

func joinURL(endpoint string, path string) string {
	return strings.TrimRight(endpoint, "/") + path
}

func doRequest(method string, requestURL string, key string, body io.Reader) (string, error) {
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("unexpected status code %d, content: %s", resp.StatusCode, string(raw))
	}

	return string(raw), nil
}

func Request(endpoint string, model string, key string, input string) (string, error) {
	requestURL := joinURL(endpoint, "/chat/completions")

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

	return doRequest(http.MethodPost, requestURL, key, bytes.NewBuffer(requestDataRaw))
}

func RequestModels(endpoint string, key string) ([]ModelData, error) {
	requestURL := joinURL(endpoint, "/models")

	raw, err := doRequest(http.MethodGet, requestURL, key, nil)
	if err != nil {
		return nil, err
	}

	var responseData ModelsResponseData
	if err := json.Unmarshal([]byte(raw), &responseData); err != nil {
		return nil, fmt.Errorf("failed json.Unmarshal, content: %v, err: %w", raw, err)
	}

	return responseData.Data, nil
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
