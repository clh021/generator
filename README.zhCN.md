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

生成器也可以作为 Go 项目中的库使用。导入 `github.com/clh021/generator/pkg/generator` 包并调用 `GenerateFiles` 函数。

### 基本示例代码

```go
package main

import (
	"log"
	"os"
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
	files, err := gen.GenerateFiles(cfg)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	// 写入生成的文件
	for _, file := range files {
		// 创建输出目录
		outputDir := filepath.Dir(file.OutputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("创建输出目录失败: %v", err)
		}

		// 创建输出文件
		if err := os.WriteFile(file.OutputPath, []byte(file.Content), 0644); err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}

		log.Printf("已写入文件: %s", file.OutputPath)
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

3. **执行生成**：调用 `gen.GenerateFiles(cfg)` 方法执行代码生成，返回生成的文件列表。

4. **写入文件**：遍历生成的文件列表，创建目录并写入文件内容。

生成器会自动：
- 加载所有变量文件
- 遍历模板目录中的所有模板文件
- 处理模板中的变量引用和子模板
- 返回生成的文件列表

### 高级用法：自定义生成过程

生成器提供了多个接口，允许用户自定义生成过程的各个步骤：

1. **自定义模板过滤器**：实现 `TemplateFilter` 接口，自定义模板筛选逻辑

```go
// 自定义模板过滤器
type CustomTemplateFilter struct {
	// 继承默认过滤器
	*generator.DefaultTemplateFilter
	// 添加自定义属性
	AllowedExtensions []string
}

// 实现 ShouldInclude 方法
func (f *CustomTemplateFilter) ShouldInclude(path, relativePath string) (bool, string) {
	// 首先使用默认过滤器
	include, reason := f.DefaultTemplateFilter.ShouldInclude(path, relativePath)
	if !include {
		return false, reason
	}

	// 然后应用自定义过滤逻辑
	if len(f.AllowedExtensions) > 0 {
		allowed := false
		for _, ext := range f.AllowedExtensions {
			if filepath.Ext(path) == ext || filepath.Ext(path) == ext+".tpl" {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "不允许的扩展名"
		}
	}

	return true, ""
}

// 使用自定义过滤器
customFilter := &CustomTemplateFilter{
	DefaultTemplateFilter: generator.NewDefaultTemplateFilter(true, "", "", cfg.TemplateDir),
	AllowedExtensions:     []string{".go", ".md", ".json"},
}
gen := generator.NewGenerator().WithTemplateFilter(customFilter)
```

2. **自定义模板扫描器**：实现 `TemplateScanner` 接口，自定义模板扫描过程

```go
// 自定义模板扫描器
type CustomTemplateScanner struct {
	// 可以添加自定义属性
	IncludePatterns []string
}

// 实现 ScanTemplates 方法
func (s *CustomTemplateScanner) ScanTemplates(templateDir string, filter generator.TemplateFilter) ([]generator.TemplateFile, error) {
	var templateFiles []generator.TemplateFile

	// 自定义扫描逻辑...
	// ...

	return templateFiles, nil
}

// 使用自定义扫描器
customScanner := &CustomTemplateScanner{
	IncludePatterns: []string{"model", "controller"},
}
gen := generator.NewGenerator().WithTemplateScanner(customScanner)
```

3. **自定义内容生成器**：实现 `ContentGenerator` 接口，自定义内容生成过程

```go
// 自定义内容生成器
type CustomContentGenerator struct {
	// 可以添加自定义属性
	AddGeneratedComment bool
	CommentPrefix       string
}

// 实现 GenerateContent 方法
func (g *CustomContentGenerator) GenerateContent(templateFile generator.TemplateFile, outputPath string, engine *template.Engine) (string, error) {
	// 生成内容
	content, err := engine.GenerateContent(templateFile.Path, outputPath)
	if err != nil {
		return "", err
	}

	// 添加自定义处理逻辑...
	// ...

	return content, nil
}

// 使用自定义内容生成器
customContentGenerator := &CustomContentGenerator{
	AddGeneratedComment: true,
}
gen := generator.NewGenerator().WithContentGenerator(customContentGenerator)
```

4. **自定义变量加载器**：实现 `VariableLoader` 接口，自定义变量加载过程

```go
// 自定义变量加载器
type CustomVariableLoader struct {
	// 可以添加自定义属性
	DefaultLoader *generator.DefaultVariableLoader
	ExtraVars     map[string]interface{}
}

// 实现 LoadVariables 方法
func (l *CustomVariableLoader) LoadVariables(variablesDir string, additionalFiles []string) (map[string]interface{}, error) {
	// 使用默认加载器加载变量
	vars, err := l.DefaultLoader.LoadVariables(variablesDir, additionalFiles)
	if err != nil {
		return nil, err
	}

	// 添加额外的变量
	for k, v := range l.ExtraVars {
		vars[k] = v
	}

	return vars, nil
}

// 实现 FindVariableFiles 方法
func (l *CustomVariableLoader) FindVariableFiles(variablesDir string, additionalFiles []string) ([]string, error) {
	return l.DefaultLoader.FindVariableFiles(variablesDir, additionalFiles)
}

// 使用自定义变量加载器
defaultLoader := generator.NewDefaultVariableLoader(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)
customLoader := &CustomVariableLoader{
	DefaultLoader: defaultLoader,
	ExtraVars: map[string]interface{}{
		"version": "1.0.0",
		"author":  "Your Name",
	},
}
gen := generator.NewGenerator().WithVariableLoader(customLoader)
```

5. **自定义路径处理器**：实现 `PathProcessor` 接口，自定义输出路径处理过程

```go
// 自定义路径处理器
type CustomPathProcessor struct {
	// 可以添加自定义属性
	DefaultProcessor *generator.DefaultPathProcessor
	PathPrefix       string
}

// 实现 ProcessOutputPath 方法
func (p *CustomPathProcessor) ProcessOutputPath(templateFile generator.TemplateFile, outputDir string, variables map[string]interface{}) (string, error) {
	// 使用默认处理器处理路径
	path, err := p.DefaultProcessor.ProcessOutputPath(templateFile, outputDir, variables)
	if err != nil {
		return "", err
	}

	// 添加自定义前缀
	if p.PathPrefix != "" {
		path = filepath.Join(p.PathPrefix, path)
	}

	return path, nil
}

// 使用自定义路径处理器
customPathProcessor := &CustomPathProcessor{
	DefaultProcessor: generator.NewDefaultPathProcessor(),
	PathPrefix:       "generated",
}
gen := generator.NewGenerator().WithPathProcessor(customPathProcessor)
```

通过实现这些接口，您可以自定义生成过程的各个步骤，以满足特定的需求。

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
