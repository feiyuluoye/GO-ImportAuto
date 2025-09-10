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
