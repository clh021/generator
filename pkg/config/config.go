package config

type Config struct {
	TemplateDir          string
	VariablesDir         string
	OutputDir            string
	VariableFiles        []string
	SkipTemplateSuffixes string // 要跳过的模板文件后缀，多个后缀用逗号分隔，完整路径(path)进行匹配
	SkipTemplatePrefixes string // 要跳过的模板路径前缀，多个前缀用逗号分隔，相对于模板目录，不要前置/符号
}
