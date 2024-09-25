package main

import (
	barkpush "bark_push/bark_push"
	oncereminder "bark_push/once_reminder"
	remindersender "bark_push/reminder_sender"
	taskexecutor "bark_push/task_executor"
	taskscheduler "bark_push/task_scheduler"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	scheduler := setupTaskScheduler()
	barkPush := setupBarkPushService()

	setupAddTaskAPI(r, scheduler, barkPush)
	setupDeleteTaskAPI(r, scheduler)
	setupOnceReminderAPI(r, barkPush)
	setupPushAPI(r, barkPush)

	// 启动任务调度器
	scheduler.Start()

	// 启动 Gin 服务器
	if err := r.Run(":7000"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

func setupTaskScheduler() *taskexecutor.TaskExecutor {
	return taskexecutor.NewTaskExecutor(taskscheduler.NewTaskScheduler())
}

func setupBarkPushService() *barkpush.BarkPushService {
	base_url := "http://bark.julai.fun:8080/EehxmJU8PGTBGNzcwQ8Yfc/"
	return barkpush.NewBarkPushService(base_url)
}

func setupAddTaskAPI(r *gin.Engine, scheduler *taskexecutor.TaskExecutor, barkPush *barkpush.BarkPushService) {
	r.POST("/add-task", func(c *gin.Context) {
		var task struct {
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Icon     string `json:"icon"`
			Group    string `json:"group"`
			Title    string `json:"title"`
			Body     string `json:"body"`
		}
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task.Icon = getDefaultIcon(task.Icon)
		task.Group = getDefaultGroup(task.Group, "轮询")

		reminderType, _ := remindersender.ParseReminder(task.Name)

		err := scheduler.AddTask(task.Name, task.Schedule, func() {
			if reminderType == "push" {
				barkPush.PushMessage(task.Title, task.Body, task.Icon, task.Group)
			}
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task added successfully"})
	})
}

func setupDeleteTaskAPI(r *gin.Engine, scheduler *taskexecutor.TaskExecutor) {
	r.DELETE("/delete-task", func(c *gin.Context) {
		var task struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		scheduler.RemoveTask(task.Name)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
	})
}

func setupOnceReminderAPI(r *gin.Engine, barkPush *barkpush.BarkPushService) {
	r.POST("/once-reminder", func(c *gin.Context) {
		var reminder struct {
			Time  string `json:"time"`
			Icon  string `json:"icon"`
			Group string `json:"group"`
			Title string `json:"title"`
			Body  string `json:"body"`
		}
		if err := c.ShouldBindJSON(&reminder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		reminder.Icon = getDefaultIcon(reminder.Icon)
		reminder.Group = getDefaultGroup(reminder.Group, "定时")

		reminderTime, err := time.ParseInLocation("2006-01-02 15:04:05", reminder.Time, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}

		duration := time.Until(reminderTime)
		if duration < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reminder time is in the past"})
			return
		}

		onceReminder := oncereminder.NewOnceReminder(duration, func() {
			barkPush.PushMessage(reminder.Title, reminder.Body, reminder.Icon, reminder.Group)
		})
		onceReminder.Start()

		c.JSON(http.StatusOK, gin.H{"message": "Once reminder set successfully"})
	})
}

func setupPushAPI(r *gin.Engine, barkPush *barkpush.BarkPushService) {
	r.POST("/push", func(c *gin.Context) {
		var message struct {
			Title string `json:"title"`
			Body  string `json:"body"`
			Icon  string `json:"icon"`
			Group string `json:"group"`
		}
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		message.Icon = getDefaultIcon(message.Icon)

		err := barkPush.PushMessage(message.Title, message.Body, message.Icon, message.Group)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Push message sent successfully"})
	})
}

func getDefaultIcon(icon string) string {
	if icon == "" {
		return "https://day.app/assets/images/avatar.jpg"
	}
	return icon
}

func getDefaultGroup(group string, defaultGroup string) string {
	if group == "" {
		return defaultGroup
	}
	return group
}
