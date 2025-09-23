# --- Stage 1: Builder ---
# 使用官方的 Go 语言镜像作为构建环境
FROM golang:1.24.3-alpine AS builder

# 设置工作目录，后续所有命令都将在此目录下执行
WORKDIR /app

# 复制 go.mod 和 go.sum 文件到工作目录
# 这样 Go Modules 可以利用缓存，如果这两个文件不变，则不会重新下载所有依赖
COPY go.mod .
COPY go.sum .

# 下载所有 Go 模块依赖
# 如果 go.mod 和 go.sum 没有变化，这一步会使用缓存，加快构建速度
RUN go mod download

# 复制所有项目源代码到工作目录
# 包括 main.go 和 internal/ 目录等
COPY . .

# 构建 Go 应用程序
# -o /usr/local/bin/todo-app: 指定输出的可执行文件名为 todo-app，并放置在 /usr/local/bin 目录下
# -trimpath: 移除所有文件路径前缀，减小二进制文件大小
# -ldflags "-s -w": 移除调试信息和符号表，进一步减小二进制文件大小
# ./main.go: 指定 Go 程序的入口文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /usr/local/bin/todo-app ./main.go

# --- Stage 2: Final Image ---
# 使用一个更小的 Alpine Linux 镜像作为最终运行环境
# 这样最终的 Docker 镜像会非常小，只包含运行时所需的最少组件
FROM alpine:3.18

# 设置时区（可选，但推荐）
# RUN apk add --no-cache tzdata
# ENV TZ Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的可执行文件
# 注意：这里 /usr/local/bin/todo-app 是 builder 阶段的输出路径
COPY --from=builder /usr/local/bin/todo-app .

# 复制前端静态文件到最终镜像
# 确保这个路径与你的 Gin 路由中 `router.Static("/static", "./static")` 的第二个参数相匹配
COPY static ./static

# 暴露应用程序监听的端口
EXPOSE 8000

# 定义容器启动时执行的命令
# 这里直接运行我们编译好的 todo-app 可执行文件
CMD ["./todo-app"]