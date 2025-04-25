package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Engine struct {
	templateDir     string
	variablesDir    string
	outputDir       string
	vars            map[string]interface{}
	loadedTemplates map[string]*template.Template
}

func New(templateDir, variablesDir, outputDir string) *Engine {
	return &Engine{
		templateDir:     templateDir,
		variablesDir:    variablesDir,
		outputDir:       outputDir,
		vars:            make(map[string]interface{}),
		loadedTemplates: make(map[string]*template.Template),
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
			// 如果是嵌套的映射，需要特殊处理
			if existingVal, ok := e.vars[k]; ok {
				if existingMap, ok := existingVal.(map[string]interface{}); ok {
					if newMap, ok := v.(map[string]interface{}); ok {
						// 合并嵌套映射
						for nk, nv := range newMap {
							existingMap[nk] = nv
						}
						continue
					}
				}
			}
			// 对于非嵌套映射或新键，直接赋值
			e.vars[k] = v
		}
	}

	return nil
}

// GenerateContent 生成模板内容但不写入文件
func (e *Engine) GenerateContent(tplPath, outputPath string) (string, error) {
	// 读取模板文件
	content, err := os.ReadFile(tplPath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
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
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	// 传递模板路径
	varsWithTemplatePath := make(map[string]interface{})
	for k, v := range e.vars {
		varsWithTemplatePath[k] = v
	}
	varsWithTemplatePath["__current_template_path"] = tplPath // 添加当前模板路径

	// 执行模板到字符串
	var result strings.Builder
	if err := tmpl.Execute(&result, varsWithTemplatePath); err != nil {
		// 执行模板
		log.Printf(" - 正在执行模板 %s", tplPath)
		log.Printf(" - 目标输出文件: %s", outputPath)
		log.Printf("传递给模板的变量:")
		for k, v := range varsWithTemplatePath {
			log.Printf("  %s = %v", k, v)
		}

		return "", fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
	}

	return result.String(), nil
}
