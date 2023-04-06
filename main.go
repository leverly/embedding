package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	library := newLibrary()
	err := library.Init("./files/")
	if err != nil {
		fmt.Println("init failed:", err)
		return
	}
	fmt.Println("init succ, just play for fun")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Input:")
		scanner.Scan()
		query := scanner.Text()
		if query == "exit" {
			fmt.Println("exit succ, bye")
			return
		}
		// find the related blocks
		err, score, result := library.FindSimilarBlock(query)
		if err != nil {
			fmt.Println("Find similar error:", err)
			continue
		}
		fmt.Println("find similar block:", 1)
		fmt.Println(0, score, string(result))
		err, scores, list := library.FindSimilarTopKBlock(query, 3)
		if err != nil {
			fmt.Println("Find TopK similar error:", err)
			continue
		}
		fmt.Println("find topk blocks:", len(list))
		for i, _ := range list {
			fmt.Println(i, scores[i], list[i])
		}
	}
}
