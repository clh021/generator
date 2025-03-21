package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
	"unicode"
)

// lcfirst converts the first character of a string to lowercase.
// If the string is empty, it returns an empty string.
//
// Example usage in templates:
//
//	{{ lcfirst "UserName" }} -> "userName"
func lcfirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// ucfirst converts the first character of a string to uppercase.
// If the string is empty, it returns an empty string.
//
// Example usage in templates:
//
//	{{ ucfirst "userName" }} -> "UserName"
func ucfirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// funcMap 返回模板中可用的函数映射
func (e *Engine) funcMap() template.FuncMap {
	return template.FuncMap{
		"include": e.createIncludeTemplateFunc(0, []string{}),
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"currentYear": func() int {
			return time.Now().Year()
		},
		"file": func(filePath string) (string, error) {
			content, err := os.ReadFile(filepath.Join(e.templateDir, filePath))
			if err != nil {
				return "", fmt.Errorf("读取文件失败 %s: %w", filePath, err)
			}
			return string(content), nil
		},
		"default": func(value, defaultValue interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"lcfirst": lcfirst,
		"ucfirst": ucfirst,
	}
}

// GetVariables 返回模板引擎中加载的所有变量
func (e *Engine) GetVariables() map[string]interface{} {
	return e.vars
}
