package main

import (
	"flag"
	"log"
	"os"

	"generate"
)

func main() {
	workDir := flag.String("dir", ".", "工作目录路径")
	flag.Parse()

	// 获取工作目录的绝对路径
	absPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	if *workDir != "." {
		absPath = *workDir
	}

	gen := generate.NewGenerator()
	if err := gen.Generate(absPath); err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	log.Println("生成完成")
}