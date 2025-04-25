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
	"gopkg.in/yaml.v3"
)

// VariableLoader 定义变量加载器接口
type VariableLoader interface {
	// LoadVariables 加载变量并返回变量映射
	LoadVariables() (map[string]interface{}, error)
}

// CustomVariableLoader 自定义变量加载器
type CustomVariableLoader struct {
	// 配置
	Config *config.Config
	// 额外的变量
	ExtraVariables map[string]interface{}
	// 环境变量前缀
	EnvPrefix string
}

// LoadVariables 加载变量并返回变量映射
func (l *CustomVariableLoader) LoadVariables() (map[string]interface{}, error) {
	variables := make(map[string]interface{})

	// 1. 首先加载文件中的变量
	fileVariables, err := l.loadFromFiles()
	if err != nil {
		return nil, err
	}
	for k, v := range fileVariables {
		variables[k] = v
	}

	// 2. 然后加载环境变量
	envVariables := l.loadFromEnvironment()
	for k, v := range envVariables {
		variables[k] = v
	}

	// 3. 最后加载额外的变量（优先级最高）
	for k, v := range l.ExtraVariables {
		variables[k] = v
	}

	// 添加一些内置变量
	variables["currentYear"] = time.Now().Year()
	variables["generatorVersion"] = "1.0.0"

	return variables, nil
}

// 从文件加载变量
func (l *CustomVariableLoader) loadFromFiles() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 查找变量文件
	variableFiles, err := l.findVariableFiles()
	if err != nil {
		return nil, err
	}

	// 加载每个文件
	for _, file := range variableFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "读取变量文件失败: %s", file)
		}

		var fileVariables map[string]interface{}
		if err := yaml.Unmarshal(data, &fileVariables); err != nil {
			return nil, errors.Wrapf(err, "解析变量文件失败: %s", file)
		}

		// 合并变量
		for k, v := range fileVariables {
			result[k] = v
		}
	}

	return result, nil
}

// 查找变量文件
func (l *CustomVariableLoader) findVariableFiles() ([]string, error) {
	var files []string

	// 加载目录中的文件
	if l.Config.VariablesDir != "" && dirExists(l.Config.VariablesDir) {
		yamlFiles, err := filepath.Glob(filepath.Join(l.Config.VariablesDir, "*.yaml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yaml 变量文件失败")
		}
		ymlFiles, err := filepath.Glob(filepath.Join(l.Config.VariablesDir, "*.yml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yml 变量文件失败")
		}
		files = append(yamlFiles, ymlFiles...)
	}

	// 添加额外的文件
	for _, file := range l.Config.VariableFiles {
		if fileExists(file) {
			files = append(files, file)
		}
	}

	return files, nil
}

// 从环境变量加载变量
func (l *CustomVariableLoader) loadFromEnvironment() map[string]interface{} {
	result := make(map[string]interface{})

	// 如果没有设置前缀，则不加载环境变量
	if l.EnvPrefix == "" {
		return result
	}

	// 遍历所有环境变量
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// 检查是否有指定的前缀
		if strings.HasPrefix(key, l.EnvPrefix) {
			// 移除前缀并转换为小驼峰命名
			varName := strings.TrimPrefix(key, l.EnvPrefix)
			if varName == "" {
				continue
			}

			// 转换为小驼峰命名
			varName = toCamelCase(varName)

			// 添加到结果中
			result[varName] = value
		}
	}

	return result
}

// 转换为小驼峰命名
func toCamelCase(s string) string {
	// 将下划线分隔的字符串转换为小驼峰命名
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// 简化的生成过程
func generateFiles(cfg *config.Config, variableLoader VariableLoader) ([]generator.GeneratedFile, error) {
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

	// 使用自定义变量加载器加载变量
	variables, err := variableLoader.LoadVariables()
	if err != nil {
		return nil, errors.Wrap(err, "加载变量失败")
	}

	// 查找变量文件
	var variableFiles []string
	if dirExists(variablesDir) {
		yamlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yaml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yaml 变量文件失败")
		}
		ymlFiles, err := filepath.Glob(filepath.Join(variablesDir, "*.yml"))
		if err != nil {
			return nil, errors.Wrap(err, "查找 *.yml 变量文件失败")
		}
		variableFiles = append(yamlFiles, ymlFiles...)
	}

	// 添加额外的文件
	for _, file := range cfg.VariableFiles {
		if fileExists(file) {
			variableFiles = append(variableFiles, file)
		}
	}

	// 加载变量到引擎
	if err := engine.LoadVariables(variableFiles); err != nil {
		return nil, errors.Wrap(err, "加载变量到引擎失败")
	}

	// 手动设置额外的变量
	for k, v := range variables {
		engine.GetVariables()[k] = v
	}

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
		TemplateDir:   "./templates",
		VariablesDir:  "./variables",
		OutputDir:     "./output",
		VariableFiles: []string{},
	}

	// 创建自定义变量加载器
	customVariableLoader := &CustomVariableLoader{
		Config: cfg,
		ExtraVariables: map[string]interface{}{
			"author":      "Your Name",
			"projectName": "Custom Project",
			"version":     "1.0.0",
		},
		EnvPrefix: "GEN_",
	}

	// 直接使用简化的生成过程
	files, err := generateFiles(cfg, customVariableLoader)
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
