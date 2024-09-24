package taskscheduler

import (
	"sync"
	"time"

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
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		))),
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

// AddOneTimeTask adds a task that will be executed only once at the specified time
func (ts *TaskScheduler) AddOneTimeTask(taskTime time.Time, task func()) {
	ts.cron.Schedule(cron.Schedule(cron.Every(taskTime.Sub(time.Now()))), cron.FuncJob(func() {
		task()
		// Remove the task after execution
		ts.cron.Remove(cron.EntryID(ts.cron.Entries()[len(ts.cron.Entries())-1].ID))
	}))
}
