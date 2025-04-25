package generator

// GeneratedFile 表示一个生成的文件
type GeneratedFile struct {
	// 模板文件路径
	TemplatePath string
	// 目标文件路径
	OutputPath string
	// 生成的内容
	Content string
}
