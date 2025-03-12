## 改造过程

请帮助我重新规划下这个软件的升级计划：
我希望改造配置文件，只需要配置模板目录，配置文件目录，输出目录。请帮助我推理，看看这样改造会有什么问题。我们逐步确认改造计划。


模板的依赖关系暂时不考虑，暂时只支持单一模板解析和生成。
模板到输出的映射应该确定为，所有的模板解析配置文件，可以应用于所有模板文件。
模板的合并顺序，因为只考虑单一模板解析和生成了，所以无须合并了。
配置文件会读取，配置目录下的所有 yaml 文件作为配置文件，区别只是不同文件配置方便管理而已。
输出文件名的确定，输出文件名取自模板目录下，具体模板的相对路径，在输出目录下，同样的相对路径中生成。区别只是去掉了模板名后缀。请再次帮助我推理，看看这样改造会有什么问题。我们逐步确认改造计划。


## 改造计划
基于我们之前的讨论和您提供的当前 `main.go` 文件，我们可以制定以下详细的项目改造计划：

1. 更新配置结构
   **File: /home/lee/Projects/generator/internal/config/config.go**
   ```go
   type Config struct {
       TemplateDir string `yaml:"templateDir"`
       ConfigDir   string `yaml:"configDir"`
       OutputDir   string `yaml:"outputDir"`
   }
   ```

2. 实现配置加载函数
   **File: /home/lee/Projects/generator/internal/config/config.go**
   ```go
   func LoadConfig(workDir string) (*Config, error) {
       // 实现从工作目录加载配置的逻辑
       // 尝试读取 generator.yaml, generator.yml, config.yaml, config.yml
       // 如果找不到配置文件，返回默认配置
   }
   ```

3. 更新主程序 (main.go)
   **File: /home/lee/Projects/generator/cmd/v1/main.go**
   ```go
   func main() {
       workDir := flag.String("dir", ".", "工作目录路径")
       templateDir := flag.String("template", "", "模板目录路径")
       configDir := flag.String("config-dir", "", "配置目录路径")
       outputDir := flag.String("output", "", "输出目录路径")
       flag.Parse()

       absPath, err := os.Getwd()
       if err != nil {
           log.Fatalf("获取工作目录失败: %v", err)
       }

       if *workDir != "." {
           absPath = *workDir
       }

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

       gen := generate.NewGenerator()
       if err := gen.Generate(cfg); err != nil {
           log.Fatalf("生成失败: %v", err)
       }

       log.Println("生成完成")
   }
   ```

4. 更新 Generator 结构体和 Generate 方法
   **File: /home/lee/Projects/generator/generator.go**
   ```go
   type Generator struct {}

   func NewGenerator() *Generator {
       return &Generator{}
   }

   func (g *Generator) Generate(cfg *config.Config) error {
       // 实现生成逻辑
       // 1. 加载所有配置文件
       // 2. 遍历模板目录
       // 3. 对每个模板文件执行生成
       return nil
   }
   ```

5. 更新模板引擎
   **File: /home/lee/Projects/generator/internal/template/template.go**
   ```go
   type Engine struct {
       workDir string
       vars    map[string]interface{}
       templates map[string]*template.Template
   }

   func New(templateDir string) *Engine {
       // 初始化模板引擎
   }

   func (e *Engine) LoadConfig(configDir string) error {
       // 加载配置目录下的所有 YAML 文件
   }

   func (e *Engine) Execute(tplPath, outputPath string) error {
       // 执行单个模板生成
   }
   ```

6. 实现辅助函数
   - 在 `internal/utils` 目录下创建辅助函数，如文件操作、路径处理等。

7. 更新错误处理
   **File: /home/lee/Projects/generator/internal/errors/errors.go**
   ```go
   // 定义自定义错误类型和错误处理函数
   ```

8. 添加单元测试
   - 为 `config`, `template`, `generator` 等包添加单元测试。

9. 更新文档
   - 更新 README.md，说明新的使用方法和配置格式。

10. 创建示例
    - 在 `examples` 目录下创建示例项目，展示如何使用新的生成器。

具体步骤：

1. 更新 `config.go`
2. 更新 `main.go`
3. 更新 `generator.go`
4. 更新 `template.go`
5. 创建并实现辅助函数
6. 更新错误处理
7. 编写单元测试
8. 更新文档
9. 创建示例项目

每完成一个步骤，都应该进行测试，确保功能正常。完成所有步骤后，进行整体测试，确保生成器能够正常工作。

这个计划涵盖了我们讨论的所有要点，包括配置简化、命令行参数支持、默认配置文件查找等。它提供了一个清晰的路径来重构和改进您的代码生成器项目。您对这个计划有什么想法或需要进一步调整的地方吗？