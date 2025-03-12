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
    templateDir := flag.String("template", "", "模板目录路径")
    configDir := flag.String("config", "", "配置目录路径")
    outputDir := flag.String("output", "", "输出目录路径")
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

    // 加载配置
    cfg, err := config.LoadConfig(absPath)
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }

    // 如果提供了命令行参数，覆盖配置
    if *templateDir != "" {
        cfg.TemplateDir = *templateDir
    }
    if *configDir != "" {
        cfg.ConfigDir = *configDir
    }
    if *outputDir != "" {
        cfg.OutputDir = *outputDir
    }

    // 确保所有路径都是绝对路径
    cfg.TemplateDir = filepath.Join(absPath, cfg.TemplateDir)
    cfg.ConfigDir = filepath.Join(absPath, cfg.ConfigDir)
    cfg.OutputDir = filepath.Join(absPath, cfg.OutputDir)

    gen := generate.NewGenerator()
    if err := gen.Generate(cfg); err != nil {
        log.Fatalf("生成失败: %v", err)
    }

    log.Println("生成完成")
}

func generateQuickStartExample() error {
    files := map[string]string{
        "generator.yaml": `templates:
  - path: "templates/example.tpl"
    output: "generated/example.txt"
    dependencies:
      - "configs/example.yaml"`,
        "configs/example.yaml": `greeting: "Hello"
name: "World"`,
        "templates/example.tpl": `{{ .greeting }}, {{ .name }}!`,
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

    fmt.Println("快速开始示例已生成。请运行 'generator' 命令来测试。")
    return nil
}