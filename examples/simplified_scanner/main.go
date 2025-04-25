package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/clh021/generator/internal/template"
	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
	"github.com/pkg/errors"
)

// TemplateScanner 定义模板扫描器接口
type TemplateScanner interface {
	// ScanTemplates 扫描模板目录，返回符合条件的模板文件列表
	ScanTemplates(templateDir string) ([]generator.TemplateFile, error)
}

// CustomScanner 自定义模板扫描器
type CustomScanner struct {
	// 可以添加自定义属性
	IncludePatterns []string
	ExcludePatterns []string
}

// ScanTemplates 扫描模板目录，返回符合条件的模板文件列表
func (s *CustomScanner) ScanTemplates(templateDir string) ([]generator.TemplateFile, error) {
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

		// 自定义逻辑：检查是否匹配排除模式
		if len(s.ExcludePatterns) > 0 {
			for _, pattern := range s.ExcludePatterns {
				if strings.Contains(relativePath, pattern) {
					log.Printf("跳过匹配排除模式的模板: %s (模式: %s)", path, pattern)
					return nil
				}
			}
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

// 简化的生成过程
func generateFiles(cfg *config.Config, scanner TemplateScanner) ([]generator.GeneratedFile, error) {
	var generatedFiles []generator.GeneratedFile

	// 确保路径是绝对路径
	templateDir, err := filepath.Abs(cfg.TemplateDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取模板目录的绝对路径: %s", cfg.TemplateDir)
	}

	variablesDir, err := filepath.Abs(cfg.VariablesDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取变量目录的绝对路径: %s", cfg.VariablesDir)
	}

	outputDir, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "无法获取输出目录的绝对路径: %s", cfg.OutputDir)
	}

	// 创建模板引擎
	engine := template.New(templateDir, variablesDir, outputDir)

	// 加载变量文件
	variableFiles, err := findVariableFiles(variablesDir, cfg.VariableFiles)
	if err != nil {
		return nil, errors.Wrap(err, "查找变量文件失败")
	}

	// 加载变量到引擎
	if err := engine.LoadVariables(variableFiles); err != nil {
		return nil, errors.Wrap(err, "加载变量到引擎失败")
	}

	// 获取变量
	variables := engine.GetVariables()

	// 使用自定义扫描器扫描模板文件
	templateFiles, err := scanner.ScanTemplates(templateDir)
	if err != nil {
		return nil, errors.Wrap(err, "扫描模板文件失败")
	}

	// 处理每个模板文件
	for _, templateFile := range templateFiles {
		// 处理输出路径
		outputPath := processOutputPath(templateFile.RelativePath, outputDir, variables)

		// 生成内容
		content, err := engine.GenerateContent(templateFile.Path, outputPath)
		if err != nil {
			return nil, errors.Wrapf(err, "生成内容失败 (%s)", templateFile.Path)
		}

		// 添加到生成的文件列表
		generatedFiles = append(generatedFiles, generator.GeneratedFile{
			TemplatePath: templateFile.Path,
			OutputPath:   outputPath,
			Content:      content,
		})
	}

	return generatedFiles, nil
}

// 查找变量文件
func findVariableFiles(variablesDir string, additionalFiles []string) ([]string, error) {
	var files []string

	// 加载目录中的文件
	if variablesDir != "" && dirExists(variablesDir) {
		yamlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yaml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yaml 变量文件失败")
		}
		ymlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yml 变量文件失败")
		}
		files = append(yamlFiles, ymlFiles...)
	}

	// 添加额外的文件
	for _, file := range additionalFiles {
		if fileExists(file) {
			files = append(files, file)
		}
	}

	return files, nil
}

// 处理输出路径
func processOutputPath(relativePath, outputDir string, variables map[string]interface{}) string {
	// 移除模板扩展名
	relPathWithoutExt := strings.TrimSuffix(relativePath, ".tpl")

	// 构建初始输出路径
	outputPath := filepath.Join(outputDir, relPathWithoutExt)

	// 处理路径中的变量
	re := regexp.MustCompile(`__([^_]+)__`)
	matches := re.FindAllStringSubmatch(outputPath, -1)

	result := outputPath
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		fullMatch := match[0] // __variable__
		varName := match[1]   // variable

		// 查找变量值
		value, ok := variables[varName]
		if !ok {
			log.Printf("警告: 在路径 %s 中找不到变量 %s，保留原始字符串", outputPath, varName)
			continue
		}

		// 将变量值转换为字符串并替换
		strValue, ok := value.(string)
		if !ok {
			log.Printf("警告: 在路径 %s 中变量 %s 的值不是字符串，保留原始字符串", outputPath, varName)
			continue
		}

		result = strings.Replace(result, fullMatch, strValue, -1)
	}

	// 规范化路径
	return filepath.Clean(result)
}

// 辅助函数：检查目录是否存在
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// 辅助函数：检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func main() {
	// 创建配置
	cfg := &config.Config{
		TemplateDir:   "/tmp/generator_test/templates",
		VariablesDir:  "/tmp/generator_test/variables",
		OutputDir:     "/tmp/generator_test/output/simplified_scanner",
		VariableFiles: []string{},
	}

	// 创建自定义扫描器
	customScanner := &CustomScanner{
		IncludePatterns: []string{"model", "controller"},
		ExcludePatterns: []string{"test", "mock"},
	}

	// 直接使用简化的生成过程
	files, err := generateFiles(cfg, customScanner)
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
