package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func generateWithDeepSeek(cardName string) (string, error) {
	client := openai.NewClient(os.Getenv("DEEPSEEK_API_KEY"))
	
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(
						"Сгенерируй прогноз на день для карты %s. Формат:\nПрогноз: [текст]\nТеги: [3 тега]\nДелать: [3 совета]\nНе делать: [3 совета]", 
						cardName,
					),
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}