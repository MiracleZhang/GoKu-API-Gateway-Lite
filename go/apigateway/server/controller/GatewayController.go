package controller
import (
	_ "apigateway/utils"
	"apigateway/server/module"
	"github.com/gin-gonic/gin"
	"strings"
	"strconv"
	"apigateway/conf"
)

// 新增网关
func AddGateway(c *gin.Context){
	var userID int
	
	if module.CheckLogin(c) == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"100000",})
		return 
	}else{
		result,_ := c.Request.Cookie("userID")
		userID,_ = strconv.Atoi(result.Value)
	}
    gatewayName := c.PostForm("gatewayName")
	gatewayDesc := c.PostForm("gatewayDesc")
	gatewayArea := c.DefaultPostForm("gatewayArea","0")
	gatewayNameLen := strings.Count(gatewayName,"")-1
	if gatewayNameLen<1 || gatewayNameLen > 32 {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130001",})
		return 
	}
	flag,gatewayHashkey := module.Addgateway(gatewayName,gatewayDesc,gatewayArea,userID)
	if flag == false {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130000",})
		return 
	}else{
		c.JSON(200,gin.H{"type":"gateway","statusCode":"000000","gatewayHashKey":gatewayHashkey,})
		return 
	}
}

// 修改网关
func EditGateway(c *gin.Context){
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
			c.JSON(200,gin.H{"statusCode":"100005","type":"gateway",})
			return 
		}
	}
	gatewayName := c.PostForm("gatewayName")
	gatewayDesc := c.PostForm("gatewayDesc")
	gatewayArea := c.DefaultPostForm("gatewayArea","0")
	gatewayNameLen := strings.Count(gatewayName,"")-1
	if gatewayNameLen<1 || gatewayNameLen > 32 {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130001",})
		return 
	}
	flag := module.EditGateway(gatewayName,gatewayArea,gatewayDesc,gatewayHashKey,userID)
	if flag == false {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130000",})
		return 
	}else{
		c.JSON(200,gin.H{"type":"gateway","statusCode":"000000"})
		return 
	}
}

// 删除网关
func DeleteGateway(c *gin.Context){
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
			c.JSON(200,gin.H{"statusCode":"100005","type":"gateway",})
			return 
		}
	}
	flag := module.DeleteGateway(gatewayHashKey,userID)
	if flag == false {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130000",})
		return 
	}else{
		c.JSON(200,gin.H{"type":"gateway","statusCode":"000000"})
		return 
	}
}

// 获取网关信息
func GetGatewayInfo(c *gin.Context){
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
			c.JSON(200,gin.H{"statusCode":"100005","type":"gateway",})
			return 
		}
	}
	flag,result := module.GetGatewayInfo(gatewayHashKey,userID)
	if flag == false {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130000",})
		return 
	}else{
		gatewayPort := conf.Configure["eotest_port"]
		result.GatewayPort = gatewayPort
		c.JSON(200,gin.H{"type":"gateway","statusCode":"000000","gatewayInfo":result})
		return 
	}
}


// 获取网关列表
func GetGatewayList(c *gin.Context){
	var userID int
	if module.CheckLogin(c) == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"100000",})
		return 
	}else{
		result,_ := c.Request.Cookie("userID")
		userID,_ = strconv.Atoi(result.Value)
	}
	gatewayArea := c.DefaultPostForm("gatewayArea","0")
	area,err := strconv.Atoi(gatewayArea)
	if err != nil{
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130002",})
		return 
	}
	flag,gatewayList := module.GetGatewayList(area,userID)
	if flag == false {
		c.JSON(200,gin.H{"type":"gateway","statusCode":"130000",})
		return 
	}else{
		c.JSON(200,gin.H{"type":"gateway","statusCode":"000000","gatewayList":gatewayList,})
		return 
	}
}


