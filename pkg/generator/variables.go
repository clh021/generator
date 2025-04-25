package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/clh021/generator/internal/template"
	"github.com/pkg/errors"
)

// VariableLoader 定义变量加载器接口
type VariableLoader interface {
	// LoadVariables 加载变量文件并返回变量映射
	LoadVariables(variablesDir string, additionalFiles []string) (map[string]interface{}, error)
	// FindVariableFiles 查找变量文件
	FindVariableFiles(variablesDir string, additionalFiles []string) ([]string, error)
}

// DefaultVariableLoader 默认的变量加载器实现
type DefaultVariableLoader struct {
	TemplateDir  string
	VariablesDir string
	OutputDir    string
}

// NewDefaultVariableLoader 创建默认的变量加载器
func NewDefaultVariableLoader(templateDir, variablesDir, outputDir string) *DefaultVariableLoader {
	return &DefaultVariableLoader{
		TemplateDir:  templateDir,
		VariablesDir: variablesDir,
		OutputDir:    outputDir,
	}
}

// LoadVariables 加载变量文件并返回变量映射
func (l *DefaultVariableLoader) LoadVariables(variablesDir string, additionalFiles []string) (map[string]interface{}, error) {
	// 检查变量目录是否存在
	if _, err := os.Stat(variablesDir); os.IsNotExist(err) && len(additionalFiles) == 0 {
		return nil, errors.Wrapf(err, "变量目录不存在且未指定变量文件: %s", variablesDir)
	}

	// 查找变量文件
	variableFiles, err := l.FindVariableFiles(variablesDir, additionalFiles)
	if err != nil {
		return nil, errors.Wrap(err, "查找变量文件失败")
	}

	// 检查是否有可用的变量文件
	if len(variableFiles) == 0 {
		return nil, errors.New("没有找到可用的变量文件")
	}

	// 创建模板引擎
	engine := template.New(l.TemplateDir, l.VariablesDir, l.OutputDir)

	// 加载变量
	if err := engine.LoadVariables(variableFiles); err != nil {
		return nil, errors.Wrap(err, "加载变量失败")
	}

	// 获取变量
	return engine.GetVariables(), nil
}

// FindVariableFiles 查找变量文件
func (l *DefaultVariableLoader) FindVariableFiles(variablesDir string, additionalFiles []string) ([]string, error) {
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
