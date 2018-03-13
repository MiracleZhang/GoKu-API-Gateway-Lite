package dao

import (
	"apigateway/dao/cache"
	"github.com/farseer810/yawf"
	"github.com/garyburd/redigo/redis"
)

func GetGatewayToken(context yawf.Context, gatewayHashkey string) string {
	var redisKey string = "gatewayToken:" + gatewayHashkey
	conn := cache.GetConnection(context)
	token, err := redis.String(conn.Do("GET", redisKey))

	if err == redis.ErrNil {
		return ""
	} else if err != nil {
		panic(err)
	}
	return token
}
