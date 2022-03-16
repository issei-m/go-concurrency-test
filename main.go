package main

import (
	"flag"
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

	var (
		isVerbose   = flag.Bool("verbose", false, "whether outputs debug log or not")
		concurrency = flag.Int("concurrency", 0, "the number of concurrency to process items (2 ~ 100)")
	)

	flag.Parse()

	if *concurrency != 0 && (1 >= *concurrency || *concurrency > 100) {
		logger.Error("--concurrency must be either 0 or between 2 and 100")
		os.Exit(1)
	}

	if *isVerbose {
		logger.EnableDebugLog()
	}

	// 100個分の item.Item
	targetItems := item.CreateItems(100)

	fmt.Println("Processing start!")

	if *concurrency == 0 {
		app.ProcessItems(targetItems)
	} else {
		app.ProcessItemsConcurrently(targetItems, *concurrency)
	}

	fmt.Println("DONE!")
}
