#!/bin/bash

# OpenClaw 监控系统启动脚本

echo "正在启动 OpenClaw 监控系统..."

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: 找不到 config.yaml 配置文件"
    exit 1
fi

# 创建数据目录
mkdir -p data

# 进入后端目录
cd backend

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未安装 Go"
    exit 1
fi

# 下载依赖
echo "正在下载依赖..."
go mod download

# 运行服务
echo "正在启动服务..."
go run cmd/server/main.go -config ../config.yaml
