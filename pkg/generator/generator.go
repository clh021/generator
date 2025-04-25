package generator

import (
	"log"
	"path/filepath"

	"github.com/clh021/generator/internal/template"
	"github.com/clh021/generator/pkg/config"

	"github.com/pkg/errors"
)

// Generator 代码生成器
type Generator struct {
	variables        map[string]interface{}
	templateScanner  TemplateScanner
	variableLoader   VariableLoader
	pathProcessor    PathProcessor
	contentGenerator ContentGenerator
	templateFilter   TemplateFilter
}

// NewGenerator 创建新的生成器实例
func NewGenerator() *Generator {
	return &Generator{
		variables:        make(map[string]interface{}),
		templateScanner:  NewDefaultTemplateScanner(),
		variableLoader:   nil, // 将在 GenerateFiles 中初始化
		pathProcessor:    NewDefaultPathProcessor(),
		contentGenerator: NewDefaultContentGenerator(),
		templateFilter:   nil, // 将在 GenerateFiles 中初始化
	}
}

// WithTemplateScanner 设置模板扫描器
func (g *Generator) WithTemplateScanner(scanner TemplateScanner) *Generator {
	g.templateScanner = scanner
	return g
}

// WithVariableLoader 设置变量加载器
func (g *Generator) WithVariableLoader(loader VariableLoader) *Generator {
	g.variableLoader = loader
	return g
}

// WithPathProcessor 设置路径处理器
func (g *Generator) WithPathProcessor(processor PathProcessor) *Generator {
	g.pathProcessor = processor
	return g
}

// WithContentGenerator 设置内容生成器
func (g *Generator) WithContentGenerator(generator ContentGenerator) *Generator {
	g.contentGenerator = generator
	return g
}

// WithTemplateFilter 设置模板过滤器
func (g *Generator) WithTemplateFilter(filter TemplateFilter) *Generator {
	g.templateFilter = filter
	return g
}

// GenerateFiles 执行生成过程但不写入文件，而是返回生成的文件列表
func (g *Generator) GenerateFiles(cfg *config.Config) ([]GeneratedFile, error) {
	var generatedFiles []GeneratedFile

	// 确保所有路径都是绝对路径
	var err error
	cfg.TemplateDir, err = filepath.Abs(cfg.TemplateDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取模板目录的绝对路径: %s", cfg.TemplateDir)
	}
	cfg.VariablesDir, err = filepath.Abs(cfg.VariablesDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取变量目录的绝对路径: %s", cfg.VariablesDir)
	}
	cfg.OutputDir, err = filepath.Abs(cfg.OutputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取输出目录的绝对路径: %s", cfg.OutputDir)
	}

	// 初始化变量加载器（如果未设置）
	if g.variableLoader == nil {
		g.variableLoader = NewDefaultVariableLoader(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)
	}

	// 初始化模板过滤器（如果未设置）
	if g.templateFilter == nil {
		g.templateFilter = NewDefaultTemplateFilter(true, cfg.SkipTemplateSuffixes, cfg.SkipTemplatePrefixes, cfg.TemplateDir)
	}

	// 加载变量
	variables, err := g.variableLoader.LoadVariables(cfg.VariablesDir, cfg.VariableFiles)
	if err != nil {
		return nil, errors.Wrap(err, "加载变量失败")
	}
	g.variables = variables

	// 创建模板引擎
	engine := template.New(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)

	// 加载变量文件
	variableFiles, err := g.variableLoader.FindVariableFiles(cfg.VariablesDir, cfg.VariableFiles)
	if err != nil {
		return nil, errors.Wrap(err, "查找变量文件失败")
	}

	// 加载变量到引擎
	if err := engine.LoadVariables(variableFiles); err != nil {
		return nil, errors.Wrap(err, "加载变量到引擎失败")
	}

	// 扫描模板
	templateFiles, err := g.templateScanner.ScanTemplates(cfg.TemplateDir, g.templateFilter)
	if err != nil {
		return nil, errors.Wrap(err, "扫描模板失败")
	}

	// 处理每个模板文件
	for _, templateFile := range templateFiles {
		// 处理输出路径
		outputPath, err := g.pathProcessor.ProcessOutputPath(templateFile, cfg.OutputDir, g.variables)
		if err != nil {
			log.Printf("警告: 处理输出路径失败: %v, 使用默认路径", err)
		}

		// 生成文件内容
		content, err := g.contentGenerator.GenerateContent(templateFile, outputPath, engine)
		if err != nil {
			return nil, errors.Wrapf(err, "生成内容失败 (%s)", templateFile.Path)
		}

		// 添加到生成的文件列表
		generatedFiles = append(generatedFiles, GeneratedFile{
			TemplatePath: templateFile.Path,
			OutputPath:   outputPath,
			Content:      content,
		})
	}

	return generatedFiles, nil
}

// 以下函数已移至各自的文件中，这里保留注释以便于理解代码结构
// loadVariableFiles -> variables.go: DefaultVariableLoader.FindVariableFiles
// removeTemplateExtension -> path.go
// processTemplatePath -> path.go: DefaultPathProcessor.processTemplatePath
