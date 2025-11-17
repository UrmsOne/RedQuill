// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: stream.go
/@Description: Stream processing for providers
/*/

package providers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ProcessSSEStream 处理SSE流
func (sp *StreamProcessor) ProcessSSEStream(ctx context.Context, resp *http.Response) (<-chan StreamChunk, error) {
	stream := make(chan StreamChunk, 100)
	
	go func() {
		defer close(stream)
		defer resp.Body.Close()
		
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					break
				}
				
				var chunk StreamChunk
				if err := json.Unmarshal([]byte(data), &chunk); err == nil {
					select {
					case stream <- chunk:
					case <-ctx.Done():
						return
					}
				}
			}
		}
		
		if err := scanner.Err(); err != nil {
			select {
			case stream <- StreamChunk{
				Error: &LLMError{
					Type:    string(ErrorTypeNetwork),
					Message: fmt.Sprintf("stream read error: %v", err),
				},
			}:
			case <-ctx.Done():
			}
		}
	}()
	
	return stream, nil
}

// ProcessJSONStream 处理JSON流
func (sp *StreamProcessor) ProcessJSONStream(ctx context.Context, resp *http.Response) (<-chan StreamChunk, error) {
	stream := make(chan StreamChunk, 100)
	
	go func() {
		defer close(stream)
		defer resp.Body.Close()
		
		decoder := json.NewDecoder(resp.Body)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			
			var chunk StreamChunk
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				
				select {
				case stream <- StreamChunk{
					Error: &LLMError{
						Type:    string(ErrorTypeNetwork),
						Message: fmt.Sprintf("json decode error: %v", err),
					},
				}:
				case <-ctx.Done():
				}
				return
			}
			
			select {
			case stream <- chunk:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return stream, nil
}
