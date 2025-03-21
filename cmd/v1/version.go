package main

import (
	"fmt"
	"runtime/debug"
)

// 这些变量会在编译时通过 -ldflags 注入
var (
	Version    string // 版本号，例如 0.0.123
	CommitID   string // Git Commit ID
	CommitTime string // Git Commit Time
	BuildTime  string // 构建时间
)

func printVersion() {
	fmt.Println("Version:", Version)
	fmt.Println("Commit ID:", CommitID)
	fmt.Println("Commit Time:", CommitTime)
	fmt.Println("Build Time:", BuildTime)

	// 打印 Go 版本信息
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("Go Version:", info.GoVersion)
		fmt.Println("Go Modules:")
		for _, dep := range info.Deps {
			fmt.Printf("  %s %s\n", dep.Path, dep.Version)
		}
	} else {
		fmt.Println("无法获取 Go 版本信息")
	}
}
