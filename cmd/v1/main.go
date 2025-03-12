package main

import (
	"flag"
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
    flag.Parse()

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