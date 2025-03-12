# 模板生成器

这是一个基于配置的代码生成器，支持从模板文件生成代码，具有以下特点：

- 基于工作目录的配置文件管理
- 支持多个模板文件和配置文件
- 模板依赖管理
- 详细的错误报告
- 支持作为库引入使用

## 普通使用者

### 参考目录树结构

> 关键是在当前目录下能找到 `generator.yaml` 配置文件

```
project/
├── generator.yaml  # 主配置文件
├── configs/        # 其他配置文件目录
│   ├── base.yaml
│   └── extra.yaml
├── templates/      # 模板文件目录
│   ├── main.tpl
│   └── sub.tpl
└── generated/      # 生成的代码目录
```
### 命令行使用

```bash
# 在当前目录执行生成
generator

# 指定工作目录
generator -dir /path/to/workdir

# 指定模板目录
generator -template /path/to/template/dir

# 指定配置目录
generator -config /path/to/config/dir

# 指定输出目录
generator -output /path/to/output/dir

# 组合使用多个参数
generator -dir /path/to/workdir -template /path/to/template/dir -config /path/to/config/dir -output /path/to/output/dir
```

### 使用示例

```bash
cd examples/go-projects
go run ../../cmd/v1/main.go
```

## 开发者

### 配置文件格式

主配置文件使用 YAML 格式，需要命名为 `generator.yaml`，示例：

```yaml
templates:
  - path: "templates/main.tpl"
    output: "generated/main.go"
    dependencies:
      - "configs/base.yaml"
      - "configs/extra.yaml"
  - path: "templates/sub.tpl"
    output: "generated/sub.go"
    dependencies:
      - "configs/sub.yaml"
```

### 作为库使用

```go
import "generator"

func main() {
    gen := generator.New()
    if err := gen.Generate("/path/to/workdir"); err != nil {
        log.Fatal(err)
    }
}
```

## 错误处理

程序会在以下情况报错：

- 配置文件不存在或格式错误
- 模板文件不存在
- 依赖的配置文件不存在
- 模板变量未定义
- 模板语法错误

错误信息会包含详细的文件路径和行号信息，便于快速定位问题。

## 测试

运行单元测试：

```bash
go test ./...
```

运行特定包的测试：

```bash
go test ./internal/template
```

## 将来支持的特性

- 支持使用 go/template 特性，直接引用子模板
- 模板的依赖关系：为生成一个大模板文件而拆分设计，可以通过指定不同小模板解析生成后合并而来
- 模板合并顺序：不同小模板解析后合并生成可以定义顺序，保证模板生成结果符合预期
- 自定义辅助函数：允许用户定义和注册自己的模板辅助函数
- 增量生成：只更新发生变化的文件，提高生成效率
- 插件系统：支持通过插件扩展生成器的功能

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。在提交 PR 之前，请确保您的代码通过了所有的测试，并且符合项目的代码风格。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。