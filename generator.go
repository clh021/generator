package generate

import (
	"generate/internal/config"
	"generate/internal/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Generator 代码生成器
type Generator struct {
	variables map[string]interface{}
}

// NewGenerator 创建新的生成器实例
func NewGenerator() *Generator {
	return &Generator{
		variables: make(map[string]interface{}),
	}
}

// Generate 执行生成过程
func (g *Generator) Generate(cfg *config.Config) error {
	// 确保所有路径都是绝对路径
	var err error
	cfg.TemplateDir, err = filepath.Abs(cfg.TemplateDir)
	if err != nil {
		return errors.Wrapf(err, "无法获取模板目录的绝对路径: %s", cfg.TemplateDir)
	}
	cfg.VariablesDir, err = filepath.Abs(cfg.VariablesDir)
	if err != nil {
		return errors.Wrapf(err, "无法获取变量目录的绝对路径: %s", cfg.VariablesDir)
	}
	cfg.OutputDir, err = filepath.Abs(cfg.OutputDir)
	if err != nil {
		return errors.Wrapf(err, "无法获取输出目录的绝对路径: %s", cfg.OutputDir)
	}

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

	// 获取变量供路径处理使用
	g.variables = engine.GetVariables()

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

		// 移除模板扩展名
		relPathWithoutExt := removeTemplateExtension(relativePath)

		// 构建初始输出路径
		outputPath := filepath.Join(cfg.OutputDir, relPathWithoutExt)

		// 处理路径中的变量
		processedPath, err := g.processTemplatePath(outputPath, g.variables)
		if err != nil {
			log.Printf("警告: 处理目标路径时遇到错误: %v. 使用原始路径", err)
			// 打印警告，但不终断流程
		} else {
			outputPath = processedPath
		}

		log.Printf("正在处理模板: %s", path)
		log.Printf("目标输出路径: %s", outputPath)

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
	var files []string
	yamlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yaml"))
	if err != nil {
		return nil, errors.Wrap(err, "查找 *.yaml 变量文件失败")
	}
	ymlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yml"))
	if err != nil {
		return nil, errors.Wrap(err, "查找 *.yml 变量文件失败")
	}
	files = append(yamlFiles, ymlFiles...)
	return files, nil
}

// removeTemplateExtension 移除模板文件的扩展名
func removeTemplateExtension(path string) string {
	return strings.TrimSuffix(path, ".tpl")
}

// processTemplatePath 处理路径中的变量引用
// 查找形如 __variable__ 的模板变量并尝试替换
// 如果变量不存在，则输出警告并保留原始字符串
func (g *Generator) processTemplatePath(path string, variables map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`__([^_]+)__`)
	matches := re.FindAllStringSubmatch(path, -1)

	result := path
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		fullMatch := match[0] // __variable__
		varName := match[1]   // variable

		// 查找变量值
		value, ok := variables[varName]
		if !ok {
			log.Printf("警告: 在路径 %s 中找不到变量 %s，保留原始字符串", path, varName)
			continue
		}

		// 将变量值转换为字符串并替换
		strValue, ok := value.(string)
		if !ok {
			log.Printf("警告: 在路径 %s 中变量 %s 的值不是字符串，保留原始字符串", path, varName)
			continue
		}

		result = strings.Replace(result, fullMatch, strValue, -1)
	}

	// 规范化路径
	return filepath.Clean(result), nil
}
