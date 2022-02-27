package main

import (
	"fmt"
	"go_concurrency_test/item"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 100個分の item.Item
	targetItems := item.CreateItems(100)

	fmt.Println("Processing start!")

	for _, targetItem := range targetItems {
		if result, err := item.ProcessItem(targetItem); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("\033[31mProcessing failed: %s\033[0m", err.Error()))
			break // 1個でも失敗したら処理を抜ける
		} else {
			fmt.Println(result)
		}
	}

	fmt.Println("DONE!")
}
