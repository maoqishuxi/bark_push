package main

import (
	taskscheduler "bark_push/task_scheduler"
	"fmt"
	"time"
)

func main() {
	scheduler := taskscheduler.NewTaskScheduler()

	// 添加一个每分钟执行一次的任务
	err := scheduler.AddTask("task1", "*/1 * * * *", func() {
		fmt.Println("Task 1 executed at", time.Now())
	})
	if err != nil {
		fmt.Println("Error adding task1:", err)
	}

	// 启动任务调度器
	scheduler.Start()

	// 运行一段时间后停止任务调度器
	time.Sleep(5 * time.Minute)
	scheduler.Stop()
}
