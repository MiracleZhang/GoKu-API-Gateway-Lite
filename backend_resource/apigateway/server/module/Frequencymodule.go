package module

import (
	"apigateway/utils"
	"apigateway/server/dao"
)

func AddFrequencyLimit(gatewayHashKey string,gatewayID,limitCount,intervalType int) bool{
	if dao.CheckFrequencyLimitIsExist(gatewayID,intervalType) == false{
		return dao.AddFrequencyLimit(gatewayID,limitCount,intervalType)
	}else{
		return dao.EditFrequencyLimit(gatewayHashKey,gatewayID,limitCount,intervalType)
	}
}

// 修改频率限制
func EditFrequencyLimit(gatewayHashKey string,gatewayID,limitCount,intervalType int) bool{
	if dao.CheckFrequencyLimitIsExist(gatewayID,intervalType) == true{
		return dao.EditFrequencyLimit(gatewayHashKey,gatewayID,limitCount,intervalType)
	}else{
		return false
	}
}

// 删除频率限制
func DeleteFrequencyLimit(gatewayHashKey string,gatewayID,intervalType int) bool{
	if dao.CheckFrequencyLimitIsExist(gatewayID,intervalType) == true{
		return dao.DeleteFrequencyLimit(gatewayHashKey,gatewayID,intervalType)
	}else{
		return false
	}
}

// 获取频率限制列表
func GetFrequencyLimitList(gatewayID int) (bool,[]*utils.FrequencyInfo){
	return dao.GetFrequencyLimitList(gatewayID)
}

// 检查频率限制是否存在
func CheckFrequencyLimitIsExist(gatewayID,intervalType int) bool{
	return dao.CheckFrequencyLimitIsExist(gatewayID,intervalType)
}