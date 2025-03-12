package generate

import (
	"generate/internal/config"
	"generate/internal/template"
	"os"
	"path/filepath"
)

// Generator 代码生成器
type Generator struct{}

// NewGenerator 创建新的生成器实例
func NewGenerator() *Generator {
    return &Generator{}
}

// Generate 执行生成过程
func (g *Generator) Generate(cfg *config.Config) error {
    // 创建模板引擎
    engine := template.New(cfg.TemplateDir, cfg.ConfigDir, cfg.OutputDir)

    // 加载配置目录下的所有配置文件
    configFiles, err := loadConfigFiles(cfg.ConfigDir)
    if err != nil {
        return err
    }

    // 加载所有配置文件
    if err := engine.LoadConfig(configFiles); err != nil {
        return err
    }

    // 遍历模板目录
    return filepath.Walk(cfg.TemplateDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 跳过目录
        if info.IsDir() {
            return nil
        }

        // 计算相对路径
        relPath, err := filepath.Rel(cfg.TemplateDir, path)
        if err != nil {
            return err
        }

        // 构造输出文件路径
        outputPath := filepath.Join(cfg.OutputDir, relPath)
        // 移除模板文件扩展名
        outputPath = removeTemplateExtension(outputPath)

        // 确保输出目录存在
        if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
            return err
        }

        // 执行模板生成
        return engine.Execute(path, outputPath)
    })
}

// loadConfigFiles 加载配置目录下的所有 YAML 文件
func loadConfigFiles(configDir string) ([]string, error) {
    var configFiles []string
    err := filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
            configFiles = append(configFiles, path)
        }
        return nil
    })
    return configFiles, err
}

// removeTemplateExtension 移除模板文件的扩展名
func removeTemplateExtension(path string) string {
    ext := filepath.Ext(path)
    if ext == ".tpl" || ext == ".tmpl" {
        return path[:len(path)-len(ext)]
    }
    return path
}
