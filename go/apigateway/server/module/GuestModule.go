package module

import (
	"apigateway/server/dao"
	"github.com/gin-gonic/gin"
)


func Login(loginCall,loginPassword string) (bool,int){
	return dao.Login(loginCall,loginPassword)
}

func CheckLogin(c *gin.Context) bool{
	_,err := c.Request.Cookie("userID")
	if err != nil{
		return false
	}else{
		return true
	}
}


func Register(loginCall,loginPassword string) (bool){
	return dao.Register(loginCall,loginPassword)
}
