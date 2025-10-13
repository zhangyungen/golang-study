package tssync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ==================== 对比实现：原生全局锁 ====================
type globalMutex struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewGlobalMutex() *globalMutex {
	return &globalMutex{
		data: make(map[string]interface{}),
	}
}

func (g *globalMutex) Lock(key string) {
	g.mu.Lock()
}

func (g *globalMutex) Unlock(key string) {
	g.mu.Unlock()
}

func (g *globalMutex) RLock(key string) {
	g.mu.RLock()
}

func (g *globalMutex) RUnlock(key string) {
	g.mu.RUnlock()
}

// ==================== 对比实现：分段锁 ====================
type shardedMutex struct {
	shards []*shard
	count  uint32
}

type shard struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewShardedMutex(shardCount uint32) *shardedMutex {
	sm := &shardedMutex{
		shards: make([]*shard, shardCount),
		count:  shardCount,
	}
	for i := range sm.shards {
		sm.shards[i] = &shard{
			data: make(map[string]interface{}),
		}
	}
	return sm
}

func (s *shardedMutex) getShard(key string) *shard {
	// 简单的哈希函数
	hash := uint32(0)
	for i := 0; i < len(key); i++ {
		hash = hash*31 + uint32(key[i])
	}
	return s.shards[hash%s.count]
}

func (s *shardedMutex) Lock(key string) {
	s.getShard(key).mu.Lock()
}

func (s *shardedMutex) Unlock(key string) {
	s.getShard(key).mu.Unlock()
}

func (s *shardedMutex) RLock(key string) {
	s.getShard(key).mu.RLock()
}

func (s *shardedMutex) RUnlock(key string) {
	s.getShard(key).mu.RUnlock()
}

// ==================== 性能测试框架 ====================
type testConfig struct {
	name        string
	goroutines  int
	operations  int
	keyCount    int
	readRatio   float64 // 读操作比例
	hotspotRate float64 // 热点key访问比例
}

type testResult struct {
	name       string
	totalOps   int64
	successOps int64
	totalTime  time.Duration
	avgLatency time.Duration
	minLatency time.Duration
	maxLatency time.Duration
	qps        float64
}

// 性能测试接口
type locker interface {
	Lock(key string)
	Unlock(key string)
	RLock(key string)
	RUnlock(key string)
}

func runPerformanceTest(locker locker, config testConfig) testResult {
	var totalOps, successOps int64
	var minLatency, maxLatency time.Duration
	var totalLatency time.Duration

	startTime := time.Now()
	var wg sync.WaitGroup

	// 生成测试keys
	keys := generateKeys(config.keyCount)

	// 用于统计的channel
	results := make(chan time.Duration, config.goroutines*config.operations)

	for i := 0; i < config.goroutines; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			localOps := int64(0)
			localSuccess := int64(0)

			for j := 0; j < config.operations; j++ {
				// 选择key（模拟热点访问）
				key := selectKey(keys, config.hotspotRate)

				opStart := time.Now()

				// 根据读写比例执行操作
				if float64(j%100) < config.readRatio*100 {
					// 读操作
					locker.RLock(key)
					// 模拟读操作耗时
					time.Sleep(10 * time.Nanosecond)
					locker.RUnlock(key)
					localSuccess++
				} else {
					// 写操作
					locker.Lock(key)
					// 模拟写操作耗时
					time.Sleep(20 * time.Nanosecond)
					locker.Unlock(key)
					localSuccess++
				}

				latency := time.Since(opStart)
				results <- latency
				localOps++
			}

			atomic.AddInt64(&totalOps, localOps)
			atomic.AddInt64(&successOps, localSuccess)
		}(i)
	}

	wg.Wait()
	close(results)
	totalTime := time.Since(startTime)

	// 处理延迟统计
	var count int64
	for latency := range results {
		count++
		totalLatency += latency

		if count == 1 {
			minLatency = latency
			maxLatency = latency
		} else {
			if latency < minLatency {
				minLatency = latency
			}
			if latency > maxLatency {
				maxLatency = latency
			}
		}
	}

	avgLatency := time.Duration(0)
	if count > 0 {
		avgLatency = time.Duration(int64(totalLatency) / count)
	}

	return testResult{
		name:       config.name,
		totalOps:   totalOps,
		successOps: successOps,
		totalTime:  totalTime,
		avgLatency: avgLatency,
		minLatency: minLatency,
		maxLatency: maxLatency,
		qps:        float64(totalOps) / totalTime.Seconds(),
	}
}

func generateKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	return keys
}

func selectKey(keys []string, hotspotRate float64) string {
	// 模拟热点访问：部分key被更频繁地访问
	if hotspotRate > 0 && len(keys) > 10 {
		hotspotCount := int(float64(len(keys)) * hotspotRate)
		if hotspotCount < 1 {
			hotspotCount = 1
		}
		// 80%的请求访问前hotspotCount个key
		return keys[int(float64(hotspotCount))]
	}
	// 均匀分布
	return keys[len(keys)/2]
}

// ==================== 基准测试函数 ====================
func BenchmarkLockers(b *testing.B) {
	configs := []testConfig{
		{
			name:        "低并发-多key",
			goroutines:  10,
			operations:  1000,
			keyCount:    100,
			readRatio:   0.8,
			hotspotRate: 0.1,
		},
		{
			name:        "高并发-少key",
			goroutines:  400,
			operations:  400,
			keyCount:    50,
			readRatio:   0.8,
			hotspotRate: 0.5,
		},
		{
			name:        "高并发-多key",
			goroutines:  100,
			operations:  400,
			keyCount:    100,
			readRatio:   0.9,
			hotspotRate: 0.1,
		},
		{
			name:        "写密集-多key",
			goroutines:  30,
			operations:  300,
			keyCount:    50,
			readRatio:   0.2,
			hotspotRate: 0.2,
		},
	}

	for _, config := range configs {
		b.Run(fmt.Sprintf("MiniMutex-%s", config.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// 包装miniMutex实现locker接口
				miniLocker := &miniMutexLocker{}
				runPerformanceTest(miniLocker, config)
			}
		})

		b.Run(fmt.Sprintf("GlobalMutex-%s", config.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				globalLocker := NewGlobalMutex()
				runPerformanceTest(globalLocker, config)
			}
		})

		b.Run(fmt.Sprintf("ShardedMutex-%s", config.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shardedLocker := NewShardedMutex(32)
				runPerformanceTest(shardedLocker, config)
			}
		})
	}
}

// miniMutex的locker接口适配器
type miniMutexLocker struct{}

func (m *miniMutexLocker) Lock(key string) {
	Lock(key)
}

func (m *miniMutexLocker) Unlock(key string) {
	UnLock(key)
}

func (m *miniMutexLocker) RLock(key string) {
	// miniMutex没有区分读写锁，使用写锁代替
	RLock(key)
}

func (m *miniMutexLocker) RUnlock(key string) {
	RUnlock(key)
}

// ==================== 性能对比测试 ====================
func TestPerformanceComparison(t *testing.T) {
	testScenarios := []testConfig{
		{
			name:        "读多写少-低竞争",
			goroutines:  20,
			operations:  500,
			keyCount:    100,
			readRatio:   0.9,
			hotspotRate: 0.1,
		},
		{
			name:        "读写均衡-高竞争",
			goroutines:  4000,
			operations:  1000,
			keyCount:    20,
			readRatio:   0.6,
			hotspotRate: 0.5,
		},
		{
			name:        "写密集-中等竞争",
			goroutines:  2000,
			operations:  5000,
			keyCount:    30,
			readRatio:   0,
			hotspotRate: 0.5,
		},
	}

	fmt.Println("🔬 锁性能对比测试结果")
	fmt.Println("=" + repeatString("=", 78))

	for _, scenario := range testScenarios {
		fmt.Printf("\n📊 测试场景: %s\n", scenario.name)
		fmt.Printf("配置: goroutines=%d, operations=%d, keys=%d, 读比例=%.1f, 热点率=%.1f\n",
			scenario.goroutines, scenario.operations, scenario.keyCount,
			scenario.readRatio, scenario.hotspotRate)

		// 测试三种锁实现
		miniLocker := &miniMutexLocker{}
		globalLocker := NewGlobalMutex()
		shardedLocker := NewShardedMutex(32)

		miniResult := runPerformanceTest(miniLocker, scenario)
		globalResult := runPerformanceTest(globalLocker, scenario)
		shardedResult := runPerformanceTest(shardedLocker, scenario)

		printResult("您的 MiniMutex", miniResult)
		printResult("全局锁 GlobalMutex", globalResult)
		printResult("分段锁 ShardedMutex", shardedResult)

		// 性能对比分析
		miniQPS := miniResult.qps
		globalQPS := globalResult.qps
		shardedQPS := shardedResult.qps

		fmt.Printf("\n📈 性能对比分析:\n")
		fmt.Printf("MiniMutex 相对于 GlobalMutex 的性能提升: %.2f%%\n",
			(miniQPS/globalQPS-1)*100)
		fmt.Printf("MiniMutex 相对于 ShardedMutex 的性能差异: %.2f%%\n",
			(miniQPS/shardedQPS-1)*100)

		fmt.Println("─" + repeatString("─", 78))
	}
}

func printResult(name string, result testResult) {
	fmt.Printf("%-20s: QPS=%.0f, 平均延迟=%-10v, 成功率=%.1f%%, 总操作=%d\n",
		name, result.qps, result.avgLatency,
		float64(result.successOps)/float64(result.totalOps)*100,
		result.totalOps)
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// ==================== 并发安全测试 ====================
func TestConcurrentSafety(t *testing.T) {
	fmt.Println("\n🔒 并发安全测试")

	// 测试数据竞争
	const iterations = 10000
	key := "test_key"

	// 使用miniMutex保护共享计数器
	counter := 0
	var wg sync.WaitGroup
	wg.Add(2)

	start := time.Now()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			Lock(key)
			counter++
			UnLock(key)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			Lock(key)
			counter--
			UnLock(key)
		}
	}()

	wg.Wait()
	duration := time.Since(start)

	if counter != 0 {
		t.Errorf("❌ 并发安全测试失败: 期望counter=0, 实际counter=%d", counter)
	} else {
		fmt.Printf("✅ 并发安全测试通过: 操作数=%d, 耗时=%v, QPS=%.0f\n",
			iterations*2, duration, float64(iterations*2)/duration.Seconds())
	}
}

func main() {
	fmt.Println("🚀 Go锁性能压测工具")
	fmt.Println()

	// 运行性能测试
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{Name: "TestPerformanceComparison", F: TestPerformanceComparison},
			{Name: "TestConcurrentSafety", F: TestConcurrentSafety},
		},
		[]testing.InternalBenchmark{
			{Name: "BenchmarkLockers", F: BenchmarkLockers},
		},
		nil)
}
