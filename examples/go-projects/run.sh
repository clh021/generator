#!/bin/bash

# 启用错误检查和命令回显
set -e # 遇到错误立即退出
# set -x # 打印执行的每一行命令

# 设置Go环境变量
export GO111MODULE=on

# 定义目录路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
EXAMPLE_DIR="${SCRIPT_DIR}"
CMD_DIR="${PROJECT_ROOT}/cmd/v1"

echo "示例目录: ${EXAMPLE_DIR}"

# 进入项目根目录
cd "${PROJECT_ROOT}"

# 下载依赖
go mod tidy
go mod vendor

# 进入当前 example 目录
cd "${EXAMPLE_DIR}"

# 运行生成器
go run "${CMD_DIR}/main.go" -template "${EXAMPLE_DIR}/templates" -variables "${EXAMPLE_DIR}/variables" -output "${EXAMPLE_DIR}/output"

echo "项目生成完成！"
echo "生成的项目位于: ${EXAMPLE_DIR}/output"
echo ""
echo "要运行生成的项目，请执行以下步骤："
echo "1. cd ${EXAMPLE_DIR}/output"
echo "2. go mod tidy"
echo "3. go run main.go"
echo ""
echo "服务器将在 http://localhost:8080 上运行"
