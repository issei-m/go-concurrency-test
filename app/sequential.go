package app

import (
	"fmt"
	"go_concurrency_test/item"
	"os"
)

func ProcessItems(items []item.Item) {
	for _, targetItem := range items {
		if result, err := item.ProcessItem(targetItem); err != nil {
			if _, err := fmt.Fprintln(os.Stderr, fmt.Sprintf("\033[31mProcessing failed: %s\033[0m", err.Error())); err != nil {
				panic(err)
			}
			break // 1個でも失敗したら処理を抜ける
		} else {
			fmt.Println(result)
		}
	}
}
