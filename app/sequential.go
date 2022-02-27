package app

import (
	"fmt"
	"go_concurrency_test/item"
	"os"
)

var bufferSize = 20

// ProcessItems 引数の items を全件 item.ProcessItem で処理し、成功時も失敗時もログを出す.
// 1個でも失敗した場合は処理を辞める.
// また何らかのバッファリングしている想定で、20回に1回バッファをクリアする.
func ProcessItems(items []item.Item) {
	failed := false
	processedCount := 0

	for _, targetItem := range items {
		result, err := item.ProcessItem(targetItem)
		if err != nil {
			if _, err := fmt.Fprintln(os.Stderr, fmt.Sprintf("\033[31mProcessing failed: %s\033[0m", err.Error())); err != nil {
				panic(err)
			}
			failed = true
			break // 1個でも失敗したら処理を抜ける
		}

		fmt.Println(result)

		processedCount++

		if processedCount%bufferSize == 0 {
			fmt.Println("Flush buffer!")
		}
	}

	if !failed && processedCount%bufferSize > 0 {
		fmt.Println("Flush buffer!")
	}
}
