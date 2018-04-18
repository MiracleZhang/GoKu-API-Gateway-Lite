package module

import (
	"apigateway/server/dao"
)

func EditPassword(userID int,oldPassword,newPassword string) bool{
	return dao.EditPassword(userID,oldPassword,newPassword)
}