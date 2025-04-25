package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/clh021/generator/internal/template"
)

// MockEngine is a mock implementation of the template engine for testing
type MockEngine struct {
	generateContentFunc func(templatePath, outputPath string) (string, error)
}

func (m *MockEngine) GenerateContent(templatePath, outputPath string) (string, error) {
	return m.generateContentFunc(templatePath, outputPath)
}

func (m *MockEngine) LoadVariables(variableFiles []string) error {
	return nil
}

func (m *MockEngine) GetVariables() map[string]interface{} {
	return nil
}

func TestDefaultContentGenerator_GenerateContent(t *testing.T) {
	// 跳过这个测试，因为它依赖于模板引擎的实现
	t.Skip("Skipping test that depends on template engine implementation")
}

func TestNewDefaultContentGenerator(t *testing.T) {
	generator := NewDefaultContentGenerator()
	if generator == nil {
		t.Error("NewDefaultContentGenerator() returned nil")
	}
}

// Integration test with real template engine
func TestDefaultContentGenerator_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "content_generator_test")
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

	// Create a simple template file
	templatePath := filepath.Join(templateDir, "hello.txt.tpl")
	templateContent := "Hello, {{.name}}!"
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create a variables file
	variablesPath := filepath.Join(variablesDir, "variables.yaml")
	variablesContent := "name: World"
	if err := os.WriteFile(variablesPath, []byte(variablesContent), 0644); err != nil {
		t.Fatalf("Failed to write variables file: %v", err)
	}

	// Create template engine
	engine := template.New(templateDir, variablesDir, outputDir)

	// Load variables
	if err := engine.LoadVariables([]string{variablesPath}); err != nil {
		t.Fatalf("Failed to load variables: %v", err)
	}

	// Create content generator
	generator := NewDefaultContentGenerator()

	// Generate content
	templateFile := TemplateFile{
		Path:         templatePath,
		RelativePath: "hello.txt.tpl",
	}
	outputPath := filepath.Join(outputDir, "hello.txt")

	content, err := generator.GenerateContent(templateFile, outputPath, engine)
	if err != nil {
		t.Fatalf("Failed to generate content: %v", err)
	}

	// Check generated content
	expectedContent := "Hello, World!"
	if content != expectedContent {
		t.Errorf("Generated content = %q, want %q", content, expectedContent)
	}
}
