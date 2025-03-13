package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Engine struct {
	templateDir  string
	variablesDir string
	outputDir    string
	vars         map[string]interface{}
}

func New(templateDir, variablesDir, outputDir string) *Engine {
	return &Engine{
		templateDir:  templateDir,
		variablesDir: variablesDir,
		outputDir:    outputDir,
		vars:         make(map[string]interface{}),
	}
}

func (e *Engine) LoadVariables(variableFiles []string) error {
	for _, path := range variableFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取变量文件 %s 失败: %w", path, err)
		}

		var vars map[string]interface{}
		if err := yaml.Unmarshal(data, &vars); err != nil {
			return fmt.Errorf("解析变量文件 %s 失败: %w", path, err)
		}

		// 合并变量
		for k, v := range vars {
			e.vars[k] = v
		}
	}

	return nil
}

// Execute 执行单个模板生成
// Execute 执行单个模板生成
func (e *Engine) Execute(tplPath, outputPath string) error {
	// 读取模板文件
	content, err := os.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 检查是否允许未定义的变量
	allowUndefined, ok := e.vars["$config.allowUndefinedVariables"].(bool)

	// 创建模板
	tmpl := template.New(filepath.Base(tplPath)).Funcs(e.funcMap())

	// 根据配置设置 missingkey 选项
	if ok && allowUndefined {
		tmpl = tmpl.Option("missingkey=zero")
	} else {
		tmpl = tmpl.Option("missingkey=error")
	}

	// 解析模板
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 创建输出目录
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建输出文件
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer out.Close()

	// 执行模板
	log.Printf(" - 正在执行模板 %s", tplPath)
	log.Printf(" - 目标输出文件: %s", outputPath)
	log.Printf("传递给模板的变量:")
	for k, v := range e.vars {
		log.Printf("  %s = %v", k, v)
	}

	if err := tmpl.Execute(out, e.vars); err != nil {
		return fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
	}

	return nil
}

// lcfirst converts the first character of a string to lowercase.
// If the string is empty, it returns an empty string.
//
// Example usage in templates:
//
//	{{ lcfirst "UserName" }} -> "userName"
func lcfirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// ucfirst converts the first character of a string to uppercase.
// If the string is empty, it returns an empty string.
//
// Example usage in templates:
//
//	{{ ucfirst "userName" }} -> "UserName"
func ucfirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// funcMap 返回模板中可用的函数映射
func (e *Engine) funcMap() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"currentYear": func() int {
			return time.Now().Year()
		},
		"file": func(filePath string) (string, error) {
			content, err := os.ReadFile(filepath.Join(e.templateDir, filePath))
			if err != nil {
				return "", fmt.Errorf("读取文件失败 %s: %w", filePath, err)
			}
			return string(content), nil
		},
		"default": func(value, defaultValue interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"lcfirst": lcfirst,
		"ucfirst": ucfirst,
	}
}
