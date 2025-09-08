#!/bin/bash

echo "启动 GoPass 密码管理器..."

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 请先安装 Go 1.21+"
    echo "下载: https://golang.org/dl/"
    exit 1
fi

# 下载依赖并构建
go mod tidy
go build -o gopass ./cmd/server

# 启动
echo "访问: http://localhost:8080"
./gopass
