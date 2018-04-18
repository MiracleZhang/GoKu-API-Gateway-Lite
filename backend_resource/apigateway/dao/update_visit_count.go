package dao

import (
	"apigateway/dao/cache"
	"apigateway/utils"
	"strconv"
	"time"

	"github.com/farseer810/requests"
	"github.com/farseer810/yawf"
)

func UpdateVisitCount(context yawf.Context, info *utils.MappingInfo, response requests.Response,remoteIP string) {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	var redisKey string = "gatewayDayCount:" + info.GatewayHashKey + ":" + dateStr
	redisConn,err := utils.GetRedisConnection()
	defer redisConn.Close()
	if err != nil{
		panic(err)
	}
	redisConn.Do("INCR", redisKey)

	timeStr := dateStr + "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute())
	redisKey = "gatewayMinuteCount:" + info.GatewayHashKey + ":" + timeStr
	redisConn.Do("INCR", redisKey)
	
	redisKey = "gatewaySecondCount:" + info.GatewayHashKey + ":" + timeStr + "-" + strconv.Itoa(now.Second())
	redisConn.Do("INCR", redisKey)

	redisKey = "gatewayDayThroughput:" + info.GatewayHashKey + ":" + dateStr
	redisConn.Do("INCRBY", redisKey, response.ContentLength())

	redisKey = "IPMinuteCount:" + remoteIP + ":" + timeStr
	redisConn.Do("INCR", redisKey)
}

func UpdateIPVisitCount(context yawf.Context, remoteIP string, response requests.Response) {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	timeStr := dateStr + "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute())
	var redisKey string = "IPMinuteCount:" + remoteIP + ":" + timeStr
	conn := cache.GetConnection(context)
	conn.Do("INCR", redisKey)
}
