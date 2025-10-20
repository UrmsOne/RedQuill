# LLM 统一客户端

基于设计文档实现的统一大模型推理客户端，支持 OpenAI、Azure OpenAI、Ollama 等主流厂商。

## 功能特性

- **统一接口**: 所有厂商使用相同的 API 调用方式
- **多厂商支持**: OpenAI、Azure OpenAI、Ollama
- **流式响应**: 支持实时流式输出
- **类型安全**: 强类型定义，编译时检查
- **配置灵活**: 支持环境变量和配置文件
- **错误处理**: 统一的错误处理和重试机制

## 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "redquill-backend/pkg/utils/llm"
)

func main() {
    // 创建客户端
    client, err := llm.NewClient(llm.LLMConfig{
        Provider: "openai",
        BaseURL:  "https://api.openai.com/v1",
        APIKey:   "sk-xxx",
        Model:    "gpt-3.5-turbo",
        Timeout:  30 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 同步调用
    resp, err := client.Chat(context.Background(), llm.ChatRequest{
        Messages: []llm.Message{
            {Role: "user", Content: "Hello, how are you?"},
        },
        Temperature: 0.7,
        MaxTokens:   1000,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 2. 流式调用

```go
// 流式调用
stream, err := client.ChatStream(context.Background(), llm.ChatRequest{
    Messages: []llm.Message{
        {Role: "user", Content: "Tell me a story"},
    },
    Stream: true,
})
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Error != nil {
        log.Printf("Error: %v", chunk.Error)
        break
    }
    
    for _, choice := range chunk.Choices {
        if choice.Delta.Content != "" {
            fmt.Print(choice.Delta.Content)
        }
    }
}
```

### 3. 多厂商配置

```go
// 从配置文件加载
client, err := llm.NewClientFromFile("config.json", "azure")
if err != nil {
    log.Fatal(err)
}

// 从环境变量加载
client, err := llm.NewClientFromEnv()
if err != nil {
    log.Fatal(err)
}
```

## 配置

### 环境变量

```bash
LLM_PROVIDER=openai
LLM_BASE_URL=https://api.openai.com/v1
LLM_API_KEY=sk-xxx
LLM_MODEL=gpt-3.5-turbo
LLM_TIMEOUT=30s
```

### 配置文件

```json
{
  "default": "openai",
  "providers": {
    "openai": {
      "provider": "openai",
      "base_url": "https://api.openai.com/v1",
      "api_key": "sk-xxx",
      "model": "gpt-3.5-turbo",
      "timeout": "30s"
    },
    "azure": {
      "provider": "azure",
      "base_url": "https://your-resource.openai.azure.com",
      "api_key": "your-api-key",
      "model": "gpt-35-turbo",
      "timeout": "30s"
    },
    "ollama": {
      "provider": "ollama",
      "base_url": "http://localhost:11434",
      "model": "llama2",
      "timeout": "60s"
    }
  }
}
```

## 支持的厂商

### OpenAI
- 官方 OpenAI API
- 支持同步和流式响应
- 需要 API Key

### Azure OpenAI
- Azure OpenAI 服务
- 需要 API Key 和 Endpoint
- 支持自定义部署名称

### Ollama
- 本地 Ollama 服务
- 默认端口 11434
- 无需 API Key

## 错误处理

```go
resp, err := client.Chat(ctx, req)
if err != nil {
    if llmErr, ok := err.(*llm.LLMError); ok {
        switch llmErr.Type {
        case string(llm.ErrorTypeAuth):
            // 认证错误
        case string(llm.ErrorTypeRateLimit):
            // 限流错误
        case string(llm.ErrorTypeServer):
            // 服务器错误
        }
    }
}
```

## 扩展

### 添加新厂商

1. 在 `providers/` 目录下创建新的提供商文件
2. 实现 `LLMClient` 接口
3. 在 `client.go` 中添加新的 case

### 自定义配置

```go
config := llm.LLMConfig{
    Provider:   "custom",
    BaseURL:    "https://api.custom.com",
    APIKey:     "custom-key",
    Model:      "custom-model",
    Headers:    map[string]string{
        "X-Custom-Header": "value",
    },
    Timeout:    30 * time.Second,
    MaxRetries: 3,
    RetryDelay: 1 * time.Second,
}
```

## 测试

```go
func TestLLMClient(t *testing.T) {
    client, err := llm.NewClient(testConfig)
    require.NoError(t, err)
    
    resp, err := client.Chat(ctx, llm.ChatRequest{
        Messages: []llm.Message{
            {Role: "user", Content: "Hello"},
        },
    })
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Choices[0].Message.Content)
}
```
