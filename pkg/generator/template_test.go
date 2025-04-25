package generator

import (
	"testing"
)

func TestDefaultTemplateFilter_ShouldInclude(t *testing.T) {
	tests := []struct {
		name         string
		filter       *DefaultTemplateFilter
		path         string
		relativePath string
		wantInclude  bool
		wantReason   string
	}{
		{
			name: "include normal template",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "",
				TemplateDir:          "/templates",
			},
			path:         "/templates/normal.tpl",
			relativePath: "normal.tpl",
			wantInclude:  true,
			wantReason:   "",
		},
		{
			name: "exclude child template",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "",
				TemplateDir:          "/templates",
			},
			path:         "/templates/__child__/child.tpl",
			relativePath: "__child__/child.tpl",
			wantInclude:  false,
			wantReason:   "子模板",
		},
		{
			name: "exclude by suffix",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: ".go.tpl,.vue.tpl",
				SkipTemplatePrefixes: "",
				TemplateDir:          "/templates",
			},
			path:         "/templates/file.go.tpl",
			relativePath: "file.go.tpl",
			wantInclude:  false,
			wantReason:   "后缀匹配: .go.tpl",
		},
		{
			name: "exclude by prefix",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "web,server/config",
				TemplateDir:          "/templates",
			},
			path:         "/templates/web/index.tpl",
			relativePath: "web/index.tpl",
			wantInclude:  false,
			wantReason:   "前缀匹配: web",
		},
		{
			name: "exclude by nested prefix",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "web,server/config",
				TemplateDir:          "/templates",
			},
			path:         "/templates/server/config/app.tpl",
			relativePath: "server/config/app.tpl",
			wantInclude:  false,
			wantReason:   "前缀匹配: server/config",
		},
		{
			name: "include server template but not config",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "web,server/config",
				TemplateDir:          "/templates",
			},
			path:         "/templates/server/main.tpl",
			relativePath: "server/main.tpl",
			wantInclude:  true,
			wantReason:   "",
		},
		{
			name: "include when child templates not skipped",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   false,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "",
				TemplateDir:          "/templates",
			},
			path:         "/templates/__child__/child.tpl",
			relativePath: "__child__/child.tpl",
			wantInclude:  true,
			wantReason:   "",
		},
		{
			name: "multiple suffixes with spaces",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: ".go.tpl, .vue.tpl, .test.tpl",
				SkipTemplatePrefixes: "",
				TemplateDir:          "/templates",
			},
			path:         "/templates/file.test.tpl",
			relativePath: "file.test.tpl",
			wantInclude:  false,
			wantReason:   "后缀匹配: .test.tpl",
		},
		{
			name: "multiple prefixes with spaces",
			filter: &DefaultTemplateFilter{
				SkipChildTemplates:   true,
				SkipTemplateSuffixes: "",
				SkipTemplatePrefixes: "web, server/config, test",
				TemplateDir:          "/templates",
			},
			path:         "/templates/test/unit.tpl",
			relativePath: "test/unit.tpl",
			wantInclude:  false,
			wantReason:   "前缀匹配: test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInclude, gotReason := tt.filter.ShouldInclude(tt.path, tt.relativePath)
			if gotInclude != tt.wantInclude {
				t.Errorf("DefaultTemplateFilter.ShouldInclude() include = %v, want %v", gotInclude, tt.wantInclude)
			}
			if gotReason != tt.wantReason {
				t.Errorf("DefaultTemplateFilter.ShouldInclude() reason = %v, want %v", gotReason, tt.wantReason)
			}
		})
	}
}

func TestNewDefaultTemplateFilter(t *testing.T) {
	filter := NewDefaultTemplateFilter(true, ".go.tpl", "web", "/templates")
	
	if filter.SkipChildTemplates != true {
		t.Errorf("NewDefaultTemplateFilter() SkipChildTemplates = %v, want %v", filter.SkipChildTemplates, true)
	}
	
	if filter.SkipTemplateSuffixes != ".go.tpl" {
		t.Errorf("NewDefaultTemplateFilter() SkipTemplateSuffixes = %v, want %v", filter.SkipTemplateSuffixes, ".go.tpl")
	}
	
	if filter.SkipTemplatePrefixes != "web" {
		t.Errorf("NewDefaultTemplateFilter() SkipTemplatePrefixes = %v, want %v", filter.SkipTemplatePrefixes, "web")
	}
	
	if filter.TemplateDir != "/templates" {
		t.Errorf("NewDefaultTemplateFilter() TemplateDir = %v, want %v", filter.TemplateDir, "/templates")
	}
}
