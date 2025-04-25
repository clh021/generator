# 代码生成器

这是一个基于模板的代码生成器，支持从配置文件生成代码。

## 特性

- 支持多个模板和配置文件
- 变量替换
- 错误报告（包括文件路径和行号）
- 支持在输出路径中使用变量
- 快速启动示例生成
- 支持多个变量文件
- **支持子模板：允许模板包含其他模板，实现模板复用和模块化。**
- **模板路径变量：支持在输出路径中使用变量引用，例如 `__variable__`，实现更灵活的文件组织。**
- **跳过子模板生成：自动跳过模板文件路径中包含 `__child__` 的文件，避免生成多余的子模板文件。**
- **按后缀跳过模板：跳过具有特定后缀的模板文件，例如 `.go.tpl.tpl`，以选择性地生成特定类型的文件。**
- **按前缀跳过模板：跳过具有特定路径前缀的模板文件，例如 `web/`，以选择性地生成服务端或客户端代码。**

## 开始使用

1.  确保已安装 Go 和 Go 工具链。
2.  克隆项目仓库到本地机器。
3.  进入 `generator` 目录。
4.  运行 `go build -o generator cmd/v1/main.go` 来编译生成器。
5.  运行 `./generator -quickstart` 来生成快速启动示例。

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
  -skip-suffixes string
        跳过特定后缀的模板文件，多个后缀用逗号分隔
        完整路径(path)进行匹配
        例如: -skip-suffixes=.go.tpl.tpl,.vue.tpl
  -skip-prefixes string
        跳过特定前缀路径的模板文件，多个前缀用逗号分隔
        相对于模板目录，不要前置/符号
        例如: -skip-prefixes=web,server/config
```

### 示例

1.  生成快速启动示例：

    ```
    ./generator -quickstart
    ```

2.  使用默认配置生成代码：

    ```
    ./generator
    ```

3.  指定工作目录：

    ```
    ./generator -dir /path/to/workdir
    ```

4.  自定义模板、变量和输出目录：

    ```
    ./generator -template /path/to/templates -variables /path/to/variables -output /path/to/output
    ```

5.  使用多个变量文件：

    ```
    ./generator -varfiles file1.yaml,file2.yaml
    ```

6.  跳过特定后缀的模板文件：

    ```
    ./generator -skip-suffixes=.go.tpl.tpl,.vue.tpl
    ```

7.  只生成服务端代码（跳过 web 模板）：

    ```
    ./generator -skip-prefixes=web
    ```

8.  只生成客户端代码（跳过 server 模板）：

    ```
    ./generator -skip-prefixes=server
    ```

## 配置文件

生成器使用 YAML 格式的配置文件来定义模板和它们的依赖关系。

## 作为库使用

生成器也可以作为 Go 项目中的库使用。导入 `github.com/clh021/generator/pkg/generator` 包并调用 `Generate` 函数。

### 示例代码

```go
package main

import (
	"log"
	"path/filepath"

	"github.com/clh021/generator/pkg/generator"
	"github.com/clh021/generator/pkg/config"
)

func main() {
	// 创建生成器实例
	gen := generator.NewGenerator()

	// 配置生成器
	cfg := &config.Config{
		TemplateDir:  "./templates",      // 模板目录
		VariablesDir: "./variables",     // 变量目录
		OutputDir:    "./output",        // 输出目录
		VariableFiles: []string{         // 可选：指定额外的变量文件
			"./custom_variables.yaml",
		},
		SkipTemplateSuffixes: ".go.tpl.tpl,.vue.tpl",  // 可选：跳过这些后缀的文件
		SkipTemplatePrefixes: "web",                   // 可选：跳过这些路径前缀的文件
	}

	// 执行生成
	if err := gen.Generate(cfg); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	log.Println("生成完成")
}
```

### 使用说明

1. **创建生成器实例**：使用 `generator.NewGenerator()` 创建一个新的生成器实例。

2. **配置生成器**：创建 `config.Config` 结构体并设置以下字段：
   - `TemplateDir`：模板目录路径
   - `VariablesDir`：变量目录路径
   - `OutputDir`：输出目录路径
   - `VariableFiles`：（可选）额外的变量文件路径列表
   - `SkipTemplateSuffixes`：（可选）跳过特定后缀的模板文件，多个后缀用逗号分隔，完整路径(path)进行匹配
   - `SkipTemplatePrefixes`：（可选）跳过特定前缀路径的模板文件，多个前缀用逗号分隔，相对于模板目录，不要前置/符号

3. **执行生成**：调用 `gen.Generate(cfg)` 方法执行代码生成。

生成器会自动：
- 加载所有变量文件
- 遍历模板目录中的所有模板文件
- 处理模板中的变量引用和子模板
- 将生成的文件保存到输出目录

## 模板特性

- 内置字符串处理函数 (`lcfirst`, `ucfirst`, `default`, `file`, `currentYear`, `dict`)
- 支持在输出路径中使用变量，例如 `__variableName__`。
- **支持子模板： 使用 `{{ include "path/to/sub_template.tpl" . }}` 在模板中包含其他模板文件。 子模板可以访问父模板中的变量。 限制最大嵌套层数为 2 层，防止循环引用。**

## 子模板使用说明

1.  **路径查找：** 子模板的路径优先作为相对于父模板的相对路径查找。如果指定了绝对路径，则直接使用绝对路径。

2.  **循环引用：** 限制子模板的嵌套层数为 2 层，超过则会报错。

3.  **变量传递：** 子模板可以访问父模板中定义的变量。

4.  **子模板命名:** 为了避免子模板被独立生成，请在子模板文件名或路径中包含 `__child__` 字符串。 包含 `__child__` 的模板文件将被自动跳过生成，并给出提示。  例如：`child__child__.tpl` 或者 `__child__/template.tpl`。

## 错误处理

生成器提供详细的错误报告，包括文件路径和行号。

## 贡献

欢迎提交问题和拉取请求。

## 许可证

本项目采用 MIT 许可证。
