#!/bin/bash

set -e  # 遇到错误立即退出
# set -x  # 打印执行的每一行命令

# 设置Go环境变量
export GO111MODULE=on

# 进入项目根目录
cd "$(dirname "$0")/"

# 确保 dist/bin 目录存在
mkdir -p dist/bin


# 下载依赖
go mod tidy
go mod vendor

# 编译项目，添加更多的编译信息
go build -v -o dist/bin/gcode ./cmd/v1/main.go

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

# 确保输出目录存在
mkdir -p dist/bin/output

# 运行 astgen
dist/bin/astgen
