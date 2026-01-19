# Go 核心特性使用文档 - make、切片和 Channel

## 📋 概述

项目中现在已经集成了 Go 的核心特性：

- **make** - 用于创建切片、map 和 Channel
- **切片（Slice）** - 动态数组
- **Channel** - 用于 Goroutine 之间的通信
- **Goroutine** - 轻量级并发

## 📂 新增文件

### 1. `utils/batch.go` - 批量处理工具

演示了以下 Go 特性：

#### ✅ make 创建切片

```go
// 创建空切片，预分配容量
errors := make([]error, 0)

// 创建指定长度和容量的切片
results := make([]interface{}, len(items))

// 创建二维切片
chunks := make([][]interface{}, 0, totalChunks)
```

#### ✅ make 创建 Channel

```go
// 创建带缓冲的 Channel
itemsChan := make(chan interface{}, batchSize)
errorsChan := make(chan error, len(items))

// 创建无缓冲的 Channel
doneChan := make(chan bool)
```

#### ✅ Goroutine 并发处理

```go
// 启动工作协程池
for i := 0; i < workers; i++ {
    go func(workerID int) {
        // 从 Channel 读取任务
        for item := range itemsChan {
            processFunc(item)
        }
    }(i)
}
```

#### ✅ Channel 通信

```go
// 发送数据到 Channel
itemsChan <- item

// 从 Channel 接收数据
item := <-itemsChan

// 关闭 Channel
close(itemsChan)
```

### 2. `services/async_task_service.go` - 异步任务服务

演示了高级并发模式：

#### ✅ make 创建 map

```go
tasks := make(map[string]*Task)
```

#### ✅ 带缓冲的 Channel

```go
taskQueue := make(chan *Task, queueSize)
```

#### ✅ select 语句

```go
select {
case task := <-s.taskQueue:
    // 处理任务
    s.processTask(task)
case <-s.stopChan:
    // 停止信号
    return
}
```

#### ✅ 超时处理

```go
select {
case result := <-resultChan:
    return result, nil
case <-time.After(timeout):
    return nil, errors.New("超时")
}
```

### 3. `services/user_service.go` - 用户服务批量操作

在实际业务中应用这些特性：

#### ✅ 批量获取用户

```go
func (s *userService) GetUsersByIDs(ids []uint) ([]models.UserResponse, error) {
    // 使用 make 创建切片
    results := make([]models.UserResponse, 0, len(ids))

    // 使用 make 创建 Channel
    resultChan := make(chan userResult, len(ids))

    // 并发获取
    for _, id := range ids {
        go func(userID uint) {
            user, err := s.GetUserByID(userID)
            resultChan <- userResult{user: user, err: err}
        }(id)
    }

    // 收集结果
    for i := 0; i < len(ids); i++ {
        result := <-resultChan
        if result.err == nil {
            results = append(results, *result.user)
        }
    }

    return results, nil
}
```

#### ✅ 批量创建用户

```go
func (s *userService) BatchCreateUsers(requests []*models.UserCreateRequest) ([]models.UserResponse, []error) {
    // 使用 make 创建结果切片
    successUsers := make([]models.UserResponse, 0)
    errors := make([]error, 0)

    // 使用 Channel 并发创建
    resultChan := make(chan createResult, len(requests))

    for i, req := range requests {
        go func(index int, request *models.UserCreateRequest) {
            user, err := s.CreateUser(request)
            resultChan <- createResult{user: user, err: err, index: index}
        }(i, req)
    }

    // 收集结果
    for i := 0; i < len(requests); i++ {
        result := <-resultChan
        if result.err != nil {
            errors = append(errors, result.err)
        } else {
            successUsers = append(successUsers, *result.user)
        }
    }

    return successUsers, errors
}
```

## 🎯 核心概念详解

### 1. make 函数

`make` 用于创建切片、map 和 Channel：

```go
// 切片
s1 := make([]int, 0)        // 长度0，容量0
s2 := make([]int, 5)        // 长度5，容量5
s3 := make([]int, 0, 10)    // 长度0，容量10

// map
m := make(map[string]int)

// Channel
ch1 := make(chan int)       // 无缓冲
ch2 := make(chan int, 10)   // 缓冲大小10
```

### 2. 切片（Slice）

动态数组，可以自动扩容：

```go
// 创建切片
slice := make([]int, 0, 10)

// 添加元素
slice = append(slice, 1, 2, 3)

// 遍历
for i, v := range slice {
    fmt.Printf("索引: %d, 值: %d\n", i, v)
}

// 切片操作
subSlice := slice[1:3]  // 获取子切片
```

### 3. Channel

用于 Goroutine 之间的通信：

```go
// 创建 Channel
ch := make(chan int, 5)

// 发送数据
ch <- 42

// 接收数据
value := <-ch

// 关闭 Channel
close(ch)

// 遍历 Channel
for value := range ch {
    fmt.Println(value)
}
```

### 4. Goroutine

轻量级线程：

```go
// 启动 Goroutine
go func() {
    fmt.Println("并发执行")
}()

// 带参数的 Goroutine
go func(msg string) {
    fmt.Println(msg)
}("Hello")
```

## 📊 使用场景

### 场景 1: 批量数据处理

```go
processor := utils.NewBatchProcessor(100, 5)

items := []interface{}{1, 2, 3, 4, 5}
errors := processor.ProcessItems(items, func(item interface{}) error {
    // 处理每个项目
    return nil
})
```

### 场景 2: 异步任务

```go
taskService := services.NewAsyncTaskService(10, 100)

// 提交任务
taskService.SubmitTask("task-1")

// 等待任务完成
task, err := taskService.WaitForTask("task-1", 30*time.Second)
```

### 场景 3: 并发获取数据

```go
userService := services.NewUserService(userRepo)

// 批量获取用户
ids := []uint{1, 2, 3, 4, 5}
users, err := userService.GetUsersByIDs(ids)
```

## 🧪 测试示例

### 测试批量处理

```go
package main

import (
    "fmt"
    "gin-backend/utils"
)

func main() {
    processor := utils.NewBatchProcessor(10, 3)

    // 准备数据
    items := make([]interface{}, 100)
    for i := 0; i < 100; i++ {
        items[i] = i
    }

    // 批量处理
    errors := processor.ProcessItems(items, func(item interface{}) error {
        fmt.Printf("处理: %v\n", item)
        return nil
    })

    fmt.Printf("完成，错误数: %d\n", len(errors))
}
```

### 测试异步任务

```go
package main

import (
    "fmt"
    "gin-backend/services"
    "time"
)

func main() {
    taskService := services.NewAsyncTaskService(5, 50)

    // 批量提交任务
    taskIDs := []string{"task-1", "task-2", "task-3"}
    results := taskService.BatchSubmitTasks(taskIDs)

    for taskID, err := range results {
        if err != nil {
            fmt.Printf("任务 %s 提交失败: %v\n", taskID, err)
        } else {
            fmt.Printf("任务 %s 提交成功\n", taskID)
        }
    }

    // 等待任务完成
    time.Sleep(3 * time.Second)

    // 获取所有任务状态
    tasks := taskService.GetAllTasks()
    for _, task := range tasks {
        fmt.Printf("任务 %s: %s\n", task.ID, task.Status)
    }
}
```

## 💡 最佳实践

### 1. 使用 make 预分配容量

```go
// ✅ 好 - 预分配容量，减少内存分配
slice := make([]int, 0, 100)

// ❌ 不好 - 频繁扩容
slice := []int{}
```

### 2. 使用带缓冲的 Channel

```go
// ✅ 好 - 带缓冲，减少阻塞
ch := make(chan int, 100)

// ❌ 不好 - 无缓冲，容易阻塞
ch := make(chan int)
```

### 3. 记得关闭 Channel

```go
// ✅ 好
ch := make(chan int, 10)
// ... 发送数据
close(ch)

// ❌ 不好 - 忘记关闭，可能导致 Goroutine 泄漏
```

### 4. 使用 sync.WaitGroup 等待 Goroutine

```go
// ✅ 好
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 处理任务
    }()
}
wg.Wait()

// ❌ 不好 - 使用 time.Sleep 等待
time.Sleep(time.Second)
```

### 5. 避免 Goroutine 泄漏

```go
// ✅ 好 - 使用 context 或 done channel
done := make(chan bool)
go func() {
    for {
        select {
        case <-done:
            return
        default:
            // 处理任务
        }
    }
}()

// ❌ 不好 - 无法停止的 Goroutine
go func() {
    for {
        // 处理任务
    }
}()
```

## 🎉 总结

现在项目中已经包含了 Go 的核心特性：

1. ✅ **make** - 创建切片、map、Channel
2. ✅ **切片** - 动态数组操作
3. ✅ **Channel** - Goroutine 通信
4. ✅ **Goroutine** - 并发处理
5. ✅ **select** - 多路复用
6. ✅ **sync.WaitGroup** - 等待协程完成

这些特性在以下场景中得到应用：

- 批量数据处理
- 异步任务执行
- 并发 API 调用
- 实时数据处理

通过这些实现，你可以更好地理解和使用 Go 的并发特性！🚀
