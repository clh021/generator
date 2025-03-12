package template

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
		e := New("/tmp/template", "/tmp/config", "/tmp/output")
		if e.templateDir != "/tmp/template" {
				t.Errorf("Expected templateDir to be /tmp/template, got %s", e.templateDir)
		}
		if e.configDir != "/tmp/config" {
				t.Errorf("Expected configDir to be /tmp/config, got %s", e.configDir)
		}
		if e.outputDir != "/tmp/output" {
				t.Errorf("Expected outputDir to be /tmp/output, got %s", e.outputDir)
		}
}

func TestLoadConfig(t *testing.T) {
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", "config_test")
		if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// 创建测试配置文件
    configContent := []byte(`
key1: value1
key2:
  nestedKey: nestedValue
`)
		configPath := filepath.Join(tempDir, "config.yaml")
		if err := os.WriteFile(configPath, configContent, 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
		}

		e := New("/tmp/template", tempDir, "/tmp/output")
		err = e.LoadConfig([]string{configPath})
		if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
		}

		if e.vars["key1"] != "value1" {
				t.Errorf("Expected key1 to be 'value1', got %v", e.vars["key1"])
		}

		nestedMap, ok := e.vars["key2"].(map[string]interface{})
		if !ok {
				t.Fatalf("Expected key2 to be a map")
		}
		if nestedMap["nestedKey"] != "nestedValue" {
				t.Errorf("Expected nestedKey to be 'nestedValue', got %v", nestedMap["nestedKey"])
		}
}

func TestExecute(t *testing.T) {
		// 创建临时目录
		tempDir, err := os.MkdirTemp("", "template_test")
		if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// 创建测试模板文件
		templateContent := []byte("Hello, {{.Name}}!")
		templatePath := filepath.Join(tempDir, "test.tpl")
		if err := os.WriteFile(templatePath, templateContent, 0644); err != nil {
				t.Fatalf("Failed to write template file: %v", err)
		}

		// 创建输出目录
		outputDir := filepath.Join(tempDir, "output")
		if err := os.Mkdir(outputDir, 0755); err != nil {
				t.Fatalf("Failed to create output dir: %v", err)
		}

		e := New(tempDir, "/tmp/config", outputDir)
		e.vars["Name"] = "World"

		outputPath := filepath.Join(outputDir, "result.txt")
		err = e.Execute(templatePath, outputPath)
		if err != nil {
				t.Fatalf("Execute failed: %v", err)
		}

		// 检查输出文件内容
		content, err := os.ReadFile(outputPath)
		if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
		}

		expectedContent := "Hello, World!"
		if string(content) != expectedContent {
				t.Errorf("Expected content to be '%s', got '%s'", expectedContent, string(content))
		}
}

func TestFuncMap(t *testing.T) {
		e := New("/tmp/template", "/tmp/config", "/tmp/output")
		funcMap := e.funcMap()

		// 测试 dict 函数
		dictResult, err := funcMap["dict"].(func(...interface{}) (map[string]interface{}, error))("key1", "value1", "key2", 2)
		if err != nil {
				t.Fatalf("dict function failed: %v", err)
		}
		if dictResult["key1"] != "value1" || dictResult["key2"] != 2 {
				t.Errorf("dict function returned unexpected result: %v", dictResult)
		}

		// 测试 currentYear 函数
		year := funcMap["currentYear"].(func() int)()
		if year != time.Now().Year() {
				t.Errorf("currentYear function returned %d, expected %d", year, time.Now().Year())
		}

		// 测试 default 函数
		defaultResult := funcMap["default"].(func(interface{}, interface{}) interface{})("", "defaultValue")
		if defaultResult != "defaultValue" {
				t.Errorf("default function returned %v, expected 'defaultValue'", defaultResult)
		}
}