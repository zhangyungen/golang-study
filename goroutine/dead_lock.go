package goroutine

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// DeadlockDetector 死锁检测器
type DeadlockDetector struct {
	checkInterval      time.Duration
	maxGoroutines      int
	maxBlockingTime    time.Duration
	enabled            bool
	goroutineSnapshots map[string]*GoroutineSnapshot
	snapshotMutex      sync.RWMutex
}

// GoroutineSnapshot goroutine 快照
type GoroutineSnapshot struct {
	StackTrace string
	FirstSeen  time.Time
	LastSeen   time.Time
	Count      int
}

// NewDeadlockDetector 创建死锁检测器
func NewDeadlockDetector(checkInterval, maxBlockingTime time.Duration, maxGoroutines int) *DeadlockDetector {
	return &DeadlockDetector{
		checkInterval:      checkInterval,
		maxGoroutines:      maxGoroutines,
		maxBlockingTime:    maxBlockingTime,
		enabled:            true,
		goroutineSnapshots: make(map[string]*GoroutineSnapshot),
	}
}

// Start 启动死锁检测
func (d *DeadlockDetector) Start() {
	if !d.enabled {
		return
	}

	fmt.Printf("🚀 Deadlock detector started (check interval: %v, max goroutines: %d)\n",
		d.checkInterval, d.maxGoroutines)

	go d.monitor()
}

// Stop 停止检测
func (d *DeadlockDetector) Stop() {
	d.enabled = false
}

// monitor 监控协程状态
func (d *DeadlockDetector) monitor() {
	ticker := time.NewTicker(d.checkInterval)
	defer ticker.Stop()

	for d.enabled {
		select {
		case <-ticker.C:
			d.checkDeadlock()
		}
	}
}

// checkDeadlock 检查死锁
func (d *DeadlockDetector) checkDeadlock() {
	// 获取当前 goroutine 数量
	numGoroutines := runtime.NumGoroutine()

	if numGoroutines > d.maxGoroutines {
		fmt.Printf("⚠️  Possible deadlock detected: too many goroutines (%d > %d)\n",
			numGoroutines, d.maxGoroutines)
		d.dumpGoroutines()
	}

	// 检查是否有 goroutine 长时间阻塞
	d.checkLongRunningGoroutines()
}

// dumpGoroutines 输出所有 goroutine 堆栈
func (d *DeadlockDetector) dumpGoroutines() {
	buf := make([]byte, 1024*1024) // 1MB buffer
	n := runtime.Stack(buf, true)

	fmt.Printf("=== Goroutine Dump (%d goroutines) ===\n", runtime.NumGoroutine())
	fmt.Println(string(buf[:n]))
	fmt.Println("=== End Goroutine Dump ===")
}

// checkLongRunningGoroutines 检查长时间运行的 goroutine
func (d *DeadlockDetector) checkLongRunningGoroutines() {
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	stack := string(buf[:n])

	goroutines := strings.Split(stack, "goroutine ")[1:]

	d.snapshotMutex.Lock()
	defer d.snapshotMutex.Unlock()

	currentTime := time.Now()
	currentSnapshots := make(map[string]bool)

	for _, goroutine := range goroutines {
		lines := strings.SplitN(goroutine, "\n", 2)
		if len(lines) < 2 {
			continue
		}

		// 提取关键堆栈信息作为标识
		stackKey := d.extractStackKey(lines[1])
		if stackKey == "" {
			continue
		}

		currentSnapshots[stackKey] = true

		if snapshot, exists := d.goroutineSnapshots[stackKey]; exists {
			snapshot.LastSeen = currentTime
			snapshot.Count++

			// 检查是否长时间阻塞
			if currentTime.Sub(snapshot.FirstSeen) > d.maxBlockingTime {
				fmt.Printf("🚨 Potential deadlock detected - goroutine blocked for %v\n",
					currentTime.Sub(snapshot.FirstSeen))
				fmt.Printf("Stack trace:\n%s\n", lines[1])

				// 重置检测，避免重复报警
				delete(d.goroutineSnapshots, stackKey)
			}
		} else {
			// 新的 goroutine 堆栈
			d.goroutineSnapshots[stackKey] = &GoroutineSnapshot{
				StackTrace: lines[1],
				FirstSeen:  currentTime,
				LastSeen:   currentTime,
				Count:      1,
			}
		}
	}

	// 清理过期的快照
	for key, snapshot := range d.goroutineSnapshots {
		if !currentSnapshots[key] {
			delete(d.goroutineSnapshots, key)
		} else if currentTime.Sub(snapshot.LastSeen) > d.checkInterval*5 {
			// 长时间未更新的快照也清理掉
			delete(d.goroutineSnapshots, key)
		}
	}
}

// extractStackKey 提取堆栈关键信息作为标识
func (d *DeadlockDetector) extractStackKey(stackTrace string) string {
	lines := strings.Split(stackTrace, "\n")

	// 取前几行作为关键信息
	keyLines := []string{}
	for i, line := range lines {
		if i >= 6 { // 最多取6行
			break
		}
		if strings.Contains(line, "created by") {
			break
		}
		keyLines = append(keyLines, strings.TrimSpace(line))
	}

	return strings.Join(keyLines, "|")
}
