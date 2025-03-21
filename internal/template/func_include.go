package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// <!-- 相对路径引用 -->
// Included template (relative): {{ include "child.txt.tpl" . }}
// <!-- 绝对路径引用 （假设工作目录为/path/to/project）-->
// Included template (absolute): {{ include "/path/to/project/.gen_templates/child.txt.tpl" . }}

// createIncludeTemplateFunc 创建一个闭包函数，用于处理递归调用时的深度限制和循环引用检测
func (e *Engine) createIncludeTemplateFunc(depth int, templateStack []string) func(tplName string, data interface{}) (string, error) {
	return func(tplName string, data interface{}) (string, error) {
		// 循环引用检测，限制嵌套层数为2
		if depth > 2 {
			return "", fmt.Errorf("超过最大嵌套层数 (2)，禁止循环引用")
		}

		// 获取当前模板的完整路径
		currentTplPath, ok := data.(map[string]interface{})["__current_template_path"].(string)
		if !ok {
			return "", fmt.Errorf("无法获取当前模板路径，请确保在顶级模板执行时传递了 __current_template_path 变量")
		}
		currentDir := filepath.Dir(currentTplPath)

		var tplPath string
		// 如果 tplName 是绝对路径，直接使用
		if filepath.IsAbs(tplName) {
			tplPath = tplName
		} else {
			// 否则，作为相对于父模板的相对路径查找
			tplPath = filepath.Join(currentDir, tplName)
		}
		tplPath = filepath.Clean(tplPath) // 清理路径

		// 循环引用检测
		for _, path := range templateStack {
			if path == tplPath {
				return "", fmt.Errorf("发现循环引用: %s -> %s", strings.Join(templateStack, " -> "), tplPath)
			}
		}

		// 检查模板是否已经加载
		tmpl, ok := e.loadedTemplates[tplPath]
		if !ok {
			// 读取模板文件
			content, err := os.ReadFile(tplPath)
			if err != nil {
				return "", fmt.Errorf("读取子模板文件 %s 失败: %w", tplPath, err)
			}

			// 创建模板
			tmpl = template.New(filepath.Base(tplName)).Funcs(e.funcMap())

			// 解析模板
			tmpl, err = tmpl.Parse(string(content))
			if err != nil {
				return "", fmt.Errorf("解析子模板 %s 失败: %w", tplPath, err)
			}

			// 缓存模板
			e.loadedTemplates[tplPath] = tmpl
		}

		// 创建新的模板栈
		newTemplateStack := append(templateStack, tplPath)

		// 创建新的 include 函数，增加深度
		newIncludeFunc := e.createIncludeTemplateFunc(depth+1, newTemplateStack)

		// 创建包含新的 include 函数的 FuncMap
		funcMapWithNewInclude := template.FuncMap{
			"include": newIncludeFunc,
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

		// 克隆模板并添加新的 FuncMap
		clonedTmpl, err := tmpl.Clone()
		if err != nil {
			return "", fmt.Errorf("克隆模板失败: %w", err)
		}
		clonedTmpl = clonedTmpl.Funcs(funcMapWithNewInclude)

		// 执行模板
		var buf strings.Builder
		err = clonedTmpl.Execute(&buf, data)
		if err != nil {
			return "", fmt.Errorf("执行子模板 %s 失败: %w", tplPath, err)
		}

		return buf.String(), nil
	}
}
