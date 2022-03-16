package app

import (
	"fmt"
	"go_concurrency_test/item"
	"go_concurrency_test/logger"
)

var bufferSize = 20

// ProcessItems processes all the given items passing to item.ProcessItem and logging its result (success or failure).
// If an error is detected while processing, remaining processes should be aborted.
// In addition, some buffering is done in each processing, and it will be flushed every 20 processing.
func ProcessItems(items []item.Item) {
	failed := false
	processedCount := 0

	for _, targetItem := range items {
		result, err := item.ProcessItem(targetItem)
		if err != nil {
			logger.Error(fmt.Sprintf("Processing failed: %s", err.Error()))
			failed = true
			break // break the loop when caught an error
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

	logger.Debug(fmt.Sprintf("Succeeded count: %d", processedCount))
}
