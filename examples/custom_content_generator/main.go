package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/clh021/generator/internal/template"
	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
	"github.com/pkg/errors"
)

// ContentGenerator 定义内容生成器接口
// 这个接口可以直接在项目中定义，不需要依赖 generator 包
type ContentGenerator interface {
	// GenerateContent 生成单个文件的内容
	GenerateContent(templatePath, relativePath, outputPath string, engine *template.Engine) (string, error)
}

// CustomContentGenerator 自定义内容生成器
type CustomContentGenerator struct {
	// 可以添加自定义属性
	AddGeneratedComment bool
	CommentPrefix       string
}

// GenerateContent 生成内容并添加自定义注释
func (g *CustomContentGenerator) GenerateContent(templatePath, relativePath, outputPath string, engine *template.Engine) (string, error) {
	log.Printf("正在处理模板: %s", templatePath)
	log.Printf("相对路径: %s", relativePath)
	log.Printf("目标输出路径: %s", outputPath)

	// 生成内容
	content, err := engine.GenerateContent(templatePath, outputPath)
	if err != nil {
		return "", errors.Wrapf(err, "执行模板失败 (%s)", templatePath)
	}

	// 添加生成注释
	if g.AddGeneratedComment {
		ext := filepath.Ext(outputPath)
		commentPrefix := g.CommentPrefix

		// 根据文件类型选择注释前缀
		if commentPrefix == "" {
			switch ext {
			case ".go", ".java", ".js", ".ts", ".c", ".cpp", ".cs":
				commentPrefix = "//"
			case ".py", ".rb", ".sh":
				commentPrefix = "#"
			case ".html", ".xml":
				commentPrefix = "<!--"
			default:
				commentPrefix = "#"
			}
		}

		// 生成注释
		comment := commentPrefix + " 此文件由生成器自动生成于 " + time.Now().Format("2006-01-02 15:04:05") + "\n"
		comment += commentPrefix + " 模板文件: " + templatePath + "\n"
		comment += commentPrefix + " 相对路径: " + relativePath + "\n\n"

		// 添加注释到内容开头
		if ext == ".html" || ext == ".xml" {
			comment += " -->\n"
		}

		content = comment + content
	}

	return content, nil
}

// 简化的生成过程
func generateFiles(cfg *config.Config, contentGenerator ContentGenerator) ([]generator.GeneratedFile, error) {
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

	// 扫描模板文件
	templateFiles, err := scanTemplateFiles(templateDir, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "扫描模板文件失败")
	}

	// 处理每个模板文件
	for _, templateFile := range templateFiles {
		// 处理输出路径
		outputPath := processOutputPath(templateFile.RelativePath, outputDir, variables)

		// 生成内容
		content, err := contentGenerator.GenerateContent(
			templateFile.Path,
			templateFile.RelativePath,
			outputPath,
			engine,
		)
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

// 扫描模板文件
func scanTemplateFiles(templateDir string, cfg *config.Config) ([]generator.TemplateFile, error) {
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

		// 检查是否为子模板
		if strings.Contains(path, "__child__") {
			log.Printf("跳过子模板: %s", path)
			return nil
		}

		// 检查是否应该跳过此模板（基于后缀）
		if cfg.SkipTemplateSuffixes != "" {
			suffixes := strings.Split(cfg.SkipTemplateSuffixes, ",")
			for _, suffix := range suffixes {
				suffix = strings.TrimSpace(suffix)
				if suffix != "" && strings.HasSuffix(path, suffix) {
					log.Printf("跳过后缀匹配的模板: %s (后缀: %s)", path, suffix)
					return nil
				}
			}
		}

		// 检查是否应该跳过此模板（基于前缀）
		if cfg.SkipTemplatePrefixes != "" {
			prefixes := strings.Split(cfg.SkipTemplatePrefixes, ",")
			for _, prefix := range prefixes {
				prefix = strings.TrimSpace(prefix)
				if prefix != "" {
					if strings.HasPrefix(relativePath, prefix) {
						log.Printf("跳过前缀匹配的模板: %s (前缀: %s)", path, prefix)
						return nil
					}
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
		OutputDir:     "/tmp/generator_test/output/custom_content_generator",
		VariableFiles: []string{},
	}

	// 创建自定义内容生成器
	customContentGenerator := &CustomContentGenerator{
		AddGeneratedComment: true,
	}

	// 直接使用简化的生成过程
	files, err := generateFiles(cfg, customContentGenerator)
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
