package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
	Key      string `json:"key"`
}

func (receiver *Config) String() string {
	raw, _ := json.Marshal(receiver)
	return string(raw)
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
