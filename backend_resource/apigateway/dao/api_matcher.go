package dao

import (
	"apigateway/dao/cache"
	"github.com/farseer810/yawf"
	"github.com/garyburd/redigo/redis"
)

func GetAllAPIPaths(context yawf.Context, gatewayHashKey string) []string {
	result := getAllAPIPathsFromCache(context, gatewayHashKey)
	return result
}

func getAllAPIPathsFromCache(context yawf.Context, gatewayHashKey string) []string {
	conn := cache.GetConnection(context)

	var redisKey string = "apiList:" + gatewayHashKey
	result, err := redis.Strings(conn.Do("LRANGE", redisKey, 0, -1))
	if err != nil {
		panic(err)
	}
	return result
}
