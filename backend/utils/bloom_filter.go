package utils

import (
	"hash/fnv"
	"math"
	"sync"
)

// BloomFilter 布隆过滤器接口
type BloomFilter interface {
	Add(data string)           // 添加元素
	Contains(data string) bool // 判断元素是否可能存在
}

// MemoryBloomFilter 内存版布隆过滤器
type MemoryBloomFilter struct {
	bitset []uint64
	size   uint
	k      uint // 哈希函数的个数
	mu     sync.RWMutex
}

// NewMemoryBloomFilter 创建一个内存布隆过滤器
// n: 预计存储的元素个数
// p: 允许的误报率 (0 < p < 1)
func NewMemoryBloomFilter(n uint, p float64) *MemoryBloomFilter {
	// 计算所需的二进制位数 m = -(n * ln(p)) / (ln(2)^2)
	m := uint(math.Ceil(-(float64(n) * math.Log(p)) / math.Pow(math.Log(2), 2)))

	// 计算所需的哈希函数个数 k = (m/n) * ln(2)
	k := uint(math.Round(float64(m) / float64(n) * math.Log(2)))
	if k == 0 {
		k = 1
	}

	// 将 m 向上取整到 64 的倍数，方便用 uint64 数组存储
	numWords := (m + 63) / 64
	return &MemoryBloomFilter{
		bitset: make([]uint64, numWords),
		size:   m,
		k:      k,
	}
}

// Add 向过滤器中添加数据
func (b *MemoryBloomFilter) Add(data string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := uint(0); i < b.k; i++ {
		index := b.hash(data, i) % b.size
		wordIndex := index / 64
		bitOffset := index % 64
		b.bitset[wordIndex] |= (1 << bitOffset)
	}
}

// Contains 判断数据是否可能存在
func (b *MemoryBloomFilter) Contains(data string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for i := uint(0); i < b.k; i++ {
		index := b.hash(data, i) % b.size
		wordIndex := index / 64
		bitOffset := index % 64
		if (b.bitset[wordIndex] & (1 << bitOffset)) == 0 {
			return false // 只要有一位是0，一定不存在
		}
	}
	return true // 可能存在
}

// hash 使用 FNV-1a 算法并结合不同的种子生成多个哈希索引
func (b *MemoryBloomFilter) hash(data string, seed uint) uint {
	h := fnv.New64a()
	h.Write([]byte(data))
	// 简单的加盐或者是多次哈希，这里采用种子作为额外输入
	result := h.Sum64()

	// 为了得到 k 个不同的哈希值，常用双重哈希策略: hash(i) = h1 + i * h2
	// 这里简化处理：将 seed 混入
	return uint(result ^ uint64(seed)*0xBF58476D1CE4E5B9)
}
