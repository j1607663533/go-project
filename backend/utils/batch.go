package utils

import (
	"sync"
)

// BatchProcessor 批量处理器
type BatchProcessor struct {
	batchSize int
	workers   int
}

// NewBatchProcessor 创建批量处理器
func NewBatchProcessor(batchSize, workers int) *BatchProcessor {
	return &BatchProcessor{
		batchSize: batchSize,
		workers:   workers,
	}
}

// ProcessItems 批量处理项目（使用 Channel 和 Goroutine）
// 演示：make 创建切片、Channel、并发处理
func (bp *BatchProcessor) ProcessItems(items []interface{}, processFunc func(interface{}) error) []error {
	// 使用 make 创建切片存储错误
	errors := make([]error, 0)
	errorsMutex := &sync.Mutex{}

	// 使用 make 创建 Channel
	itemsChan := make(chan interface{}, bp.batchSize)
	errorsChan := make(chan error, len(items))
	doneChan := make(chan bool)

	// 启动工作协程池
	var wg sync.WaitGroup
	for i := 0; i < bp.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// 从 Channel 中读取任务
			for item := range itemsChan {
				if err := processFunc(item); err != nil {
					errorsChan <- err
				}
			}
		}(i)
	}

	// 发送任务到 Channel
	go func() {
		for _, item := range items {
			itemsChan <- item
		}
		close(itemsChan) // 关闭 Channel，通知工作协程没有更多任务
	}()

	// 等待所有工作协程完成
	go func() {
		wg.Wait()
		close(errorsChan)
		doneChan <- true
	}()

	// 收集错误
	go func() {
		for err := range errorsChan {
			errorsMutex.Lock()
			errors = append(errors, err)
			errorsMutex.Unlock()
		}
	}()

	// 等待完成
	<-doneChan

	return errors
}

// BatchResult 批量操作结果
type BatchResult struct {
	SuccessCount int
	FailureCount int
	Errors       []error
}

// ProcessWithResult 批量处理并返回详细结果
func (bp *BatchProcessor) ProcessWithResult(items []interface{}, processFunc func(interface{}) error) *BatchResult {
	// 使用 make 创建结果切片
	result := &BatchResult{
		Errors: make([]error, 0),
	}

	// 使用 Channel 进行并发处理
	itemsChan := make(chan interface{}, bp.batchSize)
	resultChan := make(chan error, len(items))

	// 工作协程池
	var wg sync.WaitGroup
	for i := 0; i < bp.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range itemsChan {
				err := processFunc(item)
				resultChan <- err
			}
		}()
	}

	// 发送任务
	go func() {
		for _, item := range items {
			itemsChan <- item
		}
		close(itemsChan)
	}()

	// 等待完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for err := range resultChan {
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, err)
		} else {
			result.SuccessCount++
		}
	}

	return result
}

// ChunkSlice 将切片分块（演示 make 创建二维切片）
func ChunkSlice(items []interface{}, chunkSize int) [][]interface{} {
	// 计算需要多少个块
	totalChunks := (len(items) + chunkSize - 1) / chunkSize

	// 使用 make 创建二维切片
	chunks := make([][]interface{}, 0, totalChunks)

	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		// 使用 make 创建每个块的切片
		chunk := make([]interface{}, end-i)
		copy(chunk, items[i:end])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// ParallelMap 并行映射函数（演示 Channel 和 Goroutine）
func ParallelMap(items []interface{}, mapFunc func(interface{}) interface{}, workers int) []interface{} {
	// 使用 make 创建结果切片
	results := make([]interface{}, len(items))

	// 使用 Channel 分发任务
	type task struct {
		index int
		item  interface{}
	}

	tasksChan := make(chan task, len(items))
	resultsChan := make(chan task, len(items))

	// 启动工作协程
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasksChan {
				result := mapFunc(t.item)
				resultsChan <- task{index: t.index, item: result}
			}
		}()
	}

	// 发送任务
	go func() {
		for i, item := range items {
			tasksChan <- task{index: i, item: item}
		}
		close(tasksChan)
	}()

	// 等待完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集结果（保持顺序）
	for t := range resultsChan {
		results[t.index] = t.item
	}

	return results
}
