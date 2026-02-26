package ai

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/genai" // 新的导入路径
)

// Client 包装了新的 genai.Client
type Client struct {
	client *genai.Client
}

// NewClient 创建一个新的 Gen AI 客户端
func NewClient(ctx context.Context) (*Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("环境变量 GEMINI_API_KEY 未设置")
	}

	// --- 主要修改点 1: 客户端初始化 ---
	// 使用 ClientConfig 来配置客户端，指定使用 Gemini API 后端和 API Key [citation:2][citation:3]
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI, // 明确指定使用 Gemini API
		// HTTPOptions: 如果需要自定义 HTTP 客户端，可以在这里设置 [citation:2]
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Gen AI 客户端失败: %w", err)
	}

	return &Client{client: client}, nil
}

// Ask 发送提示词给模型并返回响应
func (c *Client) Ask(prompt string) (string, error) {
	if c.client == nil {
		return "", errors.New("客户端未初始化")
	}

	// --- 主要修改点 2: 模型调用方式 ---
	// 新 SDK 使用 client.Models.GenerateContent 方法，需要指定模型名称 [citation:2]
	// 模型名称字符串与旧 SDK 中使用的保持一致，例如 "gemini-2.5-flash"
	resp, err := c.client.Models.GenerateContent(
		context.Background(), // 注意：这里需要传入一个 context，可以从函数参数传入，或使用接收器中的 ctx
		"gemini-2.5-flash",   // 模型 ID [citation:2]
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{Text: prompt},
				},
				Role: "user", // 可以指定角色
			},
		},
		nil, // 这里可以传入 *genai.GenerateContentConfig 来设置可选参数，如 Temperature、MaxTokens 等 [citation:2]
	)

	if err != nil {
		return "", fmt.Errorf("AI 生成内容失败: %w", err)
	}

	// --- 主要修改点 3: 响应解析方式 ---
	// 新 SDK 的响应结构略有不同，需要遍历 Candidates 和 Parts 来提取文本
	if resp == nil || len(resp.Candidates) == 0 {
		return "", errors.New("AI 返回了空响应")
	}

	var fullResponse string
	for _, candidate := range resp.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			// Part 可能是文本、内联数据或其他类型，我们只关心文本
			if part.Text != "" {
				fullResponse += part.Text
			}
			// 如果需要处理其他类型（如函数调用），可以在此扩展
		}
	}

	if fullResponse == "" {
		// 可能被安全过滤或其他原因
		return "AI 没有生成有效的文本响应。", nil
	}

	return fullResponse, nil
}

// --- 主要修改点 4: 移除 Close 方法 ---
// 新 SDK 的客户端不需要显式调用 Close() 来释放资源 [citation:3]
// 因此可以删除 Close 方法，或者保留一个空方法以防旧代码调用
// func (c *Client) Close() error {
//     return nil
// }