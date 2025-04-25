package generator

import (
	"os"
	"path/filepath"
	"testing"
)

// MockTemplateFilter is a mock implementation of TemplateFilter for testing
type MockTemplateFilter struct {
	shouldIncludeFunc func(path, relativePath string) (bool, string)
}

func (m *MockTemplateFilter) ShouldInclude(path, relativePath string) (bool, string) {
	return m.shouldIncludeFunc(path, relativePath)
}

func TestDefaultTemplateScanner_ScanTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test template files
	files := map[string]string{
		"file1.tpl":           "content1",
		"file2.tpl":           "content2",
		"subdir/file3.tpl":    "content3",
		"__child__/child.tpl": "child content",
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
		name           string
		filter         TemplateFilter
		expectedCount  int
		expectedPaths  []string
		expectedError  bool
		nonExistentDir bool
	}{
		{
			name: "include all files",
			filter: &MockTemplateFilter{
				shouldIncludeFunc: func(path, relativePath string) (bool, string) {
					return true, ""
				},
			},
			expectedCount: 4,
			expectedPaths: []string{"file1.tpl", "file2.tpl", "subdir/file3.tpl", "__child__/child.tpl"},
			expectedError: false,
		},
		{
			name: "exclude child templates",
			filter: &MockTemplateFilter{
				shouldIncludeFunc: func(path, relativePath string) (bool, string) {
					return !filepath.HasPrefix(relativePath, "__child__"), "child template"
				},
			},
			expectedCount: 3,
			expectedPaths: []string{"file1.tpl", "file2.tpl", "subdir/file3.tpl"},
			expectedError: false,
		},
		{
			name: "exclude subdirectory",
			filter: &MockTemplateFilter{
				shouldIncludeFunc: func(path, relativePath string) (bool, string) {
					return !filepath.HasPrefix(relativePath, "subdir/"), "subdirectory"
				},
			},
			expectedCount: 3,
			expectedPaths: []string{"file1.tpl", "file2.tpl", "__child__/child.tpl"},
			expectedError: false,
		},
		{
			name:           "non-existent directory",
			filter:         &MockTemplateFilter{shouldIncludeFunc: func(path, relativePath string) (bool, string) { return true, "" }},
			expectedCount:  0,
			expectedPaths:  []string{},
			expectedError:  true,
			nonExistentDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewDefaultTemplateScanner()
			
			dir := tempDir
			if tt.nonExistentDir {
				dir = filepath.Join(tempDir, "non-existent")
			}
			
			templateFiles, err := scanner.ScanTemplates(dir, tt.filter)
			
			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("DefaultTemplateScanner.ScanTemplates() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			
			if tt.expectedError {
				return
			}
			
			// Check count
			if len(templateFiles) != tt.expectedCount {
				t.Errorf("DefaultTemplateScanner.ScanTemplates() returned %d files, expected %d", len(templateFiles), tt.expectedCount)
			}
			
			// Check paths
			pathMap := make(map[string]bool)
			for _, file := range templateFiles {
				pathMap[file.RelativePath] = true
			}
			
			for _, expectedPath := range tt.expectedPaths {
				if !pathMap[expectedPath] {
					t.Errorf("Expected path %s not found in result", expectedPath)
				}
			}
		})
	}
}

func TestNewDefaultTemplateScanner(t *testing.T) {
	scanner := NewDefaultTemplateScanner()
	if scanner == nil {
		t.Error("NewDefaultTemplateScanner() returned nil")
	}
}
