package generator

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// PathProcessor 定义路径处理器接口
type PathProcessor interface {
	// ProcessOutputPath 处理模板的输出路径
	ProcessOutputPath(templateFile TemplateFile, outputDir string, variables map[string]interface{}) (string, error)
}

// DefaultPathProcessor 默认的路径处理器实现
type DefaultPathProcessor struct{}

// NewDefaultPathProcessor 创建默认的路径处理器
func NewDefaultPathProcessor() *DefaultPathProcessor {
	return &DefaultPathProcessor{}
}

// ProcessOutputPath 处理模板的输出路径
func (p *DefaultPathProcessor) ProcessOutputPath(templateFile TemplateFile, outputDir string, variables map[string]interface{}) (string, error) {
	// 移除模板扩展名
	relPathWithoutExt := removeTemplateExtension(templateFile.RelativePath)

	// 构建初始输出路径
	outputPath := filepath.Join(outputDir, relPathWithoutExt)

	// 处理路径中的变量
	processedPath, err := p.processTemplatePath(outputPath, variables)
	if err != nil {
		return outputPath, err
	}

	return processedPath, nil
}

// processTemplatePath 处理路径中的变量引用
// 查找形如 __variable__ 的模板变量并尝试替换
// 如果变量不存在，则输出警告并保留原始字符串
func (p *DefaultPathProcessor) processTemplatePath(path string, variables map[string]interface{}) (string, error) {
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

// removeTemplateExtension 移除模板文件的扩展名
func removeTemplateExtension(path string) string {
	return strings.TrimSuffix(path, ".tpl")
}
