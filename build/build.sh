#!/bin/bash

set -e  # 遇到错误立即退出
# set -x  # 打印执行的每一行命令

# 设置Go环境变量
export GO111MODULE=on
export CGO_ENABLED=0  # 禁用CGO，生成静态二进制文件

# 进入项目根目录
cd "$(dirname "$0")/../"

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

# 确保 build/dist/bin 目录存在
mkdir -p build/dist/bin

# 下载依赖
go mod tidy
go mod vendor

# 编译项目，添加更多的编译信息
go build -v -trimpath -ldflags="${LDFLAGS} -s -w" -o build/dist/bin/generator cmd/v1/*.go

# 显示编译后的文件大小
echo "编译后的文件大小:"
ls -lh build/dist/bin/generator

# 如果安装了 UPX，使用它来压缩二进制文件
if command -v upx &> /dev/null; then
    echo "使用 UPX 压缩二进制文件..."
    upx --best --lzma build/dist/bin/generator
    
    # 显示压缩后的文件大小
    echo "压缩后的文件大小:"
    ls -lh build/dist/bin/generator
else
    echo "UPX 未安装，跳过压缩步骤。如需压缩二进制文件，请安装 UPX。"
fi

echo "构建完成，版本号: ${VERSION}"

# 检查编译是否成功
if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

# 确保输出目录存在
mkdir -p build/dist/bin/

# 运行 astgen
# build/dist/bin/astgen
