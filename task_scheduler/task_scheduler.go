package taskscheduler

import (
	"sync"

	"github.com/robfig/cron/v3"
)

// TaskScheduler 结构体
type TaskScheduler struct {
	cron *cron.Cron
	mu   sync.Mutex
	jobs map[string]cron.EntryID
}

// NewTaskScheduler 创建一个新的任务调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		cron: cron.New(),
		jobs: make(map[string]cron.EntryID),
	}
}

// Start 启动任务调度器
func (ts *TaskScheduler) Start() {
	ts.cron.Start()
}

// Stop 停止任务调度器
func (ts *TaskScheduler) Stop() {
	ts.cron.Stop()
}

// AddTask 添加一个新任务
func (ts *TaskScheduler) AddTask(name string, spec string, cmd func()) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	id, err := ts.cron.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	ts.jobs[name] = id
	return nil
}

// RemoveTask 删除一个任务
func (ts *TaskScheduler) RemoveTask(name string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if id, ok := ts.jobs[name]; ok {
		ts.cron.Remove(id)
		delete(ts.jobs, name)
	}
}

// UpdateTask 更新一个任务
func (ts *TaskScheduler) UpdateTask(name string, spec string, cmd func()) error {
	ts.RemoveTask(name)
	return ts.AddTask(name, spec, cmd)
}
