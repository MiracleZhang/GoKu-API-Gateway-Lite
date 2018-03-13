package controller
import (
	_ "apigateway/utils"
    "apigateway/server/module"
	"github.com/gin-gonic/gin"
	"strconv"
)


func EditIPList(c *gin.Context){
	var userID int
	gatewayHashKey := c.PostForm("gatewayHashKey")
	
	if module.CheckLogin(c) == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"100000",})
		return 
	}else{
		result,_ := c.Request.Cookie("userID")
		userID,_ = strconv.Atoi(result.Value)
		flag := module.CheckGatewayPermission(gatewayHashKey,userID)
		if flag == false{
			c.JSON(200,gin.H{"statusCode":"100005","type":"guest",})
			return 
		}
	}
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	ipList := c.PostForm("ipList")
	chooseType := c.PostForm("chooseType")
	cType,_ := strconv.Atoi(chooseType)
	flag := module.EditIPList(gatewayID,cType,gatewayHashKey,ipList)
	if flag == true{
		c.JSON(200,gin.H{"type":"ip","statusCode":"000000"})
		return
	}else{
		c.JSON(200,gin.H{"type":"ip","statusCode":"180000"})
		return
	}
}

// 获取IP名单列表
func GetIPInfo(c *gin.Context){
	var userID int
	gatewayHashKey := c.PostForm("gatewayHashKey")
	
	if module.CheckLogin(c) == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"100000",})
		return 
	}else{
		result,_ := c.Request.Cookie("userID")
		userID,_ = strconv.Atoi(result.Value)
		flag := module.CheckGatewayPermission(gatewayHashKey,userID)
		if flag == false{
			c.JSON(200,gin.H{"statusCode":"100005","type":"guest",})
			return 
		}
	}
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)

	flag,ipList := module.GetIPList(gatewayID)
	if flag == true{
		c.JSON(200,gin.H{"type":"ip","statusCode":"000000","ipInfo":ipList})
		return
	}else{
		c.JSON(200,gin.H{"type":"ip","statusCode":"180000"})
	}
}

// 检查IP名单是否存在
func CheckIPListIsExist(c *gin.Context) {
	var userID int
	gatewayHashKey := c.PostForm("gatewayHashKey")
	
	if module.CheckLogin(c) == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"100000",})
		return 
	}else{
		result,_ := c.Request.Cookie("userID")
		userID,_ = strconv.Atoi(result.Value)
		flag := module.CheckGatewayPermission(gatewayHashKey,userID)
		if flag == false{
			c.JSON(200,gin.H{"statusCode":"100005","type":"guest",})
			return 
		}
	}
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	flag := module.CheckIPListIsExist(gatewayID)
	if flag == true{
		c.JSON(200,gin.H{"type":"ip","statusCode":"000000"})
		return
	}else{
		c.JSON(200,gin.H{"type":"ip","statusCode":"180000"})
	}
}
 