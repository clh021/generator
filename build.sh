#!/bin/bash

set -e  # 遇到错误立即退出
# set -x  # 打印执行的每一行命令

# 设置Go环境变量
export GO111MODULE=on

# 进入项目根目录
cd "$(dirname "$0")/"

# 获取 Git Commit ID
COMMIT_ID=$(git rev-parse --short HEAD)

# 获取 Git Commit Time (RFC3339 格式)
COMMIT_TIME=$(git log -1 --format=%aI)

# 获取 Commit Count
COMMIT_COUNT=$(git rev-list --count HEAD)

# 构建版本号
VERSION="0.0.${COMMIT_COUNT}"

# 获取构建时间
BUILD_TIME=$(date -Iseconds)

# 设置 -ldflags 参数
LDFLAGS="-X main.Version=${VERSION} -X main.CommitID=${COMMIT_ID} -X main.CommitTime=${COMMIT_TIME} -X main.BuildTime=${BUILD_TIME}"

# 确保 dist/bin 目录存在
mkdir -p dist/bin

# 下载依赖
go mod tidy
go mod vendor

# 编译项目，添加更多的编译信息
go build -v -ldflags="${LDFLAGS}" -o dist/bin/generator cmd/v1/*.go

echo "构建完成，版本号: ${VERSION}"

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

# 确保输出目录存在
mkdir -p dist/bin/

# 运行 astgen
# dist/bin/astgen
