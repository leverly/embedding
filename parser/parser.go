package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type FileParser struct {
	Content  []string
	fileName string
}

func NewParser(file string) *FileParser {
	return &FileParser{fileName: file}
}

func (p *FileParser) ParseFile() error {
	file, err := os.Open(p.fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var block string
	for scanner.Scan() {
		line := scanner.Bytes()
		// new block start
		if bytes.HasPrefix(line, []byte("###")) {
			if len(block) > 0 {
				p.Content = append(p.Content, string(block))
			}
			// reset to find a new block
			block = ""
		} else {
			if len(block) != 0 {
				block += ("\n" + string(line))
			} else {
				block = string(line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
