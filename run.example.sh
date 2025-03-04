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

# 运行 antdv-book-fields 示例
echo "Running antdv-book-fields example..."
cd ./examples/antdv-book-fields/
go run ../../cmd/v1/main.go

# 检查生成结果
if [ $? -eq 0 ]; then
    echo "生成成功！"
    echo "生成的文件位于：examples/antdv-book-fields/src/"
else
    echo "生成失败，请检查错误信息"
    exit 1
fi
