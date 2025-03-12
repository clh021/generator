#!/usr/bin/env bash
# leehom Chen clh021@gmail.com

# 使用说明：
# 1. 运行格式化: ./fmt-modify.sh
# 2. 安装 pre-commit 钩子: ./fmt-modify.sh --install-hook
# 3. 格式化特定目录: ./fmt-modify.sh /path/to/directory

# 设置项目根目录
PROJECT_ROOT=$(git rev-parse --show-toplevel)
PRETTIER_BIN="$PROJECT_ROOT/rjgc/node_modules/.bin/prettier"

# 开关变量
FORMAT_BACKEND=true
FORMAT_FRONTEND=false  # 设置为 true 以启用前端格式化

# 检查项目根目录是否存在
if [ ! -d "$PROJECT_ROOT" ]; then
  echo "无法找到项目根目录，请确保在项目仓库内运行此脚本。"
  exit 1
fi

# 安装 pre-commit 钩子的函数
install_pre_commit() {
    local pre_commit_path="$PROJECT_ROOT/.git/hooks/pre-commit"
    local pre_commit_content='#!/usr/bin/env bash
# 自动生成的 pre-commit 钩子

# 获取 git 仓库根目录
REPO_ROOT=$(git rev-parse --show-toplevel)

# 运行格式化脚本
"$REPO_ROOT/fmt-modify.sh"

# 检查格式化脚本是否成功执行
if [ $? -ne 0 ]; then
    echo "格式化失败，commit 操作已取消。"
    exit 1
fi'

    echo "$pre_commit_content" > "$pre_commit_path"
    chmod +x "$pre_commit_path"
    echo "pre-commit 钩子已成功安装到 $pre_commit_path"
}

# 处理命令行参数
if [ "$1" = "--install-hook" ]; then
    install_pre_commit
    exit 0
fi

# 格式化后端文件的函数
format_backend_files() {
    local target_path=$1
    echo "开始格式化 $target_path 下的所有 .go 文件..."
    find "$target_path" -type f -name "*.go" -exec gofmt -w {} \;
}

# 格式化前端文件的函数
format_frontend_files() {
    local target_path=$1
    if [ -f "$PRETTIER_BIN" ]; then
        echo "开始格式化 $target_path 下的所有前端文件..."
        $PRETTIER_BIN --write "$target_path/**/*.{js,jsx,ts,vue,tsx,json}"
    else
        echo "无法找到 prettier 程序，跳过前端文件格式化。"
    fi
}

# 处理特定类型文件的函数
format_files() {
    local files=$1
    local format_cmd=$2
    local file_pattern=$3

    if [ -z "$files" ]; then
        echo "没有需要格式化的 $file_pattern 文件。"
    else
        echo "开始格式化 $file_pattern 文件..."
        for file in $files; do
            if [ -f "$file" ]; then
                echo "正在格式化: $file"
                $format_cmd "$file"
                git add "$file"
                echo "已格式化并添加到暂存区: $file"
            else
                echo "文件 $file 不存在，跳过格式化。"
            fi
        done
    fi
}

# 主要逻辑
if [ $# -gt 0 ]; then
    TARGET_PATH="$1"
    if [ -d "$TARGET_PATH" ]; then
        [ "$FORMAT_BACKEND" = true ] && format_backend_files "$TARGET_PATH"
        [ "$FORMAT_FRONTEND" = true ] && format_frontend_files "$TARGET_PATH"
        echo "格式化完成！"
        exit 0
    else
        echo "提供的路径 $TARGET_PATH 不存在。"
        exit 1
    fi
fi

# 处理修改的文件
if [ "$FORMAT_BACKEND" = true ]; then
    MODIFIED_GO_FILES=$(git diff --name-only --cached | grep '\.go$')
    format_files "$MODIFIED_GO_FILES" "gofmt -w" ".go"
fi

if [ "$FORMAT_FRONTEND" = true ]; then
    if [ -f "$PRETTIER_BIN" ]; then
        MODIFIED_FRONTEND_FILES=$(git diff --name-only --cached | grep -E '\.(js|jsx|ts|vue|tsx|json)$')
        format_files "$MODIFIED_FRONTEND_FILES" "$PRETTIER_BIN --write" "前端"
    else
        echo "无法找到 prettier 程序，跳过前端文件格式化。"
    fi
fi
echo "格式化完成！"
