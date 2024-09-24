package barkpush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BarkPushService 结构体
type BarkPushService struct {
	client  *http.Client
	baseURL string
}

// NewBarkPushService 创建一个新的 Bark 推送服务
func NewBarkPushService(baseURL string) *BarkPushService {
	return &BarkPushService{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
	}
}

// PushMessage 发送推送消息
func (b *BarkPushService) PushMessage(title, body, icon, group string) error {
	url := b.baseURL

	payload := map[string]interface{}{
		"title": title,
		"body":  body,
		"badge": 1,
		"icon":  icon,
		"group": group,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
