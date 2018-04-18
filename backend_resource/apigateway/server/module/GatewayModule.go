package module

import (
	"apigateway/utils"
	"apigateway/server/dao"
	"time"
)
// 新增网关
func Addgateway(gatewayName,gatewayDesc,gatewayArea string,userID int) (bool,string){
	createTime := time.Now().Format("2006-01-02 15:04:05")
	hashKey := utils.GetHashKey(gatewayName,gatewayDesc,gatewayArea)
	token := utils.GetHashKey(gatewayName)
	if flag,_ :=dao.Addgateway(gatewayName,gatewayDesc,gatewayArea,createTime,hashKey,token,userID);flag{
		return true,hashKey
	}else{
		return true,""
	}
}

// 修改网关
func EditGateway(gatewayName,gatewayArea,gatewayDesc,gatewayHashKey string,userID int) (bool){
	flag,gatewayID := dao.GetIDFromHashKey(gatewayHashKey)
	if flag{
		if dao.CheckGatewayPermission(gatewayID,userID){
			return dao.EditGateway(gatewayName,gatewayArea,gatewayDesc,gatewayHashKey)
		}else{
			return false
		}
	}else{
		return false
	}
}

// 删除网关
func DeleteGateway(gatewayHashKey string,userID int) (bool){
	flag,gatewayID := dao.GetIDFromHashKey(gatewayHashKey)
	if flag{
		if dao.CheckGatewayPermission(gatewayID,userID){
			return dao.DeleteGateway(gatewayHashKey)
		}else{
			return false
		}
	}else{
		return false
	}
}

// 获取网关信息
func GetGatewayInfo(gatewayHashKey string,userID int) (bool,utils.GatewayInfo){
	var gatewayInfo utils.GatewayInfo
	flag,gatewayID := dao.GetIDFromHashKey(gatewayHashKey)
	if flag{
		if dao.CheckGatewayPermission(gatewayID,userID){
			return dao.GetGatewayInfo(gatewayHashKey)
		}else{
			return false,gatewayInfo
		}
	}else{
		return false,gatewayInfo
	}
}

// 获取网关列表
func GetGatewayList(gatewayArea int,userID int) (bool,[]*utils.GatewayInfo){
	return dao.GetGatewayList(userID,gatewayArea)
}
/**
 * 判断用户是否拥有对网关的操作权限
 * @param $hash_key 网关的hash_key
 * @param $user_id 用户的数字ID
 */
func CheckGatewayPermission(gatewayHashKey string,userID int) bool{
	flag,gatewayID := dao.GetIDFromHashKey(gatewayHashKey)
	
	if flag{
		return dao.CheckGatewayPermission(gatewayID,userID)
	}else{
		return false
	}
}

// 从hashKey获取ID
func GetIDFromHashKey(gatewayHashKey string) (bool,int){
	return dao.GetIDFromHashKey(gatewayHashKey)
}