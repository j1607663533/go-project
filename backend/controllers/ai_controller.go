package controllers

import (
	"context"
	"fmt"
	"gin-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type AIChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type AIChatResponse struct {
	Reply string `json:"reply"`
}

type AIController struct {
	client *openai.Client
}

func NewAIController() *AIController {
	// 初始化 OpenAI 客户端 (兼容 SiliconFlow 等提供商)
	conf := openai.DefaultConfig(config.AppConfig.AI.APIKey)
	if config.AppConfig.AI.BaseURL != "" {
		conf.BaseURL = config.AppConfig.AI.BaseURL
	}

	client := openai.NewClientWithConfig(conf)
	return &AIController{
		client: client,
	}
}

func (ctrl *AIController) Chat(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}

	// 检查 API Key 是否配置
	if config.AppConfig.AI.APIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 服务未配置 API Key，请在 .env 文件中配置 AI_API_KEY"})
		return
	}

	// 调用 AI 接口
	resp, err := ctrl.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: config.AppConfig.AI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI 接口调用失败: %v", err)})
		return
	}

	reply := resp.Choices[0].Message.Content
	c.JSON(http.StatusOK, AIChatResponse{
		Reply: reply,
	})
}
