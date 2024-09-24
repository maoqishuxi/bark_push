package remindersender

import (
	taskexecutor "bark_push/task_executor"
	"fmt"
	"strings"
	"time"
)

// ReminderSender 结构体
type ReminderSender struct {
	executor *taskexecutor.TaskExecutor
}

// NewReminderSender 创建一个新的提醒发送器
func NewReminderSender(executor *taskexecutor.TaskExecutor) *ReminderSender {
	return &ReminderSender{
		executor: executor,
	}
}

// SendReminders 发送提醒
func (rs *ReminderSender) SendReminders() {
	go func() {
		for reminder := range rs.executor.GetQueue() {
			var sendFunc func(string)
			reminderType, message := ParseReminder(reminder)
			switch reminderType {
			case "push":
				sendFunc = sendPushNotification
			default:
				fmt.Printf("Unknown reminder type: %s\n", reminderType)
				continue
			}
			sendFunc(message)
		}
	}()
}

// parseReminder 解析提醒
func ParseReminder(reminder string) (string, string) {
	parts := strings.SplitN(reminder, ":", 2) // 只分割一次
	if len(parts) != 2 {
		fmt.Printf("Failed to parse reminder: %s, format error\n", reminder)
		return "", ""
	}
	reminderType := strings.TrimSpace(parts[0])
	message := strings.TrimSpace(parts[1])
	return reminderType, message
}

// sendPushNotification 发送推送通知
func sendPushNotification(message string) {
	fmt.Printf("%s Sending push notification: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	// 调用外部推送通知服务
}
