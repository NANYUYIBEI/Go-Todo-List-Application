package database

import (
	// 导入 fmt 以便在日志中使用格式化字符串
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"todo/internal/database/models"
)

var DB *gorm.DB // 全局数据库连接变量

// InitMySQL 初始化 MySQL 数据库连接
func InitMySQL(dsn string) {
	var err error
	// 使用 gorm.Open 连接到 MySQL
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}

	log.Println("Connected to MySQL database!")

	// 自动迁移模式，创建或更新表结构
	// 这里会创建 todos 表，如果它不存在的话
	err = DB.AutoMigrate(&models.Todo{})
	if err != nil {
		// 使用 fmt.Sprintf 格式化错误消息
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Println("Database migration completed.")
}

// GetDB 返回全局数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
