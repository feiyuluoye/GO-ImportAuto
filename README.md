## Go Gin 基于配置文件自动导入模块包



[TOC]



1. **配置文件驱动**：在配置文件中定义要启用的功能模块。
2. **模块注册机制**：Go 本身没有像 Python 那样的动态 import，所以一般通过 **工厂模式** 或 **注册表模式** 来“自动”加载。

------

### 实现思路

#### 1. 配置文件（例如 config.yaml）

```yaml
modules:
  - user
  - auth
  - order
```

#### 2. 定义模块接口

每个模块要有统一的初始化接口，比如：

```go
package module

import "github.com/gin-gonic/gin"

type Module interface {
    RegisterRoutes(r *gin.Engine)
}
```

#### 3. 模块实现（例如 user 模块）

```go
package user

import (
    "github.com/gin-gonic/gin"
    "myapp/module"
)

type UserModule struct{}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
    r.GET("/user", func(c *gin.Context) {
        c.JSON(200, gin.H{"msg": "user module"})
    })
}

// 工厂函数
func New() module.Module {
    return &UserModule{}
}
```

#### 4. 模块注册表

可以通过 `map[string]func() Module` 来管理模块：

```go
package registry

import (
    "myapp/module"
    "myapp/module/user"
    "myapp/module/auth"
    "myapp/module/order"
)

var Modules = map[string]func() module.Module{
    "user":  user.New,
    "auth":  auth.New,
    "order": order.New,
}
```

#### 5. 主程序根据配置文件动态加载

```go
package main

import (
    "fmt"
    "log"

    "github.com/gin-gonic/gin"
    "gopkg.in/yaml.v3"
    "os"

    "myapp/registry"
)

type Config struct {
    Modules []string `yaml:"modules"`
}

func main() {
    // 读取配置
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatal(err)
    }

    r := gin.Default()

    // 按配置加载模块
    for _, m := range cfg.Modules {
        if newFn, ok := registry.Modules[m]; ok {
            mod := newFn()
            mod.RegisterRoutes(r)
            fmt.Println("Loaded module:", m)
        } else {
            fmt.Println("Unknown module:", m)
        }
    }

    r.Run(":8080")
}
```

这样就实现了 **“根据配置文件自动导入模块包”**。好的 👍 我给你一个最小可运行的 **Gin + 配置驱动模块加载** 项目结构。
 你可以复制整个目录结构，`go run main.go` 就能跑起来。

------

## 实验项目目录结构

```
myapp/
├── go.mod
├── main.go
├── config.yaml
├── module/
│   └── module.go
├── registry/
│   └── registry.go
└── modules/
    ├── user/
    │   └── user.go
    ├── auth/
    │   └── auth.go
    └── order/
        └── order.go
```

### 测试代码

### go.mod

```go
module myapp

go 1.21

require (
    github.com/gin-gonic/gin v1.10.0
    gopkg.in/yaml.v3 v3.0.1
)
```

------

### config.yaml

```yaml
modules:
  - user
  - auth
```

------

### module/module.go

```go
package module

import "github.com/gin-gonic/gin"

// 所有模块必须实现的接口
type Module interface {
	RegisterRoutes(r *gin.Engine)
}
```

------

### modules/user/user.go

```go
package user

import (
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type UserModule struct{}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from user module"})
	})
}

func New() module.Module {
	return &UserModule{}
}
```

------

### modules/auth/auth.go

```go
package auth

import (
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type AuthModule struct{}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from auth module"})
	})
}

func New() module.Module {
	return &AuthModule{}
}
```

------

### modules/order/order.go

```go
package order

import (
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type OrderModule struct{}

func (m *OrderModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/order", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from order module"})
	})
}

func New() module.Module {
	return &OrderModule{}
}
```

------

### registry/registry.go

```go
package registry

import (
	"myapp/module"
	"myapp/modules/auth"
	"myapp/modules/order"
	"myapp/modules/user"
)

// 注册表：模块名 -> 工厂函数
var Modules = map[string]func() module.Module{
	"user":  user.New,
	"auth":  auth.New,
	"order": order.New,
}
```

------

### main.go

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"myapp/registry"
)

type Config struct {
	Modules []string `yaml:"modules"`
}

func main() {
	// 读取配置文件
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// 根据配置加载模块
	for _, m := range cfg.Modules {
		if newFn, ok := registry.Modules[m]; ok {
			mod := newFn()
			mod.RegisterRoutes(r)
			fmt.Println("Loaded module:", m)
		} else {
			fmt.Println("Unknown module:", m)
		}
	}

	// 启动服务
	r.Run(":8080")
}
```

------

## 运行

```bash
cd myapp
go mod tidy
go run main.go
```

然后访问：

- `http://localhost:8080/user` → 返回 `"Hello from user module"`
- `http://localhost:8080/auth` → 返回 `"Hello from auth module"`
- `http://localhost:8080/order` → 不会生效（因为没在 `config.yaml` 里启用）

------

好嘞 👍 我帮你把 **热加载 config.yaml 动态启停模块** 的版本写出来。
 我们用 `fsnotify` 监听配置文件变化，当文件被修改时，重新加载配置，重建路由。

------

## 目录结构（和之前差不多）

```
myapp/
├── go.mod
├── main.go
├── config.yaml
├── module/
│   └── module.go
├── registry/
│   └── registry.go
└── modules/
    ├── user/user.go
    ├── auth/auth.go
    └── order/order.go
```

------

## go.mod

```go
module myapp

go 1.21

require (
    github.com/fsnotify/fsnotify v1.7.0
    github.com/gin-gonic/gin v1.10.0
    gopkg.in/yaml.v3 v3.0.1
)
```

------

## config.yaml

```yaml
modules:
  - user
  - auth
```

------

## main.go（核心改造）

```go
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"myapp/registry"
)

type Config struct {
	Modules []string `yaml:"modules"`
}

var (
	currentModules []string
	router         *gin.Engine
	lock           sync.Mutex
)

// 读取配置
func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// 根据配置重新构建路由
func rebuildRouter(cfg Config) {
	lock.Lock()
	defer lock.Unlock()

	r := gin.Default()

	for _, m := range cfg.Modules {
		if newFn, ok := registry.Modules[m]; ok {
			mod := newFn()
			mod.RegisterRoutes(r)
			fmt.Println("Loaded module:", m)
		} else {
			fmt.Println("Unknown module:", m)
		}
	}

	router = r
	currentModules = cfg.Modules
	fmt.Println("Router rebuilt with modules:", currentModules)
}

// Gin Handler 委托：始终转发到最新 router
func handler() *gin.Engine {
	lock.Lock()
	defer lock.Unlock()
	return router
}

func main() {
	// 初始加载配置
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	rebuildRouter(cfg)

	// 启动监听 config.yaml 的协程
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		if err := watcher.Add("config.yaml"); err != nil {
			log.Fatal(err)
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					fmt.Println("Config file changed, reloading...")
					newCfg, err := loadConfig()
					if err != nil {
						fmt.Println("Error reloading config:", err)
						continue
					}
					rebuildRouter(newCfg)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Watcher error:", err)
			}
		}
	}()

	// 使用单一 http.Server，始终从 handler() 获取路由
	ginEngine := gin.New()
	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
```

------

## 使用效果

1. 启动服务：

   ```bash
   go run main.go
   ```

2. 初始访问：

   - `http://localhost:8080/user` ✅
   - `http://localhost:8080/auth` ✅
   - `http://localhost:8080/order` ❌（未启用）

3. 修改 `config.yaml`：

   ```yaml
   modules:
     - user
     - order
   ```

4. 保存文件 → 控制台输出：

   ```
   Config file changed, reloading...
   Loaded module: user
   Loaded module: order
   Router rebuilt with modules: [user order]
   ```

5. 再次访问：

   - `http://localhost:8080/user` ✅
   - `http://localhost:8080/order` ✅
   - `http://localhost:8080/auth` ❌（已禁用）

------

要不要我帮你升级一下，让 **模块支持初始化/销毁钩子**（比如数据库连接、定时任务，在热加载时能优雅关闭旧模块）？

好，这一步我们就让每个模块支持 **初始化 (Init)** 和 **销毁 (Shutdown)** 钩子，这样在热加载时就能优雅释放资源（比如数据库连接、定时任务、goroutine）。

------

## 改造思路

1. **模块接口扩展**：
   - `Init()`：模块启动时调用。
   - `Shutdown()`：模块被卸载时调用。
2. **模块管理器 (ModuleManager)**：
   - 维护一个当前已加载的模块集合。
   - 支持增量更新（新增模块 → Init，移除模块 → Shutdown，保留的不动）。
3. **热加载时**：
   - 解析新配置。
   - 比对新旧模块列表。
   - 调用对应钩子。
   - 重建 Gin 路由。

------

## 修改后的代码

### module/module.go

```go
package module

import "github.com/gin-gonic/gin"

// 模块接口：支持生命周期钩子
type Module interface {
	Init() error                  // 模块初始化
	RegisterRoutes(r *gin.Engine) // 注册路由
	Shutdown() error              // 模块销毁（关闭资源）
}
```

------

### modules/user/user.go

```go
package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type UserModule struct{}

func (m *UserModule) Init() error {
	fmt.Println("[user] Init")
	return nil
}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from user module"})
	})
}

func (m *UserModule) Shutdown() error {
	fmt.Println("[user] Shutdown")
	return nil
}

func New() module.Module {
	return &UserModule{}
}
```

------

### modules/auth/auth.go

```go
package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type AuthModule struct{}

func (m *AuthModule) Init() error {
	fmt.Println("[auth] Init")
	return nil
}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from auth module"})
	})
}

func (m *AuthModule) Shutdown() error {
	fmt.Println("[auth] Shutdown")
	return nil
}

func New() module.Module {
	return &AuthModule{}
}
```

------

### modules/order/order.go

```go
package order

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type OrderModule struct{}

func (m *OrderModule) Init() error {
	fmt.Println("[order] Init (e.g. open DB connection)")
	return nil
}

func (m *OrderModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/order", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from order module"})
	})
}

func (m *OrderModule) Shutdown() error {
	fmt.Println("[order] Shutdown (e.g. close DB connection)")
	return nil
}

func New() module.Module {
	return &OrderModule{}
}
```

------

### main.go

```go
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"myapp/module"
	"myapp/registry"
)

type Config struct {
	Modules []string `yaml:"modules"`
}

// 模块管理器
type ModuleManager struct {
	active map[string]module.Module
	lock   sync.Mutex
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{active: make(map[string]module.Module)}
}

func (m *ModuleManager) Update(modNames []string) *gin.Engine {
	m.lock.Lock()
	defer m.lock.Unlock()

	newActive := make(map[string]module.Module)
	r := gin.Default()

	// 启用新模块
	for _, name := range modNames {
		if old, exists := m.active[name]; exists {
			// 已存在的模块，保留
			newActive[name] = old
			old.RegisterRoutes(r)
		} else if newFn, ok := registry.Modules[name]; ok {
			mod := newFn()
			if err := mod.Init(); err != nil {
				fmt.Println("Failed to init module:", name, err)
				continue
			}
			mod.RegisterRoutes(r)
			newActive[name] = mod
			fmt.Println("Started module:", name)
		} else {
			fmt.Println("Unknown module:", name)
		}
	}

	// 停止不再需要的模块
	for name, mod := range m.active {
		if _, stillActive := newActive[name]; !stillActive {
			if err := mod.Shutdown(); err != nil {
				fmt.Println("Error shutting down module:", name, err)
			} else {
				fmt.Println("Stopped module:", name)
			}
		}
	}

	m.active = newActive
	return r
}

// ===============
// 配置 & 热加载
// ===============
var (
	router       *gin.Engine
	manager      = NewModuleManager()
	globalRouter sync.Mutex
)

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func rebuildRouter(cfg Config) {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	router = manager.Update(cfg.Modules)
}

func handler() *gin.Engine {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	return router
}

func main() {
	// 初始加载
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
	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
```

------

## 效果演示

1. 启动：

   ```
   go run main.go
   ```

   控制台输出：

   ```
   [user] Init
   Started module: user
   [auth] Init
   Started module: auth
   ```

2. 修改 `config.yaml` → 去掉 `auth`，加上 `order`：

   ```yaml
   modules:
     - user
     - order
   ```

   控制台输出：

   ```
   Config changed, reloading...
   [order] Init (e.g. open DB connection)
   Started module: order
   [auth] Shutdown
   Stopped module: auth
   ```

3. 路由也随之变化。

------

要不要我再帮你升级一下，让 **模块支持依赖注入**（比如 `order` 依赖 `auth`，自动按依赖顺序 Init/Shutdown）？



好嘞 🚀 我们来给模块系统加上 **依赖注入 (DI)**，让模块能够声明自己的依赖关系，然后由 **ModuleManager** 自动解析依赖顺序，在热加载时正确 Init / Shutdown。

------

## 改造思路

1. **模块接口增加 `Deps()` 方法**
   - 每个模块可以返回一个字符串数组，表示它依赖哪些模块。
   - 例如：`order` 模块 → `[]string{"auth"}`
2. **ModuleManager 初始化逻辑修改**
   - 读取用户配置的模块列表。
   - 递归解析依赖关系（拓扑排序）。
   - 自动保证先初始化依赖，再初始化当前模块。
3. **Shutdown 逻辑修改**
   - 停止模块时，按 **依赖的逆序** 关闭。
   - 保证被依赖的模块最后关闭。

------

## 修改后的代码

### module/module.go

```go
package module

import "github.com/gin-gonic/gin"

// 模块接口：支持生命周期 & 依赖声明
type Module interface {
	Deps() []string               // 模块依赖哪些其他模块
	Init() error                  // 模块初始化
	RegisterRoutes(r *gin.Engine) // 注册路由
	Shutdown() error              // 模块销毁（释放资源）
}
```

------

### modules/auth/auth.go

```go
package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type AuthModule struct{}

func (m *AuthModule) Deps() []string { return nil }

func (m *AuthModule) Init() error {
	fmt.Println("[auth] Init")
	return nil
}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from auth module"})
	})
}

func (m *AuthModule) Shutdown() error {
	fmt.Println("[auth] Shutdown")
	return nil
}

func New() module.Module {
	return &AuthModule{}
}
```

------

### modules/user/user.go

```go
package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type UserModule struct{}

func (m *UserModule) Deps() []string { return nil }

func (m *UserModule) Init() error {
	fmt.Println("[user] Init")
	return nil
}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from user module"})
	})
}

func (m *UserModule) Shutdown() error {
	fmt.Println("[user] Shutdown")
	return nil
}

func New() module.Module {
	return &UserModule{}
}
```

------

### modules/order/order.go

```go
package order

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type OrderModule struct{}

func (m *OrderModule) Deps() []string { return []string{"auth"} }

func (m *OrderModule) Init() error {
	fmt.Println("[order] Init (depends on auth)")
	return nil
}

func (m *OrderModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/order", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello from order module"})
	})
}

func (m *OrderModule) Shutdown() error {
	fmt.Println("[order] Shutdown")
	return nil
}

func New() module.Module {
	return &OrderModule{}
}
```

------

### main.go（重点改造 ModuleManager）

```go
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"myapp/module"
	"myapp/registry"
)

type Config struct {
	Modules []string `yaml:"modules"`
}

// 模块管理器
type ModuleManager struct {
	active map[string]module.Module
	lock   sync.Mutex
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{active: make(map[string]module.Module)}
}

// 解析依赖顺序（拓扑排序）
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

func (m *ModuleManager) Update(modNames []string) *gin.Engine {
	m.lock.Lock()
	defer m.lock.Unlock()

	ordered, err := resolveDependencies(modNames)
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
			if err := mod.Init(); err != nil {
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

// ===============
// 配置 & 热加载
// ===============
var (
	router       *gin.Engine
	manager      = NewModuleManager()
	globalRouter sync.Mutex
)

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func rebuildRouter(cfg Config) {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	router = manager.Update(cfg.Modules)
}

func handler() *gin.Engine {
	globalRouter.Lock()
	defer globalRouter.Unlock()
	return router
}

func main() {
	// 初始加载
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
	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
```

------

## 测试效果

### config.yaml 初始：

```yaml
modules:
  - order
```

### 启动：

```
[auth] Init
Started module: auth
[order] Init (depends on auth)
Started module: order
```

> 即使配置里只写了 `order`，也会自动加载 `auth`。

### 修改 config.yaml：

```yaml
modules:
  - user
```

### 输出：

```
Config changed, reloading...
[user] Init
Started module: user
[order] Shutdown
Stopped module: order
[auth] Shutdown
Stopped module: auth
```

> 先关掉依赖 `order`，再关掉被依赖的 `auth`。

------

要不要我再帮你升级一下，让模块支持 **配置注入**（比如 `order` 模块可以从 `config.yaml` 里拿到数据库连接参数）？

好！接下来我们给模块系统加上 **配置注入** 🛠️，这样每个模块不仅可以声明依赖，还可以从 `config.yaml` 中读取自己需要的配置（例如数据库连接参数、缓存设置、API Key 等）。

------

## 改造思路

1. **扩展模块接口**

   - 在 `Init(cfg ModuleConfig)` 里传入模块自己的配置块。
   - 模块可以解析并使用这些参数。

2. **配置文件结构升级**

   - 在 `config.yaml` 中，除了模块启用列表，还支持 `configs` 区块，里面写每个模块的配置。

   例如：

   ```yaml
   modules:
     - user
     - order
   
   configs:
     order:
       dsn: "root:123456@tcp(localhost:3306)/shop"
     user:
       greeting: "Hi from user module"
   ```

3. **ModuleManager 负责把对应配置传给模块**

   - `Update` 时，把 `configs[name]` 传给 `Init(cfg)`。
   - 如果没有配置，就传空 map。

------

## 修改后的代码

### module/module.go

```go
package module

import "github.com/gin-gonic/gin"

// 通用配置类型
type ModuleConfig map[string]any

// 模块接口：支持生命周期、依赖 & 配置注入
type Module interface {
	Deps() []string
	Init(cfg ModuleConfig) error
	RegisterRoutes(r *gin.Engine)
	Shutdown() error
}
```

------

### config.yaml（示例）

```yaml
modules:
  - user
  - order

configs:
  user:
    greeting: "Hi from user module"
  order:
    dsn: "root:123456@tcp(localhost:3306)/shop"
```

------

### modules/user/user.go

```go
package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type UserModule struct {
	greeting string
}

func (m *UserModule) Deps() []string { return nil }

func (m *UserModule) Init(cfg module.ModuleConfig) error {
	if g, ok := cfg["greeting"].(string); ok {
		m.greeting = g
	} else {
		m.greeting = "Hello from user (default)"
	}
	fmt.Println("[user] Init with greeting =", m.greeting)
	return nil
}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": m.greeting})
	})
}

func (m *UserModule) Shutdown() error {
	fmt.Println("[user] Shutdown")
	return nil
}

func New() module.Module {
	return &UserModule{}
}
```

------

### modules/order/order.go

```go
package order

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"myapp/module"
)

type OrderModule struct {
	dsn string
}

func (m *OrderModule) Deps() []string { return []string{"auth"} }

func (m *OrderModule) Init(cfg module.ModuleConfig) error {
	if dsn, ok := cfg["dsn"].(string); ok {
		m.dsn = dsn
	} else {
		m.dsn = "memory://default"
	}
	fmt.Println("[order] Init with DSN =", m.dsn)
	return nil
}

func (m *OrderModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/order", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Order module using DSN: " + m.dsn})
	})
}

func (m *OrderModule) Shutdown() error {
	fmt.Println("[order] Shutdown")
	return nil
}

func New() module.Module {
	return &OrderModule{}
}
```

------

### main.go（只贴关键改动）

```go
type Config struct {
	Modules []string                       `yaml:"modules"`
	Configs map[string]map[string]any      `yaml:"configs"`
}

func (m *ModuleManager) Update(cfg Config) *gin.Engine {
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
	for oldName, oldMod := range m.active {
		if _, still := newActive[oldName]; !still {
			if err := oldMod.Shutdown(); err != nil {
				fmt.Println("Error shutting down module:", oldName, err)
			} else {
				fmt.Println("Stopped module:", oldName)
			}
		}
	}

	m.active = newActive
	return r
}
```

------

## 效果演示

1. 启动：

   ```
   go run main.go
   ```

   控制台输出：

   ```
   [auth] Init
   Started module: auth
   [order] Init with DSN = root:123456@tcp(localhost:3306)/shop
   Started module: order
   [user] Init with greeting = Hi from user module
   Started module: user
   ```

2. 访问接口：

   - `http://localhost:8080/user` → `{ "msg": "Hi from user module" }`
   - `http://localhost:8080/order` → `{ "msg": "Order module using DSN: root:123456@tcp(localhost:3306)/shop" }`

3. 修改配置热加载：

   - 改 `user.greeting` → `"Hello, World!"`
   - 保存文件 → 控制台输出 `[user] Init with greeting = Hello, World!`

------

要不要我再帮你升级一下，让配置支持 **环境变量替换**（比如 `${DB_DSN}` 自动从环境变量注入到模块配置）？

好👌 我们来给配置系统加一个 **环境变量替换** 功能：

- 在 `config.yaml` 里写 `${ENV_NAME}`，加载时自动替换为系统环境变量的值。
- 如果环境变量不存在，可以保留原值，或者支持默认值语法 `${ENV_NAME:default}`。

------

## 改造思路

1. **增加一个函数 `expandEnv`**
   - 用正则匹配 `${VAR}` 或 `${VAR:default}`。
   - 优先取 `os.Getenv(VAR)`，若为空则用 `default`，再为空就保留原样。
2. **在 `loadConfig()` 后统一展开**
   - 将整个配置结构（map、slice、string）递归展开。
   - 保证所有 string 类型的配置值都支持环境变量替换。

------

## 代码实现

### utils/env.go

```go
package utils

import (
	"os"
	"regexp"
)

// ${VAR} 或 ${VAR:default}
var envPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(?::([^}]+))?\}`)

func ExpandEnv(s string) string {
	return envPattern.ReplaceAllStringFunc(s, func(m string) string {
		groups := envPattern.FindStringSubmatch(m)
		if len(groups) < 2 {
			return m
		}
		key := groups[1]
		def := ""
		if len(groups) > 2 {
			def = groups[2]
		}
		if val := os.Getenv(key); val != "" {
			return val
		}
		if def != "" {
			return def
		}
		return m
	})
}

// 递归展开配置
func ExpandConfig(v any) any {
	switch val := v.(type) {
	case string:
		return ExpandEnv(val)
	case map[string]any:
		newMap := make(map[string]any)
		for k, v2 := range val {
			newMap[k] = ExpandConfig(v2)
		}
		return newMap
	case []any:
		newSlice := make([]any, len(val))
		for i, v2 := range val {
			newSlice[i] = ExpandConfig(v2)
		}
		return newSlice
	default:
		return v
	}
}
```

------

### main.go 中改造 loadConfig

```go
import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
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

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	// 环境变量展开
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
```

------

### config.yaml 示例

```yaml
modules:
  - order
  - user

configs:
  order:
    dsn: "${DB_DSN:mysql://root:123@tcp(localhost:3306)/shop}"
  user:
    greeting: "${USER_GREETING:Hello, Default User!}"
```

------

### 运行效果

假设环境变量：

```bash
export DB_DSN="postgres://admin:pwd@localhost:5432/mydb"
```

启动：

```
[auth] Init
Started module: auth
[order] Init with DSN = postgres://admin:pwd@localhost:5432/mydb
Started module: order
[user] Init with greeting = Hello, Default User!
Started module: user
```

如果 `USER_GREETING` 没有设置，则使用默认值 `"Hello, Default User!"`。

------

要不要我再帮你做一个小工具，比如 `config dump`，可以打印出 **展开环境变量后的完整配置**，方便调试？

好 👌 我们来做一个 `config dump` 工具，可以打印出 **展开环境变量后的完整配置**，方便调试。

------

## 实现思路

1. 在 `main.go` 的入口加一个命令行参数解析（比如 `go run main.go dump`）。
2. 如果参数是 `dump`，就加载配置并打印 JSON（带环境变量替换的结果）。
3. 否则正常启动 Gin 服务。

------

## 修改 main.go（关键部分）

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
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

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	// 环境变量展开
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
	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
```

------

## 使用方法

### 1. 启动服务

```bash
go run main.go
```

### 2. 打印展开后的配置

```bash
go run main.go dump
```

输出示例：

```json
{
  "Modules": [
    "user",
    "order"
  ],
  "Configs": {
    "order": {
      "dsn": "postgres://admin:pwd@localhost:5432/mydb"
    },
    "user": {
      "greeting": "Hello, Default User!"
    }
  }
}
```

------

要不要我帮你加个 `makefile`，让你能直接用 `make run` / `make dump` 来跑，而不用手写 `go run main.go ...`？

好 👌 我们来写一个简单的 `Makefile`，支持：

- `make run` → 启动服务
- `make dump` → 打印展开后的配置
- 另外顺手加上 `make tidy`、`make build`

------

## Makefile 示例

```makefile
APP_NAME := myapp
MAIN := main.go

.PHONY: run dump tidy build

# 启动服务
run:
	go run $(MAIN)

# 打印展开后的配置
dump:
	go run $(MAIN) dump

# 整理依赖
tidy:
	go mod tidy

# 编译可执行文件
build:
	go build -o $(APP_NAME) $(MAIN)
```

------

## 使用方法

1. **启动服务**

   ```bash
   make run
   ```

2. **打印展开后的配置**

   ```bash
   make dump
   ```

3. **整理 go.mod**

   ```bash
   make tidy
   ```

4. **编译**

   ```bash
   make build
   ./myapp
   ```

------

要不要我再帮你加一个 `make dev`，结合 `air` 或 `reflex` 实现 **代码热重载**（改 go 文件就自动重启服务）？

好👌 我们加一个 `make dev`，用 **air** 或 **reflex** 来做 Go 代码热重载。

------

## 方案 1：用 [**air**](https://github.com/cosmtrek/air)

`air` 是 Go 社区常用的热重载工具，配置文件简单，体验好。

### 1. 安装 air

```bash
go install github.com/cosmtrek/air@latest
```

> 安装后 `$(go env GOPATH)/bin` 要在 PATH 里。

------

### 2. 新增 `.air.toml`（可选）

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/myapp main.go"
  bin = "tmp/myapp"
  full_bin = "./tmp/myapp"
  include_ext = ["go", "yaml"]
  exclude_dir = ["tmp", "vendor"]

[log]
  time = true
```

------

### 3. 修改 `Makefile`

```makefile
APP_NAME := myapp
MAIN := main.go

.PHONY: run dump tidy build dev

run:
	go run $(MAIN)

dump:
	go run $(MAIN) dump

tidy:
	go mod tidy

build:
	go build -o $(APP_NAME) $(MAIN)

# 开发热重载（用 air）
dev:
	air -c .air.toml
```

------

## 方案 2：用 reflex

如果不想用 air，可以用轻量级的 [reflex](https://github.com/cespare/reflex)。

### 1. 安装 reflex

```bash
go install github.com/cespare/reflex@latest
```

------

### 2. 新增 `reflex.conf`

```conf
# 监听 go 和 yaml 文件变化，重启服务
-r '\.go$$' -r '\.yaml$$' go run main.go
```

------

### 3. 修改 Makefile

```makefile
dev:
	reflex -c reflex.conf
```

------

## 使用方法

```bash
make dev   # 开发模式，修改 .go / config.yaml 会自动重启
```

------

要不要我帮你在 **热加载 config.yaml** 的功能里，也加上 `make dev` 的集成提示？比如控制台打印 `[dev mode] Watching config.yaml ...`，让人一眼就知道当前是开发模式。

明白了 ✅
 我们来把 **`make dev` 热重载** 和 **配置文件热加载** 结合起来，加一个清晰的提示：

- 如果是用 `make dev` 启动（即 `air` / `reflex` 监控），控制台在启动时打印：

  ```
  [dev mode] Watching config.yaml ...
  ```

- 这样一眼就能区分出是生产还是开发模式。

------

## 实现方式

可以通过环境变量来区分运行模式：

- 在 `make dev` 里设置 `APP_ENV=dev`
- 在 `main.go` 判断 `APP_ENV` 是否为 `dev`，然后打印提示

------

### 修改 Makefile

```makefile
APP_NAME := myapp
MAIN := main.go

.PHONY: run dump tidy build dev

run:
	go run $(MAIN)

dump:
	go run $(MAIN) dump

tidy:
	go mod tidy

build:
	go build -o $(APP_NAME) $(MAIN)

# 开发热重载（用 air）
dev:
	APP_ENV=dev air -c .air.toml
```

------

### 修改 main.go（监控部分）

```go
func main() {
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
		if os.Getenv("APP_ENV") == "dev" {
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
	ginEngine.Any("/*path", func(c *gin.Context) {
		handler().ServeHTTP(c.Writer, c.Request)
	})

	ginEngine.Run(":8080")
}
```

------

## 效果

### 开发模式（`make dev`）：

```
[dev mode] Watching config.yaml ...
Config changed, reloading...
```

### 生产模式（`make run`）：

```
Watching config.yaml ...
Config changed, reloading...
```

------

要不要我帮你把 **生产/开发模式** 的区别再做彻底点？比如：

- 开发模式下开启 Gin 的 `DebugMode`
- 生产模式默认 `ReleaseMode`



好主意 👍
 在 **开发模式** 下开启 `pprof`，可以用来分析性能。Gin 官方其实提供了一个 [pprof 中间件](https://pkg.go.dev/github.com/gin-contrib/pprof)，用起来非常方便。

------

## 改造步骤

1. 安装依赖：

```bash
go get github.com/gin-contrib/pprof
```

1. 在 `main.go` 里开发模式时注册 `pprof` 路由。

------

### 修改 main.go（关键部分）

```go
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
```

------

### main 函数里增加 pprof

```go
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
		// ... 原有的 fsnotify 逻辑 ...
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
```

------

## 效果

### 开发模式（`make dev`）

控制台：

```
[dev mode] Gin running in DebugMode
[dev mode] Watching config.yaml ...
[dev mode] pprof enabled at /debug/pprof
```

浏览器访问：

- `http://localhost:8080/debug/pprof/` → 性能剖析工具页面
- 可以结合 `go tool pprof` 来分析

### 生产模式（`make run`）

```
Gin running in ReleaseMode
Watching config.yaml ...
```

👉 `pprof` 不会启用。

------

要不要我顺便帮你加一个 **make profile** 命令，自动跑 `go tool pprof` 连到 `/debug/pprof/profile`？

好 ✅ 我来帮你整理一份完整的 **工程项目文档**，涵盖我们上面一步步搭建的所有功能：

- 模块化加载
- 热加载 config.yaml
- 模块依赖管理
- 生命周期管理（Init/Shutdown）
- 配置注入 + 环境变量替换
- 开发/生产模式区分
- pprof 调试
- Makefile 脚本

------

# 📖 Go Gin 模块化工程模板

## 目录结构

```
myapp/
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── .air.toml         # air 配置（可选）
├── config.yaml       # 应用配置
├── module/           # 模块接口定义
│   └── module.go
├── registry/         # 模块注册表
│   └── registry.go
├── utils/            # 工具函数
│   └── env.go
└── modules/          # 模块实现
    ├── auth/
    │   └── auth.go
    ├── user/
    │   └── user.go
    └── order/
        └── order.go
```

------

## go.mod（示例）

```go
module myapp

go 1.22

require (
    github.com/fsnotify/fsnotify v1.7.0
    github.com/gin-contrib/pprof v1.5.1
    github.com/gin-gonic/gin v1.10.0
    gopkg.in/yaml.v3 v3.0.1
)
```

------

## config.yaml（示例）

```yaml
modules:
  - auth
  - user
  - order

configs:
  user:
    greeting: "${USER_GREETING:Hello, Default User!}"
  order:
    dsn: "${DB_DSN:mysql://root:123@tcp(localhost:3306)/shop}"
```

> - `modules`：定义启用哪些模块
> - `configs`：每个模块的独立配置，支持 `${ENV:default}` 环境变量替换

------

## module/module.go

```go
package module

import "github.com/gin-gonic/gin"

type ModuleConfig map[string]any

// 模块接口：支持依赖声明、配置注入、生命周期
type Module interface {
	Deps() []string
	Init(cfg ModuleConfig) error
	RegisterRoutes(r *gin.Engine)
	Shutdown() error
}
```

------

## utils/env.go

```go
package utils

import (
	"os"
	"regexp"
)

// ${VAR} 或 ${VAR:default}
var envPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_]+)(?::([^}]+))?\}`)

func ExpandEnv(s string) string {
	return envPattern.ReplaceAllStringFunc(s, func(m string) string {
		groups := envPattern.FindStringSubmatch(m)
		if len(groups) < 2 {
			return m
		}
		key := groups[1]
		def := ""
		if len(groups) > 2 {
			def = groups[2]
		}
		if val := os.Getenv(key); val != "" {
			return val
		}
		if def != "" {
			return def
		}
		return m
	})
}

func ExpandConfig(v any) any {
	switch val := v.(type) {
	case string:
		return ExpandEnv(val)
	case map[string]any:
		newMap := make(map[string]any)
		for k, v2 := range val {
			newMap[k] = ExpandConfig(v2)
		}
		return newMap
	case []any:
		newSlice := make([]any, len(val))
		for i, v2 := range val {
			newSlice[i] = ExpandConfig(v2)
		}
		return newSlice
	default:
		return v
	}
}
```

------

## registry/registry.go

```go
package registry

import "myapp/module"

var Modules = map[string]func() module.Module{}
```

------

## modules/auth/auth.go

```go
package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type AuthModule struct{}

func (m *AuthModule) Deps() []string { return nil }

func (m *AuthModule) Init(cfg module.ModuleConfig) error {
	fmt.Println("[auth] Init")
	return nil
}

func (m *AuthModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Auth OK"})
	})
}

func (m *AuthModule) Shutdown() error {
	fmt.Println("[auth] Shutdown")
	return nil
}

func New() module.Module { return &AuthModule{} }
```

------

## modules/user/user.go

```go
package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type UserModule struct {
	greeting string
}

func (m *UserModule) Deps() []string { return nil }

func (m *UserModule) Init(cfg module.ModuleConfig) error {
	if g, ok := cfg["greeting"].(string); ok {
		m.greeting = g
	} else {
		m.greeting = "Hello from user (default)"
	}
	fmt.Println("[user] Init with greeting =", m.greeting)
	return nil
}

func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": m.greeting})
	})
}

func (m *UserModule) Shutdown() error {
	fmt.Println("[user] Shutdown")
	return nil
}

func New() module.Module { return &UserModule{} }
```

------

## modules/order/order.go

```go
package order

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type OrderModule struct {
	dsn string
}

func (m *OrderModule) Deps() []string { return []string{"auth"} }

func (m *OrderModule) Init(cfg module.ModuleConfig) error {
	if dsn, ok := cfg["dsn"].(string); ok {
		m.dsn = dsn
	} else {
		m.dsn = "memory://default"
	}
	fmt.Println("[order] Init with DSN =", m.dsn)
	return nil
}

func (m *OrderModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/order", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Order module using DSN: " + m.dsn})
	})
}

func (m *OrderModule) Shutdown() error {
	fmt.Println("[order] Shutdown")
	return nil
}

func New() module.Module { return &OrderModule{} }
```

------

## main.go（核心）

太长了就不全贴了，关键点：

- 通过 `fsnotify` 监控 `config.yaml`
- `APP_ENV=dev` → `DebugMode` + pprof
- `APP_ENV` 不是 dev → `ReleaseMode`
- `make dump` → 打印完整配置

👉 这一部分我们之前已经整理过，你只要把上面版本粘贴到 `main.go` 就行。

------

## Makefile

```makefile
APP_NAME := myapp
MAIN := main.go

.PHONY: run dump tidy build dev

run:
	go run $(MAIN)

dump:
	go run $(MAIN) dump

tidy:
	go mod tidy

build:
	go build -o $(APP_NAME) $(MAIN)

# 开发热重载（用 air）
dev:
	APP_ENV=dev air -c .air.toml
```

------

## 开发 & 调试

### 启动服务

```bash
make run
```

### 打印配置（含环境变量替换结果）

```bash
make dump
```

### 开发模式（热重载 + Debug + pprof）

```bash
make dev
```

输出示例：

```
[dev mode] Gin running in DebugMode
[dev mode] Watching config.yaml ...
[dev mode] pprof enabled at /debug/pprof
```

### 访问路由

- `GET /auth` → `{ "msg": "Auth OK" }`
- `GET /user` → `{ "msg": "Hi from user module" }`
- `GET /order` → `{ "msg": "Order module using DSN: ..." }`

------

## 总结

这个工程模板已经具备：

- ✅ 模块化加载，支持依赖声明
- ✅ 配置注入，支持 `${ENV:default}` 环境变量替换
- ✅ 热加载 config.yaml，动态启停模块
- ✅ 生命周期管理（Init/Shutdown）
- ✅ 开发/生产模式区分
- ✅ 开发模式开启 pprof 调试路由
- ✅ Makefile 脚本简化常用操作

------

要不要我再帮你补一个 **README.md**，把这些功能整理成对外文档，方便团队成员直接使用？

好主意 👍
 这样新同事就能直接按照 README 步骤扩展模块。
 我在 README 里加一节 **新增模块示例（Product 模块）**，展示完整流程。

------

## 更新版 README.md

```markdown
# MyApp - Go Gin 模块化工程模板

这是一个基于 **Gin** 的模块化工程模板，支持：

- ✅ 模块化加载（支持依赖）
- ✅ 模块生命周期管理（Init / Shutdown）
- ✅ 配置注入（支持 `${ENV:default}` 环境变量替换）
- ✅ 热加载 `config.yaml`（动态启停模块）
- ✅ 开发 / 生产模式区分
- ✅ 开发模式自动开启 `pprof` 调试
- ✅ Makefile 快捷命令（run / dump / dev / build）

---

## 📂 项目结构
```

myapp/
 ├── go.mod
 ├── go.sum
 ├── main.go
 ├── Makefile
 ├── .air.toml
 ├── config.yaml
 ├── module/
 │   └── module.go
 ├── registry/
 │   └── registry.go
 ├── utils/
 │   └── env.go
 └── modules/
 ├── auth/
 │   └── auth.go
 ├── user/
 │   └── user.go
 └── order/
 └── order.go

```
---

## ⚙️ 配置示例 `config.yaml`

```yaml
modules:
  - auth
  - user
  - order

configs:
  user:
    greeting: "${USER_GREETING:Hello, Default User!}"
  order:
    dsn: "${DB_DSN:mysql://root:123@tcp(localhost:3306)/shop}"
```

------

## 🚀 快速开始

### 依赖安装

```bash
go mod tidy
```

### 启动服务

```bash
make run
```

### 打印配置（含环境变量替换结果）

```bash
make dump
```

### 开发模式（热重载 + DebugMode + pprof）

```bash
make dev
```

### 编译

```bash
make build
./myapp
```

------

## 🔥 调试

- 开发模式下自动启用 **pprof**：
  - `http://localhost:8080/debug/pprof/`
  - 可用 `go tool pprof` 进行性能分析

------

## 📡 API 示例

- `GET /auth` → `{ "msg": "Auth OK" }`
- `GET /user` → `{ "msg": "Hi from user module" }`
- `GET /order` → `{ "msg": "Order module using DSN: ..." }`

------

## 🆕 新增模块示例：Product 模块

### 1. 新建目录 `modules/product/product.go`

```go
package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myapp/module"
)

type ProductModule struct {
	prefix string
}

func (m *ProductModule) Deps() []string { return nil }

// 支持配置注入，例如 configs.product.prefix
func (m *ProductModule) Init(cfg module.ModuleConfig) error {
	if p, ok := cfg["prefix"].(string); ok {
		m.prefix = p
	} else {
		m.prefix = "DefaultProduct"
	}
	fmt.Println("[product] Init with prefix =", m.prefix)
	return nil
}

func (m *ProductModule) RegisterRoutes(r *gin.Engine) {
	r.GET("/product", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": m.prefix + " module works"})
	})
}

func (m *ProductModule) Shutdown() error {
	fmt.Println("[product] Shutdown")
	return nil
}

func New() module.Module { return &ProductModule{} }
```

------

### 2. 在 `registry/registry.go` 注册模块

```go
import (
    "myapp/module"
    "myapp/modules/auth"
    "myapp/modules/order"
    "myapp/modules/user"
    "myapp/modules/product"   // 新增
)

var Modules = map[string]func() module.Module{
    "auth":    auth.New,
    "user":    user.New,
    "order":   order.New,
    "product": product.New,   // 新增
}
```

------

### 3. 修改配置 `config.yaml`

```yaml
modules:
  - auth
  - product   # 新增模块
  - user
  - order

configs:
  product:
    prefix: "Awesome"
```

------

### 4. 启动后访问

```bash
make run
```

输出：

```
[product] Init with prefix = Awesome
```

访问接口：

```
GET http://localhost:8080/product
→ { "msg": "Awesome module works" }
```