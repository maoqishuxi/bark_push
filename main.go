package main

import (
	remindersender "bark_push/reminder_sender"
	taskexecutor "bark_push/task_executor"
	taskscheduler "bark_push/task_scheduler"
	"fmt"
	"time"
)

func main() {
	scheduler := taskexecutor.NewTaskExecutor(taskscheduler.NewTaskScheduler())

	// 添加任务
	err := scheduler.AddTask("push: 北溟鱼发微博了", "*/30 * * * * *", func() {
		// fmt.Println("Add Task 1 executed at", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Println("Error adding task1:", err)
	}

	err = scheduler.AddTask("push: 叶清眉发微博了", "*/5 * * * * *", func() {
		// fmt.Println("Add Task 2 executed at", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Println("Error adding task2:", err)
	}

	err = scheduler.AddTask("push: 测试推送通知", "*/10 * * * * *", func() {
		// fmt.Println("Add Task 3 executed at", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Println("Error adding task3:", err)
	}

	reminderSender := remindersender.NewReminderSender(scheduler)
	reminderSender.SendReminders()

	// 启动任务调度器
	scheduler.Start()

	// 运行一段时间后停止任务调度器
	time.Sleep(5 * time.Minute)
	scheduler.Stop()
}
