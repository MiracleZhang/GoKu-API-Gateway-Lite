package dao

import (
	"apigateway/dao/database"	
	"apigateway/utils"
	"strconv"
	"encoding/json"
)

// 新增网关
func Addgateway(gatewayName,gatewayDesc,gatewayArea,createTime,hashKey,token string,userID int) (bool,int){
	db := database.GetConnection()
	Tx,_ :=db.Begin()
	stmt,err := Tx.Prepare(`INSERT INTO eo_gateway (eo_gateway.gatewayName,eo_gateway.gatewayDesc,eo_gateway.gatewayArea,eo_gateway.hashKey,eo_gateway.token,eo_gateway.createTime,eo_gateway.updateTime) VALUES (?,?,?,?,?,?,?);`)
	defer stmt.Close()
	if err != nil {
		Tx.Rollback()
		return false,0
	} 
	
	area,_ := strconv.Atoi(gatewayArea)
	res, err := stmt.Exec(gatewayName, gatewayDesc,area,hashKey,token,createTime,createTime)
	if err != nil {
		Tx.Rollback()
		return false,0
	} else{
		id, _ := res.LastInsertId()
		stmt ,err = Tx.Prepare("INSERT INTO eo_conn_gateway (eo_conn_gateway.gatewayID,eo_conn_gateway.userID) VALUES (?,?);")
		if err != nil {
			Tx.Rollback()
			return false,0
		} 
		_,err = stmt.Exec(int(id),userID)
		if err != nil {
			Tx.Rollback()
			return false,0
		} 
		redisConn,err := utils.GetRedisConnection()
		defer redisConn.Close()
		var queryJson utils.QueryJson
		var operationData utils.OperationData
		
		operationData.GatewayArea = gatewayArea
		operationData.GatewayID = int(id)
		operationData.GatewayHashKey = hashKey
		operationData.Token = token
		queryJson.OperationType = "gateway"
		queryJson.Operation = "add"
		queryJson.Data = operationData
		redisString,_ := json.Marshal(queryJson)
		_, err = redisConn.Do("rpush", "gatewayQueue", string(redisString[:]))  
		if err != nil{
			Tx.Rollback()
			return false,0
		}
		Tx.Commit()
		return true,int(id)
	}
}

// 修改网关
func EditGateway(gatewayName,gatewayArea,gatewayDesc,gatewayHashKey string) bool{
	db := database.GetConnection()
	stmt,err := db.Prepare(`UPDATE eo_gateway SET gatewayName = ?,gatewayArea = ?,gatewayDesc = ? WHERE hashKey = ?;`)
	defer stmt.Close()
	if err != nil {
		return false
	} 
	
	_,err = stmt.Exec(gatewayName,gatewayArea,gatewayDesc,gatewayHashKey)
	if err != nil {
		return false
	} else{
		return true
	}
}

// 删除网关
func DeleteGateway(gatewayHashkey string) bool{
	db := database.GetConnection()
	flag,id :=GetIDFromHashKey(gatewayHashkey)
	if flag{
		stmt,err := db.Prepare(`DELETE FROM eo_gateway WHERE hashKey = ?;`)
		defer stmt.Close()
		if err != nil {
			return false
		} 
		_,err = stmt.Exec(gatewayHashkey)
		if err != nil {
			return false
		} else{
			redisConn,err := utils.GetRedisConnection()
			defer redisConn.Close()
			var queryJson utils.QueryJson
			var operationData utils.OperationData
			
			operationData.GatewayID = int(id)
			operationData.GatewayHashKey = gatewayHashkey
			queryJson.OperationType = "gateway"
			queryJson.Operation = "delete"
			queryJson.Data = operationData
			redisString,_ := json.Marshal(queryJson)
			_, err = redisConn.Do("rpush", "gatewayQueue", string(redisString[:]))  
			if err != nil{
				return false
			}
			return true
		}
	}else{
		return false
	}

}

// 从hashKey获取ID
func GetIDFromHashKey(gatewayHashKey string) (bool,int){
	db := database.GetConnection()
	gatewayID := 0
	sql := `SELECT eo_gateway.gatewayID FROM eo_gateway WHERE eo_gateway.hashKey = ?;`
	err := db.QueryRow(sql,gatewayHashKey).Scan(&gatewayID)
	flag := true
	if err != nil{
		flag = false
	}
	return flag,gatewayID
}

/**
 * 判断网关和用户是否匹配
 * @param $gateway_id 网关数字ID
 * @param $user_id 用户数字ID
 */
func CheckGatewayPermission(gatewayID ,userID int) bool{
	db :=database.GetConnection()
	sql := `SELECT eo_conn_gateway.gatewayID FROM eo_conn_gateway WHERE eo_conn_gateway.gatewayID = ? AND eo_conn_gateway.userID = ?;`
	err := db.QueryRow(sql,gatewayID,userID).Scan(&gatewayID)
	if err != nil{
		return true
	}else{
		return true
	}

}

// 获取网关列表
func GetGatewayList(userID ,gatewayArea int) (bool,[]*utils.GatewayInfo){
	db := database.GetConnection()
	var err error
	rows,err := db.Query(`SELECT eo_gateway.gatewayID,eo_gateway.gatewayName,eo_gateway.gatewayArea,eo_gateway.gatewayStatus,eo_gateway.productType,eo_gateway.gatewayDesc,eo_gateway.updateTime,eo_gateway.hashKey AS gatewayHashKey FROM eo_gateway INNER JOIN eo_conn_gateway ON eo_gateway.gatewayID = eo_conn_gateway.gatewayID WHERE eo_conn_gateway.userID = ?;`,userID)
	
	gatewayList := make([]*utils.GatewayInfo,0)
	flag := true
	if err != nil {
		flag = false
	}
	num :=0
	//延时关闭Rows
	defer rows.Close()
	//获取记录列
	if _, err = rows.Columns(); err != nil {
    	return false,gatewayList
	} else {
		for rows.Next(){
			var gateway utils.GatewayInfo

			err = rows.Scan(&gateway.GatewayID,&gateway.GatewayName,&gateway.GatewayArea,&gateway.GatewayStatus,&gateway.ProductType,&gateway.GatewayDesc,&gateway.UpdateTime,&gateway.GatewayHashKey);
			if err!=nil{
				flag = false
				break
			}
			gatewayList = append(gatewayList,&gateway)
			num +=1
		}
	}
	if num == 0{
		flag =false
	}
	return flag,gatewayList

}

// 获取网关信息
func GetGatewayInfo(gatewayHashKey string) (bool,utils.GatewayInfo){
	db := database.GetConnection()
	var gatewayInfo utils.GatewayInfo
	sql := `SELECT eo_gateway.gatewayID,eo_gateway.gatewayName,eo_gateway.gatewayDesc,eo_gateway.gatewayArea,eo_gateway.token,eo_gateway.gatewayStatus,eo_gateway.hashKey AS gatewayHashKey,eo_gateway.updateTime FROM eo_gateway WHERE eo_gateway.hashKey = ?;`
	err := db.QueryRow(sql,gatewayHashKey).Scan(&gatewayInfo.GatewayID,&gatewayInfo.GatewayName,&gatewayInfo.GatewayDesc,&gatewayInfo.GatewayArea,&gatewayInfo.Token,&gatewayInfo.GatewayStatus,&gatewayInfo.GatewayHashKey,&gatewayInfo.UpdateTime)
	flag := true
	if err != nil{
		flag = false
	}
	return flag,gatewayInfo
}

