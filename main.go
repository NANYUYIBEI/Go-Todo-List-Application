package main

import (
	"fmt"       // 用于格式化字符串
	"log"       // 用于日志输出
	"os"        // 用于操作系统接口，例如信号处理
	"os/signal" // 用于捕获操作系统信号
	"syscall"   // 用于系统调用，例如 SIGINT 和 SIGTERM

	"todo/internal/config"   // 导入配置模块
	"todo/internal/database" // 导入数据库模块
	"todo/internal/handlers" // 导入处理器模块
	"todo/internal/kafka"    // 导入 Kafka 模块
	"todo/internal/redis"    // 导入 Redis 模块
	"todo/internal/routers"  // 导入路由模块
)

func main() {
	// 1. 加载应用程序配置
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully.")

	// 2. 初始化 MySQL 数据库连接
	// 注意：这里的 cfg.MySQLDSN 可能会包含敏感信息，生产环境应考虑更安全的管理方式
	database.InitMySQL(cfg.MySQLDSN)
	// defer database.GetDB().Close() // Gorm v2 推荐使用连接池，通常无需手动关闭

	// 3. 初始化 Redis 客户端连接
	redisClient := redis.InitRedis(cfg.RedisAddr)
	// defer redisClient.Close() // go-redis 客户端通常会自动管理连接池，但在某些场景下，你可能希望显式关闭

	// 4. 初始化 Kafka 消息生产者
	kafkaProducer := kafka.InitProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	defer kafkaProducer.Close() // 确保在应用退出时关闭 Kafka 生产者

	// 5. (可选) 启动 Kafka 消费者作为一个 Goroutine
	// 实际生产环境，Kafka 消费者通常会是独立的微服务，或者在单独的进程中运行。
	// 这里为了演示目的，将其作为主应用的一个 Goroutine 启动。
	kafka.StartConsumer(cfg.KafkaBroker, cfg.KafkaTopic, "todo-event-group")
	log.Println("Kafka consumer started in background.")

	// 6. 初始化 Gin 处理器，注入所有必要的依赖
	// 这样处理器函数就可以访问数据库、Redis 客户端和 Kafka 生产者
	todoHandler := handlers.NewTodoHandler(database.GetDB(), redisClient, kafkaProducer)
	log.Println("TodoHandler initialized with dependencies.")

	// 7. 设置 Gin 路由，包括静态文件服务和 API 接口
	router := routers.SetupRouter(todoHandler)
	log.Println("Gin router configured.")

	// 8. 启动 Gin HTTP 服务器，并监听指定端口
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Gin server failed to start on %s: %v", serverAddr, err)
		}
	}()
	log.Printf("Gin server started on %s", serverAddr)
	log.Printf("Access the web application at http://localhost%s/static/index.html", serverAddr)

	// --- 优雅关机处理 ---
	// 创建一个通道来接收操作系统信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT (Ctrl+C) 和 SIGTERM (终止信号)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞主 Goroutine，直到接收到退出信号
	<-quit
	log.Println("Shutting down server gracefully...")

	// 在这里可以执行其他的清理工作，例如：
	// - 等待 Gin 服务器的所有活动请求完成
	// - 对于 Kafka 消费者，如果其 StartConsumer 中有 stop channel，可以在此发送信号以优雅停止。
	//   （目前示例消费者是无限循环，但其所在的 Goroutine 会随主进程退出而停止）

	log.Println("Server exited.")
}
