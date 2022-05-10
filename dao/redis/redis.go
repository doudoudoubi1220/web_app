package redis

import (
	"fmt"
	"web_app/settings"

	"github.com/go-redis/redis"
)

var client *redis.Client

//Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return
}
func Close() {
	_ = client.Close()
}
