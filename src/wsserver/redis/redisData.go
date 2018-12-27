package redisPool

import (
	"github.com/gomodule/redigo/redis"
	"strings"
)

func GetRedisData(conn redis.Conn, page string) (data map[string]string, err error) {
	conn.Do("SELECT", 1)
	redisKeys, _ := redis.String(conn.Do("GET", page))
	keys := strings.Split(redisKeys,",");

	result:=make(map[string]string)

	conn.Do("SELECT", 0)
	for _, key := range keys{
		value,_ := redis.String(conn.Do("GET", key))
		result[key]=value
	}
	return result, nil
}