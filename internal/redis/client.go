package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8" // 确保能导入这个包

	"todo/internal/database/models" // 导入你的 models 包，路径为 "todo/internal/database/models"
)

// 定义一个 context 上下文，用于 Redis 操作。在实际项目中，通常会从请求上下文传递。
var ctx = context.Background()

const (
	TodosCacheKey = "all_todos" // 用于缓存所有 Todos 列表的键
)

// Client 结构体封装了 Redis 客户端
type Client struct {
	*redis.Client
}

// InitRedis 初始化 Redis 客户端连接
func InitRedis(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // Redis 服务器地址 (例如 "localhost:6379")
		Password: "",   // 如果 Redis 有密码，请在此处填写
		DB:       0,    // 默认数据库
	})

	// 尝试 ping Redis 服务器，测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis!")
	return &Client{rdb}
}

// GetTodosCache 从 Redis 获取所有 Todos 的缓存
func (r *Client) GetTodosCache() ([]models.Todo, bool) {
	val, err := r.Get(ctx, TodosCacheKey).Result()
	if err != nil {
		if err == redis.Nil { // 缓存未命中
			return nil, false
		}
		log.Printf("Error getting todos from Redis: %v", err) // 其他错误
		return nil, false
	}

	var todos []models.Todo
	err = json.Unmarshal([]byte(val), &todos)
	if err != nil {
		log.Printf("Error unmarshaling todos from Redis: %v", err)
		return nil, false
	}
	log.Println("Todos loaded from Redis cache.")
	return todos, true
}

// SetTodosCache 将所有 Todos 列表存入 Redis 缓存
func (r *Client) SetTodosCache(todos []models.Todo, expiration time.Duration) {
	data, err := json.Marshal(todos)
	if err != nil {
		log.Printf("Error marshaling todos for Redis: %v", err)
		return
	}
	err = r.Set(ctx, TodosCacheKey, data, expiration).Err()
	if err != nil {
		log.Printf("Error setting todos in Redis: %v", err)
	} else {
		log.Println("Todos cached in Redis.")
	}
}

// InvalidateTodosCache 清除所有 Todos 的缓存
// 在任务被创建、更新或删除时调用，以确保数据一致性
func (r *Client) InvalidateTodosCache() {
	err := r.Del(ctx, TodosCacheKey).Err()
	if err != nil && err != redis.Nil { // redis.Nil 表示键不存在，不是错误
		log.Printf("Error invalidating todos cache in Redis: %v", err)
	} else {
		log.Println("Todos cache invalidated in Redis.")
	}
}
