package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/maypok86/otter/v2/stats"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maypok86/otter/v2"
)

// 缓存操作类型
type OperationType string

const (
	OperationDelete OperationType = "delete"
	OperationSet    OperationType = "set"
)

// 缓存同步消息结构
type SyncMessage struct {
	Operation  OperationType `json:"operation"`
	Key        string        `json:"key"`
	Value      string        `json:"value,omitempty"` // 仅set操作需要
	InstanceID string        `json:"instance_id"`     // 消息来源实例ID
	Timestamp  int64         `json:"timestamp"`       // 消息时间戳
}

// 二级缓存配置
type Config struct {
	// Otter 一级缓存配置
	OtterMaxSize int           `json:"otterMaxSize"`
	OtterTTL     time.Duration `json:"otterTTL"`

	// Redis 二级缓存配置
	RedisAddr     string        `json:"redisAddr"`
	RedisPassword string        `json:"redisPassword"`
	RedisDB       int           `json:"redisDB"`
	RedisTTL      time.Duration `json:"redisTTL"`

	// Pub/Sub 配置
	PubSubChannel string `json:"pubSubChannel"`
	InstanceID    string `json:"instanceID"` // 当前实例标识
}

// 分布式二级缓存主结构
type DistributedCache struct {
	primaryCache   *otter.Cache[string, string] // 一级缓存 (Otter)
	secondaryCache *redis.Client                // 二级缓存 (Redis)
	pubSub         *redis.PubSub                // Redis Pub/Sub
	config         Config
	cacheCtx       context.Context
	cancel         context.CancelFunc
	syncMutex      sync.RWMutex    // 同步操作锁
	localOnlyKeys  map[string]bool // 仅本地操作标记
	MsgSendCount   atomic.Uint64
	MsgRecvdCount  atomic.Uint64
}

// 创建分布式缓存实例
func NewDistributedCache(ctx context.Context, config Config) (*DistributedCache, error) {
	// 创建带取消的上下文
	cacheCtx, cancel := context.WithCancel(ctx)

	// 初始化一级缓存 (Otter)
	primaryCache := otter.Must(&otter.Options[string, string]{
		MaximumSize:     config.OtterMaxSize,
		InitialCapacity: 100,
		Weigher: func(key string, value string) uint32 {
			return uint32(len(key) + len(value))
		},
		ExpiryCalculator:  otter.ExpiryAccessing[string, string](config.OtterTTL),
		RefreshCalculator: otter.RefreshWriting[string, string](config.OtterTTL),
		StatsRecorder:     stats.NewCounter(),
	})

	// 初始化二级缓存 (Redis)
	secondaryCache := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// 测试Redis连接
	if err := secondaryCache.Ping(cacheCtx).Err(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	// 创建Pub/Sub订阅
	pubSub := secondaryCache.Subscribe(cacheCtx, config.PubSubChannel)

	cache := &DistributedCache{
		primaryCache:   primaryCache,
		secondaryCache: secondaryCache,
		pubSub:         pubSub,
		config:         config,
		cacheCtx:       cacheCtx,
		cancel:         cancel,
		localOnlyKeys:  make(map[string]bool),
	}

	// 启动消息监听goroutine
	go cache.listenForSyncMessages()

	return cache, nil
}

// 关闭缓存实例
func (dc *DistributedCache) Close() error {
	dc.cancel()

	if err := dc.pubSub.Close(); err != nil {
		log.Printf("Error closing pubsub: %v", err)
	}
	//  关闭缓存
	dc.primaryCache.CleanUp()
	dc.primaryCache = nil
	return dc.secondaryCache.Close()
}

// 监听同步消息
func (dc *DistributedCache) listenForSyncMessages() {
	channel := dc.pubSub.Channel()
	for {
		select {
		case msg, ok := <-channel:
			if !ok {
				return
			}
			dc.handleSyncMessage(msg)
		case <-dc.cacheCtx.Done():
			return
		}
	}
}

// 处理同步消息
func (dc *DistributedCache) handleSyncMessage(msg *redis.Message) {
	var syncMsg SyncMessage
	if err := json.Unmarshal([]byte(msg.Payload), &syncMsg); err != nil {
		log.Printf("Failed to unmarshal sync message: %v", err)
		return
	}
	dc.MsgRecvdCount.Add(1)
	// 忽略自己发送的消息
	if syncMsg.InstanceID == dc.config.InstanceID {
		return
	}
	// 处理消息
	switch syncMsg.Operation {
	case OperationDelete:
		dc.silentDelete(syncMsg.Key)
	case OperationSet:
		dc.silentSet(syncMsg.Key, syncMsg.Value)
	}
}

// 静默删除（不触发消息广播）
func (dc *DistributedCache) silentDelete(key string) {
	dc.syncMutex.Lock()
	defer dc.syncMutex.Unlock()

	// 标记为本地操作，避免循环广播
	dc.localOnlyKeys[key] = true
	dc.primaryCache.Invalidate(key)
	delete(dc.localOnlyKeys, key)
}

// 静默设置（不触发消息广播）
func (dc *DistributedCache) silentSet(key, value string) {
	dc.syncMutex.Lock()
	defer dc.syncMutex.Unlock()

	// 标记为本地操作，避免循环广播
	dc.localOnlyKeys[key] = true
	dc.primaryCache.Set(key, value)
	delete(dc.localOnlyKeys, key)
}

// 发送同步消息
func (dc *DistributedCache) sendSyncMessage(operation OperationType, key, value string) error {
	syncMsg := SyncMessage{
		Operation:  operation,
		Key:        key,
		Value:      value,
		InstanceID: dc.config.InstanceID,
		Timestamp:  time.Now().UnixNano(),
	}

	message, err := json.Marshal(syncMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal sync message: %v", err)
	}

	dc.MsgSendCount.Add(1)
	return dc.secondaryCache.Publish(dc.cacheCtx, dc.config.PubSubChannel, message).Err()
}

// get 方法：二级缓存读取
func (dc *DistributedCache) get(key string) (string, error) {
	// 1. 尝试从一级缓存获取
	if value, err := dc.primaryCache.Get(dc.cacheCtx, key, nil); err != nil {
		return value, nil
	}

	// 2. 尝试从二级缓存获取
	value, err := dc.secondaryCache.Get(dc.cacheCtx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("redis error: %v", err)
	}

	// 3. 回填一级缓存（静默操作，不广播）
	dc.silentSet(key, value)

	return value, nil
}

// Set 方法：二级缓存写入
func (dc *DistributedCache) set(key, value string) error {
	// 1. 先写入二级缓存（确保数据持久化）
	if err := dc.secondaryCache.Set(dc.cacheCtx, key, value, dc.config.RedisTTL).Err(); err != nil {
		return fmt.Errorf("failed to set secondary cache: %v", err)
	}

	// 2. 写入本地一级缓存
	dc.primaryCache.Set(key, value)

	// 3. 广播设置消息到其他实例（如果不是本地操作触发的）
	dc.syncMutex.RLock()
	isLocalOnly := dc.localOnlyKeys[key]
	dc.syncMutex.RUnlock()

	if !isLocalOnly {
		// 发送同步消息建议是-删除消息，缓存重新取二级缓存数据
		if err := dc.sendSyncMessage(OperationDelete, key, value); err != nil {
			log.Printf("Failed to send sync message: %v", err)
		}
	}

	return nil
}

// Delete 方法：二级缓存删除
func (dc *DistributedCache) Delete(key string) error {
	// 1. 先删除二级缓存
	if err := dc.secondaryCache.Del(dc.cacheCtx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from secondary cache: %v", err)
	}

	// 2. 删除本地一级缓存
	dc.primaryCache.Invalidate(key)

	// 3. 广播删除消息到其他实例（如果不是本地操作触发的）
	dc.syncMutex.RLock()
	isLocalOnly := dc.localOnlyKeys[key]
	dc.syncMutex.RUnlock()

	if !isLocalOnly {
		if err := dc.sendSyncMessage(OperationDelete, key, ""); err != nil {
			log.Printf("Failed to send sync message: %v", err)
		}
	}
	return nil
}

// 带加载器的Get方法
func (dc *DistributedCache) GetWithLoader(
	key string,
	loader func(context.Context, string) (string, error),
) (string, error) {

	// 1. 尝试从缓存获取
	if value, err := dc.get(key); err == nil {
		return value, nil
	}

	// 2. 从数据源加载
	value, err := loader(dc.cacheCtx, key)
	if err != nil {
		return "", err
	}

	// 3. 设置缓存
	if err := dc.set(key, value); err != nil {
		log.Printf("Failed to set cache after loading: %v", err)
	}

	return value, nil
}

//// 批量操作支持
//func (dc *DistributedCache) MGet(keys []string) (map[string]string, error) {
//	//result := make(map[string]string)
//
//	// 1. 先从一级缓存获取可用数据
//	result, err := dc.primaryCache.BulkGet(dc.cacheCtx, keys, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	// 2. 检查未命中的键
//	var missingKeys []string
//	for _, key := range keys {
//		if _, found := result[key]; !found {
//			missingKeys = append(missingKeys, key)
//		}
//	}
//
//	// 3. 从二级缓存获取缺失的键
//	if len(missingKeys) > 0 {
//		values, err := dc.secondaryCache.MGet(dc.cacheCtx, missingKeys...).Result()
//		if err != nil {
//			return nil, err
//		}
//
//		for i, key := range missingKeys {
//			if values[i] != nil {
//				value := values[i].(string)
//				result[key] = value
//				// 回填一级缓存
//				dc.silentSet(key, value)
//			}
//		}
//	}
//
//	return result, nil
//}

func (dc *DistributedCache) MGetWithLoader(keys []string, loader func(context context.Context, key []string) ([]string, error)) (map[string]string, error) {
	// 1. 先从一级缓存获取可用数据
	result, err := dc.primaryCache.BulkGet(dc.cacheCtx, keys, nil)
	if err != nil {
		return nil, err
	}

	// 2. 检查未命中的键
	var missingKeysPrimary []string
	var missingKeysSecond []string
	for _, key := range keys {
		if _, found := result[key]; !found {
			missingKeysPrimary = append(missingKeysPrimary, key)
		}
	}

	// 3. 从二级缓存获取缺失的键
	if len(missingKeysPrimary) > 0 {
		values, err := dc.secondaryCache.MGet(dc.cacheCtx, missingKeysPrimary...).Result()
		if err != nil {
			return nil, err
		}

		for i, key := range missingKeysPrimary {
			if values[i] != nil {
				value := values[i].(string)
				result[key] = value
				// 回填一级缓存
				dc.silentSet(key, value)
			} else {
				missingKeysSecond = append(missingKeysSecond, key)
			}
		}
	}

	strings, err := loader(dc.cacheCtx, missingKeysSecond)
	if err != nil {
		return nil, err
	}
	for i, key := range missingKeysSecond {
		err := dc.set(key, strings[i])
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// 缓存统计信息
type CacheStats struct {
	PrimaryHits     uint64 `json:"primaryHits"`
	PrimaryMisses   uint64 `json:"primaryMisses"`
	SecondaryHits   uint64 `json:"secondaryHits"`
	SecondaryMisses uint64 `json:"secondaryMisses"`
	MessagesSent    uint64 `json:"syncMessagesSent"`
	MessagesRecvd   uint64 `json:"syncMessagesReceived"`
}

func (dc *DistributedCache) GetStats() *CacheStats {
	// 实现统计信息收集
	return &CacheStats{
		// 填充统计信息
		PrimaryHits:   dc.primaryCache.Stats().Hits,
		PrimaryMisses: dc.primaryCache.Stats().Misses,
		MessagesSent:  dc.MsgSendCount.Load(),
		MessagesRecvd: dc.MsgRecvdCount.Load(),
	}
}
