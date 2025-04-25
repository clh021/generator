package generator

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// TemplateScanner 定义模板扫描器接口
type TemplateScanner interface {
	// ScanTemplates 扫描模板目录，返回符合条件的模板文件列表
	ScanTemplates(templateDir string, filter TemplateFilter) ([]TemplateFile, error)
}

// DefaultTemplateScanner 默认的模板扫描器实现
type DefaultTemplateScanner struct{}

// NewDefaultTemplateScanner 创建默认的模板扫描器
func NewDefaultTemplateScanner() *DefaultTemplateScanner {
	return &DefaultTemplateScanner{}
}

// ScanTemplates 扫描模板目录，返回符合条件的模板文件列表
func (s *DefaultTemplateScanner) ScanTemplates(templateDir string, filter TemplateFilter) ([]TemplateFile, error) {
	var templateFiles []TemplateFile

	// 检查模板目录是否存在
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "模板目录不存在: %s", templateDir)
	}

	// 遍历模板目录
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "遍历模板目录时出错")
		}

		if info.IsDir() {
			return nil
		}

		// 获取相对路径
		relativePath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return errors.Wrap(err, "获取相对路径失败")
		}

		// 检查是否应该包含此模板
		include, reason := filter.ShouldInclude(path, relativePath)
		if !include {
			// 可以添加日志记录，但为了保持函数纯净，这里不添加日志
			return nil
		}

		templateFiles = append(templateFiles, TemplateFile{
			Path:         path,
			RelativePath: relativePath,
		})

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "扫描模板目录时出错")
	}

	return templateFiles, nil
}
