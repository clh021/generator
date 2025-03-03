package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 定义主配置结构
type Config struct {
    Templates []Template `yaml:"templates"`
}

// Template 定义模板配置
type Template struct {
    Path    string   `yaml:"path"`
    Configs []string `yaml:"configs"`
}

// LoadConfig 从指定目录加载配置
func LoadConfig(workDir string) (*Config, error) {
    configPath := filepath.Join(workDir, "config.yaml")

    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %w", err)
    }

    // 验证配置
    if err := cfg.validate(workDir); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// validate 验证配置的有效性
func (c *Config) validate(workDir string) error {
    if len(c.Templates) == 0 {
        return fmt.Errorf("配置文件中未找到模板定义")
    }

    for _, tmpl := range c.Templates {
        // 验证模板文件存在
        templatePath := filepath.Join(workDir, tmpl.Path)
        if _, err := os.Stat(templatePath); err != nil {
            return fmt.Errorf("模板文件不存在: %s", tmpl.Path)
        }

        // 验证配置文件存在
        for _, cfgPath := range tmpl.Configs {
            configPath := filepath.Join(workDir, cfgPath)
            if _, err := os.Stat(configPath); err != nil {
                return fmt.Errorf("配置文件不存在: %s", cfgPath)
            }
        }
    }

    return nil
}

// GetTemplateConfigs 获取指定模板的所有配置内容
func (c *Config) GetTemplateConfigs(workDir string, templatePath string) (map[string]interface{}, error) {
    var template *Template
    for _, t := range c.Templates {
        if t.Path == templatePath {
            template = &t
            break
        }
    }

    if template == nil {
        return nil, fmt.Errorf("未找到模板配置: %s", templatePath)
    }

    // 合并所有配置文件
    mergedConfig := make(map[string]interface{})
    for _, cfgPath := range template.Configs {
        configPath := filepath.Join(workDir, cfgPath)
        data, err := os.ReadFile(configPath)
        if err != nil {
            return nil, fmt.Errorf("读取配置文件失败 %s: %w", cfgPath, err)
        }

        var cfg map[string]interface{}
        if err := yaml.Unmarshal(data, &cfg); err != nil {
            return nil, fmt.Errorf("解析配置文件失败 %s: %w", cfgPath, err)
        }

        // 合并配置
        for k, v := range cfg {
            mergedConfig[k] = v
        }
    }

    return mergedConfig, nil
}
