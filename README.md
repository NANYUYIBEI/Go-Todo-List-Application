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

推荐使用 Docker Compose 部署整个应用栈（Go应用、MySQL、Redis、Kafka、Zookeeper）。

### 前提条件

*   **Git：** 用于克隆代码仓库。
*   **Docker Desktop：** (包含 Docker Engine 和 Docker Compose) 确保其已安装并正在运行。

### 步骤

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
    *   如果遇到端口冲突（例如 3306 或 8000 端口被占用），请使用 `netstat -ano | findstr :<端口号>` 查找占用进程的 PID，然后在任务管理器中结束该进程。
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

你可以通过以下命令直接拉取并运行 Go 应用容器（需先配置好其他 Docker Compose 服务）：

```bash
docker pull nanyui/go-todo-app:1.0
