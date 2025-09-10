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
