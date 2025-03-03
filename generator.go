package generate

import (
	"generate/internal/config"
	"generate/internal/template"
)

// Generator 代码生成器
type Generator struct{}

// NewGenerator 创建新的生成器实例
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate 执行生成过程
func (g *Generator) Generate(workDir string) error {
	// 加载配置
	cfg, err := config.LoadConfig(workDir)
	if err != nil {
		return err
	}

	// 处理每个模板
	for _, tpl := range cfg.Templates {
		engine := template.New(workDir)
		
		// 加载依赖的配置文件
		if err := engine.LoadConfig(tpl.Dependencies); err != nil {
			return err
		}

		// 执行模板生成
		if err := engine.Execute(tpl.Path, tpl.Output); err != nil {
			return err
		}
	}

	return nil
}
