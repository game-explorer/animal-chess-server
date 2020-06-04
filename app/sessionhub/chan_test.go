package sessionhub

import (
	"github.com/game-explorer/animal-chess-server/internal/pkg/log"
	"testing"
	"time"
)

// 测试 在放入管道的同时修改管道
// 结果: 已经阻塞到select中的管道依然保持原样(阻塞), 等到下一次写入时才是新的管道生效.
func TestSetChanWrite(t *testing.T) {
	var c chan int

	go func() {
		for range time.Tick(1 * time.Second) {
			select {
			case c <- 1:
				log.Infof("c")
			case <-time.After(5 * time.Second):
				log.Infof("timeout")
			}
		}
	}()

	// 先 timeout
	// 后 c
	// 再 default
	go func() {
		time.Sleep(2 * time.Second)
		c = make(chan int, 10)
		time.Sleep(7 * time.Second)
		c = nil
	}()

	select {}
}
