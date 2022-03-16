package item

import (
	"fmt"
	"math/rand"
	"time"
)

// ProcessItem processes the given Item and returns the result as a string when succeeded, or an error when failed.
// And this function sleeps for random milliseconds of time to pretend to do some heavy blocking I/O process.
func ProcessItem(item Item) (string, error) {
	// Sleep time is range of 100 ~ 1000 ms
	sleepTime := rand.Intn(900) + 100
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	if rand.Intn(100) == 0 {
		return "", fmt.Errorf("%d (something went wrong)", item)
	} else {
		return fmt.Sprintf("Processed successfully: %d (elapsed: %d msec)", item, sleepTime), nil
	}
}
