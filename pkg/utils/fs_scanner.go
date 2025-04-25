package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// VarFile 表示一个变量文件
type VarFile struct {
	// 文件路径
	Path string
	// 文件名（不含扩展名）
	Name string
	// 文件描述（目前使用文件名）
	Description string
}

// ScanYAMLFiles 扫描目录，返回所有 YAML 文件
func ScanYAMLFiles(dir string) ([]VarFile, error) {
	var files []VarFile

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dir)
	}

	// 扫描目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 YAML 文件
		if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
			name := strings.TrimSuffix(strings.TrimSuffix(info.Name(), ".yml"), ".yaml")
			files = append(files, VarFile{
				Path:        path,
				Name:        name,
				Description: name,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %v", err)
	}

	// 按文件名排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	return files, nil
}

// GetCommonFile 获取通用文件
func GetCommonFile(dir, filename string) (string, error) {
	filePath := filepath.Join(dir, filename)
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}
	
	return filePath, nil
}
