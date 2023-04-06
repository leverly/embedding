package database

import "github.com/sashabaranov/go-openai"

type Database interface {
	Insert(blocks []string, result []openai.Embedding) error
	Query(query openai.Embedding) (error, float32, string)
	QueryTopK(query openai.Embedding, topK int) (error, []float32, []string)
}
