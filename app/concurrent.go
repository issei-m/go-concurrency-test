package app

import (
	"context"
	"fmt"
	"go_concurrency_test/item"
	"go_concurrency_test/logger"
	"sync/atomic"
)

type taskResult struct {
	result string
	err    error
}

func ProcessItemsConcurrently(items []item.Item, concurrency int) {
	var chItems = make(chan item.Item)
	var chTaskResults = make(chan *taskResult)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		close(chTaskResults)
		cancel()
	}()

	numWorkers := int32(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			defer func() {
				atomic.AddInt32(&numWorkers, -1)
				logger.Debug(fmt.Sprintf("[WORKER %d] Terminated, bye (remaining %d worker(s))", workerID, numWorkers))
			}()

			for targetItem := range chItems {
				select {
				case <-ctx.Done():
					logger.Debug(fmt.Sprintf("[WORKER %d] Received Item(%d) to process but cancellation has been requested, so discard it and terminate this worker", workerID, targetItem))
				default:
					logger.Debug(fmt.Sprintf("[WORKER %d] Receved Item(%d) to process", workerID, targetItem))
					result, err := item.ProcessItem(targetItem)
					chTaskResults <- &taskResult{result: result, err: err}
				}
			}
		}(i + 1)
	}

	go func() {
	LOOP:
		for _, targetItem := range items {
			select {
			case <-ctx.Done():
				logger.Debug(fmt.Sprintf("Cancellation has already been requested, so no more item will be pushed"))
				break LOOP
			default:
				chItems <- targetItem
			}
		}
		close(chItems)
	}()

	processedCount := 0
	failed := false

	for numWorkers > 0 {
		select {
		case result := <-chTaskResults:
			if failed {
				break
			}

			if result.err != nil {
				// 1件でもエラーが起きたら以降の処理をキャンセルする.
				// worker が並行に動いている為、タイミングによってキャンセル後もいくつかの処理は実行されうるので、その時点でこの goroutine を閉じてしまうと deadlock になるので注意.
				// 従って、 chTaskResults の受信は worker が全ていなくなるまで続ける.

				cancel()
				logger.Error(fmt.Sprintf("Processing failed: %s", result.err.Error()))
				failed = true
				break
			}

			logger.Info(result.result)

			processedCount++

			if processedCount%bufferSize == 0 {
				logger.Info("Flush buffer!")
			}
		default:
			// `default` が無いと、タイミングによっては worker が全て終了した後に chTaskResults から受信を試みる (=deadlock) 事になるので注意
		}
	}

	if !failed && processedCount%bufferSize > 0 {
		logger.Info("Flush buffer!")
	}
}
