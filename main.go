package main

import (
	"fmt"
	"go_concurrency_test/app"
	"go_concurrency_test/item"
	"go_concurrency_test/logger"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if isDebug() {
		logger.EnableDebugLog()
	}

	// 100個分の item.Item
	targetItems := item.CreateItems(100)

	fmt.Println("Processing start!")

	app.ProcessItems(targetItems)

	fmt.Println("DONE!")
}

func isDebug() bool {
	for _, arg := range os.Args {
		if arg == "--debug" || arg == "-d" {
			return true
		}
	}

	return false
}
