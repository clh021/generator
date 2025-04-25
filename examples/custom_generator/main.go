package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
)

// 自定义模板过滤器
type CustomTemplateFilter struct {
	// 继承默认过滤器
	*generator.DefaultTemplateFilter
	// 添加自定义属性
	AllowedExtensions []string
}

// 实现 ShouldInclude 方法
func (f *CustomTemplateFilter) ShouldInclude(path, relativePath string) (bool, string) {
	// 首先使用默认过滤器
	include, reason := f.DefaultTemplateFilter.ShouldInclude(path, relativePath)
	if !include {
		return false, reason
	}

	// 然后应用自定义过滤逻辑
	if len(f.AllowedExtensions) > 0 {
		allowed := false
		for _, ext := range f.AllowedExtensions {
			if filepath.Ext(path) == ext || filepath.Ext(path) == ext+".tpl" {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, "不允许的扩展名"
		}
	}

	return true, ""
}

func main() {
	// 创建配置
	cfg := &config.Config{
		TemplateDir:   "./templates",
		VariablesDir:  "./variables",
		OutputDir:     "./output",
		VariableFiles: []string{},
	}

	// 创建自定义模板过滤器
	customFilter := &CustomTemplateFilter{
		DefaultTemplateFilter: generator.NewDefaultTemplateFilter(true, "", "", cfg.TemplateDir),
		AllowedExtensions:     []string{".go", ".md", ".json"},
	}

	// 创建生成器并设置自定义过滤器
	gen := generator.NewGenerator().WithTemplateFilter(customFilter)

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
