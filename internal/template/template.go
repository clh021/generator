package template

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
			e.vars[k] = v
		}
	}

	return nil
}

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

	// 传递模板路径
	varsWithTemplatePath := make(map[string]interface{})
	for k, v := range e.vars {
		varsWithTemplatePath[k] = v
	}
	varsWithTemplatePath["__current_template_path"] = tplPath // 添加当前模板路径

	if err := tmpl.Execute(out, varsWithTemplatePath); err != nil { // 使用新的变量
		// 执行模板
		log.Printf(" - 正在执行模板 %s", tplPath)
		log.Printf(" - 目标输出文件: %s", outputPath)
		log.Printf("传递给模板的变量:")
		for k, v := range varsWithTemplatePath {
			log.Printf("  %s = %v", k, v)
		}

		return fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
	}

	return nil
}
