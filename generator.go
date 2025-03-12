package generate

import (
	"generate/internal/config"
	"generate/internal/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Generator 代码生成器
type Generator struct{}

// NewGenerator 创建新的生成器实例
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate 执行生成过程
func (g *Generator) Generate(cfg *config.Config) error {
	// 检查模板目录是否存在
	if _, err := os.Stat(cfg.TemplateDir); os.IsNotExist(err) {
		return errors.Wrapf(err, "模板目录不存在: %s", cfg.TemplateDir)
	}

	// 检查变量目录是否存在
	if _, err := os.Stat(cfg.VariablesDir); os.IsNotExist(err) {
		return errors.Wrapf(err, "变量目录不存在: %s", cfg.VariablesDir)
	}

	// 创建输出目录
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return errors.Wrap(err, "创建输出目录失败")
	}

	// 加载变量文件
	variableFiles, err := loadVariableFiles(cfg.VariablesDir)
	if err != nil {
		return errors.Wrap(err, "加载变量文件失败")
	}

	// 创建模板引擎
	engine := template.New(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)

	// 加载变量
	if err := engine.LoadVariables(variableFiles); err != nil {
		return errors.Wrap(err, "加载变量失败")
	}

	// 遍历模板目录
	err = filepath.Walk(cfg.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "遍历模板目录时出错")
		}

		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(cfg.TemplateDir, path)
		if err != nil {
			return errors.Wrap(err, "获取相对路径失败")
		}

		outputPath := filepath.Join(cfg.OutputDir, removeTemplateExtension(relativePath))

		if err := engine.Execute(path, outputPath); err != nil {
			return errors.Wrapf(err, "执行模板失败 (%s)", path)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "生成过程中出错")
	}

	return nil
}

// loadVariableFiles 加载变量目录下的所有 YAML 文件
func loadVariableFiles(variablesDir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(variablesDir, "*.yaml"))
	if err != nil {
		return nil, errors.Wrap(err, "查找变量文件失败")
	}
	return files, nil
}

// removeTemplateExtension 移除模板文件的扩展名
func removeTemplateExtension(path string) string {
	return strings.TrimSuffix(path, ".tpl")
}
