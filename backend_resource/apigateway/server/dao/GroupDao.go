package dao

import (
	"apigateway/dao/database"	
	"apigateway/utils"
)

// 添加分组
func AddGroup(gatewayID int,groupName string) (bool,int){
	db := database.GetConnection()
	stmt,err := db.Prepare(`INSERT INTO eo_gateway_api_group (eo_gateway_api_group.groupName,eo_gateway_api_group.gatewayID) VALUES (?,?);`)
	defer stmt.Close()
	if err != nil {
		return false,0
	} 
	
	res, err := stmt.Exec(groupName,gatewayID)
	if err != nil {
		return false,0
	} else{
		id, _ := res.LastInsertId()
		return true,int(id)
	}
}

// 添加子分组
func AddChildGroup(gatewayID ,parentGroupID int,groupName string) (bool,int){
	db := database.GetConnection()
	stmt,err := db.Prepare(`INSERT INTO eo_gateway_api_group (eo_gateway_api_group.groupName,eo_gateway_api_group.gatewayID,eo_gateway_api_group.parentGroupID,eo_gateway_api_group.isChild) VALUES (?,?,?,1);`)
	defer stmt.Close()
	if err != nil {
		return false,0
	} 
	
	res, err := stmt.Exec(groupName,gatewayID,parentGroupID)
	if err != nil {
		return false,0
	} else{
		id, _ := res.LastInsertId()
		return true,int(id)
	}
}

// 删除网关api分组
func DeleteGroup(groupID int) bool{
	db := database.GetConnection()
	Tx,_ := db.Begin()
	stmt,err := Tx.Prepare(`DELETE FROM eo_gateway_api_group WHERE eo_gateway_api_group.groupID = ?;`)
	defer stmt.Close()
	if err != nil {

		return false
	} 
	
	_, err = stmt.Exec(groupID)
	if err != nil {
		Tx.Rollback()
		return false
	} else{
		stmt,_ = Tx.Prepare("DELETE FROM eo_gateway_api_group WHERE eo_gateway_api_group.parentGroupID = ?;")
		stmt.Exec(groupID)
		stmt,_ = Tx.Prepare("DELETE FROM eo_gateway_api WHERE eo_gateway_api.groupID = ?;")
		stmt.Exec(groupID)
		Tx.Commit()
		return true
	}
}

// 获取网关分组列表
func GetGroupList(gatewayID int) (bool,[]*utils.GroupInfo){
	db := database.GetConnection()
	rows,err := db.Query(`SELECT eo_gateway_api_group.groupID,eo_gateway_api_group.groupName FROM eo_gateway_api_group WHERE gatewayID = ? AND isChild = 0 ORDER BY eo_gateway_api_group.groupID DESC;`,gatewayID)
	
	defer rows.Close()
	groupList := make([]*utils.GroupInfo,0)
	flag := true
	if err != nil {
		flag = false
	}
	num :=0
	//延时关闭Rows
	defer rows.Close()
	//获取记录列
	if _, err = rows.Columns(); err != nil {
    	return false,groupList
	} else {
		for rows.Next(){
			var group utils.GroupInfo

			err:= rows.Scan(&group.GroupID,&group.GroupName);
			if err!=nil{
				flag = false
				break
			}
			childRows,err := db.Query(`SELECT eo_gateway_api_group.groupID,eo_gateway_api_group.groupName FROM eo_gateway_api_group WHERE gatewayID = ? AND isChild = 1 AND parentGroupID = ? ORDER BY eo_gateway_api_group.groupID DESC;`,gatewayID,group.GroupID)

			for childRows.Next(){
				var childGroup utils.ChildGroupInfo
				childRows.Scan(&childGroup.GroupID,&childGroup.GroupName)
				group.ChildGroupList = append(group.ChildGroupList,&childGroup)
			}
			if group.ChildGroupList == nil{
				childGroup := make([]*utils.ChildGroupInfo,0)
				group.ChildGroupList = childGroup
			}
			groupList = append(groupList,&group)
			num +=1
		}
	}
	if num == 0{
		flag =false
	}
	return flag,groupList
}

// 修改分组信息
func EditGroup(groupID,parentGroupID int,groupName string) bool{
	db := database.GetConnection()
	stmt,err := db.Prepare(`UPDATE eo_gateway_api_group SET eo_gateway_api_group.groupName = ?,eo_gateway_api_group.parentGroupID = ?,eo_gateway_api_group.isChild = ? WHERE eo_gateway_api_group.groupID = ?;`)

	defer stmt.Close()
	if err != nil {
		return false
	} 
	isChild := 1
	if parentGroupID == 0{
		isChild = 0
	}
	_, err = stmt.Exec(groupName,parentGroupID,isChild ,groupID)
	if err != nil {
		return false
	} else{
		return true
	}
}

// 判断分组与用户是否匹配
func CheckGroupPermission(groupID,userID int) bool{
	db := database.GetConnection()
	var gatewayID int
	err := db.QueryRow("SELECT eo_conn_management.gatewayID FROM eo_conn_management INNER JOIN eo_gateway_api_group ON eo_gateway_api_group.gatewayID = eo_conn_management.gatewayID WHERE userID = ? AND groupID = ?;",groupID,userID).Scan(&gatewayID)
	if err != nil {
		return false
	} else{
		return true
	}
}

// 获取分组名称
func GetGroupName(groupID int) (bool,string){
	db := database.GetConnection()
	var gatewayName string
	err := db.QueryRow("SELECT e o_gateway_api_group.groupName FROM eo_gateway_api_group WHERE eo_gateway_api_group.groupID = ?;",groupID).Scan(&gatewayName)
	if err != nil {
		return false,""
	} else{
		return true,gatewayName
	}
}
