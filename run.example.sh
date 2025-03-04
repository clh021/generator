#!/bin/bash

set -e  # 遇到错误立即退出
set -x  # 打印执行的每一行命令

# 设置Go环境变量
export GO111MODULE=on

# 进入项目根目录
cd "$(dirname "$0")/"

# 下载依赖
go mod tidy
go mod vendor

cd ./examples/element-book/
go run ../../cmd/v1/main.go