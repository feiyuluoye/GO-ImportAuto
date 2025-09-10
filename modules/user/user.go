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
