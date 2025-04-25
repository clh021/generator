package generator

import (
	"path/filepath"
	"testing"
)

func TestDefaultPathProcessor_ProcessOutputPath(t *testing.T) {
	tests := []struct {
		name       string
		templateFile TemplateFile
		outputDir  string
		variables  map[string]interface{}
		want       string
		wantErr    bool
	}{
		{
			name: "simple path without variables",
			templateFile: TemplateFile{
				Path:         "/templates/file.txt.tpl",
				RelativePath: "file.txt.tpl",
			},
			outputDir:  "/output",
			variables:  map[string]interface{}{},
			want:       filepath.Clean("/output/file.txt"),
			wantErr:    false,
		},
		{
			name: "path with single variable",
			templateFile: TemplateFile{
				Path:         "/templates/file__name__.txt.tpl",
				RelativePath: "file__name__.txt.tpl",
			},
			outputDir: "/output",
			variables: map[string]interface{}{
				"name": "test",
			},
			want:    filepath.Clean("/output/filetest.txt"),
			wantErr: false,
		},
		{
			name: "path with multiple variables",
			templateFile: TemplateFile{
				Path:         "/templates/__dir__/__name__.txt.tpl",
				RelativePath: "__dir__/__name__.txt.tpl",
			},
			outputDir: "/output",
			variables: map[string]interface{}{
				"dir":  "subdir",
				"name": "test",
			},
			want:    filepath.Clean("/output/subdir/test.txt"),
			wantErr: false,
		},
		{
			name: "path with missing variable",
			templateFile: TemplateFile{
				Path:         "/templates/__missing__.txt.tpl",
				RelativePath: "__missing__.txt.tpl",
			},
			outputDir:  "/output",
			variables:  map[string]interface{}{},
			want:       filepath.Clean("/output/__missing__.txt"),
			wantErr:    false,
		},
		{
			name: "path with non-string variable",
			templateFile: TemplateFile{
				Path:         "/templates/__number__.txt.tpl",
				RelativePath: "__number__.txt.tpl",
			},
			outputDir: "/output",
			variables: map[string]interface{}{
				"number": 123,
			},
			want:    filepath.Clean("/output/__number__.txt"),
			wantErr: false,
		},
		{
			name: "path with same variable multiple times",
			templateFile: TemplateFile{
				Path:         "/templates/__prefix__/__prefix__.txt.tpl",
				RelativePath: "__prefix__/__prefix__.txt.tpl",
			},
			outputDir: "/output",
			variables: map[string]interface{}{
				"prefix": "test",
			},
			want:    filepath.Clean("/output/test/test.txt"),
			wantErr: false,
		},
		{
			name: "path with nested directories",
			templateFile: TemplateFile{
				Path:         "/templates/a/b/c/file.txt.tpl",
				RelativePath: "a/b/c/file.txt.tpl",
			},
			outputDir:  "/output",
			variables:  map[string]interface{}{},
			want:       filepath.Clean("/output/a/b/c/file.txt"),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDefaultPathProcessor()
			got, err := p.ProcessOutputPath(tt.templateFile, tt.outputDir, tt.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultPathProcessor.ProcessOutputPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DefaultPathProcessor.ProcessOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveTemplateExtension(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "with .tpl extension",
			path: "file.txt.tpl",
			want: "file.txt",
		},
		{
			name: "without .tpl extension",
			path: "file.txt",
			want: "file.txt",
		},
		{
			name: "with multiple dots",
			path: "file.config.json.tpl",
			want: "file.config.json",
		},
		{
			name: "only .tpl extension",
			path: "file.tpl",
			want: "file",
		},
		{
			name: "no extension",
			path: "file",
			want: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeTemplateExtension(tt.path); got != tt.want {
				t.Errorf("removeTemplateExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultPathProcessor(t *testing.T) {
	processor := NewDefaultPathProcessor()
	if processor == nil {
		t.Error("NewDefaultPathProcessor() returned nil")
	}
}
