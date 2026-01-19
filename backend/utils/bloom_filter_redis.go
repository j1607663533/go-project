package utils

import (
	"context"
	"hash/fnv"
	"math"

	"github.com/redis/go-redis/v9"
)

// RedisBloomFilter 基于 Redis 的布隆过滤器（分布式）
type RedisBloomFilter struct {
	client *redis.Client
	key    string
	size   uint
	k      uint
}

// NewRedisBloomFilter 创建一个基于 Redis 的布隆过滤器
func NewRedisBloomFilter(client *redis.Client, key string, n uint, p float64) *RedisBloomFilter {
	m := uint(math.Ceil(-(float64(n) * math.Log(p)) / math.Pow(math.Log(2), 2)))
	k := uint(math.Round(float64(m) / float64(n) * math.Log(2)))
	if k == 0 {
		k = 1
	}

	return &RedisBloomFilter{
		client: client,
		key:    key,
		size:   m,
		k:      k,
	}
}

// Add 向 Redis 过滤器中添加数据
func (b *RedisBloomFilter) Add(ctx context.Context, data string) error {
	pipe := b.client.Pipeline()
	for i := uint(0); i < b.k; i++ {
		offset := int64(b.hash(data, i) % b.size)
		pipe.SetBit(ctx, b.key, offset, 1)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// Contains 判断数据是否可能存在于 Redis 过滤器中
func (b *RedisBloomFilter) Contains(ctx context.Context, data string) (bool, error) {
	pipe := b.client.Pipeline()
	for i := uint(0); i < b.k; i++ {
		offset := int64(b.hash(data, i) % b.size)
		pipe.GetBit(ctx, b.key, offset)
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	for _, cmd := range cmds {
		bit, err := cmd.(*redis.IntCmd).Result()
		if err != nil {
			return false, err
		}
		if bit == 0 {
			return false, nil
		}
	}
	return true, nil
}

// hash 生成哈希索引
func (b *RedisBloomFilter) hash(data string, seed uint) uint {
	h := fnv.New64a()
	h.Write([]byte(data))
	result := h.Sum64()
	// 使用 seed 混淆，保证得到 k 个相互独立的哈希值
	return uint(result ^ uint64(seed)*0xBF58476D1CE4E5B9)
}
