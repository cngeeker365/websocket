package redisPool

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"wsserver/config"
)

var (
	Pool *redis.Pool
)

func init()  {
	// 建立连接池
	idleTimeOut,_ := time.ParseDuration(string(config.Conf.Redis.Pool.IdleTimeout))

	Pool = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     config.Conf.Redis.Pool.MaxIdle,
		MaxActive:   config.Conf.Redis.Pool.MaxActive,
		IdleTimeout: time.Second * idleTimeOut,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Conf.Redis.Pool.Host)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
