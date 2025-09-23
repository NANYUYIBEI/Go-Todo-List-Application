# Go Todo List 全栈应用 (Gin + MySQL + Redis + Kafka)

这是一个从基础架构升级而来的全栈待办事项列表应用程序，旨在展示使用 Go (Golang) Gin 框架、MySQL 数据库、Redis 缓存和 Kafka 消息队列构建高性能、高并发分布式应用的能力。前端使用 HTML、CSS 和原生 JavaScript 实现。

---

## 🌟 功能特性

*   **任务管理：** 实现任务的完整 CRUD (创建、读取、更新、删除) 操作。
*   **状态切换：** 轻松将任务标记为已完成或未完成。
*   **数据持久化：** 所有任务数据都可靠地存储在 **MySQL** 数据库中。
*   **高性能缓存：** 任务列表查询利用 **Redis** 缓存，显著提升读取响应速度。
*   **异步事件处理：** 任务的创建、更新、删除事件通过 **Kafka** 消息队列进行异步处理，实现服务解耦。
*   **前后端分离：** 基于 RESTful API 的经典前后端交互模式。
*   **容器化部署：** 提供 Dockerfile 和 docker-compose.yml，实现多服务（Go应用、MySQL、Redis、Kafka/Zookeeper）一键部署。
*   **性能优化：** 经过 JMeter 压测验证，具备在高负载下的良好性能表现。

## 🛠️ 技术栈

**后端 (GoLang)**
*   **Go (Golang)：** 应用核心逻辑，Go 版本 1.24.3。
*   **Gin Gonic：** 高性能 Web 框架，用于构建 RESTful API。
*   **Gorm：** 开发者友好的 Go ORM 库，简化数据库交互。
*   **MySQL：** 关系型数据库，用于持久化存储任务数据。
*   **Redis：** 内存数据结构存储，用作任务列表的缓存。
*   **Kafka：** 分布式流平台，用于任务事件的异步消息队列。
*   **`go-redis/redis/v8`：** Go 语言的 Redis 客户端。
*   **`segmentio/kafka-go`：** Go 语言的 Kafka 客户端。
*   **`joho/godotenv`：** 用于本地开发时加载 .env 文件（可选）。

**前端**
*   **HTML：** (`static/index.html`) 网页结构。
*   **CSS：** (`static/style.css`) 页面样式。
*   **JavaScript (原生)：** (`static/script.js`) 异步 API 通信和 DOM 操作。

**部署与测试**
*   **Docker：** 容器化技术。
*   **Docker Compose：** 多容器应用定义与运行工具。
*   **JMeter：** 开源性能测试工具，用于API压测。

---

## 🚀 安装与运行

你可以通过多种方式来运行此项目，推荐使用 Docker Compose 以获得最顺畅的体验。

### 方式一：通过 Docker Compose (推荐)

这是最推荐的运行方式，只需几条命令即可完成所有操作。

#### 前提条件

*   **Git：** 用于克隆代码仓库。
*   **Docker Desktop：** (包含 Docker Engine 和 Docker Compose) 确保其已安装并正在运行。

#### 步骤

1.  **克隆仓库：**
    ```bash
    git clone https://github.com/NANYUYIBEI/Go-Todo-List-Application.git
    cd Go-Todo-List-Application
    ```

2.  **配置 `docker-compose.yml` (重要！)**
    *   打开项目根目录下的 `docker-compose.yml` 文件。
    *   找到 `mysql` 服务下的 `environment` 部分和 `healthcheck` 部分。
    *   将 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD` 的值 **`password_secure`** (或者你之前设置的值) **替换为你自己的强密码**。确保这个密码与 `app` 服务 `environment` 中的 `MYSQL_DSN` 字符串里的密码**完全一致**。
        *   **示例 (假设你的密码是 `mySecurePassword123`)：**
            ```yaml
            services:
              mysql:
                # ...
                environment:
                  MYSQL_ROOT_PASSWORD: mySecurePassword123
                  MYSQL_DATABASE: todo_app
                  # MYSQL_USER: root # 保持注释或删除
                  MYSQL_PASSWORD: mySecurePassword123 # 与 MYSQL_ROOT_PASSWORD 保持一致
                healthcheck:
                  test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pmySecurePassword123"] # 这里也要替换
                  # ...
              app:
                # ...
                environment:
                  MYSQL_DSN: "root:mySecurePassword123@tcp(mysql:3306)/todo_app?charset=utf8mb4&parseTime=True&loc=Local" # 这里也要替换
                  # ...
            ```
    *   **移除 `version: '3.8'` 行：** Docker Compose 较新版本中此行已过时，移除以避免警告。

3.  **启动所有服务：**
    *   在终端中，运行以下命令。这会构建你的 Go 应用镜像（如果需要），拉取其他服务镜像，并启动整个应用栈。
    ```bash
    docker-compose up -d --build
    ```
    *   **故障排除提示：**
        *   如果遇到端口冲突（例如 3306、6379、9092、2181 或 8000 端口被占用），请使用 `netstat -ano | findstr :<端口号>` 查找占用进程的 PID，然后在任务管理器中结束该进程。
        *   如果 Go 应用容器启动失败（日志中显示 MySQL 连接失败），请检查 `docker-compose.yml` 中 `app` 服务的 `MYSQL_DSN` 密码是否与 `mysql` 服务的 `MYSQL_ROOT_PASSWORD` 一致。
        *   你可以使用 `docker-compose ps` 查看所有容器的运行状态，确保它们都处于 `Up` (或 `healthy`) 状态。
        *   使用 `docker-compose logs -f` 可以实时查看所有服务的日志输出。

4.  **访问应用：**
    *   打开浏览器，访问 `http://localhost:8000/static/index.html`。
    *   尝试创建、更新、完成、删除任务。

5.  **停止服务：**
    *   在项目根目录运行：
    ```bash
    docker-compose down --volumes --remove-orphans
    ```
    *   这会停止并移除所有容器、网络和数据卷（包括 MySQL 和 Redis 的持久化数据）。

### 方式二：通过 Docker (手动构建/拉取应用镜像)

这种方式可以单独运行 Go 应用容器，但你需要确保其他依赖服务（MySQL、Redis、Kafka）已在其他地方运行并可访问。

#### 前提条件

*   **Git** 和 **Docker Desktop**。
*   **手动配置和运行 MySQL、Redis、Kafka 实例，并确保 Go 应用容器可以访问它们。**

#### 步骤 (从源码构建 Go 应用镜像并运行)

1.  **克隆仓库：**
    ```bash
    git clone https://github.com/NANYUYIBEI/Go-Todo-List-Application.git
    cd Go-Todo-List-Application
    ```
2.  **构建 Docker 镜像：**
    ```bash
    docker build -t go-todo-app:1.0 .
    ```
3.  **运行 Docker 容器：**
    *   你需要配置环境变量以指向你手动运行的 MySQL、Redis、Kafka 服务。
    *   **示例：** (假设你在本地运行了 MySQL, Redis, Kafka，或者它们在可访问的 IP 地址上)
        ```bash
        docker run -d -p 8000:8000 \
          -e MYSQL_DSN="root:mySecurePassword123@tcp(host.docker.internal:3306)/todo_app?charset=utf8mb4&parseTime=True&loc=Local" \
          -e REDIS_ADDR="host.docker.internal:6379" \
          -e KAFKA_BROKER="host.docker.internal:9092" \
          -e KAFKA_TOPIC="todo_events" \
          --name my-go-app go-todo-app:1.0
        ```
        *   **`host.docker.internal`** 是 Docker Desktop 在 Windows/macOS 上提供的一个特殊 DNS 名称，用于容器访问宿主机上的服务。如果你在 Linux 上，可能需要使用宿主机的实际 IP 地址。
        *   请将 `mySecurePassword123` 替换为你的 MySQL 密码，并将 `3306`, `6379`, `9092` 替换为你的本地服务端口。
4.  **访问应用：** `http://localhost:8000/static/index.html`。

### 方式三：在本地直接运行 (传统方式)

这种方式无需 Docker，但需要你在本地环境中安装并手动启动所有依赖服务。

#### 前提条件

*   **Go (版本 1.24.3 或更高)**：确保 Go 环境已配置。
*   **Git**：用于克隆代码仓库。
*   **MySQL 服务器**：已安装并运行，监听在 `localhost:3306`。
*   **Redis 服务器**：已安装并运行，监听在 `localhost:6379`。
*   **Kafka 集群 (包含 Zookeeper)**：已安装并运行，Kafka 监听在 `localhost:9092`。

#### 步骤

1.  **克隆仓库：**
    ```bash
    git clone https://github.com/NANYUYIBEI/Go-Todo-List-Application.git
    cd Go-Todo-List-Application
    ```
2.  **配置 Go 环境变量 (如果你的 .env 文件不存在)：**
    *   在项目根目录下创建一个 `.env` 文件，内容如下：
        ```
        MYSQL_DSN="root:your_mysql_password@tcp(localhost:3306)/todo_app?charset=utf8mb4&parseTime=True&loc=Local"
        REDIS_ADDR="localhost:6379"
        KAFKA_BROKER="localhost:9092"
        KAFKA_TOPIC="todo_events"
        SERVER_PORT="8000"
        ```
    *   请将 `your_mysql_password` 替换为你的本地 MySQL root 用户密码。
3.  **下载 Go 依赖：**
    ```bash
    go mod tidy
    ```
4.  **运行应用程序：**
    ```bash
    go run main.go
    ```
    *   服务器将在 8000 端口启动。观察终端输出，确保 MySQL、Redis 和 Kafka 连接成功。

5.  **访问应用：** 打开浏览器访问 `http://localhost:8000/static/index.html`。

---

## 📈 性能测试 (使用 JMeter)

本项目通过 JMeter 进行了性能压测，以评估在高并发场景下的 API 响应时间、吞吐量及稳定性。

### 测试工具

*   **Apache JMeter**

### 测试场景

模拟100个并发用户，每个用户循环执行10次任务创建和任务获取操作。

### 测试计划

1.  **Thread Group 配置：**
    *   Number of Threads (users)：100
    *   Ramp-up period (seconds)：10
    *   Loop Count：10
2.  **HTTP Request Defaults：** Protocol: `http`, Server Name: `localhost`, Port: `8000`
3.  **`POST /api/create` (创建任务)：**
    *   Method: `POST`, Path: `/api/create`
    *   Headers: `Content-Type: application/json`
    *   Body Data:
        ```json
        {
            "name": "JMeter Test Task ${__RandomString(10,abcdefghijklmnopqrstuvwxyz,)}",
            "description": "Description for JMeter test task ${__time(yyyy-MM-dd HH:mm:ss,)}"
        }
        ```
4.  **`GET /api/get-all-todos` (获取所有任务)：**
    *   Method: `GET`, Path: `/api/get-all-todos`
    *   Headers: `Content-Type: application/json`
5.  **监听器：** View Results Tree, Summary Report, Aggregate Report。

### 性能测试结果概览

| Label        | # Samples | Average (ms) | Median (ms) | 90% Line (ms) | Error % | Throughput (req/sec) |
| :----------- | :-------- | :----------- | :---------- | :------------ | :------ | :------------------- |
| Create Todo  | 4229      | 1761         | -           | -             | 0.00%   | 124.5/sec            |
| Get All Todos| 3524      | 3272         | -           | -             | 0.00%   | 107.1/sec            |
| **TOTAL**    | **7753**  | **2448**     | -           | -             | **0.00%** | **228.2/sec**        |

*   **注：** 上述表格数据来源于 Summary Report，Aggregate Report 提供更详细的 Median, 90%, 95%, 99% Line 数据。
*   **总吞吐量 (Total QPS)：** **228.2 req/sec**
*   **错误率 (Error Rate)：** **0.00%**
*   **解读：** 在100个并发用户场景下，应用程序展现了卓越的稳定性，未发生任何错误。总吞吐量达到每秒228.2个请求，其中任务创建和获取操作分别达到了124.5 req/sec和107.1 req/sec，显示了Go Gin结合Redis缓存和Kafka异步处理带来的高性能。

---

## 📦 Docker Hub 镜像

本项目已将构建好的应用镜像发布至 Docker Hub，方便快速部署和验证。

*   **镜像名称：** `nanyui/go-todo-app:1.0`
*   **Docker Hub 地址：** [https://hub.docker.com/r/nanyui/go-todo-app](https://hub.docker.com/r/nanyui/go-todo-app) (请确保这个链接是正确的)

你可以通过以下命令直接拉取此镜像并集成到你自己的 `docker-compose.yml` 中（将 `app` 服务的 `build` 部分替换为 `image: nanyui/go-todo-app:1.0`）。

---

## 📂 项目结构

Go-Todo-List-Application/
├── .dockerignore
├── .gitignore
├── Dockerfile # Go 应用的 Dockerfile
├── docker-compose.yml # 包含所有服务的 Docker Compose 文件
├── go.mod # Go 模块依赖定义
├── go.sum
├── main.go # 应用入口，初始化配置、数据库、Redis、Kafka、Gin 路由
├── README.md # 本文档
├── internal/ # 内部模块
│ ├── config/ # 应用程序配置
│ │ └── config.go
│ ├── database/ # 数据库相关
│ │ ├── mysql.go # MySQL 连接和 Gorm 初始化
│ │ └── models/
│ │ └── todo.go # Gorm 数据模型
│ ├── handlers/ # Gin 路由处理函数
│ │ └── todo_handler.go
│ ├── kafka/ # Kafka 客户端相关
│ │ ├── common.go # Kafka 消息结构定义
│ │ ├── consumer.go # Kafka 消息消费者
│ │ └── producer.go # Kafka 消息生产者
│ ├── redis/ # Redis 客户端相关
│ │ └── client.go # Redis 连接和缓存操作
│ └── routers/ # Gin 路由定义
│ └── router.go
└── static/ # 前端静态文件
├── index.html
├── script.js
└── style.css


---

## 📡 API 端点

本应用后端提供以下 API 接口供前端调用：

*   **`POST /api/create`：** 创建一个新任务
    *   请求体 (Body -> raw -> JSON): `{"name": "我的新任务", "description": "这是一个描述"}`
*   **`GET /api/get-all-todos`：** 获取所有任务
*   **`POST /api/update`：** 更新指定任务
    *   请求体 (Body -> raw -> JSON): `{"id": 任务ID, "name": "新名字", "description": "新描述", "completed": true}`
*   **`POST /api/delete`：** 删除指定任务
    *   请求体 (Body -> raw -> JSON): `{"id": 要删除的任务ID}`

（**注意：** 本项目为简化起见，更新和删除操作仍使用 POST 方法，而不是标准的 PUT/DELETE 方法。）

---
