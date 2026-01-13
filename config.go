package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Endpoint  string `json:"endpoint"`
	Model     string `json:"model"`
	Key       string `json:"key"`
	PromptDir string `json:"prompt_dir"`
}

func (receiver *Config) String() string {
	raw, _ := json.MarshalIndent(receiver, "", "  ")
	return string(raw)
}

func (receiver *Config) GetPromptDir() string {
	if receiver.PromptDir != "" {
		return receiver.PromptDir
	}
	return filepath.Join(os.Getenv("HOME"), ".aicc", "prompt")
}

func Load(content string) (*Config, error) {
	// 解析配置文件
	var config Config
	err := json.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadFromFile(filePath string) (*Config, error) {
	// 读取配置文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return Load(string(content))
}
