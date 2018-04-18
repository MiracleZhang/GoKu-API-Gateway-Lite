package controller
import (
	"apigateway/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"apigateway/server/dao"
	"strconv"
	"regexp"
)

// 用户登录
func Login(c *gin.Context){
	loginCall := c.PostForm("loginCall")
	loginPassword := c.PostForm("loginPassword")
	if match, _ := regexp.MatchString("^[0-9a-zA-Z][0-9a-zA-Z_]{3,63}$", loginCall);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"120002"})
		return 
	}else if match, _ := regexp.MatchString("^[0-9a-zA-Z]{32}$", loginPassword);match == false{
		c.JSON(200,gin.H{"type":"guest","statusCode":"120004"})
		return 
	}else{
		flag,userID := dao.Login(loginCall,utils.Md5(loginPassword))
		if flag == true{
			userCookie := http.Cookie{Name: "userID", Value:strconv.Itoa(userID), Path: "/", MaxAge: 86400}
			http.SetCookie(c.Writer, &userCookie)
			c.JSON(200,gin.H{"type":"guest","statusCode":"000000",})
			return
		}else{
			c.JSON(200,gin.H{"type":"guest","statusCode":"120000",})
			return 
		}
	}
}


