package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go" // 确保能导入这个包
)

// Producer 结构体封装了 Kafka 消息生产者
type Producer struct {
	writer *kafka.Writer
}

// InitProducer 初始化 Kafka 生产者
func InitProducer(broker, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),   // Kafka Broker 地址 (例如 "localhost:9092")
		Topic:    topic,               // 消息发送到的 Topic
		Balancer: &kafka.LeastBytes{}, // 负载均衡策略：发送到字节数最少的分区
		// 可选配置，根据实际需求调整性能：
		// Async:    true, // 异步发送，提高吞吐量但可能丢失消息
		// BatchTimeout: 10 * time.Millisecond, // 批量发送的超时时间
		// BatchSize:    100, // 批量发送的消息数量
	}
	log.Println("Kafka producer initialized for topic:", topic)
	return &Producer{writer: w}
}

// ProduceMessage 发送 TodoEvent 消息到 Kafka
func (p *Producer) ProduceMessage(event TodoEvent) {
	// 将事件结构体编码为 JSON 字节数组
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal Kafka event: %v", err)
		return
	}

	// 构建 Kafka 消息
	msg := kafka.Message{
		// Key: []byte(strconv.FormatUint(uint64(event.TodoID), 10)), // 可选：根据 TodoID 作为 Key，保证同一 Todo 的消息顺序
		Value: data, // 消息内容
	}

	// 使用带超时的上下文发送消息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保上下文在函数结束时被取消

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to write Kafka message for event type %s, Todo ID %d: %v", event.EventType, event.TodoID, err)
	} else {
		log.Printf("Kafka message sent for event type %s, Todo ID %d", event.EventType, event.TodoID)
	}
}

// Close 关闭 Kafka 生产者
func (p *Producer) Close() {
	if p.writer != nil {
		if err := p.writer.Close(); err != nil {
			log.Printf("Failed to close Kafka writer: %v", err)
		} else {
			log.Println("Kafka producer closed.")
		}
	}
}
