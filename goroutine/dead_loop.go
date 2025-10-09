package goroutine

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LoopDetector 死循环检测器
type LoopDetector struct {
	checkInterval  time.Duration
	maxLoopTime    time.Duration
	enabled        bool
	monitoredFuncs map[string]time.Time
	mutex          sync.RWMutex
}

// NewLoopDetector 创建死循环检测器
func NewLoopDetector(checkInterval, maxLoopTime time.Duration) *LoopDetector {
	return &LoopDetector{
		checkInterval:  checkInterval,
		maxLoopTime:    maxLoopTime,
		enabled:        true,
		monitoredFuncs: make(map[string]time.Time),
	}
}

// Start 启动死循环检测
func (l *LoopDetector) Start() {
	if !l.enabled {
		return
	}

	fmt.Printf("🔍 Loop detector started (check interval: %v, max loop time: %v)\n",
		l.checkInterval, l.maxLoopTime)

	go l.monitor()
}

// Start 启动死循环检测
func (l *LoopDetector) Stop() {
	l.enabled = false
}

// MonitorFunction 监控函数执行
func (l *LoopDetector) MonitorFunction(funcName string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.monitoredFuncs[funcName] = time.Now()
}

// UpdateFunction 更新函数执行时间
func (l *LoopDetector) UpdateFunction(funcName string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.monitoredFuncs[funcName] = time.Now()
}

// monitor 监控循环
func (l *LoopDetector) monitor() {
	ticker := time.NewTicker(l.checkInterval)
	defer ticker.Stop()

	for l.enabled {
		select {
		case <-ticker.C:
			l.checkInfiniteLoops()
		}
	}
}

// checkInfiniteLoops 检查死循环
func (l *LoopDetector) checkInfiniteLoops() {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	currentTime := time.Now()

	for funcName, lastUpdate := range l.monitoredFuncs {
		if currentTime.Sub(lastUpdate) > l.maxLoopTime {
			fmt.Printf("🚨 Potential infinite loop detected in function: %s (running for %v)\n",
				funcName, currentTime.Sub(lastUpdate))

			// 输出相关 goroutine 堆栈
			l.dumpRelevantStacks(funcName)
		}
	}
}

// dumpRelevantStacks 输出相关堆栈
func (l *LoopDetector) dumpRelevantStacks(funcName string) {
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	stack := string(buf[:n])

	goroutines := strings.Split(stack, "goroutine ")

	for _, goroutine := range goroutines[1:] {
		if strings.Contains(goroutine, funcName) {
			lines := strings.SplitN(goroutine, "\n", 10) // 取前10行
			fmt.Printf("Relevant stack trace:\n%s\n", strings.Join(lines, "\n"))
			break
		}
	}
}
