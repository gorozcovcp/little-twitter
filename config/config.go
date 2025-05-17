package config

import "syscall"

type Config struct {
	MongoURI   string
	DBName     string
	RedisAddr  string
	ServerAddr string
}

func LoadConfig() Config {
	return Config{
		MongoURI:   getEnv("MONGO_URI", "mongodb://mongo:27017/?directConnection=true"),
		DBName:     getEnv("DB_NAME", "uala"),
		RedisAddr:  getEnv("REDIS_ADDR", "redis:6379"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
