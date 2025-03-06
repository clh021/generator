# 模板生成器

这是一个基于配置的代码生成器，支持从模板文件生成代码，具有以下特点：

- 基于工作目录的配置文件管理
- 支持多个模板文件和配置文件
- 模板依赖管理
- 详细的错误报告
- 支持作为库引入使用

## 普通使用者

###参的目录树结构

> 关键是在当前目录下能找到 `generator.yaml` 配置文件

```
project/
├── generator.yaml  # 配置文件
├── templates/      # 模板文件目录
│   └── main.tpl
└── generated/      # 生成的代码目录
```

### 命令行使用

```bash
# 在当前目录执行生成
generator

# 指定工作目录
generator -dir /path/to/workdir
```

### 使用示例

```bash
cd examples/go-projects
go run ../../cmd/v1/main.go
```

## 开发者

### 配置文件格式

配置文件使用 YAML 格式，需要命名为 `generator.yaml`，示例：

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

- 配置文件不存在
- 模板文件不存在
- 依赖的配置文件不存在
- 模板变量未定义
- 模板语法错误

错误信息会包含详细的文件路径和行号信息。