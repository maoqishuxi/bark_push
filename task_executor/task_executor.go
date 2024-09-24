package taskexecutor

import (
	taskscheduler "bark_push/task_scheduler"
	"sync"
)

// TaskExecutor 结构体
type TaskExecutor struct {
	scheduler *taskscheduler.TaskScheduler
	mu        sync.Mutex
	queue     chan string
}

// NewTaskExecutor 创建一个新的任务执行器
func NewTaskExecutor(scheduler *taskscheduler.TaskScheduler) *TaskExecutor {
	return &TaskExecutor{
		scheduler: scheduler,
		queue:     make(chan string, 100), // 假设队列大小为100
	}
}

// ExecuteTask 执行任务逻辑
func (te *TaskExecutor) ExecuteTask(name string, cmd func()) {
	te.mu.Lock()
	defer te.mu.Unlock()

	// 执行任务
	cmd()

	// 将需要发送的提醒放入消息队列
	te.queue <- name
}

// Start 启动任务执行器
func (te *TaskExecutor) Start() {
	te.scheduler.Start()
}

// Stop 停止任务执行器
func (te *TaskExecutor) Stop() {
	te.scheduler.Stop()
	close(te.queue)
}

// AddTask 添加一个新任务
func (te *TaskExecutor) AddTask(name string, spec string, cmd func()) error {
	return te.scheduler.AddTask(name, spec, func() {
		te.ExecuteTask(name, cmd)
	})
}

// RemoveTask 删除一个任务
func (te *TaskExecutor) RemoveTask(name string) {
	te.scheduler.RemoveTask(name)
}

// UpdateTask 更新一个任务
func (te *TaskExecutor) UpdateTask(name string, spec string, cmd func()) error {
	return te.scheduler.UpdateTask(name, spec, func() {
		te.ExecuteTask(name, cmd)
	})
}

// GetQueue 获取消息队列
func (te *TaskExecutor) GetQueue() chan string {
	return te.queue
}
