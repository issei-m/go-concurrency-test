package app

import (
	"context"
	"fmt"
	"go_concurrency_test/item"
	"go_concurrency_test/logger"
	"runtime"
	"sync"
)

type taskResult struct {
	result string
	err    error
}

func ProcessItemsConcurrently(items []item.Item, concurrency int) {
	numItems := len(items)
	var chItems = make(chan item.Item)
	var chTaskResults = make(chan *taskResult)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		close(chTaskResults)
		cancel()
	}()

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer func() {
				logger.Debug(fmt.Sprintf("[WORKER %d] Terminated, bye (remaining %d worker(s))", workerID, runtime.NumGoroutine()-2)) // これから削除される自分自身と、 main goroutine の分を除く.
				wg.Done()
			}()

		LOOP:
			for targetItem := range chItems {
				select {
				case <-ctx.Done():
					logger.Debug(fmt.Sprintf("[WORKER %d] Received Item(%d) to process but cancellation has been requested, so discard it and terminate this worker", workerID, targetItem))
					break LOOP
				default:
				}

				logger.Debug(fmt.Sprintf("[WORKER %d] Receved Item(%d) to process", workerID, targetItem))
				result, err := item.ProcessItem(targetItem)
				taskResult := &taskResult{result: result, err: err}

				select {
				case <-ctx.Done():
					logger.Debug(fmt.Sprintf("[WORKER %d] Processed Item(%d) but cancellation has been requested, so discard it", workerID, targetItem))
				default:
					logger.Debug(fmt.Sprintf("[WORKER %d] Processed Item(%d), pushing the result", workerID, targetItem))
					chTaskResults <- taskResult
				}
			}
		}(i + 1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		processedCount := 0
		failed := false

	LOOP:
		for processedCount < numItems {
			select {
			case <-ctx.Done():
				logger.Debug("Has already detected an error, so terminate the tasks")
				break LOOP
			case result := <-chTaskResults:
				if result.err != nil {
					cancel()
					logger.Error(fmt.Sprintf("Processing failed: %s", result.err.Error()))
					failed = true
					break LOOP // 1個でも失敗したら処理を抜ける
				}

				logger.Info(result.result)

				processedCount++

				if processedCount%bufferSize == 0 {
					logger.Info("Flush buffer!")
				}
			}
		}

		if !failed && processedCount%bufferSize > 0 {
			logger.Info("Flush buffer!")
		}
	}()

	for _, targetItem := range items {
		select {
		case <-ctx.Done():
		default:
			chItems <- targetItem
		}
	}
	close(chItems)

	wg.Wait()
}
