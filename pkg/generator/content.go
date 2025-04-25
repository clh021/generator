package generator

import (
	"log"

	"github.com/clh021/generator/internal/template"
	"github.com/pkg/errors"
)

// ContentGenerator 定义内容生成器接口
type ContentGenerator interface {
	// GenerateContent 生成单个文件的内容
	GenerateContent(templateFile TemplateFile, outputPath string, engine interface{}) (string, error)
}

// DefaultContentGenerator 默认的内容生成器实现
type DefaultContentGenerator struct{}

// NewDefaultContentGenerator 创建默认的内容生成器
func NewDefaultContentGenerator() *DefaultContentGenerator {
	return &DefaultContentGenerator{}
}

// GenerateContent 生成单个文件的内容
func (g *DefaultContentGenerator) GenerateContent(templateFile TemplateFile, outputPath string, engine interface{}) (string, error) {
	log.Printf("正在处理模板: %s", templateFile.Path)
	log.Printf("目标输出路径: %s", outputPath)

	// 类型断言
	templateEngine, ok := engine.(*template.Engine)
	if !ok {
		return "", errors.New("引擎类型不是 *template.Engine")
	}

	// 生成内容
	content, err := templateEngine.GenerateContent(templateFile.Path, outputPath)
	if err != nil {
		return "", errors.Wrapf(err, "执行模板失败 (%s)", templateFile.Path)
	}

	return content, nil
}
