package handlers

import (
	"log" // 添加 log 导入
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo/internal/database/models"
	"todo/internal/kafka"
	"todo/internal/redis"
)

type TodoHandler struct {
	DB     *gorm.DB
	RedisC *redis.Client
	KafkaP *kafka.Producer
}

func NewTodoHandler(db *gorm.DB, redisC *redis.Client, kafkaP *kafka.Producer) *TodoHandler {
	return &TodoHandler{
		DB:     db,
		RedisC: redisC,
		KafkaP: kafkaP,
	}
}

// CreateTodo handles POST /api/create
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		log.Printf("Error binding JSON for CreateTodo: %v", err) // 添加日志
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if todo.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task name cannot be empty"})
		return
	}

	result := h.DB.Create(&todo)
	if result.Error != nil {
		log.Printf("Error creating todo in DB: %v", result.Error) // 添加日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	// 任务创建后，清理相关 Redis 缓存
	h.RedisC.InvalidateTodosCache()
	// 发送 Kafka 消息
	h.KafkaP.ProduceMessage(kafka.TodoEvent{
		EventType: "CREATE",
		TodoID:    todo.ID,
		TodoName:  todo.Name,
		Timestamp: time.Now(),
	})

	log.Printf("Successfully created todo: %+v", todo) // ❌ 新增日志：查看创建后的 todo 对象
	c.JSON(http.StatusCreated, todo)
}

// GetAllTodos handles GET /api/get-all-todos
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, found := h.RedisC.GetTodosCache()
	if found {
		log.Printf("Todos loaded from Redis cache: %d items", len(todos)) // 添加日志
		c.JSON(http.StatusOK, todos)
		return
	}

	var todosDB []models.Todo
	result := h.DB.Order("id desc").Find(&todosDB)
	if result.Error != nil {
		log.Printf("Error fetching todos from DB: %v", result.Error) // 添加日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	h.RedisC.SetTodosCache(todosDB, 5*time.Minute)                                       // 缓存 5 分钟
	log.Printf("Todos fetched from DB and cached: %d items, %+v", len(todosDB), todosDB) // ❌ 新增日志：查看从 DB 获取的 todo 列表
	c.JSON(http.StatusOK, todosDB)
}

// UpdateTodo handles POST /api/update
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	var req models.Todo
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for UpdateTodo: %v", err) // 添加日志
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("UpdateTodo request received: %+v", req) // ❌ 新增日志：查看收到的更新请求

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID cannot be empty"})
		return
	}

	var existingTodo models.Todo
	result := h.DB.First(&existingTodo, req.ID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			log.Printf("Error finding todo %d in DB: %v", req.ID, result.Error) // 添加日志
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find todo"})
		}
		return
	}

	if req.Name != "" {
		existingTodo.Name = req.Name
	}
	// description 可以为空字符串，所以直接赋值
	existingTodo.Description = req.Description
	existingTodo.Completed = req.Completed

	result = h.DB.Save(&existingTodo)
	if result.Error != nil {
		log.Printf("Error updating todo %d in DB: %v", req.ID, result.Error) // 添加日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	h.RedisC.InvalidateTodosCache()
	h.KafkaP.ProduceMessage(kafka.TodoEvent{
		EventType: "UPDATE",
		TodoID:    existingTodo.ID,
		TodoName:  existingTodo.Name,
		Completed: existingTodo.Completed,
		Timestamp: time.Now(),
	})

	log.Printf("Successfully updated todo: %+v", existingTodo) // ❌ 新增日志：查看更新后的 todo 对象
	c.JSON(http.StatusOK, existingTodo)
}

// DeleteTodo handles POST /api/delete
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	var req struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for DeleteTodo: %v", err) // 添加日志
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("DeleteTodo request received for ID: %d", req.ID) // ❌ 新增日志：查看收到的删除请求 ID

	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID cannot be empty"})
		return
	}

	result := h.DB.Delete(&models.Todo{}, req.ID)
	if result.Error != nil {
		log.Printf("Error deleting todo %d from DB: %v", req.ID, result.Error) // 添加日志
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	h.RedisC.InvalidateTodosCache()
	h.KafkaP.ProduceMessage(kafka.TodoEvent{
		EventType: "DELETE",
		TodoID:    req.ID,
		Timestamp: time.Now(),
	})

	log.Printf("Successfully deleted todo ID: %d", req.ID) // 添加日志
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

// GetTodoByID handles GET /api/todos/:id (新添加的，目前未在前端使用)
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var todo models.Todo
	result := h.DB.First(&todo, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo"})
		}
		return
	}
	c.JSON(http.StatusOK, todo)
}
