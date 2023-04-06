package database

import (
	"container/heap"
	"embedding/util"
	"errors"
	"github.com/sashabaranov/go-openai"
)

type MemoryDB struct {
	blocks          []string
	blocksEmbddings []openai.Embedding
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{}
}

func (m *MemoryDB) Query(query openai.Embedding) (error, float32, string) {
	max := float32(0.0)
	if len(m.blocks) == 0 {
		return errors.New("no records"), max, ""
	}
	index := -1
	for i, item := range m.blocksEmbddings {
		score := util.CosineSimilarity(query.Embedding, item.Embedding)
		if score > max {
			max = score
			index = i
		}
	}
	return nil, max, m.blocks[index]
}

func (m *MemoryDB) QueryTopK(query openai.Embedding, topK int) (error, []float32, []string) {
	if len(m.blocks) == 0 {
		return errors.New("no records"), nil, nil
	}
	if topK <= 0 {
		return errors.New("param invalid"), nil, nil
	}
	var temp util.MinHeap
	heap.Init(&temp)
	for i, item := range m.blocksEmbddings {
		score := util.CosineSimilarity(query.Embedding, item.Embedding)
		if temp.Len() < topK {
			heap.Push(&temp, util.Item{Score: score, Index: i})
		} else if score > temp[0].Score {
			heap.Pop(&temp)
			heap.Push(&temp, util.Item{Score: score, Index: i})
		}
	}
	var result []string
	var scores []float32
	for temp.Len() > 0 {
		result = append([]string{m.blocks[temp[0].Index]}, result...)
		scores = append([]float32{temp[0].Score}, scores...)
		heap.Pop(&temp)
	}
	return nil, scores, result
}

func (m *MemoryDB) Insert(blocks []string, result []openai.Embedding) error {
	if len(blocks) != len(result) {
		return errors.New("check param failed")
	}

	// store the original blocks
	for _, block := range blocks {
		m.blocks = append(m.blocks, block)
	}

	// store the embedding result
	for _, res := range result {
		m.blocksEmbddings = append(m.blocksEmbddings, res)
	}
	return nil
}
