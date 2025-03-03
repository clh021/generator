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
}

// Template 模板配置
type Template struct {
	Path         string   `yaml:"path"`           // 模板文件路径
	Dependencies []string `yaml:"dependencies"`   // 依赖的配置文件列表
	Output      string   `yaml:"output"`         // 输出文件路径
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

	for _, t := range c.Templates {
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

	return nil
}
