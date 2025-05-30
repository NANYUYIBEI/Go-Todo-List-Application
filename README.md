这是一个使用 Go (Golang) 作为后端、HTML, CSS, 和 JavaScript 作为前端构建的简单待办事项列表应用程序。它允许用户创建、查看、更新、删除任务，并将任务标记为已完成。

## 功能特性

*   **添加任务**: 用户可以添加包含名称和描述的新任务。
*   **查看任务**: 以列表形式展示所有任务。
*   **更新任务**: 修改现有任务的详细信息。
*   **标记完成**: 将任务标记为已完成或未完成。
*   **删除任务**: 从列表中移除任务。
*   **响应式用户界面**: 简单的用户界面，适配不同设备。

## 技术栈

本项目旨在展示一个全栈应用的基础构成，主要技术栈如下：

### 前端

*   **HTML**: (`static/index.html`) 用于构建网页结构。
*   **CSS**: (`static/style.css`) 用于美化页面样式。
*   **JavaScript (原生)**: (`static/script.js`) 用于处理用户交互、DOM 操作以及与后端 API 的异步通信 (例如使用 `fetch` API)。

### 后端

*   **Go (Golang)**: (`main.go`)
    *   **原生 `net/http` 包**: 用于构建 Web 服务器，处理 HTTP 请求、路由、中间件等。这是学习 Go Web 开发的基础。
*   **Gorm**:
    *   一个优秀且对开发者友好的 Go ORM (Object-Relational Mapper) 库。
    *   用于简化数据库操作，如模型的定义 (`db/model.go`)、数据库的增删改查 (CRUD)、事务处理等。
*   **Sqlite**:
    *   一个轻量级的、基于文件的 SQL 数据库引擎。
    *   无需单独的数据库服务器进程，非常适合小型项目和本地开发。数据通常存储在单个 `.db` 文件中 (例如 `todo.db`，会在首次运行时自动创建)。

## 项目内容

*   **Go 原生开发服务器**:
    *   如何使用 Go 的 `net/http` 包搭建 HTTP 服务器。
    *   理解路由 (Routing) 的概念以及如何将不同的 URL 路径映射到相应的处理函数 (Handler)。
    *   如何处理 HTTP 请求 (GET, POST, PUT, DELETE) 并发送响应。
    *   如何提供静态文件服务 (HTML, CSS, JS from the `static` directory)。
*   **Gorm 的使用**:
    *   如何定义 Gorm 模型 (Structs in `db/model.go`) 来映射数据库表。
    *   如何配置 Gorm 连接到 Sqlite 数据库。
    *   使用 Gorm 执行常见的数据库操作：
        *   创建 (Create) 记录
        *   读取 (Read) 记录 (单个和列表)
        *   更新 (Update) 记录
        *   删除 (Delete) 记录
*   **Go 语言业务逻辑开发 (CRUD)**:
    *   如何将业务需求（如创建待办事项、标记完成等）转化为具体的 Go 代码逻辑。
    *   设计和实现针对 "Task" 资源的基本 CRUD (Create, Read, Update, Delete) API 接口。
*   **RESTful API 设计基础**:
    *   了解如何设计简单的 RESTful API 端点。
    *   理解 HTTP 方法 (GET, POST, PUT, DELETE) 在 REST API 中的语义。
    *   如何使用 JSON 在前后端之间传输数据。
*   **前后端分离概念**:
    *   前端 (`static/script.js`) 通过 JavaScript 发起 API 请求，后端 Go 程序 (`main.go`) 处理请求并返回数据。
*   **Postman (或类似 API 测试工具) 的使用**:
    *   如何使用 Postman 测试你开发的后端 API 接口，验证其功能是否符合预期。
*   **数据库管理工具 (DBMS) 的使用**:
    *   例如 DB Browser for SQLite，用于查看和管理 Sqlite 数据库文件 (`todo.db`)，检查表结构和数据。
*   **Go 项目基础结构**:
    *   了解一个简单 Go Web 项目可能的文件和目录组织方式。

## 安装与运行

### 先决条件

*   [Go](https://golang.org/dl/) (推荐最新稳定版本)
*   [Git](https://git-scm.com/) (用于克隆仓库)
*   一个文本编辑器或 IDE (如 VS Code)
*   (可选) [DB Browser for SQLite](https://sqlitebrowser.org/) 用于查看数据库内容。
*   (可选) [Postman](https://www.postman.com/downloads/) 用于测试 API。

### 步骤

1.  **克隆仓库**:
    ```bash
    git clone https://github.com/NANYUYIBEI/Go-Todo-List-Application.git
    cd Go-Todo-List-Application
    ```

2.  **安装依赖**:
    Go Modules 会自动处理依赖。`go.mod` 文件已存在，此步骤确保它们被下载和验证。
    ```bash
    go mod tidy
    ```

3.  **运行应用程序**:
    这将编译并运行 `main.go` 文件中的服务器。
    ```bash
    go run main.go
    ```
    服务器通常会监听在 `8000` 端口。你会在控制台看到类似 `Server started at :8000` 的输出。

4.  **在浏览器中打开**:
    打开你的网页浏览器，访问:
    [http://localhost:8000/static/index.html](http://localhost:8000/static/index.html)

## 项目结构

```
Go-Todo-List-Application/
├── .gitignore
├── go.mod
├── go.sum             # Go模块依赖项校验和 (由 go mod tidy 生成)
├── main.go            # Go程序入口，HTTP服务器设置和路由
├── README.md          # 本说明文件
├── db/
│   └── model.go       # Gorm数据模型定义
└── static/
    ├── index.html     # 前端HTML结构
    ├── script.js      # 前端JavaScript逻辑
    └── style.css      # 前端CSS样式
└── todo.db            # (运行时自动生成) Sqlite数据库文件
```

## API 端点 (请根据 `main.go` 中的实际实现进行核对和修改)

你的应用可能会暴露以下 API 端点供前端调用：

*   `GET /api/tasks` - 获取所有任务
*   `POST /api/tasks` - 创建一个新任务
    *   请求体 (JSON): `{"name": "Task Name", "description": "Task Description"}`
*   `GET /api/tasks/{id}` - 获取指定 ID 的任务
*   `PUT /api/tasks/{id}` - 更新指定 ID 的任务
    *   请求体 (JSON): `{"name": "Updated Name", "description": "Updated Description", "completed": true}`
*   `DELETE /api/tasks/{id}` - 删除指定 ID 的任务
*   `PUT /api/tasks/{id}/toggle` - (示例) 切换指定 ID 任务的完成状态 (也可能直接通过 `PUT /api/tasks/{id}` 更新 `completed` 字段实现)

**注意**: 上述 API 端点是常见的 RESTful 示例。请检查你的 `main.go` 文件以确认实际使用的路由和请求/响应格式。

## 使用 Postman 测试 API

1.  启动你的 Go 服务器 (`go run main.go`)。
2.  打开 Postman。
3.  根据你的实际 API 端点创建请求：
    *   **获取所有任务 (示例)**:
        *   方法: `GET`
        *   URL: `http://localhost:8000/api/tasks`
    *   **创建新任务 (示例)**:
        *   方法: `POST`
        *   URL: `http://localhost:8000/api/tasks`
        *   Body -> raw -> JSON:
          ```json
          {
              "name": "学习 Go",
              "description": "完成 Todo List 项目的后端 API"
          }
          ```
    *   以此类推，测试其他 PUT 和 DELETE 端点。

## 数据库查看

你可以使用 [DB Browser for SQLite](https://sqlitebrowser.org/) 打开项目根目录下的 `todo.db` (或 Gorm 配置中指定的数据库文件名) 来查看 `tasks` 表的结构和数据，验证 Gorm 操作是否正确。

## 贡献

欢迎提出改进建议！如果您想贡献，请 Fork 本仓库并发起 Pull Request。

