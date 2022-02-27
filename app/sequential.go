package app

import (
	"fmt"
	"go_concurrency_test/item"
	"go_concurrency_test/logger"
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
			logger.Error(fmt.Sprintf("Processing failed: %s", err.Error()))
			failed = true
			break // 1個でも失敗したら処理を抜ける
		}

		logger.Info(result)

		processedCount++

		if processedCount%bufferSize == 0 {
			logger.Info("Flush buffer!")
		}
	}

	if !failed && processedCount%bufferSize > 0 {
		logger.Info("Flush buffer!")
	}
}
