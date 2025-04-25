package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/clh021/generator/pkg/config"
	"github.com/clh021/generator/pkg/generator"
)

var versionCmd = flag.NewFlagSet("version", flag.ExitOnError)

func main() {
	workDir := flag.String("dir", ".", "工作目录路径")
	templateDir := flag.String("template", ".gen_templates", "模板目录路径")
	variablesDir := flag.String("variables", ".gen_variables", "变量目录路径")
	outputDir := flag.String("output", ".gen_output", "输出目录路径")
	quickStart := flag.Bool("quickstart", false, "生成快速开始示例")
	variableFiles := flag.String("varfiles", "", "变量文件路径，多个文件用逗号分隔")

	// 定义 version 子命令
	if len(os.Args) > 1 && os.Args[1] == "version" {
		versionCmd.Parse(os.Args[2:])
		printVersion()
		os.Exit(0)
	}

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
		TemplateDir:   *templateDir,
		VariablesDir:  *variablesDir,
		OutputDir:     *outputDir,
		VariableFiles: []string{},
	}

	if *variableFiles != "" {
		cfg.VariableFiles = strings.Split(*variableFiles, ",")
	}
	// 如果提供了工作目录，则将路径调整为相对于工作目录
	if *workDir != "." {
		cfg.TemplateDir = filepath.Join(*workDir, *templateDir)
		cfg.VariablesDir = filepath.Join(*workDir, *variablesDir)
		cfg.OutputDir = filepath.Join(*workDir, *outputDir)
		for i, file := range cfg.VariableFiles {
			cfg.VariableFiles[i] = filepath.Join(*workDir, file)
		}
	}

	// 打印配置信息以便调试
	log.Printf("使用的配置：\n模板目录: %s\n变量目录: %s\n输出目录: %s\n变量文件: %v",
		cfg.TemplateDir, cfg.VariablesDir, cfg.OutputDir, cfg.VariableFiles)

	gen := generator.NewGenerator()
	if err := gen.Generate(cfg); err != nil {
		log.Fatalf("生成失败: %+v", err)
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
  output_dir: ".gen_output"`,
		".gen_variables/example.yaml": `greeting: "Hello"
name: "World"
# 特殊配置：允许模板中使用未定义的变量
# 当设置为 true 时，未定义的变量将被替换为零值（如空字符串）
# 当设置为 false 或未设置时，未定义的变量将导致错误
$config.allowUndefinedVariables: true`,
		".gen_variables/additional.yaml": `additional_var: "This is an additional variable"`,
		".gen_templates/example.txt.tpl": `{{ .greeting }}, {{ .name }}!

# 这是一个演示未定义变量的例子
未定义变量示例: {{ .undefinedVariable }}

# 如果 $config.allowUndefinedVariables 为 true，上面的行将显示为空
# 如果为 false，生成过程将因错误而停止

# 这是一个来自额外变量文件的变量
额外变量: {{ .additional_var }}`,
		".gen_templates/parent.txt.tpl": `This is the parent template.

# 这是一个相对路径引用的子模板
{{ include "child__child__.txt.tpl" . }}`,
		".gen_templates/child__child__.txt.tpl": `This is the child template.

# 从父模板传递的变量：
{{ .greeting }}`,
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
	fmt.Println("6. 在 .gen_variables/example.yaml 文件中，")
	fmt.Println("   '$config.allowUndefinedVariables: true' 允许模板中使用未定义的变量。")
	fmt.Println("   您可以将其设置为 false 或删除该行以恢复严格模式。")
	fmt.Println("7. 本示例中添加了一个额外的变量文件 .gen_variables/additional.yaml，")
	fmt.Println("   展示了如何使用多个变量文件。")
	fmt.Println("8. 您可以使用 -varfiles 参数指定额外的变量文件：")
	fmt.Println("   generator -varfiles .gen_variables/example.yaml,.gen_variables/additional.yaml")
	fmt.Println("9.  本示例中添加了子模板 parent.txt.tpl 和 child__child__.txt.tpl，")
	fmt.Println("    展示了如何使用子模板以及变量传递，以及子模板如何避免被独立生成。")
	fmt.Println("10. 由于 `child__child__.txt.tpl` 文件名中包含 `__child__`，所以它不会被独立生成。")
	fmt.Println("\n示例输出将生成在 .gen_output/example.txt 和 .gen_output/parent.txt")

	return nil
}
