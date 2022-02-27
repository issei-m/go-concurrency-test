package item

import (
	"fmt"
	"math/rand"
	"time"
)

// ProcessItem 引数の Item を処理する. 処理に成功した場合は結果を表す内容を string で、失敗した場合は error を返す.
// 関数内で時間の掛かるブロッキング I/O 処理が走る事を想定している.
func ProcessItem(item Item) (string, error) {
	// ブロッキング I/O のシミュレーション. 処理時間はランダムに 100 ~ 1000 ms とする
	sleepTime := rand.Intn(900) + 100
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	if rand.Intn(100) == 0 {
		return "", fmt.Errorf("%d (something went wrong)", item)
	} else {
		return fmt.Sprintf("Processed successfully: %d (elapsed: %d msec)", item, sleepTime), nil
	}
}
