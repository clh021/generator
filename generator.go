package generate

import (
	"generate/internal/config"
	"generate/internal/template"
	"fmt"
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

	// 创建模板引擎
	engine := template.New(workDir)

	// 创建模板名称到路径的映射
	templatePaths := make(map[string]string)
	templateDeps := make(map[string][]string)
	for _, tpl := range cfg.Templates {
		templatePaths[tpl.Name] = tpl.Path
		templateDeps[tpl.Name] = tpl.Dependencies
	}

	// 处理每个输出配置
	for _, output := range cfg.Outputs {
		// 收集所有依赖的配置文件
		allDeps := make(map[string]struct{})
		for _, tplName := range output.Templates {
			for _, dep := range templateDeps[tplName] {
				allDeps[dep] = struct{}{}
			}
		}

		// 转换为切片
		deps := make([]string, 0, len(allDeps))
		for dep := range allDeps {
			deps = append(deps, dep)
		}

		// 加载所有依赖的配置文件
		if err := engine.LoadConfig(deps); err != nil {
			return err
		}

		// 收集模板路径
		tplPaths := make([]string, len(output.Templates))
		for i, tplName := range output.Templates {
			if path, ok := templatePaths[tplName]; ok {
				tplPaths[i] = path
			} else {
				return fmt.Errorf("模板 %s 未定义", tplName)
			}
		}

		// 执行模板生成
		if err := engine.ExecuteMultiple(tplPaths, output.Path, output.Order); err != nil {
			return err
		}
	}

	// 处理单个模板配置（向后兼容）
	for _, tpl := range cfg.Templates {
		if tpl.Output != "" {
			// 加载依赖的配置文件
			if err := engine.LoadConfig(tpl.Dependencies); err != nil {
				return err
			}

			// 执行模板生成
			if err := engine.Execute(tpl.Path, tpl.Output); err != nil {
				return err
			}
		}
	}

	return nil
}
