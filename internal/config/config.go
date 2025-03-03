package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 主配置文件结构
type Config struct {
	Templates []Template `yaml:"templates"`
	Outputs   []Output   `yaml:"outputs"`    // 新增：输出文件配置
}

// Template 模板配置
type Template struct {
	Name         string   `yaml:"name"`           // 模板名称
	Path         string   `yaml:"path"`           // 模板文件路径
	Dependencies []string `yaml:"dependencies"`   // 依赖的配置文件列表
	Output       string   `yaml:"output,omitempty"` // 向后兼容：直接输出路径
}

// Output 输出配置
type Output struct {
	Path      string   `yaml:"path"`      // 输出文件路径
	Templates []string `yaml:"templates"`  // 使用的模板名称列表
	Order     []int    `yaml:"order"`     // 可选：模板合并顺序
}

// LoadConfig 从指定目录加载配置文件
func LoadConfig(workDir string) (*Config, error) {
	configPath := filepath.Join(workDir, "generator.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件 %s 不存在", configPath)
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := config.validate(workDir); err != nil {
		return nil, err
	}

	return &config, nil
}

// validate 验证配置是否有效
func (c *Config) validate(workDir string) error {
	if len(c.Templates) == 0 {
		return fmt.Errorf("没有配置模板文件")
	}

	// 创建模板名称映射，用于验证
	templateMap := make(map[string]struct{})
	for _, t := range c.Templates {
		if t.Name == "" {
			return fmt.Errorf("模板必须指定名称")
		}
		if _, exists := templateMap[t.Name]; exists {
			return fmt.Errorf("模板名称 %s 重复", t.Name)
		}
		templateMap[t.Name] = struct{}{}

		// 检查模板文件是否存在
		tplPath := filepath.Join(workDir, t.Path)
		if _, err := os.Stat(tplPath); err != nil {
			return fmt.Errorf("模板文件 %s 不存在", t.Path)
		}

		// 检查依赖的配置文件是否存在
		for _, dep := range t.Dependencies {
			depPath := filepath.Join(workDir, dep)
			if _, err := os.Stat(depPath); err != nil {
				return fmt.Errorf("依赖配置文件 %s 不存在", dep)
			}
		}
	}

	// 验证输出配置
	for _, o := range c.Outputs {
		if o.Path == "" {
			return fmt.Errorf("输出文件路径不能为空")
		}
		if len(o.Templates) == 0 {
			return fmt.Errorf("输出文件 %s 未指定模板", o.Path)
		}
		
		// 验证引用的模板是否存在
		for _, tplName := range o.Templates {
			if _, exists := templateMap[tplName]; !exists {
				return fmt.Errorf("输出文件 %s 引用的模板 %s 不存在", o.Path, tplName)
			}
		}

		// 验证顺序配置
		if len(o.Order) > 0 {
			if len(o.Order) != len(o.Templates) {
				return fmt.Errorf("输出文件 %s 的顺序配置数量与模板数量不匹配", o.Path)
			}
			// 验证顺序值是否有效
			orderMap := make(map[int]bool)
			for _, ord := range o.Order {
				if ord < 0 || ord >= len(o.Templates) {
					return fmt.Errorf("输出文件 %s 的顺序配置无效", o.Path)
				}
				if orderMap[ord] {
					return fmt.Errorf("输出文件 %s 的顺序配置重复", o.Path)
				}
				orderMap[ord] = true
			}
		}
	}

	return nil
}
