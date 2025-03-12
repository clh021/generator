package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// Engine 模板引擎
type Engine struct {
    templateDir string
    configDir   string
    outputDir   string
    vars        map[string]interface{}
}

// New 创建新的模板引擎
func New(templateDir, configDir, outputDir string) *Engine {
    return &Engine{
        templateDir: templateDir,
        configDir:   configDir,
        outputDir:   outputDir,
        vars:        make(map[string]interface{}),
    }
}

// LoadConfig 加载配置文件并合并变量
func (e *Engine) LoadConfig(configFiles []string) error {
    for _, path := range configFiles {
        data, err := os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
        }

        var vars map[string]interface{}
        if err := yaml.Unmarshal(data, &vars); err != nil {
            return fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
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

    // 创建模板
    tmpl, err := template.New(filepath.Base(tplPath)).
        Option("missingkey=error").
        Funcs(e.funcMap()).
        Parse(string(content))
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
    if err := tmpl.Execute(out, e.vars); err != nil {
        return fmt.Errorf("执行模板失败 (template: %s): %w", tplPath, err)
    }

    return nil
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
    }
}
