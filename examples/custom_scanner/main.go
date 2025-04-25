package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
	"github.com/pkg/errors"
)

// 自定义模板扫描器
type CustomTemplateScanner struct {
	// 可以添加自定义属性
	IncludePatterns []string
}

// 实现 ScanTemplates 方法
func (s *CustomTemplateScanner) ScanTemplates(templateDir string, filter generator.TemplateFilter) ([]generator.TemplateFile, error) {
	var templateFiles []generator.TemplateFile

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

		// 自定义逻辑：检查是否匹配包含模式
		if len(s.IncludePatterns) > 0 {
			matched := false
			for _, pattern := range s.IncludePatterns {
				if strings.Contains(relativePath, pattern) {
					matched = true
					break
				}
			}
			if !matched {
				log.Printf("跳过不匹配的模板: %s", path)
				return nil
			}
		}

		// 使用提供的过滤器
		include, reason := filter.ShouldInclude(path, relativePath)
		if !include {
			log.Printf("跳过模板: %s (%s)", path, reason)
			return nil
		}

		templateFiles = append(templateFiles, generator.TemplateFile{
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

func main() {
	// 创建配置
	cfg := &config.Config{
		TemplateDir:   "./templates",
		VariablesDir:  "./variables",
		OutputDir:     "./output",
		VariableFiles: []string{},
	}

	// 创建自定义模板扫描器
	customScanner := &CustomTemplateScanner{
		IncludePatterns: []string{"model", "controller"},
	}

	// 创建生成器并设置自定义扫描器
	gen := generator.NewGenerator().WithTemplateScanner(customScanner)

	// 执行生成
	files, err := gen.GenerateFiles(cfg)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	// 写入生成的文件
	for _, file := range files {
		// 创建输出目录
		outputDir := filepath.Dir(file.OutputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("创建输出目录失败: %v", err)
		}

		// 创建输出文件
		if err := os.WriteFile(file.OutputPath, []byte(file.Content), 0644); err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}

		log.Printf("已写入文件: %s", file.OutputPath)
	}

	log.Println("生成完成")
}
