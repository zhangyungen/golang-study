package cache

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
)

// BloomFilter 布隆过滤器接口
type BloomFilter interface {
	// Add 添加元素到布隆过滤器
	Add(ctx context.Context, key string) error
	// Exists 判断元素是否可能存在
	Exists(ctx context.Context, key string) (bool, error)
	// BatchAdd 批量添加元素
	BatchAdd(ctx context.Context, keys []string) error
	// Close 关闭布隆过滤器
	Close() error
}

// MemoryBloomFilter 内存布隆过滤器实现
type MemoryBloomFilter struct {
	filter *bloom.BloomFilter
}

// NewMemoryBloomFilter 创建内存布隆过滤器
// expectedElements: 预期元素数量
// falsePositiveRate: 误判率
func NewMemoryBloomFilter(expectedElements uint, falsePositiveRate float64) BloomFilter {
	return &MemoryBloomFilter{
		filter: bloom.NewWithEstimates(expectedElements, falsePositiveRate),
	}
}

func (m *MemoryBloomFilter) Add(ctx context.Context, key string) error {
	m.filter.AddString(key)
	return nil
}

func (m *MemoryBloomFilter) Exists(ctx context.Context, key string) (bool, error) {
	return m.filter.TestString(key), nil
}

func (m *MemoryBloomFilter) BatchAdd(ctx context.Context, keys []string) error {
	for _, key := range keys {
		m.filter.AddString(key)
	}
	return nil
}

func (m *MemoryBloomFilter) Close() error {
	m.filter = nil
	return nil
}
