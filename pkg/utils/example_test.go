package utils_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/clh021/generator/pkg/utils"
)

// This example demonstrates how to use the file scanning function.
func ExampleScanYAMLFiles() {
	// Create a temporary directory for the example
	tempDir, err := os.MkdirTemp("", "generator-example")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some YAML files for the example
	files := map[string]string{
		"common.yml":    "key: value",
		"project1.yml":  "name: Project 1",
		"project2.yml":  "name: Project 2",
		"settings.yaml": "debug: true",
	}

	for name, content := range files {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write file %s: %v", name, err)
		}
	}

	// Scan the directory for YAML files
	varFiles, err := utils.ScanYAMLFiles(tempDir)
	if err != nil {
		log.Fatalf("Failed to scan directory: %v", err)
	}

	// Display the files to the user
	utils.DisplayFiles(varFiles, "Available Configuration Files")

	// In a real application, you would get user input here
	// For the example, we'll simulate selecting files 1 and 3
	fmt.Println("\nSimulating user input: 1,3")

	// Get the common file
	commonFile, err := utils.GetCommonFile(tempDir, "common.yml")
	if err != nil {
		log.Fatalf("Failed to get common file: %v", err)
	}

	// In a real application, you would use GetUserSelection
	// For the example, we'll manually select files
	selectedPaths := []string{
		varFiles[0].Path, // project1.yml
		varFiles[2].Path, // settings.yaml
	}

	// Add the common file to the selected paths
	allPaths := append([]string{commonFile}, selectedPaths...)

	// Print the selected files
	fmt.Println("\nSelected files:")
	for _, path := range allPaths {
		fmt.Println("-", path)
	}

	// 这个示例不会产生固定的输出，因为路径是动态生成的
	fmt.Println("Example completed successfully")
}

// This example demonstrates how to use the DisplayFiles function.
func ExampleDisplayFiles() {
	// Create some sample VarFiles
	files := []utils.VarFile{
		{
			Path:        "/path/to/file1.yml",
			Name:        "file1",
			Description: "file1",
		},
		{
			Path:        "/path/to/file2.yml",
			Name:        "file2",
			Description: "file2",
		},
	}

	// Display the files
	utils.DisplayFiles(files, "Sample Files")

	// Output:
	//
	// Sample Files:
	// ---------------------------
	// [1] file1 (/path/to/file1.yml)
	// [2] file2 (/path/to/file2.yml)
	// ---------------------------
}
