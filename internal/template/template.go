package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Engine 模板引擎
type Engine struct {
	workDir string
	vars    map[string]interface{}
	templates map[string]*template.Template
}

// New 创建新的模板引擎
func New(workDir string) *Engine {
	return &Engine{
		workDir: workDir,
		vars:    make(map[string]interface{}),
		templates: make(map[string]*template.Template),
	}
}

// LoadConfig 加载配置文件并合并变量
func (e *Engine) LoadConfig(configPaths []string) error {
	for _, path := range configPaths {
		fullPath := filepath.Join(e.workDir, path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
		}

		// 只处理YAML文件
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			continue
		}

		var vars map[string]interface{}
		if err := yaml.Unmarshal(data, &vars); err != nil {
			return fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
		}

		// 合并变量
		for k, v := range vars {
			if existing, ok := e.vars[k]; ok {
				if existingMap, ok := existing.(map[string]interface{}); ok {
					if newMap, ok := v.(map[string]interface{}); ok {
						// 合并嵌套的map
						for subk, subv := range newMap {
							existingMap[subk] = subv
						}
						continue
					}
				}
			}
			e.vars[k] = v
		}
	}

	return nil
}

// Execute 执行单个模板生成
func (e *Engine) Execute(tplPath, outputPath string) error {
	tmpl, err := e.loadTemplate(filepath.Base(tplPath), tplPath)
	if err != nil {
		return err
	}

	// 创建输出目录
	fullOutputPath := filepath.Join(e.workDir, outputPath)
	outputDir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建输出文件
	out, err := os.Create(fullOutputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer out.Close()

	// 执行模板
	if err := tmpl.Execute(out, e.vars); err != nil {
		return fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
	}

	return nil
}

// ExecuteMultiple 执行多个模板并合并结果
func (e *Engine) ExecuteMultiple(templates []string, outputPath string, order []int) error {
	// 创建输出目录
	fullOutputPath := filepath.Join(e.workDir, outputPath)
	outputDir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 创建输出文件
	out, err := os.Create(fullOutputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer out.Close()

	// 如果没有指定顺序，按照模板列表顺序执行
	if order == nil {
		order = make([]int, len(templates))
		for i := range order {
			order[i] = i
		}
	}

	// 按顺序执行模板
	for _, idx := range order {
		tplPath := templates[idx]
		tmpl, err := e.loadTemplate(filepath.Base(tplPath), tplPath)
		if err != nil {
			return err
		}

		// 执行模板并写入文件
		if err := tmpl.Execute(out, e.vars); err != nil {
			return fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
		}
	}

	return nil
}

func (e *Engine) loadTemplate(name, path string) (*template.Template, error) {
	// 如果模板已经加载过，直接返回
	if tmpl, ok := e.templates[path]; ok {
		return tmpl, nil
	}

	// 读取模板文件
	content, err := os.ReadFile(filepath.Join(e.workDir, path))
	if err != nil {
		return nil, fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 创建模板
	tmpl := template.New(name).
		Option("missingkey=error"). // 确保所有变量都已定义
		Funcs(template.FuncMap{
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
				content, err := os.ReadFile(filepath.Join(filepath.Dir(path), filePath))
				if err != nil {
					return "", fmt.Errorf("读取文件失败 %s: %w", filePath, err)
				}
				return string(content), nil
			},
		})

	// 解析模板
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("解析模板失败: %w", err)
	}

	// 缓存模板
	e.templates[path] = tmpl
	return tmpl, nil
}
