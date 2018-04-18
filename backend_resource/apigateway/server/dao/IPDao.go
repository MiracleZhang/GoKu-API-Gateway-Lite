package dao

import (
	"apigateway/dao/database"	
	"apigateway/utils"
)

func EditIPList(gatewayID,ipType int,gatewayHashKey,ipList string) bool{
	db := database.GetConnection()
	listType := "blackList"
	if ipType == 2{
		listType = "whiteList"
	}
	isExist := CheckIPListIsExist(gatewayID)
	if isExist == true{
		if ipType == 0{
			stmt,err := db.Prepare("UPDATE eo_gateway_ip_restrict SET chooseType = 0 WHERE gateWayID =?;")
			if err != nil{
				return false
			}
			_,err = stmt.Exec(gatewayID)
			if err != nil{
				return false
			}
		}else{
			stmt,err := db.Prepare("UPDATE eo_gateway_ip_restrict SET " + listType + " = ?,chooseType = ? WHERE gateWayID =?;")
			if err != nil{
				return false
			}
			_,err = stmt.Exec(ipList,ipType,gatewayID)
			if err != nil{
				return false
			}
		}
	}else{
		if ipType == 0{
			stmt,err := db.Prepare("INSERT INTO eo_gateway_ip_restrict (gatewayID,chooseType) VALUES (?,0);")
			if err != nil{
				return false
			}
			_,err = stmt.Exec(gatewayID)
			if err != nil{
				return false
			}
		}else{
			stmt,err := db.Prepare("INSERT INTO eo_gateway_ip_restrict (gatewayID," + listType + ",chooseType) VALUES (?,?,?);")
			if err != nil{
				return false
			}
			_,err = stmt.Exec(gatewayID,ipList,ipType)
			if err != nil{
				return false
			}
		}
	}

	redisConn,err := utils.GetRedisConnection()
	defer redisConn.Close()
	
	_, err = redisConn.Do("del", "IPList:"+gatewayHashKey)  
	if err != nil{
		return false
	}
	
	return true
}

// 获取IP名单列表
func GetIPList(gatewayID int) (bool,utils.IPList){
	db := database.GetConnection()
	var ipList utils.IPList
	sql := "SELECT IFNULL(blackList,''),IFNULL(whiteList,''),chooseType FROM eo_gateway_ip_restrict WHERE gatewayID = ?;"
	err := db.QueryRow(sql,gatewayID).Scan(&ipList.BlackList,&ipList.WhiteList,&ipList.ChooseType)
	if err != nil {
		return false,ipList
	}else{
		return true,ipList
	}

}

// 检查IP名单是否存在
func CheckIPListIsExist(gatewayID int) bool{
	db := database.GetConnection()
	count :=0
	
	sql := "SELECT COUNT(0) AS count FROM eo_gateway_ip_restrict WHERE gatewayID = ?;"
	db.QueryRow(sql,gatewayID).Scan(&count)
	if count == 0 {
		return false
	}else{
		return true
	}
}
 