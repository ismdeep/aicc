package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	promptDir := receiver.PromptDir
	if promptDir == "" {
		promptDir = filepath.Join(os.Getenv("HOME"), ".aicc", "prompt")
	}

	// render ${HOME}
	if strings.Contains(promptDir, "${HOME}") {
		promptDir = strings.ReplaceAll(promptDir, "${HOME}", os.Getenv("HOME"))
	}

	return promptDir
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
