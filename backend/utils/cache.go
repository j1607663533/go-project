package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-backend/config"
	"os"
	"sync"
	"time"
)

const cacheFile = "cache_persistence.json"

// memoryCacheItem 内存缓存项
type memoryCacheItem struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

var (
	memoryCache   map[string]memoryCacheItem
	memoryCacheMu sync.RWMutex
)

func init() {
	loadCache()
}

func loadCache() {
	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()

	memoryCache = make(map[string]memoryCacheItem)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &memoryCache); err != nil {
		fmt.Printf("加载缓存文件失败: %v\n", err)
	}
}

func saveCacheUnsafe() {
	data, err := json.MarshalIndent(memoryCache, "", "  ")
	if err != nil {
		fmt.Printf("序列化缓存失败: %v\n", err)
		return
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		fmt.Printf("保存缓存文件失败: %v\n", err)
	}
}

var (
	lastRedisCheck      time.Time
	redisAvailableState bool
	checkMutex          sync.Mutex
)

// isRedisAvailable 检查 Redis 是否可用
func isRedisAvailable() bool {
	if config.RedisClient == nil {
		return false
	}

	checkMutex.Lock()
	defer checkMutex.Unlock()

	// 5 秒内不再重复检查
	if time.Since(lastRedisCheck) < 5*time.Second {
		return redisAvailableState
	}

	// 尝试 Ping 一下，确保连接仍然有效
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := config.RedisClient.Ping(ctx).Result()
	redisAvailableState = (err == nil)
	lastRedisCheck = time.Now()

	return redisAvailableState
}

// CacheSet 设置缓存
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	// 将值序列化为 JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if isRedisAvailable() {
		return config.RedisClient.Set(config.GetRedisContext(), key, jsonValue, expiration).Err()
	}

	// 降级使用内存
	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()

	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	memoryCache[key] = memoryCacheItem{
		Value:      string(jsonValue),
		Expiration: exp,
	}
	saveCacheUnsafe()
	return nil
}

// CacheGet 获取缓存
func CacheGet(key string, dest interface{}) error {
	if isRedisAvailable() {
		val, err := config.RedisClient.Get(config.GetRedisContext(), key).Result()
		if err == nil {
			return json.Unmarshal([]byte(val), dest)
		}
	}

	// 从内存获取
	memoryCacheMu.RLock()
	item, ok := memoryCache[key]
	memoryCacheMu.RUnlock()

	if !ok {
		return errors.New("redis: nil")
	}

	// 检查过期
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		CacheDel(key)
		return errors.New("redis: nil")
	}

	return json.Unmarshal([]byte(item.Value), dest)
}

// CacheGetString 获取字符串缓存
func CacheGetString(key string) (string, error) {
	if isRedisAvailable() {
		return config.RedisClient.Get(config.GetRedisContext(), key).Result()
	}

	memoryCacheMu.RLock()
	item, ok := memoryCache[key]
	memoryCacheMu.RUnlock()

	if !ok {
		return "", errors.New("redis: nil")
	}

	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		CacheDel(key)
		return "", errors.New("redis: nil")
	}

	return item.Value, nil
}

// CacheSetString 设置字符串缓存
func CacheSetString(key string, value string, expiration time.Duration) error {
	if isRedisAvailable() {
		return config.RedisClient.Set(config.GetRedisContext(), key, value, expiration).Err()
	}

	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()

	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	memoryCache[key] = memoryCacheItem{
		Value:      value,
		Expiration: exp,
	}
	saveCacheUnsafe()
	return nil
}

// CacheDel 删除缓存
func CacheDel(keys ...string) error {
	if isRedisAvailable() {
		return config.RedisClient.Del(config.GetRedisContext(), keys...).Err()
	}

	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()
	for _, key := range keys {
		delete(memoryCache, key)
	}
	saveCacheUnsafe()
	return nil
}

// CacheExists 检查缓存是否存在
func CacheExists(key string) (bool, error) {
	if isRedisAvailable() {
		count, err := config.RedisClient.Exists(config.GetRedisContext(), key).Result()
		return count > 0, err
	}

	memoryCacheMu.RLock()
	item, ok := memoryCache[key]
	memoryCacheMu.RUnlock()

	if !ok {
		return false, nil
	}

	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		CacheDel(key)
		return false, nil
	}

	return true, nil
}

// CacheExpire 设置过期时间
func CacheExpire(key string, expiration time.Duration) error {
	if isRedisAvailable() {
		return config.RedisClient.Expire(config.GetRedisContext(), key, expiration).Err()
	}

	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()

	if item, ok := memoryCache[key]; ok {
		item.Expiration = time.Now().Add(expiration)
		memoryCache[key] = item
		saveCacheUnsafe()
		return nil
	}

	return errors.New("key not found")
}

// CacheTTL 获取剩余过期时间
func CacheTTL(key string) (time.Duration, error) {
	if isRedisAvailable() {
		return config.RedisClient.TTL(config.GetRedisContext(), key).Result()
	}

	memoryCacheMu.RLock()
	item, ok := memoryCache[key]
	memoryCacheMu.RUnlock()

	if !ok {
		return -2, errors.New("key not found")
	}

	if item.Expiration.IsZero() {
		return -1, nil
	}

	ttl := time.Until(item.Expiration)
	if ttl <= 0 {
		CacheDel(key)
		return -2, errors.New("key not found")
	}

	return ttl, nil
}

// CacheIncr 自增
func CacheIncr(key string) (int64, error) {
	return 0, errors.New("not implemented in memory fallback")
}

// CacheIncrBy 指定步长自增
func CacheIncrBy(key string, value int64) (int64, error) {
	return 0, errors.New("not implemented in memory fallback")
}

// CacheDecr 自减
func CacheDecr(key string) (int64, error) {
	return 0, errors.New("not implemented in memory fallback")
}

// CacheKeys 获取匹配的键列表
func CacheKeys(pattern string) ([]string, error) {
	return nil, errors.New("not implemented in memory fallback")
}

// CacheFlushDB 清空当前数据库
func CacheFlushDB() error {
	if isRedisAvailable() {
		return config.RedisClient.FlushDB(config.GetRedisContext()).Err()
	}

	memoryCacheMu.Lock()
	defer memoryCacheMu.Unlock()
	memoryCache = make(map[string]memoryCacheItem)
	saveCacheUnsafe()
	return nil
}

// CacheHSet 设置哈希字段
func CacheHSet(key string, field string, value interface{}) error {
	return errors.New("not implemented in memory fallback")
}

// CacheHGet 获取哈希字段
func CacheHGet(key string, field string) (string, error) {
	return "", errors.New("not implemented in memory fallback")
}

// CacheHGetAll 获取所有哈希字段
func CacheHGetAll(key string) (map[string]string, error) {
	return nil, errors.New("not implemented in memory fallback")
}

// CacheHDel 删除哈希字段
func CacheHDel(key string, fields ...string) error {
	return errors.New("not implemented in memory fallback")
}
