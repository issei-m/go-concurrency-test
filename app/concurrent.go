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
	defer cancel()

	numWorkers := int32(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			var pulledCount int

			defer func() {
				atomic.AddInt32(&numWorkers, -1)
				logger.Debug(fmt.Sprintf("[WORKER %d] Terminated, bye (pulled %d items, remaining %d worker(s))", workerID, pulledCount, numWorkers))
			}()

			for targetItem := range chItems {
				pulledCount++

				select {
				case <-ctx.Done():
					logger.Debug(fmt.Sprintf("[WORKER %d] Received Item(%d) to process but cancellation has been requested, so terminating this worker", workerID, targetItem))
					return
				default:
					logger.Debug(fmt.Sprintf("[WORKER %d] Received Item(%d) to process", workerID, targetItem))
					result, err := item.ProcessItem(targetItem)
					chTaskResults <- &taskResult{result: result, err: err}
				}
			}
		}(i + 1)
	}

	go func() {
		defer close(chItems)

		for _, targetItem := range items {
			select {
			case <-ctx.Done():
				logger.Debug(fmt.Sprintf("Cancellation has already been requested, so no more item will be pushed"))
				return
			default:
				logger.Debug(fmt.Sprintf("Pushing Item(%d)", targetItem))
				chItems <- targetItem
			}
		}
	}()

	processedCount := 0
	failed := false

	for numWorkers > 0 {
		select {
		case result := <-chTaskResults:
			if failed {
				logger.Debug(fmt.Sprintf("Received the result but cancellation has already been requested, so discard it and wait for all workers' termination"))
				break
			}

			if result.err != nil {
				// We will request cancellation of later tasks as soon as detected an error.
				// To avoid a deadlock, we have to wait for all workers' termination discarding the result pulled from the channel.
				// Since workers are running concurrently, "item processing" could be done even after the cancellation request but there will be no receivers of chTaskResults if we broke the loop.
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
			// Forgetting to put `default` can lead a deadlock since there might be the case of attempt to receive from chTaskResults after all workers have terminated (i.e. no senders)
		}
	}

	if !failed && processedCount%bufferSize > 0 {
		logger.Info("Flush buffer!")
	}

	logger.Debug(fmt.Sprintf("Succeeded count: %d", processedCount))
}
