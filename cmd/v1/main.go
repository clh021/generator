package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"generate"
	"generate/internal/config"
)

func main() {
	workDir := flag.String("dir", ".", "工作目录路径")
	templateDir := flag.String("template", ".gen_templates", "模板目录路径")
	variablesDir := flag.String("variables", ".gen_variables", "变量目录路径")
	outputDir := flag.String("output", ".gen_output", "输出目录路径")
	quickStart := flag.Bool("quickstart", false, "生成快速开始示例")

	flag.Parse()

	if *quickStart {
		if err := generateQuickStartExample(); err != nil {
			log.Fatalf("生成快速开始示例失败: %v", err)
		}
		return
	}

	// 获取工作目录的绝对路径
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	if *workDir != "." {
		absPath = *workDir
	}

	// 检查默认变量目录是否存在
	defaultVariablesPath := filepath.Join(absPath, *variablesDir)
	if _, err := os.Stat(defaultVariablesPath); os.IsNotExist(err) && flag.NFlag() == 0 {
		fmt.Println("未找到默认变量目录，且未提供任何参数。")
		printHelp()
		os.Exit(1)
	}

	// 创建配置
	cfg := &config.Config{
		TemplateDir:  filepath.Join(absPath, *templateDir),
		VariablesDir: filepath.Join(absPath, *variablesDir),
		OutputDir:    filepath.Join(absPath, *outputDir),
	}

	gen := generate.NewGenerator()
	if err := gen.Generate(cfg); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	log.Println("生成完成")
}

func printHelp() {
	fmt.Println("使用方法: generator [选项]")
	fmt.Println("\n选项:")
	flag.PrintDefaults()
	fmt.Println("\n示例:")
	fmt.Println("  generator -quickstart                # 生成快速开始示例")
	fmt.Println("  generator -dir /path/to/workdir      # 指定工作目录")
	fmt.Println("  generator -template /path/to/templates -variables /path/to/variables -output /path/to/output") //

}

func generateQuickStartExample() error {
	files := map[string]string{
		".gen_config.yaml": `config:
  template_dir: ".gen_templates"
  variables_dir: ".gen_variables"
  output_dir: ".gen_output"
templates:
  - path: "example.tpl"
    output: "example.txt"
    variables: "example.yaml"`,
		".gen_variables/example.yaml": `greeting: "Hello"
name: "World"`,
		".gen_templates/example.tpl": `{{ .greeting }}, {{ .name }}!`,
	}

	fmt.Println("将要生成以下文件:")
	for path := range files {
		fmt.Printf("- %s\n", path)
	}

	fmt.Print("是否继续? (y/n): ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("操作已取消")
		return nil
	}

	for path, content := range files {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", path, err)
		}
	}

	fmt.Println("\n快速开始示例已生成。")
	fmt.Println("\n使用说明:")
	fmt.Println("1. 这个示例使用了默认的目录结构:")
	fmt.Println("   - 模板目录: .gen_templates")
	fmt.Println("   - 变量目录: .gen_variables")
	fmt.Println("   - 输出目录: .gen_output (将在生成时创建)")
	fmt.Println("2. 您可以直接运行 'generator' 命令来测试，无需额外参数。")
	fmt.Println("3. 如果您想自定义配置，可以修改 .gen_config.yaml 文件。")
	fmt.Println("4. 如果您不需要自定义配置，可以安全地删除 .gen_config.yaml 文件。")
	fmt.Println("5. 你也可以通过命令行使用以下参数进行配置:")
	fmt.Println("   generator -template <模板目录> -variables <变量目录> -output <输出目录>")
	fmt.Println("\n示例输出将生成在 .gen_output/example.txt")

	return nil
}
