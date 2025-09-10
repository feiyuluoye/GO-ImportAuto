APP_NAME := myapp
MAIN := main.go

.PHONY: run dump tidy build dev

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

# 开发热重载（用 air）
dev:
	APP_ENV=dev air -c .air.toml
