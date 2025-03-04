#!/bin/bash

set -e  # 遇到错误立即退出
set -x  # 打印执行的每一行命令

# 显示帮助信息
show_help() {
    echo "Usage: $0 [example]"
    echo
    echo "Available examples:"
    echo "  dynamic     - Run antdv-dynamic example (default)"
    echo "  book        - Run antdv-book-fields example"
    echo
    exit 1
}

# 设置Go环境变量
export GO111MODULE=on

# 进入项目根目录
cd "$(dirname "$0")/"

# 下载依赖
go mod tidy
go mod vendor

# 获取示例参数
EXAMPLE=${1:-dynamic}

case $EXAMPLE in
    dynamic)
        echo "Running antdv-dynamic example..."
        cd ./examples/antdv-dynamic/
        ;;
    book)
        echo "Running antdv-book-fields example..."
        cd ./examples/antdv-book-fields/
        ;;
    -h|--help)
        show_help
        ;;
    *)
        echo "Invalid example: $EXAMPLE"
        show_help
        ;;
esac

# 运行示例
go run ../../cmd/v1/main.go
