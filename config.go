package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
	Key      string `json:"key"`
}

func (receiver *Config) String() string {
	raw, _ := json.MarshalIndent(receiver, "", "  ")
	return string(raw)
}

func (receiver *Config) GetPromptDir() string {
	return filepath.Join(os.Getenv("HOME"), ".aicc", "prompt")
}

func Load(content string) (*Config, error) {
	// 解析配置文件
	var config Config
	err := json.Unmarshal([]byte(content), &config)
	if err != nil {
		return &Config{}, err
	}
	return &config, nil
}

func LoadFromFile(filePath string) (*Config, error) {
	// 读取配置文件
	content, _ := os.ReadFile(filePath)
	c, _ := Load(string(content))
	// 从环境变量加载
	if os.Getenv("OPENAI_BASE_URL") != "" {
		c.Endpoint = os.Getenv("OPENAI_BASE_URL")
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		c.Key = os.Getenv("OPENAI_API_KEY")
	}
	if os.Getenv("OPENAI_MODEL") != "" {
		c.Model = os.Getenv("OPENAI_MODEL")
	}
	return c, nil
}
