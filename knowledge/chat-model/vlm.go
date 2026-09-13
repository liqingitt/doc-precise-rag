package chatmodel

import (
	"context"
	"doc-precise-rag/knowledge/config"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

var VlmModel *openai.ChatModel

func init() {
	model, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey:  *config.AppConfig.AiConfig.VlmConfig.Token,
		BaseURL: *config.AppConfig.AiConfig.VlmConfig.BaseURL,
		Model:   *config.AppConfig.AiConfig.VlmConfig.ModelName,
	})

	if err != nil {
		panic(err)
	}

	VlmModel = model

}
