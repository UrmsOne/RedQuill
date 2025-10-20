// Package llm
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: config.go
/@Description: LLM client configuration
/*/

package llm

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

// LLMConfig LLM配置
type LLMConfig struct {
	Provider   string            `json:"provider"`   // openai, azure, ollama, deepseek, doubao, qwen
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key"`
	Model      string            `json:"model"`
	Headers    map[string]string `json:"headers,omitempty"`
	Timeout    time.Duration     `json:"timeout"`
	MaxRetries int               `json:"max_retries"`
	RetryDelay time.Duration     `json:"retry_delay"`
}

// MultiProviderConfig 多厂商配置
type MultiProviderConfig struct {
	Default   string                `json:"default"`
	Providers map[string]LLMConfig  `json:"providers"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int         `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
	Backoff    BackoffType `json:"backoff"`
}

// BackoffType 退避类型
type BackoffType string

const (
	BackoffLinear      BackoffType = "linear"
	BackoffExponential BackoffType = "exponential"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	MaxIdleConns        int           `json:"max_idle_conns"`
	MaxIdleConnsPerHost int           `json:"max_idle_conns_per_host"`
	IdleConnTimeout     time.Duration `json:"idle_conn_timeout"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"`
	Size    int           `json:"size"`
}

// ConcurrencyConfig 并发配置
type ConcurrencyConfig struct {
	MaxConcurrentRequests int `json:"max_concurrent_requests"`
	MaxConcurrentStreams  int `json:"max_concurrent_streams"`
}

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(filename string) (*MultiProviderConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var config MultiProviderConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *LLMConfig {
	return &LLMConfig{
		Provider:   getEnv("LLM_PROVIDER", "openai"),
		BaseURL:    getEnv("LLM_BASE_URL", "https://api.openai.com/v1"),
		APIKey:     getEnv("LLM_API_KEY", ""),
		Model:      getEnv("LLM_MODEL", "gpt-3.5-turbo"),
		Timeout:    parseDuration(getEnv("LLM_TIMEOUT", "30s")),
		MaxRetries: parseInt(getEnv("LLM_MAX_RETRIES", "3")),
		RetryDelay: parseDuration(getEnv("LLM_RETRY_DELAY", "1s")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second
	}
	return d
}

func parseInt(s string) int {
	if s == "" {
		return 3
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 3
}