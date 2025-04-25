package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultVariableLoader_FindVariableFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "variable_loader_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test variable files
	files := map[string]string{
		"variables.yaml":      "key: value",
		"variables.yml":       "key2: value2",
		"additional.yaml":     "key3: value3",
		"subdir/nested.yaml": "nested: value",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name            string
		variablesDir    string
		additionalFiles []string
		expectedCount   int
		expectedFiles   []string
		expectedError   bool
	}{
		{
			name:            "find all yaml files in directory",
			variablesDir:    tempDir,
			additionalFiles: []string{},
			expectedCount:   3, // variables.yaml, variables.yml, additional.yaml (not nested.yaml)
			expectedFiles:   []string{filepath.Join(tempDir, "variables.yaml"), filepath.Join(tempDir, "variables.yml"), filepath.Join(tempDir, "additional.yaml")},
			expectedError:   false,
		},
		{
			name:            "find with additional files",
			variablesDir:    tempDir,
			additionalFiles: []string{filepath.Join(tempDir, "subdir/nested.yaml")},
			expectedCount:   4,
			expectedFiles:   []string{filepath.Join(tempDir, "variables.yaml"), filepath.Join(tempDir, "variables.yml"), filepath.Join(tempDir, "additional.yaml"), filepath.Join(tempDir, "subdir/nested.yaml")},
			expectedError:   false,
		},
		{
			name:            "non-existent directory",
			variablesDir:    filepath.Join(tempDir, "non-existent"),
			additionalFiles: []string{},
			expectedCount:   0,
			expectedFiles:   []string{},
			expectedError:   false,
		},
		{
			name:            "non-existent additional file",
			variablesDir:    tempDir,
			additionalFiles: []string{filepath.Join(tempDir, "non-existent.yaml")},
			expectedCount:   3,
			expectedFiles:   []string{filepath.Join(tempDir, "variables.yaml"), filepath.Join(tempDir, "variables.yml"), filepath.Join(tempDir, "additional.yaml")},
			expectedError:   false,
		},
		{
			name:            "only additional files",
			variablesDir:    "",
			additionalFiles: []string{filepath.Join(tempDir, "additional.yaml")},
			expectedCount:   1,
			expectedFiles:   []string{filepath.Join(tempDir, "additional.yaml")},
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewDefaultVariableLoader("", "", "")
			
			files, err := loader.FindVariableFiles(tt.variablesDir, tt.additionalFiles)
			
			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("DefaultVariableLoader.FindVariableFiles() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			
			if tt.expectedError {
				return
			}
			
			// Check count
			if len(files) != tt.expectedCount {
				t.Errorf("DefaultVariableLoader.FindVariableFiles() returned %d files, expected %d", len(files), tt.expectedCount)
			}
			
			// Check files
			fileMap := make(map[string]bool)
			for _, file := range files {
				fileMap[file] = true
			}
			
			for _, expectedFile := range tt.expectedFiles {
				if !fileMap[expectedFile] {
					t.Errorf("Expected file %s not found in result", expectedFile)
				}
			}
		})
	}
}

func TestDefaultVariableLoader_LoadVariables(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "variable_loader_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create template directory
	templateDir := filepath.Join(tempDir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create variables directory
	variablesDir := filepath.Join(tempDir, "variables")
	if err := os.MkdirAll(variablesDir, 0755); err != nil {
		t.Fatalf("Failed to create variables directory: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create test variable files
	files := map[string]string{
		"variables.yaml":  "name: Test\nversion: 1.0.0",
		"additional.yaml": "description: Test project\nauthor: Test Author",
	}

	for path, content := range files {
		fullPath := filepath.Join(variablesDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name            string
		variablesDir    string
		additionalFiles []string
		expectedVars    map[string]interface{}
		expectedError   bool
	}{
		{
			name:            "load all variables",
			variablesDir:    variablesDir,
			additionalFiles: []string{},
			expectedVars: map[string]interface{}{
				"name":        "Test",
				"version":     "1.0.0",
				"description": "Test project",
				"author":      "Test Author",
			},
			expectedError: false,
		},
		{
			name:            "load with additional file",
			variablesDir:    variablesDir,
			additionalFiles: []string{filepath.Join(tempDir, "extra.yaml")},
			expectedVars: map[string]interface{}{
				"name":        "Test",
				"version":     "1.0.0",
				"description": "Test project",
				"author":      "Test Author",
			},
			expectedError: false,
		},
		{
			name:            "non-existent directory with no additional files",
			variablesDir:    filepath.Join(tempDir, "non-existent"),
			additionalFiles: []string{},
			expectedVars:    nil,
			expectedError:   true,
		},
		{
			name:            "non-existent directory with additional files",
			variablesDir:    filepath.Join(tempDir, "non-existent"),
			additionalFiles: []string{filepath.Join(variablesDir, "variables.yaml")},
			expectedVars: map[string]interface{}{
				"name":    "Test",
				"version": "1.0.0",
			},
			expectedError: false,
		},
	}

	// Create an extra file outside the variables directory
	extraFilePath := filepath.Join(tempDir, "extra.yaml")
	extraFileContent := "extra: value"
	if err := os.WriteFile(extraFilePath, []byte(extraFileContent), 0644); err != nil {
		t.Fatalf("Failed to write extra file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewDefaultVariableLoader(templateDir, tt.variablesDir, outputDir)
			
			vars, err := loader.LoadVariables(tt.variablesDir, tt.additionalFiles)
			
			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("DefaultVariableLoader.LoadVariables() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			
			if tt.expectedError {
				return
			}
			
			// Check variables
			for key, expectedValue := range tt.expectedVars {
				value, ok := vars[key]
				if !ok {
					t.Errorf("Expected variable %s not found in result", key)
					continue
				}
				
				if !reflect.DeepEqual(value, expectedValue) {
					t.Errorf("Variable %s = %v, want %v", key, value, expectedValue)
				}
			}
		})
	}
}

func TestNewDefaultVariableLoader(t *testing.T) {
	loader := NewDefaultVariableLoader("/templates", "/variables", "/output")
	
	if loader.TemplateDir != "/templates" {
		t.Errorf("NewDefaultVariableLoader() TemplateDir = %v, want %v", loader.TemplateDir, "/templates")
	}
	
	if loader.VariablesDir != "/variables" {
		t.Errorf("NewDefaultVariableLoader() VariablesDir = %v, want %v", loader.VariablesDir, "/variables")
	}
	
	if loader.OutputDir != "/output" {
		t.Errorf("NewDefaultVariableLoader() OutputDir = %v, want %v", loader.OutputDir, "/output")
	}
}

func TestDirExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dir_exists_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing directory",
			path: tempDir,
			want: true,
		},
		{
			name: "non-existent directory",
			path: filepath.Join(tempDir, "non-existent"),
			want: false,
		},
		{
			name: "file (not a directory)",
			path: filePath,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dirExists(tt.path); got != tt.want {
				t.Errorf("dirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_exists_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: filePath,
			want: true,
		},
		{
			name: "non-existent file",
			path: filepath.Join(tempDir, "non-existent.txt"),
			want: false,
		},
		{
			name: "directory (not a file)",
			path: tempDir,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.path); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
