package client

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"time"
)

var apikeyList = [3]string{
	"sk-ziIjWUX7ZmfiVA86pSPHT",
	"sk-Fsf8MmcrssrKy5B7PqvxT",
	"sk-t3wJTgGCLSv6nwNrAmBLT",
}

type OpenaiClient struct {
	client  *openai.Client
	timeout time.Duration
}

func NewOpenAIClient(timeout time.Duration) *OpenaiClient {
	rand.Seed(time.Now().Unix())
	return &OpenaiClient{client: openai.NewClient(apikeyList[rand.Int()%len(apikeyList)]), timeout: timeout}
}

func (c *OpenaiClient) Embedding(blocks []string) (error, []openai.Embedding) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := c.client.CreateEmbeddings(ctx,
		openai.EmbeddingRequest{
			Input: blocks,
			Model: openai.AdaEmbeddingV2,
		})
	if err != nil {
		fmt.Println("Create Embeddings error:", err)
		return err, nil
	}
	return nil, resp.Data
}

func (c *OpenaiClient) ChatCompletion(messages []openai.ChatCompletionMessage) (error, *openai.ChatCompletionResponse) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := c.client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		fmt.Println("ChatCompletion error:", err)
		return err, nil
	}
	return nil, &resp
}

func (c *OpenaiClient) Transcription(file string) (error, *openai.AudioResponse) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.CreateTranscription(ctx,
		openai.AudioRequest{
			FilePath:    file,
			Model:       "whisper-1",
			Prompt:      "",
			Temperature: 0.3,
			Language:    "zh",
		},
	)

	if err != nil {
		fmt.Println("CreateTranscription error:", err)
		return err, nil
	}
	return nil, &resp
}

func (c *OpenaiClient) Translation(file, lang string) (error, *openai.AudioResponse) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.CreateTranslation(ctx,
		openai.AudioRequest{
			FilePath:    file,
			Model:       "whisper-3",
			Prompt:      "",
			Temperature: 0.5,
			Language:    lang,
		},
	)

	if err != nil {
		fmt.Println("CreateTranslation error:", err)
		return err, nil
	}
	return nil, &resp
}
