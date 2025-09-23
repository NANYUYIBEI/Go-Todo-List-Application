package routers

import (
	"net/http" // 用于 HTTP 状态码和重定向

	"github.com/gin-gonic/gin" // 导入 Gin 框架

	"todo/internal/handlers" // 导入 handlers 包，以便使用 TodoHandler
)

// SetupRouter 配置并返回一个 Gin 引擎，包含所有应用程序路由
func SetupRouter(todoHandler *handlers.TodoHandler) *gin.Engine {
	router := gin.Default() // 使用 Gin 默认的中间件 (Logger 和 Recovery)

	// 配置 CORS 中间件，允许跨域请求
	// 注意：在生产环境，建议将 Access-Control-Allow-Origin 设置为你的前端域名，而不是 "*"
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")         // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 允许携带认证信息（如 Cookies, Authorization Headers）
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE") // 允许的 HTTP 方法

		// 处理 OPTIONS 请求 (CORS 预检请求)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 返回 204 No Content，表示预检成功
			return
		}

		c.Next() // 继续处理请求
	})

	// 提供 "static" 目录下的静态文件 (HTML, CSS, JS)
	// 这使得可以通过 /static/index.html 访问前端页面
	router.Static("/static", "./static")

	// 重定向根路径 "/" 到静态文件 /static/index.html
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// 定义 API 路由组
	api := router.Group("/api") // 所有 API 接口都以 /api 前缀
	{
		// 保持与你原有的 API 路径兼容性
		api.POST("/create", todoHandler.CreateTodo)
		api.GET("/get-all-todos", todoHandler.GetAllTodos)
		api.POST("/update", todoHandler.UpdateTodo)
		api.POST("/delete", todoHandler.DeleteTodo)

		// 额外的，更 RESTful 的 API 风格（你可以选择性地启用这些，并调整前端）
		// api.GET("/todos", todoHandler.GetAllTodos)         // 获取所有任务
		// api.GET("/todos/:id", todoHandler.GetTodoByID)     // 获取单个任务
		// api.POST("/todos", todoHandler.CreateTodo)         // 创建任务
		// api.PUT("/todos/:id", todoHandler.UpdateTodo)       // 更新任务 (注意：你的原有更新是 POST /update，这里是 PUT /todos/:id)
		// api.DELETE("/todos/:id", todoHandler.DeleteTodo)    // 删除任务 (注意：你的原有删除是 POST /delete，这里是 DELETE /todos/:id)
	}

	return router
}
