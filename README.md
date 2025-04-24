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

## Configuration Files

The generator uses YAML format configuration files to define templates and their dependencies.

## Using as a Library

The generator can also be used as a library in Go projects. Import the `generate/pkg/generator` package and call the `Generate` function.

### Example Code

```go
package main

import (
	"log"
	"path/filepath"

	"generate/pkg/generator"
	"generate/pkg/config"
)

func main() {
	// Create a generator instance
	gen := generator.NewGenerator()

	// Configure the generator
	cfg := &config.Config{
		TemplateDir:  "./templates",      // Template directory
		VariablesDir: "./variables",     // Variables directory
		OutputDir:    "./output",        // Output directory
		VariableFiles: []string{         // Optional: specify additional variable files
			"./custom_variables.yaml",
		},
	}

	// Execute generation
	if err := gen.Generate(cfg); err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	log.Println("Generation completed")
}
```

### Usage Instructions

1. **Create a generator instance**: Use `generator.NewGenerator()` to create a new generator instance.

2. **Configure the generator**: Create a `config.Config` struct and set the following fields:
   - `TemplateDir`: Template directory path
   - `VariablesDir`: Variables directory path
   - `OutputDir`: Output directory path
   - `VariableFiles`: (Optional) List of additional variable file paths

3. **Execute generation**: Call the `gen.Generate(cfg)` method to execute code generation.

The generator will automatically:
- Load all variable files
- Traverse all template files in the template directory
- Process variable references and sub-templates in the templates
- Save the generated files to the output directory

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
