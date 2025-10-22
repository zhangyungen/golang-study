package test

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"testing"
	"time"
)

func TestAntsRun(t *testing.T) {
	// 延迟关闭默认池
	p, _ := ants.NewPool(10000)
	defer p.Release()
	// 提交任务
	for i := 0; i < 2; i++ {
		n := i
		// 提交任务
		err := p.Submit(func() {
			time.Sleep(time.Second * 3)
			// 故意抛出错误
			if n > 0 {
				panic("运行遇到错误~")
			}
			fmt.Println("run end time: ", time.Now().Format("2006-01-02 15:04:05"))
		})
		if err != nil {
			fmt.Println("submit err:", err)
		}
	}
	// 打印当前运行的协程数量
	fmt.Println("run go num: ", p.Running())
	fmt.Println("cap go num: ", p.Cap())
	// 主动等待协程运行完成
	time.Sleep(time.Second * 10)

	fmt.Println("run go num: ", p.Running())
	fmt.Println("cap go num: ", p.Cap())

	fmt.Println("finish")
}

/**运行输出
=== RUN   TestAntsRun
run go num:  2
time:  2022-04-28 17:58:57
2022/04/28 17:58:57 worker exits from a panic: 运行遇到错误~
2022/04/28 17:58:57 worker exits from panic: goroutine 9 [running]:
....

finish
--- PASS: TestAntsRun (10.00s)
PASS
*/
