package services

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task 任务结构
type Task struct {
	ID        string
	Status    TaskStatus
	Result    interface{}
	Error     error
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AsyncTaskService 异步任务服务（演示 Channel 和 Goroutine）
type AsyncTaskService struct {
	tasks      map[string]*Task
	tasksMutex sync.RWMutex
	taskQueue  chan *Task
	workers    int
	stopChan   chan bool
}

// NewAsyncTaskService 创建异步任务服务
func NewAsyncTaskService(workers int, queueSize int) *AsyncTaskService {
	service := &AsyncTaskService{
		tasks:     make(map[string]*Task),      // 使用 make 创建 map
		taskQueue: make(chan *Task, queueSize), // 使用 make 创建带缓冲的 Channel
		workers:   workers,
		stopChan:  make(chan bool), // 使用 make 创建 Channel
	}

	// 启动工作协程池
	service.Start()

	return service
}

// Start 启动工作协程池
func (s *AsyncTaskService) Start() {
	for i := 0; i < s.workers; i++ {
		go s.worker(i)
	}
}

// worker 工作协程（演示 Goroutine 和 Channel）
func (s *AsyncTaskService) worker(workerID int) {
	for {
		select {
		case task := <-s.taskQueue:
			// 处理任务
			s.processTask(task, workerID)
		case <-s.stopChan:
			// 停止信号
			return
		}
	}
}

// processTask 处理任务
func (s *AsyncTaskService) processTask(task *Task, workerID int) {
	// 更新任务状态为运行中
	s.updateTaskStatus(task.ID, TaskStatusRunning)

	// 模拟任务处理
	time.Sleep(time.Second * 2)

	// 更新任务状态为完成
	task.Result = fmt.Sprintf("Task processed by worker %d", workerID)
	s.updateTaskStatus(task.ID, TaskStatusCompleted)
}

// SubmitTask 提交任务（使用 Channel）
func (s *AsyncTaskService) SubmitTask(taskID string) error {
	task := &Task{
		ID:        taskID,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存任务
	s.tasksMutex.Lock()
	s.tasks[taskID] = task
	s.tasksMutex.Unlock()

	// 发送到任务队列（使用 Channel）
	select {
	case s.taskQueue <- task:
		return nil
	default:
		return errors.New("任务队列已满")
	}
}

// GetTask 获取任务状态
func (s *AsyncTaskService) GetTask(taskID string) (*Task, error) {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, errors.New("任务不存在")
	}

	return task, nil
}

// GetAllTasks 获取所有任务（演示 make 创建切片）
func (s *AsyncTaskService) GetAllTasks() []*Task {
	s.tasksMutex.RLock()
	defer s.tasksMutex.RUnlock()

	// 使用 make 创建切片
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// updateTaskStatus 更新任务状态
func (s *AsyncTaskService) updateTaskStatus(taskID string, status TaskStatus) {
	s.tasksMutex.Lock()
	defer s.tasksMutex.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Status = status
		task.UpdatedAt = time.Now()
	}
}

// Stop 停止服务
func (s *AsyncTaskService) Stop() {
	close(s.stopChan)
}

// BatchSubmitTasks 批量提交任务（演示 Channel 和并发）
func (s *AsyncTaskService) BatchSubmitTasks(taskIDs []string) map[string]error {
	// 使用 make 创建结果 map
	results := make(map[string]error)
	resultsMutex := &sync.Mutex{}

	// 使用 Channel 进行并发提交
	taskChan := make(chan string, len(taskIDs))
	var wg sync.WaitGroup

	// 启动多个协程并发提交
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskID := range taskChan {
				err := s.SubmitTask(taskID)
				resultsMutex.Lock()
				results[taskID] = err
				resultsMutex.Unlock()
			}
		}()
	}

	// 发送任务到 Channel
	go func() {
		for _, taskID := range taskIDs {
			taskChan <- taskID
		}
		close(taskChan)
	}()

	// 等待完成
	wg.Wait()

	return results
}

// WaitForTask 等待任务完成（演示 Channel 和超时）
func (s *AsyncTaskService) WaitForTask(taskID string, timeout time.Duration) (*Task, error) {
	// 使用 make 创建 Channel
	resultChan := make(chan *Task, 1)
	errorChan := make(chan error, 1)

	// 启动协程监控任务状态
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			task, err := s.GetTask(taskID)
			if err != nil {
				errorChan <- err
				return
			}
			if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed {
				resultChan <- task
				return
			}
		}
	}()

	// 等待结果或超时
	select {
	case task := <-resultChan:
		return task, nil
	case err := <-errorChan:
		return nil, err
	case <-time.After(timeout):
		return nil, errors.New("等待任务超时")
	}
}
