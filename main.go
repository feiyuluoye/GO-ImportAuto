package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"myapp/module"
	"myapp/registry"
	"myapp/utils"
)

type Config struct {
	Modules []string                  `yaml:"modules"`
	Configs map[string]map[string]any `yaml:"configs"`
}

type ModuleManager struct {
	active map[string]module.Module
	lock   sync.Mutex
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{active: make(map[string]module.Module)}
}

func resolveDependencies(modNames []string) ([]string, error) {
	visited := make(map[string]bool)
	result := []string{}
	var visit func(string) error

	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		factory, ok := registry.Modules[name]
		if !ok {
			return fmt.Errorf("unknown module: %s", name)
		}
		tmp := factory() // 创建临时实例来获取依赖
		for _, dep := range tmp.Deps() {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[name] = true
		result = append(result, name)
		return nil
	}

	for _, m := range modNames {
		if err := visit(m); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (m *ModuleManager) Update(cfg Config) *gin.Engine {
	m.lock.Lock()
	defer m.lock.Unlock()

	ordered, err := resolveDependencies(cfg.Modules)
	if err != nil {
		fmt.Println("Dependency resolution error:", err)
		return gin.Default()
	}

	newActive := make(map[string]module.Module)
	r := gin.Default()

	// 启动新模块
	for _, name := range ordered {
		if old, exists := m.active[name]; exists {
			// 已存在，保留
			newActive[name] = old
			old.RegisterRoutes(r)
		} else if newFn, ok := registry.Modules[name]; ok {
			mod := newFn()
			modCfg := module.ModuleConfig(cfg.Configs[name])
			if err := mod.Init(modCfg); err != nil {
				fmt.Println("Failed to init module:", name, err)
				continue
			}
			mod.RegisterRoutes(r)
			newActive[name] = mod
			fmt.Println("Started module:", name)
		}
	}

	// 停止不再需要的模块（逆序）
	for i := len(ordered) - 1; i >= 0; i-- {
		name := ordered[i]
		if _, stillActive := newActive[name]; stillActive {
			continue
		}
		if old, exists := m.active[name]; exists {
			if err := old.Shutdown(); err != nil {
				fmt.Println("Error shutting down module:", name, err)
			} else {
				fmt.Println("Stopped module:", name)
			}
		}
	}

	m.active = newActive
	return r
}

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	
	newCfg := Config{
		Modules: cfg.Modules,
		Configs: map[string]map[string]any{},
	}
	for k, v := range cfg.Configs {
		expanded := utils.ExpandConfig(v)
		if m, ok := expanded.(map[string]any); ok {
			newCfg.Configs[k] = m
		}
	}
	return newCfg, nil
}

var (
	router       *gin.Engine
	manager      = NewModuleManager()
	globalRouter sync.Mutex
)

func rebuildRouter(cfg Config) {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	router = manager.Update(cfg)
}

func handler() *gin.Engine {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	return router
}

func main() {
	// 设置 Gin 模式
	devMode := os.Getenv("APP_ENV") == "dev"
	if devMode {
		gin.SetMode(gin.DebugMode)
		fmt.Println("[dev mode] Gin running in DebugMode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Gin running in ReleaseMode")
	}

	// 如果是 dump 模式
	if len(os.Args) > 1 && os.Args[1] == "dump" {
		cfg, err := loadConfig()
		if err != nil {
			log.Fatal("Failed to load config:", err)
		}
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(data))
		return
	}

	// 正常启动 Gin 服务
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	rebuildRouter(cfg)

	// 文件监控
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		if err := watcher.Add("config.yaml"); err != nil {
			log.Fatal(err)
		}

		// 提示 dev 模式
		if devMode {
			fmt.Println("[dev mode] Watching config.yaml ...")
		} else {
			fmt.Println("Watching config.yaml ...")
		}

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					fmt.Println("Config changed, reloading...")
					newCfg, err := loadConfig()
					if err != nil {
						fmt.Println("Error loading config:", err)
						continue
					}
					rebuildRouter(newCfg)
				}
			case err := <-watcher.Errors:
				fmt.Println("Watcher error:", err)
			}
		}
	}()

	// HTTP server
	ginEngine := gin.New()

	// 开发模式下启用 pprof
	if devMode {
		pprof.Register(ginEngine)
		fmt.Println("[dev mode] pprof enabled at /debug/pprof")
	}

	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
