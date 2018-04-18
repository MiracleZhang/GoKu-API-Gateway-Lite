package dao

import (
	"encoding/json"
	"apigateway/dao/cache"
	"apigateway/dao/database"
	"apigateway/utils"
	"fmt"
	"strings"
	"github.com/farseer810/yawf"
	"github.com/garyburd/redigo/redis"
)

func loadGatewayIPList(context yawf.Context, hashKey string) (*utils.IPListInfo,bool) {
	var ipList []string
	var blackList string
	var whiteList string
	var chooseType int
	var info *utils.IPListInfo = &utils.IPListInfo{}
	flag := true
	db := database.GetConnection()
	sql := `SELECT chooseType,IFNULL(blackList,''),IFNULL(whiteList,'') FROM eo_gateway_ip_restrict WHERE ` +
		`gatewayID=(SELECT gatewayID FROM eo_gateway WHERE hashKey=?) `
	err := db.QueryRow(sql, hashKey).Scan(&chooseType, &blackList, &whiteList)
	if err != nil {
		flag =false
	}else{
		if chooseType == 1 {
			blackList = strings.Replace(blackList,"；",";",-1)
			ipList = strings.Split(blackList,";")
		} else if chooseType == 2 {
			whiteList = strings.Replace(whiteList,"；",";",-1)
			ipList = strings.Split(whiteList,";")
		} else {
			ipList = []string{}
		}
		info.IPList = ipList
		info.ChooseType = chooseType
	}
	
	return info,flag
}

func getGatewayIPList(context yawf.Context, hashKey string) *utils.IPListInfo {
	var redisKey string = "IPList:" + hashKey
	var info *utils.IPListInfo
	fmt.Println(hashKey)
	conn := cache.GetConnection(context)
	infoStr, err := redis.String(conn.Do("GET", redisKey))
	if err == redis.ErrNil {
		ipInfo,flag := loadGatewayIPList(context, hashKey)
		if flag{
			infoStr = ipInfo.String()
			// 缓存时间为1 hour
			conn.Do("SET", redisKey, infoStr, "EX", 3600)
		}
		info = ipInfo
	} else if err != nil {
		panic(err)
	} else {
		info = &utils.IPListInfo{}
		err = json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			panic(err)
		}
	}
	return info
}

func getAllBlackList(context yawf.Context) *utils.IPListInfo {
	var redisKey string = "Gateway:BlackList"
	var info *utils.IPListInfo

	conn := cache.GetConnection(context)
	infoStr, err := redis.String(conn.Do("GET", redisKey))
	if err == redis.ErrNil {
		info = loadBlackList(context)
		infoStr = info.String()
		// 缓存时间为1 hour
		conn.Do("SET", redisKey, infoStr, "EX", 3600)
	} else if err != nil {
		panic(err)
	} else {

		info = &utils.IPListInfo{}
		err = json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			panic(err)
		}
	}
	return info
}

//合并网关黑名单和动态黑名单
func mergeIPList(context yawf.Context, hashKey string) *utils.IPListInfo {
	personIPList := getGatewayIPList(context, hashKey)
	allBlackList := getAllBlackList(context)
	if personIPList != nil{
		if personIPList.ChooseType == 1 {
			for _, ip := range allBlackList.IPList {
				personIPList.IPList = append(personIPList.IPList, ip)
			}
		}
		fmt.Println(personIPList.IPList)
		return personIPList
	}else{
		return allBlackList
	}
}

//从数据库中加载全局黑名单
func loadBlackList(context yawf.Context) *utils.IPListInfo {
	var ipList []string
	var blackList *string

	db := database.GetConnection()
	
	sql := `SELECT ip FROM eo_gateway_ip_cache`
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		if err := rows.Scan(&blackList); err != nil {
			panic(err)
		}
		fmt.Println(*blackList)
		ipList = append(ipList, *blackList)
	}
	var info *utils.IPListInfo = &utils.IPListInfo{}

	info.IPList = ipList
	info.ChooseType = 1
	return info

}

//将新增黑名单IP更新到redis数据库中
func UpdateBlackList(context yawf.Context, remoteIP string) bool {
	var redisKey string = "Gateway:BlackList"
	var info *utils.IPListInfo

	conn := cache.GetConnection(context)
	infoStr, err := redis.String(conn.Do("GET", redisKey))
	if err == redis.ErrNil {
		info = loadBlackList(context)
		infoStr = info.String()
		// 缓存时间为1 hour
		conn.Do("SET", redisKey, infoStr, "EX", 3600)
	} else if err != nil {
		panic(err)
	} else {
		info = &utils.IPListInfo{}
		err = json.Unmarshal([]byte(infoStr), &info)
		if err != nil {
			panic(err)
		}
		for _, ip := range info.IPList {
			if ip == remoteIP {
				return true
			}
		}
		info.IPList = append(info.IPList, remoteIP)
		infoStr = info.String()
		conn.Do("SET", redisKey, infoStr, "EX", 3600)
	}
	return true
}
func GetIPList(context yawf.Context, hashKey string) *utils.IPListInfo {
	return mergeIPList(context, hashKey)
}
