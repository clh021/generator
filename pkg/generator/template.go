package generator

import (
	"strings"
)

// TemplateFile 表示一个模板文件
type TemplateFile struct {
	// 模板文件的完整路径
	Path string
	// 相对于模板目录的路径
	RelativePath string
}

// TemplateFilter 定义模板过滤器接口
// 返回值: (是否包含该模板, 排除原因)
type TemplateFilter interface {
	ShouldInclude(path, relativePath string) (bool, string)
}

// DefaultTemplateFilter 默认的模板过滤器实现
type DefaultTemplateFilter struct {
	SkipChildTemplates   bool   // 是否跳过子模板
	SkipTemplateSuffixes string // 要跳过的模板文件后缀，多个后缀用逗号分隔
	SkipTemplatePrefixes string // 要跳过的模板路径前缀，多个前缀用逗号分隔
	TemplateDir          string // 模板目录，用于计算相对路径
}

// NewDefaultTemplateFilter 创建默认的模板过滤器
func NewDefaultTemplateFilter(skipChild bool, skipSuffixes, skipPrefixes, templateDir string) *DefaultTemplateFilter {
	return &DefaultTemplateFilter{
		SkipChildTemplates:   skipChild,
		SkipTemplateSuffixes: skipSuffixes,
		SkipTemplatePrefixes: skipPrefixes,
		TemplateDir:          templateDir,
	}
}

// ShouldInclude 检查是否应该包含指定的模板文件
func (f *DefaultTemplateFilter) ShouldInclude(path, relativePath string) (bool, string) {
	// 检查是否为子模板
	if f.SkipChildTemplates && strings.Contains(path, "__child__") {
		return false, "子模板"
	}

	// 检查是否应该跳过此模板（基于后缀）
	if f.SkipTemplateSuffixes != "" {
		suffixes := strings.Split(f.SkipTemplateSuffixes, ",")
		for _, suffix := range suffixes {
			suffix = strings.TrimSpace(suffix)
			if suffix != "" && strings.HasSuffix(path, suffix) {
				return false, "后缀匹配: " + suffix
			}
		}
	}

	// 检查是否应该跳过此模板（基于前缀）
	if f.SkipTemplatePrefixes != "" {
		prefixes := strings.Split(f.SkipTemplatePrefixes, ",")
		for _, prefix := range prefixes {
			prefix = strings.TrimSpace(prefix)
			if prefix != "" {
				if strings.HasPrefix(relativePath, prefix) {
					return false, "前缀匹配: " + prefix
				}
			}
		}
	}

	return true, ""
}
