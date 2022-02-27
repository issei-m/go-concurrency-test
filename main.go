package main

import (
	"fmt"
	"go_concurrency_test/app"
	"go_concurrency_test/item"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 100個分の item.Item
	targetItems := item.CreateItems(100)

	fmt.Println("Processing start!")

	app.ProcessItems(targetItems)

	fmt.Println("DONE!")
}
