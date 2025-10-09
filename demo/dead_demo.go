package main

import (
	"fmt"
	"sync"
	"time"
	"zyj.com/golang-study/goroutine"
)

func main() {
	// 初始化检测器
	var deadlockDetector = goroutine.NewDeadlockDetector(
		10*time.Second, // 检查间隔
		30*time.Second, // 最大阻塞时间
		1000,           // 最大 goroutine 数
	)
	loopDetector := goroutine.NewLoopDetector(
		5*time.Second,  // 检查间隔
		10*time.Second, // 最大循环时间
	)
	// 启动检测器
	deadlockDetector.Start()
	loopDetector.Start()

	defer func() {
		deadlockDetector.Stop()
		loopDetector.Stop()
	}()

	// 演示各种场景
	demoNormalOperations()
	demoPotentialDeadlock()
	demoInfiniteLoop(loopDetector)

	// 保持程序运行
	select {}
	time.Sleep(1000 * time.Second)

	deadlockDetector.Stop()
	loopDetector.Stop()
}

func demoNormalOperations() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Second)
			fmt.Printf("Worker %d done\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}

func demoPotentialDeadlock() {
	var mu1, mu2 sync.Mutex

	// 潜在死锁场景
	go func() {
		mu1.Lock()
		time.Sleep(100 * time.Millisecond)
		mu2.Lock() // 这里可能会死锁
		fmt.Println("Acquired both locks in goroutine 1")
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		mu2.Lock()
		time.Sleep(100 * time.Millisecond)
		mu1.Lock() // 这里可能会死锁
		fmt.Println("Acquired both locks in goroutine 2")
		mu1.Unlock()
		mu2.Unlock()
	}()
}

func demoInfiniteLoop(loopDetector *goroutine.LoopDetector) {
	// 注册循环监控
	//loopChecker := loopMonitor.RegisterLoop("demoInfiniteLoop")

	go func() {
		loopDetector.MonitorFunction("demoInfiniteLoop")

		counter := 0
		for {
			// 模拟工作
			time.Sleep(10 * time.Millisecond)

			// 检查循环次数
			//loopChecker()

			counter++
			if counter%1000 == 0 {
				// 正常循环中定期更新检测器
				loopDetector.UpdateFunction("demoInfiniteLoop")
			}

			// 模拟长时间运行（但非死循环）
			if counter > 5000 {
				break
			}
		}
		fmt.Println("Loop completed normally")
	}()

	// 真正的潜在死循环
	go func() {
		loopDetector.MonitorFunction("realInfiniteLoop")

		for {
			// 没有 break 条件的循环
			// 这会被死循环检测器捕获
			time.Sleep(1 * time.Millisecond)

			// 注意：这里没有调用 UpdateFunction
			// 所以检测器会认为这个函数长时间未更新
		}
	}()
}
