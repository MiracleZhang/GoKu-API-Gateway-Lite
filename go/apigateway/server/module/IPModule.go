package module

import (
	"apigateway/utils"
	"apigateway/server/dao"
)

func EditIPList(gatewayID,ipType int,gatewayHashKey,ipList string) bool{
	return dao.EditIPList(gatewayID,ipType,gatewayHashKey,ipList)
}

// 获取IP名单列表
func GetIPList(gatewayID int) (bool,utils.IPList){
	return dao.GetIPList(gatewayID)
}

// 检查IP名单是否存在
func CheckIPListIsExist(gatewayID int) bool{
	return dao.CheckIPListIsExist(gatewayID)
}
 