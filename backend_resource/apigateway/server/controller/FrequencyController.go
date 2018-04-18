package controller
import (
	_ "apigateway/utils"
    "apigateway/server/module"
	"github.com/gin-gonic/gin"
	"strconv"
	"regexp"
)

func AddFrequencyLimit(c *gin.Context) {
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
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)

	count := c.PostForm("count")
	intervalType := c.PostForm("intervalType")
	if match, _ := regexp.MatchString("^[1-9][0-9]*$", count);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"170002"})
		return 
	}else if match, _ := regexp.MatchString("^[01]$", intervalType);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"170001"})
		return 
	}else{
		cou,_ := strconv.Atoi(count)
		
		iType,_ := strconv.Atoi(intervalType)
		
		flag := module.AddFrequencyLimit(gatewayHashKey,gatewayID,cou,iType)
		if flag == true{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"000000",})
			return
		}else{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"170000",})
			return
		}
	}
}

// 修改频率限制
func EditFrequencyLimit(c *gin.Context) {
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
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	count := c.PostForm("count")
	intervalType := c.PostForm("intervalType")
	if match, _ := regexp.MatchString("^[1-9][0-9]*$", count);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"170002"})
		return 
	}else if match, _ := regexp.MatchString("^[01]$", intervalType);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"170001"})
		return 
	}else{
		cou,_ := strconv.Atoi(count)
		iType,_ := strconv.Atoi(intervalType)
		flag := module.EditFrequencyLimit(gatewayHashKey,gatewayID,cou,iType)
		if flag == true{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"000000"})
			return
		}else{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"170000"})
			return
		}
	}
}

// 删除频率限制
func DeleteFrequencyLimit(c *gin.Context){
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
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	intervalType := c.PostForm("intervalType")
	iType,_ := strconv.Atoi(intervalType)
	flag := module.DeleteFrequencyLimit(gatewayHashKey,gatewayID,iType)
	if flag == true{
		c.JSON(200,gin.H{"type":"frequency","statusCode":"000000"})
		return
	}else{
		c.JSON(200,gin.H{"type":"frequency","statusCode":"170000"})
		return
	}
}

// 获取频率限制列表
func GetFrequencyLimitList(c *gin.Context) {
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
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	flag,frequencyList := module.GetFrequencyLimitList(gatewayID)
	if flag == true{
		c.JSON(200,gin.H{"type":"frequency","statusCode":"000000","frequencyList":frequencyList})
		return
	}else{
		c.JSON(200,gin.H{"type":"frequency","statusCode":"170000"})
		return
	}
	
}

// 检查频率限制是否存在
func CheckFrequencyLimitIsExist(c *gin.Context) {
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
	var gatewayID int
	_,gatewayID = module.GetIDFromHashKey(gatewayHashKey)
	intervalType := c.PostForm("intervalType")
	if match, _ := regexp.MatchString("^[01]$", intervalType);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"170001"})
		return 
	}else{
		iType,_ := strconv.Atoi(intervalType)
		flag := module.CheckFrequencyLimitIsExist(gatewayID,iType)
		if flag == true{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"000000"})
			return
		}else{
			c.JSON(200,gin.H{"type":"frequency","statusCode":"170000"})
			return
		}
	}
}