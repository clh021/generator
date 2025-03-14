# 代码生成器

这是一个基于模板的代码生成器，支持从配置文件生成代码。

## 特性

- 支持多个模板和配置文件
- 变量替换
- 错误报告（包括文件路径和行号）
- 支持在输出路径中使用变量
- 快速启动示例生成
- 支持多个变量文件

## 开始使用

1. 确保已安装 Go 和 Go 工具链。
2. 克隆项目仓库到本地机器。
3. 进入 `generator` 目录。
4. 运行 `go build -o generator cmd/v1/main.go` 来编译生成器。
5. 运行 `./generator -quickstart` 来生成快速启动示例。

## 使用方法

```
generator [选项]

选项:
  -dir string
        工作目录路径 (默认 ".")
  -output string
        输出目录路径 (默认 ".gen_output")
  -quickstart
        生成快速开始示例
  -template string
        模板目录路径 (默认 ".gen_templates")
  -variables string
        变量目录路径 (默认 ".gen_variables")
  -varfiles string
        变量文件路径，多个文件用逗号分隔
```

### 示例

1. 生成快速启动示例：
   ```
   ./generator -quickstart
   ```

2. 使用默认配置生成代码：
   ```
   ./generator
   ```

3. 指定工作目录：
   ```
   ./generator -dir /path/to/workdir
   ```

4. 自定义模板、变量和输出目录：
   ```
   ./generator -template /path/to/templates -variables /path/to/variables -output /path/to/output
   ```

5. 使用多个变量文件：
   ```
   ./generator -varfiles file1.yaml,file2.yaml
   ```

## 配置文件

生成器使用 YAML 格式的配置文件来定义模板和它们的依赖关系。

## 作为库使用

生成器也可以作为 Go 项目中的库使用。导入 `generator` 包并调用 `Generate` 函数。

## 模板特性

- 内置字符串处理函数
- 支持在输出路径中使用变量

## 错误处理

生成器提供详细的错误报告，包括文件路径和行号。

## 贡献

欢迎提交问题和拉取请求。

## 许可证

本项目采用 MIT 许可证。
```

这个更新后的 README.md 文件包含了新增功能的说明，特别是多个变量文件的支持和快速启动示例的使用方法。您可以根据需要进一步调整或补充这个 README 文件。