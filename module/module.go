package module

import "github.com/gin-gonic/gin"

type ModuleConfig map[string]any

// 模块接口：支持生命周期 & 依赖声明
type Module interface {
	Deps() []string               // 模块依赖哪些其他模块
	Init(cfg ModuleConfig) error  // 模块初始化
	RegisterRoutes(r *gin.Engine) // 注册路由
	Shutdown() error              // 模块销毁（释放资源）
}
