package dao

import (
	"apigateway/dao/cache"
	"apigateway/utils"
	"strconv"
	"time"
	"apigateway/dao/database"
	"github.com/farseer810/yawf"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func GetGatewayDayVisitCount(context yawf.Context, info *utils.MappingInfo) int {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)

	var redisKey string = "gatewayDayCount:" + info.GatewayHashKey + ":" + dateStr
	conn := cache.GetConnection(context)
	count, err := redis.Int(conn.Do("GET", redisKey))

	if err == redis.ErrNil {
		return 0
	} else if err != nil {
		panic(err)
	}
	return count
}

func GetGatewayMinuteCount(context yawf.Context, info *utils.MappingInfo) int {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	timeStr := dateStr + "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute())

	var redisKey string = "gatewayMinuteCount:" + info.GatewayHashKey + ":" + timeStr
	conn := cache.GetConnection(context)
	count, err := redis.Int(conn.Do("GET", redisKey))

	if err == redis.ErrNil {
		return 0
	} else if err != nil {
		panic(err)
	}
	return count
}

func GetGatewaySecondCount(context yawf.Context, info *utils.MappingInfo) int {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	timeStr := dateStr + "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute()) + "-" + strconv.Itoa(now.Second())

	var redisKey string = "gatewaySecondCount:" + info.GatewayHashKey + ":" + timeStr
	conn := cache.GetConnection(context)
	count, err := redis.Int(conn.Do("GET", redisKey))

	if err == redis.ErrNil {
		return 0
	} else if err != nil {
		panic(err)
	}
	return count
}

func GetGatewayDayThroughput(context yawf.Context, info *utils.MappingInfo) int {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)

	var redisKey string = "gatewayDayThroughput:" + info.GatewayHashKey + ":" + dateStr
	conn := cache.GetConnection(context)
	count, err := redis.Int(conn.Do("GET", redisKey))

	if err == redis.ErrNil {
		return 0
	} else if err != nil {
		panic(err)
	}
	return count
}

func GetIPMinuteCount(context yawf.Context, remoteIP string) int {
	now := time.Now()
	year, month, day := now.Date()
	dateStr := strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	timeStr := dateStr + "-" + strconv.Itoa(now.Hour()) + "-" + strconv.Itoa(now.Minute())
	var redisKey string = "IPMinuteCount:" + remoteIP + ":" + timeStr
	conn := cache.GetConnection(context)
	count, err := redis.Int(conn.Do("GET", redisKey))
	if err == redis.ErrNil {
		return 0
	} else if err != nil {
		panic(err)
	}
	return count
}

// 加载秒阀值
func loadSecondValve(hashKey string) (utils.FrequencyInfo,bool) {
	var intervalType,count int
	var info utils.FrequencyInfo
	db := database.GetConnection()
	sql := `SELECT intervalType,count FROM eo_gateway_frequency WHERE ` +
		`gatewayID=(SELECT gatewayID FROM eo_gateway WHERE hashKey=?) ` +
		`AND intervalType = 0`
	err := db.QueryRow(sql, hashKey).Scan(&intervalType,&count)
	flag := true
	if err != nil {
		flag = false
	}
	info.Count = count
	info.IntervalType = intervalType
	return info,flag
}

// 加载分阀值
func loadMinuteValve(hashKey string) (utils.FrequencyInfo,bool) {
	var intervalType,count int
	var info utils.FrequencyInfo
	db := database.GetConnection()
	sql := `SELECT intervalType,count FROM eo_gateway_frequency WHERE ` +
		`gatewayID=(SELECT gatewayID FROM eo_gateway WHERE hashKey=?) ` +
		`AND intervalType = 1`
	err := db.QueryRow(sql, hashKey).Scan(&intervalType,&count)
	flag := true
	if err != nil {
		flag = false
	}
	info.Count = count
	info.IntervalType = intervalType
	return info,flag
}

// 获取每分钟的阀值
func getGatewayMinuteValve(context yawf.Context,hashKey string) (utils.FrequencyInfo,bool) {
	var redisKey string = "gatewayMinuteFrequency:" + hashKey
	var info utils.FrequencyInfo
	conn := cache.GetConnection(context)
	infoStr, err := redis.String(conn.Do("GET", redisKey))
	flag := true
	if err == redis.ErrNil {
		info,flag = loadMinuteValve(hashKey)
		if flag{
			infoByte,_ := json.Marshal(info)
			infoStr =  string(infoByte[:])
			// 缓存时间为1 hour
			conn.Do("SET", redisKey, infoStr, "EX", 3600)
		}
	} else if err != nil {
		flag = false
	} else {
		err = json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			flag = false
		}
	}
	return info,flag
}

// 获取每分钟的阀值
func getGatewaySecondValve(context yawf.Context,hashKey string) (utils.FrequencyInfo,bool){
	var redisKey string = "gatewaySecondFrequency:" + hashKey
	var info utils.FrequencyInfo
	conn := cache.GetConnection(context)
	infoStr, err := redis.String(conn.Do("GET", redisKey))
	flag := true
	if err == redis.ErrNil {
		info,flag = loadSecondValve(hashKey)
		if flag{
			infoByte,_ := json.Marshal(info)
			infoStr =  string(infoByte[:])
			// 缓存时间为1 hour
			conn.Do("SET", redisKey, infoStr, "EX", 3600)
		}
	} else if err != nil {
		flag = false
	} else {
		err = json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			flag = false
		}
	}
	return info,flag
}

// 获取阀值列表
func getGatewayValve(context yawf.Context,hashKey string) []utils.FrequencyInfo{
	info := make([]utils.FrequencyInfo,0)
	minuteValve,flag := getGatewayMinuteValve(context,hashKey)
	if flag{
		info = append(info,minuteValve)
	}
	secondValve,flag := getGatewaySecondValve(context,hashKey)
	if flag{
		info = append(info,secondValve)
	}		
	return info
}

func GetGatewayValve(context yawf.Context,hashKey string) []utils.FrequencyInfo{
	return getGatewayValve(context,hashKey)
}
