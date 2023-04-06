package main

import (
	"embedding/client"
	"embedding/database"
	"embedding/parser"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Library struct {
	db   database.Database
	init bool
}

func newLibrary() *Library {
	// using memory as database and init must be false
	return &Library{db: database.NewMemoryDB(), init: false}

	// using milvus as database and init should be true or false
	return &Library{
		db: database.NewVectorDB(
			"https://vectordb.zillizcloud.com:19538",
			"db_admin",
			"password"),
		init: true}
}

func (l *Library) Init(filedir string) error {
	if l.init == false {
		visitFile := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".txt") {
				return nil
			}
			return l.embeddingFile(path)
		}
		err := filepath.Walk(filedir, visitFile)
		if err != nil {
			return err
		}
		l.init = true
	}
	return nil
}

func (l *Library) embeddingFile(filename string) error {
	fileParse := parser.NewParser(filename)
	err := fileParse.ParseFile()
	if err != nil {
		return err
	}
	if len(fileParse.Content) > 0 {
		client := client.NewOpenAIClient(time.Second * 20)
		err, result := client.Embedding(fileParse.Content)
		if err != nil {
			return err
		}
		err = l.db.Insert(fileParse.Content, result)
		if err != nil {
			return err
		}
	}
	return nil
}

// find the similarest block
func (l *Library) FindSimilarBlock(query string) (error, float32, string) {
	if l.init == true {
		client := client.NewOpenAIClient(time.Second * 20)
		err, result := client.Embedding([]string{query})
		if err != nil {
			return err, 0.0, ""
		}
		return l.db.Query(result[0])
	}
	return errors.New("not inited"), float32(0.0), ""
}

// find the topk similar blocks
func (l *Library) FindSimilarTopKBlock(query string, topK int) (error, []float32, []string) {
	if l.init == true {
		client := client.NewOpenAIClient(time.Second * 20)
		err, result := client.Embedding([]string{query})
		if err != nil {
			return err, nil, nil
		}
		return l.db.QueryTopK(result[0], topK)
	}
	return errors.New("not inited"), nil, nil
}
