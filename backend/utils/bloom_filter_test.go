package utils

import (
	"testing"
)

func TestMemoryBloomFilter(t *testing.T) {
	// 创建一个预计存储 1000 个元素，误报率为 0.01 的过滤器
	bf := NewMemoryBloomFilter(1000, 0.01)

	elements := []string{"apple", "banana", "cherry"}

	// 添加元素
	for _, e := range elements {
		bf.Add(e)
	}

	// 测试已存在的元素
	for _, e := range elements {
		if !bf.Contains(e) {
			t.Errorf("Expected %s to be contained", e)
		}
	}

	// 测试不存在的元素
	if bf.Contains("dog") {
		t.Log("Note: 'dog' might be a false positive (which is allowed)")
	} else {
		t.Log("'dog' correctly identified as not present")
	}
}
