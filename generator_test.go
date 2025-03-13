package generate

import (
	"os"
	"path/filepath"
	"testing"

	"generate/internal/config"

	"gopkg.in/yaml.v3"
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

// loadVariableFilesMap 加载 variablesDir 目录下所有 YAML 文件，并解析为 map[string]interface{} 类型的变量映射
func loadVariableFilesMap(variablesDir string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 调用原有 loadVariableFiles 获取 YAML 文件路径列表
	files, err := loadVariableFiles(variablesDir)
	if err != nil {
		return nil, err
	}

	// 导入 yaml 包用于解析 YAML 文件
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var m map[string]interface{}
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		for k, v := range m {
			result[k] = v
		}
	}
	return result, nil
}

func TestProcessTemplatePath(t *testing.T) {
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
			g := NewGenerator()
			g.variables = tt.vars

			result, err := g.processTemplatePath(tt.path, tt.vars)

			if (err != nil) != tt.wantErr {
				t.Errorf("processTemplatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("processTemplatePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadVariableFiles(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "generator-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test YAML files
	yamlContent1 := []byte(`
key1: value1
key2: value2
common: yaml1
`)
	yamlContent2 := []byte(`
key3: value3
key4: value4
common: yaml2
`)

	yamlFile1 := filepath.Join(tempDir, "vars1.yaml")
	yamlFile2 := filepath.Join(tempDir, "vars2.yml")

	if err := os.WriteFile(yamlFile1, yamlContent1, 0644); err != nil {
		t.Fatalf("Failed to write yaml file: %v", err)
	}
	if err := os.WriteFile(yamlFile2, yamlContent2, 0644); err != nil {
		t.Fatalf("Failed to write yml file: %v", err)
	}
	g := NewGenerator()

	// Test loading variable files
	var vars map[string]interface{}
	vars, err = loadVariableFilesMap(tempDir)
	if err != nil {
		t.Fatalf("loadVariableFiles() error = %v", err)
	}
	g.variables = vars

	// Verify variables were loaded and merged correctly
	expectedVars := map[string]interface{}{
		"key1":   "value1",
		"key2":   "value2",
		"key3":   "value3",
		"key4":   "value4",
		"common": "yaml2", // Last file loaded should override
	}

	for k, v := range expectedVars {
		if g.variables[k] != v {
			t.Errorf("Variable %s = %v, want %v", k, g.variables[k], v)
		}
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
			result := removeTemplateExtension(tt.path)
			if result != tt.expected {
				t.Errorf("removeTemplateExtension() = %v, want %v", result, tt.expected)
			}
		})
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
	err = g.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Verify output files
	expectedFiles := map[string]string{
		filepath.Join(outputDir, "README.md"):                          "# TestProject\n\nA test project generated by generator\n",
		filepath.Join(outputDir, "test-project", "src", "main.go"):     "package main\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
		filepath.Join(outputDir, "test-project", "src", "config.json"): "{\n  \"name\": \"testProject\",\n  \"version\": \"1.0.0\"\n}",
		filepath.Join(outputDir, "test-project", "LICENSE"):            "Copyright (c) 2025 John doe",
	}

	for path, expectedContent := range expectedFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("File %s content = %q, want %q", path, string(content), expectedContent)
		}
	}

	// Verify directory structure
	expectedDirs := []string{
		outputDir,
		filepath.Join(outputDir, "test-project"),
		filepath.Join(outputDir, "test-project", "src"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", dir)
		}
	}
}
