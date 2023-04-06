package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"time"
)

type VectorDB struct {
	uri      string
	user     string
	password string
	client   client.Client
}

func NewVectorDB(uri, user, password string) *VectorDB {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := client.NewDefaultGrpcClientWithURI(ctx, uri, user, password)
	if err != nil {
		fmt.Println("connect failed:", err)
		return nil
	}

	return &VectorDB{uri: uri, user: user, password: password, client: client}
}

func (v *VectorDB) Insert(blocks []string, result []openai.Embedding) error {
	// here is the collection name we use in this example
	collectionName := `docs`
	var docIds []int8
	var vectors [][]float32
	for i := 0; i < len(blocks); i++ {
		docIds = append(docIds, int8(rand.Int()%10))
		vectors = append(vectors, result[i].Embedding)
	}

	docIdColumn := entity.NewColumnInt8("docId", docIds)
	contentColumn := entity.NewColumnVarChar("content", blocks)
	vectorColumn := entity.NewColumnFloatVector("vector", 1536, vectors)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := v.client.Insert(ctx, collectionName, "", docIdColumn, contentColumn, vectorColumn)
	if err != nil {
		fmt.Println("failed to insert data:", err.Error())
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = v.client.Flush(ctx, collectionName, false)
	if err != nil {
		fmt.Println("failed to flush collection:", err.Error())
		return err
	}
	return nil
}

func (v *VectorDB) Query(query openai.Embedding) (error, float32, string) {
	vec2search := []entity.Vector{
		entity.FloatVector(query.Embedding),
	}
	sp, _ := entity.NewIndexFlatSearchParam()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := v.client.Search(ctx, "docs", nil, "", []string{"id", "content", "docId"},
		vec2search, "vector", entity.IP, 1, sp)
	if err != nil {
		fmt.Println("failed to search collection, err:", err)
		return err, float32(0.0), ""
	}
	if len(result) == 0 || result[0].ResultCount == 0 {
		fmt.Println("failed to search collection")
		return err, float32(0.0), ""
	}

	// find the content column convert into string result
	var content string
	for _, field := range result[0].Fields {
		if field.Name() == "content" {
			content = string(field.(*entity.ColumnVarChar).Data()[0])
			break
		}
	}
	return nil, result[0].Scores[0], content
}

func (v *VectorDB) QueryTopK(query openai.Embedding, topK int) (error, []float32, []string) {
	if topK <= 0 {
		return errors.New("param invalid"), nil, nil
	}
	vec2search := []entity.Vector{
		entity.FloatVector(query.Embedding),
	}
	sp, _ := entity.NewIndexFlatSearchParam()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := v.client.Search(ctx, "docs", nil, "", []string{"id", "content", "docId"},
		vec2search, "vector", entity.IP, topK, sp)
	if err != nil {
		fmt.Println("failed to search collection, err:", err)
		return err, nil, nil
	}

	if len(result) == 0 {
		fmt.Println("failed to search collection")
		return err, nil, nil
	}

	// find the content column convert into string result
	var contents []string
	var scores []float32
	for _, field := range result[0].Fields {
		if field.Name() == "content" {
			for i := 0; i < result[0].ResultCount; i++ {
				contents = append(contents, string(field.(*entity.ColumnVarChar).Data()[i]))
			}
		}
	}
	for _, score := range result[0].Scores {
		scores = append(scores, score)
	}
	return nil, scores, contents
}
