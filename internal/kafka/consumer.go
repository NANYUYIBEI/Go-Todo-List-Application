package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go" // 确保这个包能被导入
)

// StartConsumer 启动 Kafka 消费者
// 它会在一个独立的 Goroutine 中运行，持续从 Kafka topic 中读取消息。
// 注意：在实际生产环境中，Kafka 消费者通常会作为独立的微服务，或者有更健壮的错误处理和重试机制。
func StartConsumer(broker, topic, groupID string) {
	// 创建一个新的 Kafka Reader 配置
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker}, // Kafka Broker 地址列表
		Topic:    topic,            // 要消费的 Topic
		GroupID:  groupID,          // 消费者组ID，用于分区消费和负载均衡
		MinBytes: 10e3,             // 10KB：每次从 Kafka 获取消息的最小字节数
		MaxBytes: 10e6,             // 10MB：每次从 Kafka 获取消息的最大字节数
		MaxWait:  1 * time.Second,  // 最多等待 1 秒来获取新数据
		// CommitInterval: 1 * time.Second, // 可选：自动提交消息偏移量的时间间隔，如果设置为0则需要手动提交
		Logger:      kafka.LoggerFunc(log.Printf), // 设置用于日志输出的 logger
		ErrorLogger: kafka.LoggerFunc(log.Printf), // 设置用于错误日志输出的 logger
	})

	log.Printf("Kafka consumer started for topic %s, group %s", topic, groupID)

	// 在一个 Goroutine 中运行消费者，避免阻塞主线程
	go func() {
		for { // 无限循环，持续消费消息
			m, err := r.FetchMessage(context.Background()) // 阻塞式地获取下一条消息
			if err != nil {
				log.Printf("Error fetching Kafka message: %v", err)
				// 根据错误类型决定是否继续重试
				time.Sleep(5 * time.Second) // 避免在高频错误时频繁打印日志或空转
				continue
			}

			// 解析消息值到 TodoEvent 结构体
			var event TodoEvent
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				log.Printf("Error unmarshaling Kafka message (value: %s): %v", string(m.Value), err)
				// 即使反序列化失败，也提交消息，避免重复处理损坏的消息
				if err := r.CommitMessages(context.Background(), m); err != nil {
					log.Printf("Error committing Kafka message (unmarshal error): %v", err)
				}
				continue
			}

			// 打印接收到的消息详情
			log.Printf("Kafka message received: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s, EventType=%s, TodoID=%d",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value), event.EventType, event.TodoID)

			// -------------------------------------------------------------
			// 在这里处理你的业务逻辑，例如：
			// - 将事件写入审计日志数据库
			// - 触发其他微服务 (例如发送通知邮件、Webhook)
			// - 更新聚合服务的数据或进行数据同步
			// -------------------------------------------------------------

			// 提交消息偏移量，表示已成功处理该消息
			if err := r.CommitMessages(context.Background(), m); err != nil {
				log.Printf("Error committing Kafka message: %v", err)
			}
		}
	}()
	// 注意：在实际应用中，你可能需要一个机制来优雅地停止这个 Goroutine (例如通过一个 stop channel)
	// 但为了演示简单，这里没有实现。
}

// CloseConsumer 可以用来关闭 Kafka 消费者，但由于 StartConsumer 在 Goroutine 中无限循环，
// 且其 reader 实例不是全局可访问的，通常在 main 函数的优雅关机逻辑中，
// 对于这种后台 Goroutine 启动的消费者，需要更复杂的协调机制来停止。
// 对于本示例，我们暂不提供一个直接外部调用的 Close 方法，而是让其随应用进程退出。
