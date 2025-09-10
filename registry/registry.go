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
