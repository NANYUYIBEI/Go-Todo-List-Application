package kafka

import (
	"time"
)

// TodoEvent 定义了发送到 Kafka 的任务事件结构
// omitempty 标签表示如果字段是零值，则在 JSON 编码时省略该字段
type TodoEvent struct {
	EventType string    `json:"eventType"`           // "CREATE", "UPDATE", "DELETE"
	TodoID    uint      `json:"todoId"`              // 任务的 ID
	TodoName  string    `json:"todoName,omitempty"`  // 任务名称 (创建/更新时有，删除时可无)
	Completed bool      `json:"completed,omitempty"` // 任务完成状态 (更新时有)
	Timestamp time.Time `json:"timestamp"`           // 事件发生时间
}
