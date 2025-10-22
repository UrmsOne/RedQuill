// Package providers
/*
/@Author: urmsone urmsone@163.com
/@Date: 2025/10/20 20:44
/@Name: mock.go
/@Description: Mock provider for testing
/*/

package providers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// MockProvider Mock提供商
type MockProvider struct {
	config LLMConfig
}

// NewMockProvider 创建Mock提供商
func NewMockProvider(config LLMConfig, client *http.Client) *MockProvider {
	return &MockProvider{
		config: config,
	}
}

// Chat 同步聊天
func (p *MockProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// 模拟延迟
	time.Sleep(100 * time.Millisecond)

	// 根据模板类型返回不同的模拟响应
	var content string
	switch req.Messages[0].Content {
	case "story_core":
		content = `{
			"concepts": [
				{
					"title": "重生之逆天改命",
					"core_conflict": "主角重生回到十八岁，获得系统帮助，开始逆天改命",
					"theme": "重生逆袭，系统流爽文",
					"innovation": "系统+重生+修炼的完美结合",
					"commercial_potential": "高，符合当前网文市场趋势",
					"target_audience": "玄幻爱好者，重生文读者"
				}
			]
		}`
	case "worldview":
		content = `{
			"power_system": {
				"name": "修真体系",
				"levels": ["练气", "筑基", "金丹", "元婴", "化神", "合体", "大乘", "渡劫"],
				"cultivation_method": "吸收天地灵气，炼化真气",
				"limitations": "需要天赋和机缘"
			},
			"society_structure": {
				"hierarchy": "以修为论地位",
				"major_factions": [
					{
						"name": "正道联盟",
						"type": "正派",
						"influence": "维护修真界秩序"
					}
				],
				"economic_system": "以灵石为货币的修真经济"
			},
			"geography": {
				"major_regions": ["东域", "西域", "南域", "北域", "中域"],
				"special_locations": ["秘境", "洞府", "仙山"]
			},
			"special_rules": ["天道规则", "因果循环", "修炼瓶颈"]
		}`
	case "character":
		content = `{
			"name": "林轩",
			"soul_profile": {
				"personality": {
					"core_traits": ["坚韧", "聪明", "正义"],
					"moral_compass": "惩恶扬善",
					"internal_conflicts": ["复仇与宽恕"],
					"fears": ["失去亲人"],
					"desires": ["变强", "保护家人"]
				},
				"background": {
					"origin": "普通家庭出身",
					"defining_events": ["重生", "获得系统"],
					"hidden_secrets": ["重生者身份"]
				},
				"motivations": {
					"immediate_goal": "快速提升实力",
					"long_term_goal": "成为最强修士",
					"core_drive": "保护所爱之人"
				}
			},
			"core_attributes": {
				"cultivation_level": "练气期",
				"current_items": ["系统", "基础功法"],
				"abilities": ["系统辅助", "修炼加速"],
				"relationships": {
					"enemies": ["前世仇人"],
					"allies": ["家人", "朋友"],
					"mentors": ["系统"]
				}
			}
		}`
	case "chapter":
		content = `{
			"title": "重生归来",
			"summary": "主角重生回到十八岁，获得系统帮助",
			"plot_advancements": ["重生觉醒", "系统激活", "开始修炼"],
			"character_development": {
				"林轩": "从普通人到修炼者的转变"
			},
			"next_chapter_hook": "系统发布第一个任务",
			"outline": {
				"goal": "建立重生设定",
				"key_events": ["重生", "系统激活", "开始修炼"],
				"dramatic_points": 3
			},
			"quality_metrics": {
				"score": 8,
				"strengths": ["设定清晰", "节奏紧凑"],
				"improvement_areas": ["细节描写"]
			},
			"content": "林轩睁开眼睛，发现自己回到了十八岁。这是怎么回事？他明明记得自己已经死了，死在了那个背叛他的女人手里。\n\n'系统激活中...'一个机械的声音在他脑海中响起。\n\n'检测到宿主重生，系统开始绑定...'\n\n林轩愣住了，系统？这不是小说里才有的东西吗？\n\n'绑定成功！欢迎使用逆天改命系统！'\n\n林轩的嘴角微微上扬，既然重生了，那就让那些曾经伤害过他的人付出代价！"
		}`
	default:
		content = `{
			"title": "测试生成",
			"content": "这是一个测试生成的内容",
			"success": true
		}`
	}

	return &ChatResponse{
		ID:    "mock-response-" + fmt.Sprintf("%d", time.Now().Unix()),
		Model: req.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     100,
			CompletionTokens: 200,
			TotalTokens:      300,
		},
		Created: time.Now().Unix(),
	}, nil
}

// ChatStream 流式聊天
func (p *MockProvider) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk, 10)

	go func() {
		defer close(ch)

		// 模拟流式响应
		content := "这是一个流式响应的测试内容，用于验证流式功能是否正常工作。"

		for i, char := range content {
			select {
			case <-ctx.Done():
				return
			case ch <- StreamChunk{
				ID:    "mock-stream-" + fmt.Sprintf("%d", time.Now().Unix()),
				Model: req.Model,
				Choices: []Choice{
					{
						Index: 0,
						Delta: Message{
							Role:    "assistant",
							Content: string(char),
						},
					},
				},
				Usage: &Usage{
					PromptTokens:     100,
					CompletionTokens: int64(i + 1),
					TotalTokens:      100 + int64(i+1),
				},
			}:
				time.Sleep(50 * time.Millisecond) // 模拟延迟
			}
		}

		// 发送结束标记
		ch <- StreamChunk{
			ID:    "mock-stream-end",
			Model: req.Model,
			Choices: []Choice{
				{
					Index:        0,
					FinishReason: "stop",
				},
			},
			Usage: &Usage{
				PromptTokens:     100,
				CompletionTokens: int64(len(content)),
				TotalTokens:      100 + int64(len(content)),
			},
		}
	}()

	return ch, nil
}

// Health 健康检查
func (p *MockProvider) Health(ctx context.Context) error {
	return nil
}

// Models 获取模型列表
func (p *MockProvider) Models(ctx context.Context) ([]Model, error) {
	return []Model{
		{
			ID:      "mock-model-1",
			Name:    "Mock Model 1",
			Owner:   "mock",
			Context: 4096,
		},
		{
			ID:      "mock-model-2",
			Name:    "Mock Model 2",
			Owner:   "mock",
			Context: 8192,
		},
	}, nil
}
