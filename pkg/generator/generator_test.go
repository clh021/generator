package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/clh021/generator/pkg/config"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator()

	if g == nil {
		t.Fatal("NewGenerator returned nil")
	}

	if g.variables == nil {
		t.Error("variables map not initialized")
	}
}

func TestDefaultPathProcessor_ProcessTemplatePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		vars     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "no variables",
			path:     "path/to/file.txt",
			vars:     map[string]interface{}{},
			expected: "path/to/file.txt",
			wantErr:  false,
		},
		{
			name:     "single variable",
			path:     "path/to/__name__.txt",
			vars:     map[string]interface{}{"name": "test"},
			expected: "path/to/test.txt",
			wantErr:  false,
		},
		{
			name:     "multiple variables",
			path:     "__dir__/to/__name__.txt",
			vars:     map[string]interface{}{"dir": "path", "name": "test"},
			expected: "path/to/test.txt",
			wantErr:  false,
		},
		{
			name:     "missing variable",
			path:     "path/to/__missing__.txt",
			vars:     map[string]interface{}{},
			expected: "path/to/__missing__.txt",
			wantErr:  false,
		},
		{
			name:     "non-string variable",
			path:     "path/to/__number__.txt",
			vars:     map[string]interface{}{"number": 123},
			expected: "path/to/__number__.txt",
			wantErr:  false,
		},
		{
			name:     "same variable multiple times",
			path:     "__prefix__/to/__prefix__.txt",
			vars:     map[string]interface{}{"prefix": "test"},
			expected: "test/to/test.txt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewDefaultPathProcessor()
			templateFile := TemplateFile{
				Path:         "dummy/path",
				RelativePath: tt.path,
			}

			result, err := processor.ProcessOutputPath(templateFile, "", tt.vars)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessOutputPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("ProcessOutputPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultVariableLoader_FindVariableFiles(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "variable_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	yamlFile := filepath.Join(tempDir, "test.yaml")
	if err := os.WriteFile(yamlFile, []byte("key: value"), 0644); err != nil {
		t.Fatalf("Failed to write yaml file: %v", err)
	}

	ymlFile := filepath.Join(tempDir, "test.yml")
	if err := os.WriteFile(ymlFile, []byte("key2: value2"), 0644); err != nil {
		t.Fatalf("Failed to write yml file: %v", err)
	}

	additionalFile := filepath.Join(tempDir, "additional.yaml")
	if err := os.WriteFile(additionalFile, []byte("key3: value3"), 0644); err != nil {
		t.Fatalf("Failed to write additional file: %v", err)
	}

	// 测试加载目录和额外文件
	loader := NewDefaultVariableLoader("", "", "")
	files, err := loader.FindVariableFiles(tempDir, []string{additionalFile})
	if err != nil {
		t.Fatalf("FindVariableFiles failed: %v", err)
	}

	// 检查是否包含所有预期的文件
	expectedPaths := []string{yamlFile, ymlFile, additionalFile}
	for _, expected := range expectedPaths {
		found := false
		for _, file := range files {
			if file == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in result", expected)
		}
	}

	// 测试非存在的目录
	files, err = loader.FindVariableFiles("/non/existent/dir", nil)
	if err != nil {
		t.Errorf("Expected no error for non-existent directory, got: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected empty file list for non-existent directory, got: %v", files)
	}

	// 测试非存在的额外文件
	files, err = loader.FindVariableFiles("", []string{"/non/existent/file.yaml"})
	if err != nil {
		t.Errorf("Expected no error for non-existent additional file, got: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected empty file list for non-existent additional file, got: %v", files)
	}
}

func TestRemoveTemplateExtension(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "with .tpl extension",
			path:     "path/to/file.txt.tpl",
			expected: "path/to/file.txt",
		},
		{
			name:     "without .tpl extension",
			path:     "path/to/file.txt",
			expected: "path/to/file.txt",
		},
		{
			name:     "with multiple dots",
			path:     "path/to/file.config.json.tpl",
			expected: "path/to/file.config.json",
		},
		{
			name:     "only .tpl extension",
			path:     "file.tpl",
			expected: "file",
		},
		{
			name:     "no extension",
			path:     "file",
			expected: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 path 包中的函数
			result := removeTemplateExtension(tt.path)
			if result != tt.expected {
				t.Errorf("removeTemplateExtension() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSkipTemplateSuffixes(t *testing.T) {
	// Create temporary directories
	rootDir, err := os.MkdirTemp("", "generator-skip-suffix-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	templateDir := filepath.Join(rootDir, "templates")
	variableDir := filepath.Join(rootDir, "variables")
	outputDir := filepath.Join(rootDir, "output")

	for _, dir := range []string{templateDir, variableDir, outputDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create template files with different suffixes
	templateFiles := map[string]string{
		filepath.Join(templateDir, "normal.txt.tpl"):      "Normal template",
		filepath.Join(templateDir, "skip1.go.tpl.tpl"):    "Should be skipped by suffix .go.tpl.tpl",
		filepath.Join(templateDir, "skip2.vue.tpl"):       "Should be skipped by suffix .vue.tpl",
		filepath.Join(templateDir, "not_skip.go.tpl"):     "Should not be skipped",
	}

	for path, content := range templateFiles {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", path, err)
		}
	}

	// Create variable file
	variableContent := []byte(`
key: value
`)

	variableFile := filepath.Join(variableDir, "variables.yaml")
	if err := os.WriteFile(variableFile, variableContent, 0644); err != nil {
		t.Fatalf("Failed to write variable file: %v", err)
	}

	// Run generator with skip suffixes
	g := NewGenerator()
	cfg := &config.Config{
		TemplateDir:          templateDir,
		VariablesDir:         variableDir,
		OutputDir:            outputDir,
		SkipTemplateSuffixes: ".go.tpl.tpl,.vue.tpl",
	}
	files, err := g.GenerateFiles(cfg)
	if err != nil {
		t.Fatalf("GenerateFiles() error = %v", err)
	}

	// Verify generated files
	expectedPaths := map[string]bool{
		filepath.Join(outputDir, "normal.txt"):  true,  // Should exist
		filepath.Join(outputDir, "skip1.go.tpl"): false, // Should not exist
		filepath.Join(outputDir, "skip2.vue"):    false, // Should not exist
		filepath.Join(outputDir, "not_skip.go"):   true,  // Should exist
	}

	// Check that the expected files are generated
	for _, file := range files {
		if expectedPaths[file.OutputPath] {
			// File should exist
			delete(expectedPaths, file.OutputPath)
		} else if _, exists := expectedPaths[file.OutputPath]; exists {
			// File should not exist but was generated
			t.Errorf("File %s should not be generated, but it was", file.OutputPath)
		}
	}

	// Check that all expected files were found
	for path, shouldExist := range expectedPaths {
		if shouldExist {
			t.Errorf("Expected file %s to be generated, but it wasn't", path)
		}
	}
}

func TestSkipTemplatePrefixes(t *testing.T) {
	// Create temporary directories
	rootDir, err := os.MkdirTemp("", "generator-skip-prefix-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	templateDir := filepath.Join(rootDir, "templates")
	variableDir := filepath.Join(rootDir, "variables")
	outputDir := filepath.Join(rootDir, "output")

	for _, dir := range []string{templateDir, variableDir, outputDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create template directories with different prefixes
	webDir := filepath.Join(templateDir, "web")
	serverDir := filepath.Join(templateDir, "server")
	commonDir := filepath.Join(templateDir, "common")

	for _, dir := range []string{webDir, serverDir, commonDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create template files in different directories
	templateFiles := map[string]string{
		filepath.Join(webDir, "index.html.tpl"):      "Web template",
		filepath.Join(webDir, "style.css.tpl"):       "Web CSS",
		filepath.Join(serverDir, "main.go.tpl"):      "Server template",
		filepath.Join(serverDir, "config.json.tpl"):  "Server config",
		filepath.Join(commonDir, "README.md.tpl"):    "Common template",
	}

	for path, content := range templateFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", path, err)
		}
	}

	// Create variable file
	variableContent := []byte(`
key: value
`)

	variableFile := filepath.Join(variableDir, "variables.yaml")
	if err := os.WriteFile(variableFile, variableContent, 0644); err != nil {
		t.Fatalf("Failed to write variable file: %v", err)
	}

	// Run generator with skip prefixes
	g := NewGenerator()
	cfg := &config.Config{
		TemplateDir:          templateDir,
		VariablesDir:         variableDir,
		OutputDir:            outputDir,
		SkipTemplatePrefixes: "web,server/config",
	}
	files, err := g.GenerateFiles(cfg)
	if err != nil {
		t.Fatalf("GenerateFiles() error = %v", err)
	}

	// Verify generated files
	expectedPaths := map[string]bool{
		filepath.Join(outputDir, "web", "index.html"):     false, // Should not exist (skipped by prefix)
		filepath.Join(outputDir, "web", "style.css"):      false, // Should not exist (skipped by prefix)
		filepath.Join(outputDir, "server", "main.go"):     true,  // Should exist
		filepath.Join(outputDir, "server", "config.json"): false, // Should not exist (skipped by prefix)
		filepath.Join(outputDir, "common", "README.md"):   true,  // Should exist
	}

	// Check that the expected files are generated
	for _, file := range files {
		if expectedPaths[file.OutputPath] {
			// File should exist
			delete(expectedPaths, file.OutputPath)
		} else if _, exists := expectedPaths[file.OutputPath]; exists {
			// File should not exist but was generated
			t.Errorf("File %s should not be generated, but it was", file.OutputPath)
		}
	}

	// Check that all expected files were found
	for path, shouldExist := range expectedPaths {
		if shouldExist {
			t.Errorf("Expected file %s to be generated, but it wasn't", path)
		}
	}
}

func TestGenerateEndToEnd(t *testing.T) {
	// Create temporary directories
	rootDir, err := os.MkdirTemp("", "generator-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir)

	templateDir := filepath.Join(rootDir, "templates")
	variableDir := filepath.Join(rootDir, "variables")
	outputDir := filepath.Join(rootDir, "output")

	for _, dir := range []string{templateDir, variableDir, outputDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create nested template directory with variable in path
	nestedTemplateDir := filepath.Join(templateDir, "__project__", "src")
	if err := os.MkdirAll(nestedTemplateDir, 0755); err != nil {
		t.Fatalf("Failed to create nested template directory: %v", err)
	}

	// Create template files
	templateFiles := map[string]string{
		filepath.Join(templateDir, "README.md.tpl"):              "# {{.projectName}}\n\n{{.description}}\n",
		filepath.Join(nestedTemplateDir, "main.go.tpl"):          "package main\n\nfunc main() {\n\tfmt.Println(\"{{.greeting}}\")\n}\n",
		filepath.Join(nestedTemplateDir, "config.json.tpl"):      "{\n  \"name\": \"{{.projectName | lcfirst}}\",\n  \"version\": \"{{.version}}\"\n}",
		filepath.Join(templateDir, "__project__", "LICENSE.tpl"): "Copyright (c) {{.year}} {{.author | ucfirst}}",
	}

	for path, content := range templateFiles {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", path, err)
		}
	}

	// Create variable file
	variableContent := []byte(`
projectName: TestProject
project: test-project
description: A test project generated by generator
greeting: Hello, World!
version: 1.0.0
year: 2025
author: john doe
`)

	variableFile := filepath.Join(variableDir, "variables.yaml")
	if err := os.WriteFile(variableFile, variableContent, 0644); err != nil {
		t.Fatalf("Failed to write variable file: %v", err)
	}

	// Run generator
	g := NewGenerator()
	cfg := &config.Config{
		TemplateDir:  templateDir,
		VariablesDir: variableDir,
		OutputDir:    outputDir,
	}
	files, err := g.GenerateFiles(cfg)
	if err != nil {
		t.Fatalf("GenerateFiles() error = %v", err)
	}

	// Verify generated files
	expectedContents := map[string]string{
		filepath.Join(outputDir, "README.md"):                          "# TestProject\n\nA test project generated by generator\n",
		filepath.Join(outputDir, "test-project", "src", "main.go"):     "package main\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
		filepath.Join(outputDir, "test-project", "src", "config.json"): "{\n  \"name\": \"testProject\",\n  \"version\": \"1.0.0\"\n}",
		filepath.Join(outputDir, "test-project", "LICENSE"):            "Copyright (c) 2025 John doe",
	}

	// Check that all expected files are generated with correct content
	for _, file := range files {
		if expectedContent, ok := expectedContents[file.OutputPath]; ok {
			if file.Content != expectedContent {
				t.Errorf("File %s content = %q, want %q", file.OutputPath, file.Content, expectedContent)
			}
			delete(expectedContents, file.OutputPath)
		}
	}

	// Check that all expected files were found
	for path := range expectedContents {
		t.Errorf("Expected file %s to be generated, but it wasn't", path)
	}

	// Verify that output paths include expected directories
	expectedDirs := []string{
		outputDir,
		filepath.Join(outputDir, "test-project"),
		filepath.Join(outputDir, "test-project", "src"),
	}

	// Create a map of all parent directories from generated files
	generatedDirs := make(map[string]bool)
	for _, file := range files {
		dir := filepath.Dir(file.OutputPath)
		for dir != "." && dir != "/" {
			generatedDirs[dir] = true
			dir = filepath.Dir(dir)
		}
	}

	// Check that all expected directories would be created
	for _, dir := range expectedDirs {
		if !generatedDirs[dir] {
			t.Errorf("Expected directory %s would not be created", dir)
		}
	}
}
