package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DisplayFiles 显示文件列表
func DisplayFiles(files []VarFile, title string) {
	if title == "" {
		title = "Available Files"
	}
	
	fmt.Printf("\n%s:\n", title)
	fmt.Println("---------------------------")
	for i, file := range files {
		fmt.Printf("[%d] %s (%s)\n", i+1, file.Name, file.Path)
	}
	fmt.Println("---------------------------")
}

// GetUserSelection 获取用户选择
func GetUserSelection(files []VarFile, prompt string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	
	if prompt == "" {
		prompt = "Please select file numbers (comma separated, or 'all' to select all)"
	}
	
	fmt.Printf("\n%s: ", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %v", err)
	}
	
	input = strings.TrimSpace(input)
	
	// 如果用户输入 "all"，选择所有文件
	if strings.ToLower(input) == "all" {
		var paths []string
		for _, file := range files {
			paths = append(paths, file.Path)
		}
		return paths, nil
	}
	
	// 解析用户输入的编号
	var selectedPaths []string
	selections := strings.Split(input, ",")
	
	for _, s := range selections {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		
		index, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid selection: %s", s)
		}
		
		// 检查索引是否有效
		if index < 1 || index > len(files) {
			return nil, fmt.Errorf("invalid number: %d, valid range: 1-%d", index, len(files))
		}
		
		selectedPaths = append(selectedPaths, files[index-1].Path)
	}
	
	if len(selectedPaths) == 0 {
		return nil, fmt.Errorf("no files selected")
	}
	
	return selectedPaths, nil
}
