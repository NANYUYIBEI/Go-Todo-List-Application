package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a single todo item
type Todo struct {
	// Gorm v2 会将 ID、CreatedAt、UpdatedAt、DeletedAt 默认序列化为 JSON。
	// 如果你希望明确控制，可以像下面这样：
	ID        uint           `gorm:"primaryKey" json:"id"`             // 明确定义 ID 字段并添加 json 标签
	CreatedAt time.Time      `json:"createdAt"`                        // 明确定义时间戳并添加 json 标签
	UpdatedAt time.Time      `json:"updatedAt"`                        // 明确定义时间戳并添加 json 标签
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"` // 明确定义软删除并添加 json 标签

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Completed   bool   `gorm:"default:false" json:"completed"`
}
