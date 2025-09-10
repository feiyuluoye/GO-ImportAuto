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
