package template

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	e := New("/tmp/template", "/tmp/variables", "/tmp/output")
	if e.templateDir != "/tmp/template" {
		t.Errorf("Expected templateDir to be /tmp/template, got %s", e.templateDir)
	}
	if e.outputDir != "/tmp/output" {
		t.Errorf("Expected outputDir to be /tmp/output, got %s", e.outputDir)
	}
}

func TestLoadVariables(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "variables_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试变量文件
	variablesContent1 := []byte(`
key1: value1
key2:
  nestedKey: nestedValue
`)
	variablesPath1 := filepath.Join(tempDir, "variables1.yaml")
	if err := os.WriteFile(variablesPath1, variablesContent1, 0644); err != nil {
		t.Fatalf("Failed to write variables file: %v", err)
	}

	variablesContent2 := []byte(`
key3: value3
key2:
  nestedKey2: nestedValue2
`)
	variablesPath2 := filepath.Join(tempDir, "variables2.yaml")
	if err := os.WriteFile(variablesPath2, variablesContent2, 0644); err != nil {
		t.Fatalf("Failed to write variables file: %v", err)
	}

	e := New("/tmp/template", tempDir, "/tmp/output")
	err = e.LoadVariables([]string{variablesPath1, variablesPath2})
	if err != nil {
		t.Fatalf("LoadVariables failed: %v", err)
	}

	// 检查变量是否正确加载
	if e.vars["key1"] != "value1" {
		t.Errorf("Expected key1 to be 'value1', got %v", e.vars["key1"])
	}

	if e.vars["key3"] != "value3" {
		t.Errorf("Expected key3 to be 'value3', got %v", e.vars["key3"])
	}

	nestedMap, ok := e.vars["key2"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected key2 to be a map")
	}

	// 检查嵌套的键值对
	if v, ok := nestedMap["nestedKey"].(string); !ok || v != "nestedValue" {
		t.Errorf("Expected nestedKey to be 'nestedValue', got %v", nestedMap["nestedKey"])
	}
	if v, ok := nestedMap["nestedKey2"].(string); !ok || v != "nestedValue2" {
		t.Errorf("Expected nestedKey2 to be 'nestedValue2', got %v", nestedMap["nestedKey2"])
	}

	// 测试加载不存在的文件
	err = e.LoadVariables([]string{"/non/existent/file.yaml"})
	if err == nil {
		t.Errorf("Expected error when loading non-existent file, got nil")
	}
}

func TestGenerateContent(t *testing.T) {
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
	content, err := e.GenerateContent(templatePath, outputPath)
	if err != nil {
		t.Fatalf("GenerateContent failed: %v", err)
	}

	expectedContent := "Hello, World!"
	if content != expectedContent {
		t.Errorf("Expected content to be '%s', got '%s'", expectedContent, content)
	}
}

func TestGenerateContentWithFunctions(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "template_functions_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试模板文件，包含条件语句、循环和内置函数
	templateContent := []byte(`
Name: {{.Name}}
Lowercase first: {{lcfirst .Name}}
Uppercase first: {{ucfirst .Title}}
{{if .ShowGreeting}}Greeting: Hello!{{else}}No greeting{{end}}
{{if .Items}}Items:{{range .Items}}
- {{.}}{{end}}
{{else}}No items{{end}}
Current Year: {{currentYear}}
Default value: {{default .MissingValue "default"}}
`)
	templatePath := filepath.Join(tempDir, "advanced.tpl")
	if err := os.WriteFile(templatePath, templateContent, 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// 创建输出目录
	outputDir := filepath.Join(tempDir, "output")
	if err := os.Mkdir(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	e := New(tempDir, "/tmp/config", outputDir)
	e.vars = map[string]interface{}{
		"Name":         "World",
		"Title":        "project",
		"ShowGreeting": true,
		"Items":        []string{"Item1", "Item2", "Item3"},
		"MissingValue": "",
	}

	outputPath := filepath.Join(outputDir, "advanced_result.txt")
	content, err := e.GenerateContent(templatePath, outputPath)
	if err != nil {
		t.Fatalf("GenerateContent failed: %v", err)
	}

	// 验证输出内容包含预期的模板函数处理结果
	expectedParts := []string{
		"Name: World",
		"Lowercase first: world",
		"Uppercase first: Project",
		"Greeting: Hello!",
		"Items:",
		"- Item1",
		"- Item2",
		"- Item3",
		"Current Year: " + strconv.Itoa(time.Now().Year()),
		"Default value: default",
	}

	for _, part := range expectedParts {
		if !strings.Contains(content, part) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nOutput: %s", part, content)
		}
	}
}

func TestLcfirst(t *testing.T) {
	e := New("/tmp/template", "/tmp/config", "/tmp/output")
	funcMap := e.funcMap()
	lcfirstFunc := funcMap["lcfirst"].(func(string) string)

	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},                       // 空字符串
		{"A", "a"},                     // 单字符
		{"ABC", "aBC"},                 // 普通字符串
		{"Hello World", "hello World"}, // 带空格的字符串
		{"Über", "über"},               // 非ASCII字符
		{"123", "123"},                 // 数字开头
		{"_test", "_test"},             // 特殊字符开头
		{" space", " space"},           // 空格开头
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := lcfirstFunc(tc.input)
			if result != tc.expected {
				t.Errorf("lcfirst(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestUcfirst(t *testing.T) {
	e := New("/tmp/template", "/tmp/config", "/tmp/output")
	funcMap := e.funcMap()
	ucfirstFunc := funcMap["ucfirst"].(func(string) string)

	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},                       // 空字符串
		{"a", "A"},                     // 单字符
		{"abc", "Abc"},                 // 普通字符串
		{"hello world", "Hello world"}, // 带空格的字符串
		{"über", "Über"},               // 非ASCII字符
		{"123", "123"},                 // 数字开头
		{"_test", "_test"},             // 特殊字符开头
		{" space", " space"},           // 空格开头
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ucfirstFunc(tc.input)
			if result != tc.expected {
				t.Errorf("ucfirst(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
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
