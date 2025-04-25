# Code Generator

A template-based code generator that supports generating code from configuration files.

## Features

- Support for multiple templates and configuration files
- Variable substitution
- Error reporting (including file paths and line numbers)
- Support for variables in output paths
- Quick start example generation
- Support for multiple variable files
- **Sub-templates: Allow templates to include other templates, enabling template reuse and modularization.**
- **Template path variables: Support for variable references in output paths, e.g., `__variable__`, for more flexible file organization.**
- **Skip child template generation: Automatically skip files with `__child__` in the template file path to avoid generating unnecessary child template files.**
- **Skip templates by suffix: Skip template files with specific suffixes, e.g., `.go.tpl.tpl`, to selectively generate certain types of files.**
- **Skip templates by prefix: Skip template files with specific path prefixes, e.g., `web/`, to selectively generate server-side or client-side code.**

## Getting Started

1. Ensure you have Go and the Go toolchain installed.
2. Clone the repository to your local machine.
3. Navigate to the `generator` directory.
4. Run `go build -o generator cmd/v1/main.go` to compile the generator.
5. Run `./generator -quickstart` to generate a quick start example.

## Usage

```
generator [options]

Options:
  -dir string
        Working directory path (default ".")
  -output string
        Output directory path (default ".gen_output")
  -quickstart
        Generate quick start example
  -template string
        Template directory path (default ".gen_templates")
  -variables string
        Variables directory path (default ".gen_variables")
  -varfiles string
        Variable files path, multiple files separated by commas
  -skip-suffixes string
        Skip template files with specific suffixes, multiple suffixes separated by commas
        Full path (path) is used for matching
        Example: -skip-suffixes=.go.tpl.tpl,.vue.tpl
  -skip-prefixes string
        Skip template files with specific path prefixes, multiple prefixes separated by commas
        Relative to the template directory, do not include leading / character
        Example: -skip-prefixes=web,server/config
```

### Examples

1. Generate a quick start example:

    ```
    ./generator -quickstart
    ```

2. Generate code with default configuration:

    ```
    ./generator
    ```

3. Specify a working directory:

    ```
    ./generator -dir /path/to/workdir
    ```

4. Customize template, variables, and output directories:

    ```
    ./generator -template /path/to/templates -variables /path/to/variables -output /path/to/output
    ```

5. Use multiple variable files:

    ```
    ./generator -varfiles file1.yaml,file2.yaml
    ```

6. Skip template files with specific suffixes:

    ```
    ./generator -skip-suffixes=.go.tpl.tpl,.vue.tpl
    ```

7. Generate only server-side code (skip web templates):

    ```
    ./generator -skip-prefixes=web
    ```

8. Generate only client-side code (skip server templates):

    ```
    ./generator -skip-prefixes=server
    ```

## Configuration Files

The generator uses YAML format configuration files to define templates and their dependencies.

## Using as a Library

The generator can be used as a library in Go projects. Import the `github.com/clh021/generator/pkg/generator` package and use the provided interfaces and functions.

### Basic Example

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
)

func main() {
	// Configure the generator
	cfg := &config.Config{
		TemplateDir:   "./templates",      // Template directory
		VariablesDir:  "./variables",      // Variables directory
		OutputDir:     "./output",         // Output directory
		VariableFiles: []string{           // Optional: specify additional variable files
			"./custom_variables.yaml",
		},
		SkipTemplateSuffixes: ".go.tpl.tpl,.vue.tpl",  // Optional: skip files with these suffixes
		SkipTemplatePrefixes: "web",                   // Optional: skip files with these path prefixes
	}

	// Create a generator instance with default components
	scanner := generator.NewDefaultTemplateScanner()
	filter := generator.NewDefaultTemplateFilter(true, cfg.SkipTemplateSuffixes, cfg.SkipTemplatePrefixes, cfg.TemplateDir)
	pathProcessor := generator.NewDefaultPathProcessor()
	contentGenerator := generator.NewDefaultContentGenerator()

	// Scan templates
	templateFiles, err := scanner.ScanTemplates(cfg.TemplateDir, filter)
	if err != nil {
		log.Fatalf("Failed to scan templates: %v", err)
	}

	// Create template engine
	engine := template.New(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)

	// Load variables
	variableLoader := generator.NewDefaultVariableLoader(cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir)
	variableFiles, err := variableLoader.FindVariableFiles(cfg.VariablesDir, cfg.VariableFiles)
	if err != nil {
		log.Fatalf("Failed to find variable files: %v", err)
	}

	if err := engine.LoadVariables(variableFiles); err != nil {
		log.Fatalf("Failed to load variables: %v", err)
	}

	variables := engine.GetVariables()

	// Process each template
	var generatedFiles []generator.GeneratedFile
	for _, templateFile := range templateFiles {
		// Process output path
		outputPath, err := pathProcessor.ProcessOutputPath(templateFile, cfg.OutputDir, variables)
		if err != nil {
			log.Printf("Warning: Failed to process output path: %v, using default path", err)
		}

		// Generate content
		content, err := contentGenerator.GenerateContent(templateFile, outputPath, engine)
		if err != nil {
			log.Fatalf("Failed to generate content: %v", err)
		}

		// Add to generated files list
		generatedFiles = append(generatedFiles, generator.GeneratedFile{
			TemplatePath: templateFile.Path,
			OutputPath:   outputPath,
			Content:      content,
		})
	}

	// Write generated files
	for _, file := range generatedFiles {
		// Create output directory
		outputDir := filepath.Dir(file.OutputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		// Create output file
		if err := os.WriteFile(file.OutputPath, []byte(file.Content), 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

		log.Printf("Generated file: %s", file.OutputPath)
	}

	log.Println("Generation completed")
}
```

### Modular Components

The generator provides several interfaces that can be implemented to customize the generation process:

#### 1. Template Scanner

The `TemplateScanner` interface is responsible for scanning the template directory and finding template files:

```go
// TemplateScanner defines the template scanner interface
type TemplateScanner interface {
	// ScanTemplates scans the template directory and returns a list of template files
	ScanTemplates(templateDir string, filter TemplateFilter) ([]TemplateFile, error)
}
```

Example of a custom template scanner:

```go
// CustomScanner is a custom template scanner
type CustomScanner struct {
	IncludePatterns []string
}

// ScanTemplates scans the template directory and returns a list of template files
func (s *CustomScanner) ScanTemplates(templateDir string, filter TemplateFilter) ([]generator.TemplateFile, error) {
	var templateFiles []generator.TemplateFile

	// Custom scanning logic...

	return templateFiles, nil
}
```

#### 2. Template Filter

The `TemplateFilter` interface is responsible for filtering template files:

```go
// TemplateFilter defines the template filter interface
type TemplateFilter interface {
	// ShouldInclude checks if a template file should be included
	// Returns: (should include, reason for exclusion)
	ShouldInclude(path, relativePath string) (bool, string)
}
```

Example of a custom template filter:

```go
// CustomFilter is a custom template filter
type CustomFilter struct {
	*generator.DefaultTemplateFilter
	AllowedExtensions []string
}

// ShouldInclude checks if a template file should be included
func (f *CustomFilter) ShouldInclude(path, relativePath string) (bool, string) {
	// First use the default filter
	include, reason := f.DefaultTemplateFilter.ShouldInclude(path, relativePath)
	if !include {
		return false, reason
	}

	// Then apply custom filtering logic
	// ...

	return true, ""
}
```

#### 3. Path Processor

The `PathProcessor` interface is responsible for processing output paths:

```go
// PathProcessor defines the path processor interface
type PathProcessor interface {
	// ProcessOutputPath processes the output path for a template file
	ProcessOutputPath(templateFile TemplateFile, outputDir string, variables map[string]interface{}) (string, error)
}
```

Example of a custom path processor:

```go
// CustomPathProcessor is a custom path processor
type CustomPathProcessor struct {
	*generator.DefaultPathProcessor
	PathPrefix string
}

// ProcessOutputPath processes the output path for a template file
func (p *CustomPathProcessor) ProcessOutputPath(templateFile generator.TemplateFile, outputDir string, variables map[string]interface{}) (string, error) {
	// Use the default processor
	path, err := p.DefaultPathProcessor.ProcessOutputPath(templateFile, outputDir, variables)
	if err != nil {
		return "", err
	}

	// Add custom prefix
	if p.PathPrefix != "" {
		path = filepath.Join(p.PathPrefix, path)
	}

	return path, nil
}
```

#### 4. Content Generator

The `ContentGenerator` interface is responsible for generating content from templates:

```go
// ContentGenerator defines the content generator interface
type ContentGenerator interface {
	// GenerateContent generates content for a template file
	GenerateContent(templateFile TemplateFile, outputPath string, engine *template.Engine) (string, error)
}
```

Example of a custom content generator:

```go
// CustomContentGenerator is a custom content generator
type CustomContentGenerator struct {
	AddGeneratedComment bool
	CommentPrefix       string
}

// GenerateContent generates content for a template file
func (g *CustomContentGenerator) GenerateContent(templateFile generator.TemplateFile, outputPath string, engine *template.Engine) (string, error) {
	// Generate content
	content, err := engine.GenerateContent(templateFile.Path, outputPath)
	if err != nil {
		return "", err
	}

	// Add generated comment
	if g.AddGeneratedComment {
		// Add comment based on file type
		// ...
	}

	return content, nil
}
```

#### 5. Variable Loader

The `VariableLoader` interface is responsible for loading variables:

```go
// VariableLoader defines the variable loader interface
type VariableLoader interface {
	// LoadVariables loads variables from the variables directory and additional files
	LoadVariables(variablesDir string, additionalFiles []string) (map[string]interface{}, error)
	// FindVariableFiles finds variable files in the variables directory and additional files
	FindVariableFiles(variablesDir string, additionalFiles []string) ([]string, error)
}
```

Example of a custom variable loader:

```go
// CustomVariableLoader is a custom variable loader
type CustomVariableLoader struct {
	*generator.DefaultVariableLoader
	ExtraVariables map[string]interface{}
}

// LoadVariables loads variables from the variables directory and additional files
func (l *CustomVariableLoader) LoadVariables(variablesDir string, additionalFiles []string) (map[string]interface{}, error) {
	// Load variables using the default loader
	variables, err := l.DefaultVariableLoader.LoadVariables(variablesDir, additionalFiles)
	if err != nil {
		return nil, err
	}

	// Add extra variables
	for k, v := range l.ExtraVariables {
		variables[k] = v
	}

	return variables, nil
}
```

### Simplified Approach

For a more simplified approach, you can create custom functions that encapsulate the generation process:

```go
// generateFiles generates files using a custom content generator
func generateFiles(cfg *config.Config, contentGenerator ContentGenerator) ([]generator.GeneratedFile, error) {
	var generatedFiles []generator.GeneratedFile

	// Template scanning, variable loading, and content generation logic...

	return generatedFiles, nil
}
```

This approach allows you to focus on the specific part of the generation process that you want to customize, while reusing the rest of the logic.

## Template Features

- Built-in string processing functions (`lcfirst`, `ucfirst`, `default`, `file`, `currentYear`, `dict`)
- Support for variables in output paths, e.g., `__variableName__`.
- **Sub-templates: Use `{{ include "path/to/sub_template.tpl" . }}` to include other template files in a template. Sub-templates can access variables from the parent template. Maximum nesting depth is limited to 2 levels to prevent circular references.**

## Sub-template Usage Instructions

1. **Path lookup:** Sub-template paths are first looked up as relative paths to the parent template. If an absolute path is specified, it is used directly.

2. **Circular references:** The nesting depth of sub-templates is limited to 2 levels, beyond which an error will be reported.

3. **Variable passing:** Sub-templates can access variables defined in the parent template.

4. **Sub-template naming:** To prevent sub-templates from being generated independently, include the string `__child__` in the sub-template file name or path. Template files containing `__child__` will be automatically skipped during generation with a notification. For example: `child__child__.tpl` or `__child__/template.tpl`.

## Error Handling

The generator provides detailed error reports, including file paths and line numbers.

## Contributing

Issues and pull requests are welcome.

## License

This project is licensed under the [MIT License](LICENSE).
