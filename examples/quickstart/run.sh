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

# 打印目录信息
# echo "脚本目录: ${SCRIPT_DIR}"
# echo "项目根目录: ${PROJECT_ROOT}"
echo "示例目录: ${EXAMPLE_DIR}"
echo "命令目录: ${CMD_DIR}"
# 进入项目根目录
cd "${PROJECT_ROOT}"

# 下载依赖
go mod tidy
go mod vendor

# 进入当前 example 目录
cd "${EXAMPLE_DIR}"

# 清理之前的输出
rm -rf .gen_templates .gen_variables .gen_output .gen_config.yaml

# 运行生成器生成示例文件
go build -o generator "${CMD_DIR}/main.go"
./generator -quickstart

# 运行生成器处理示例文件
./generator

echo "快速开始示例已生成并执行。"
echo "生成的文件位于: ${EXAMPLE_DIR}/.gen_output"
