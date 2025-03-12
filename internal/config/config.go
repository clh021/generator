package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 主配置文件结构
type Config struct {
    TemplateDir string `yaml:"templateDir"`
    ConfigDir   string `yaml:"configDir"`
    OutputDir   string `yaml:"outputDir"`
}

// LoadConfig 从指定目录加载配置文件
func LoadConfig(workDir string) (*Config, error) {
    configPath := filepath.Join(workDir, "generator.yaml")

    data, err := os.ReadFile(configPath)
    if err != nil {
        if os.IsNotExist(err) {
            // 如果找不到配置文件，返回默认配置
            return &Config{
                TemplateDir: "templates",
                ConfigDir:   "config",
                OutputDir:   "output",
            }, nil
        }
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %w", err)
    }

    // 验证配置
    if err := config.validate(); err != nil {
        return nil, err
    }

    return &config, nil
}

// validate 验证配置是否有效
func (c *Config) validate() error {
    if c.TemplateDir == "" {
        return fmt.Errorf("模板目录不能为空")
    }
    if c.ConfigDir == "" {
        return fmt.Errorf("配置目录不能为空")
    }
    if c.OutputDir == "" {
        return fmt.Errorf("输出目录不能为空")
    }
    return nil
}
