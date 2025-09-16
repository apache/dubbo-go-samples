# Dubbo-Go 微服务商城项目 - Makefile 使用指南

## 📋 项目概述

这是一个基于 Dubbo-Go 的微服务商城系统，包含5个核心服务：
- **Frontend** - 商城前端 (端口: 8080)
- **User** - 用户服务 (端口: 20001)  
- **Detail** - 商品详情服务 (端口: 20002)
- **Order** - 订单服务 (端口: 20003)
- **Comment** - 评论服务 (端口: 20004)

## 🚀 快速开始

### 一键启动项目
```bash
make quick-start
```
这个命令会自动完成：初始化目录 → 安装依赖 → 构建服务 → 启动所有服务

### 分步操作
```bash
# 1. 初始化项目目录
make init

# 2. 安装依赖
make install-deps

# 3. 构建所有服务
make build

# 4. 启动服务 (V1版本)
make start

# 或启动 V2 版本服务
make start-v2
```

## 🔧 常用命令

### 服务管理
```bash
make start          # 启动所有服务 (V1版本)
make start-v2       # 启动所有服务 (V2版本)
make stop           # 停止所有服务
make restart        # 重启所有服务
make status         # 查看服务运行状态
```

### 开发工具
```bash
make build          # 构建所有服务
make clean          # 清理构建文件
make clean-all      # 深度清理 (包括日志和PID文件)
make logs           # 查看所有服务日志
make fmt            # 格式化Go代码
make vet            # 代码静态检查
```

### 测试功能
```bash
make test           # 运行项目测试
make test-user      # 测试用户服务客户端
make test-detail    # 测试商品详情服务客户端
make test-order     # 测试订单服务客户端
make test-comment   # 测试评论服务客户端
```

### 开发工具
```bash
make proto          # 生成 Protocol Buffers 代码
```

## 📁 项目结构

启动后会创建以下目录：
```
shop/
├── bin/            # 编译后的可执行文件
├── logs/           # 服务日志文件
├── pids/           # 进程ID文件
└── Makefile        # 构建脚本
```

## 🌐 访问地址

启动成功后，可以通过以下地址访问：
- **前端页面**: http://localhost:8080
- **用户服务**: localhost:20001
- **商品详情服务**: localhost:20002  
- **订单服务**: localhost:20003
- **评论服务**: localhost:20004

## 🔍 监控和调试

### 查看服务状态
```bash
make status
```

### 查看服务日志
```bash
make logs
```

### 查看特定服务日志
```bash
tail -f logs/frontend.log    # 前端服务日志
tail -f logs/user.log        # 用户服务日志
tail -f logs/detail-v1.log   # 商品详情服务日志
tail -f logs/order-v1.log    # 订单服务日志
tail -f logs/comment-v1.log  # 评论服务日志
```

## ⚙️ 高级配置

### 服务版本切换
项目支持V1和V2两个版本的服务：
```bash
make start     # 启动V1版本
make start-v2  # 启动V2版本
```

### 开发模式
安装热重载工具：
```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Protocol Buffers
重新生成protobuf代码：
```bash
make proto
```

## 📝 常见问题

### Q: 服务启动失败怎么办？
A: 检查端口是否被占用，查看日志文件排查错误：
```bash
make status
make logs
```

### Q: 如何完全重置项目？
A: 使用深度清理命令：
```bash
make stop
make clean-all
make quick-start
```

### Q: 如何只启动部分服务？
A: 可以手动启动特定服务：
```bash
make build
./bin/frontend &        # 只启动前端
./bin/user &           # 只启动用户服务
```

## 🛠️ 开发建议

1. **开发时使用热重载**: `make dev`
2. **提交前检查代码**: `make fmt && make vet`
3. **定期更新依赖**: `make mod-update`
4. **运行测试**: `make test`

## 📞 技术支持

如遇到问题，请查看：
1. 服务日志: `make logs`
2. 服务状态: `make status`
3. 项目文档: `READEME_CN.md`
