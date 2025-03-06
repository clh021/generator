package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// convertJSONToYAML 将 JSON 文件转换为 YAML 文件。
func convertJSONToYAML(jsonFilePath, yamlOutputPath string) error {
	// 读取 JSON 文件
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return fmt.Errorf("打开 JSON 文件时出错: %w", err)
	}
	defer jsonFile.Close()

	// 使用 io.ReadAll 替代 ioutil.ReadAll
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("读取 JSON 文件时出错: %w", err)
	}

	// 将 JSON 数据反序列化到一个 interface{} 中
	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return fmt.Errorf("反序列化 JSON 时出错: %w", err)
	}

	// 将数据序列化为 YAML 格式
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化为 YAML 时出错: %w", err)
	}

	// 将 YAML 数据写入文件
	err = os.WriteFile(yamlOutputPath, yamlData, 0644) // 使用 os.WriteFile 替代 ioutil.WriteFile
	if err != nil {
		return fmt.Errorf("写入 YAML 文件时出错: %w", err)
	}

	fmt.Printf("成功将 %s 转换为 %s\n", jsonFilePath, yamlOutputPath)
	return nil
}

func main() {
	// 检查是否通过命令行参数提供了 JSON 文件路径和 YAML 输出路径
	if len(os.Args) < 3 {
		fmt.Println("错误：未提供 JSON 文件路径或 YAML 输出路径")
		fmt.Println("用法：", os.Args[0], "<json文件路径> <yaml输出路径>")
		fmt.Println("功能：将指定的 JSON 文件转换为 YAML 文件")
		return
	}

	jsonFilePath := os.Args[1]
	yamlOutputPath := os.Args[2]

	// 如果 yamlOutputPath 不带后缀，则识别为目录，并在目录下生成同名 YAML 文件
	if filepath.Ext(yamlOutputPath) == "" {
		baseName := filepath.Base(jsonFilePath)
		yamlOutputPath = filepath.Join(yamlOutputPath, strings.TrimSuffix(baseName, filepath.Ext(baseName))+".yaml")
	}

	err := convertJSONToYAML(jsonFilePath, yamlOutputPath)
	if err != nil {
		fmt.Println("转换 JSON 到 YAML 时出错:", err)
		return
	}
}
