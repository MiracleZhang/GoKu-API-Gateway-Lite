package dao

import (
	"apigateway/dao/database"	
	"apigateway/utils"
)

func AddFrequencyLimit(gatewayID,limitCount,intervalType int) bool{
	db := database.GetConnection()
	stmt,err := db.Prepare("INSERT INTO eo_gateway_frequency (gatewayID,count,intervalType) VALUES (?,?,?);")
	defer stmt.Close()
	if err != nil{
		return false
	}
	_,err = stmt.Exec(gatewayID,limitCount,intervalType)
	if err != nil{
		return false
	}else{
		return true
	}
}

// 修改频率限制
func EditFrequencyLimit(gatewayHashKey string,gatewayID,limitCount,intervalType int) bool{
	db := database.GetConnection()
	stmt,err := db.Prepare("UPDATE eo_gateway_frequency SET count = ? WHERE intervalType =? AND gatewayID = ?;")
	defer stmt.Close()
	if err != nil{
		return false
	}
	_,err = stmt.Exec(limitCount,intervalType,gatewayID)
	if err != nil{
		return false
	}else{
		redisConn,_ := utils.GetRedisConnection()
		defer redisConn.Close()
		
		redisConn.Do("del", "gatewaySecondFrequency:"+gatewayHashKey)  
		redisConn.Do("del", "gatewayMinuteFrequency:"+gatewayHashKey)  
		return true
	}
}

// 删除频率限制
func DeleteFrequencyLimit(gatewayHashKey string,gatewayID,intervalType int) bool{
	db := database.GetConnection()
	stmt,err := db.Prepare("DELETE FROM eo_gateway_frequency WHERE gatewayID = ? AND intervalType = ?;")
	defer stmt.Close()
	if err != nil{
		return false
	}
	_,err = stmt.Exec(gatewayID,intervalType)
	if err != nil{
		return false
	}else{
		redisConn,_ := utils.GetRedisConnection()
		defer redisConn.Close()
		
		redisConn.Do("del", "gatewaySecondFrequency:"+gatewayHashKey)  
		redisConn.Do("del", "gatewayMinuteFrequency:"+gatewayHashKey)  
		return true
	}
}

// 获取频率限制列表
func GetFrequencyLimitList(gatewayID int) (bool,[]*utils.FrequencyInfo){
	db := database.GetConnection()
	rows,err :=db.Query(`SELECT intervalType,count FROM eo_gateway_frequency WHERE gatewayID = ? ORDER BY intervalType;`,gatewayID)
	frequencyList := make([]*utils.FrequencyInfo,0)

	flag := true
	num :=0
	if _, err = rows.Columns(); err != nil {
		return false,frequencyList
	} else {
		for rows.Next(){
			var frequency utils.FrequencyInfo
			err = rows.Scan(&frequency.IntervalType,&frequency.Count)
			if err != nil{
				flag = false
				break
			}
			frequencyList = append(frequencyList,&frequency)
			num +=1
		}
	}
	if num == 0{
		flag =false
	}
	
	return flag,frequencyList
}

// 检查频率限制是否存在
func CheckFrequencyLimitIsExist(gatewayID,intervalType int) bool{
	db := database.GetConnection()
	count := 0
	db.QueryRow("SELECT COUNT(0) AS count FROM eo_gateway_frequency WHERE gatewayID = ? AND intervalType = ?;",gatewayID,intervalType).Scan(&count)
	if count == 0{
		return false 
	}else{
		return true
	}
}