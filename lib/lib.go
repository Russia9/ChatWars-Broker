package lib

import "os"

func GetEnv(key string, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}
