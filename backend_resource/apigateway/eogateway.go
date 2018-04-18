package main

import (
	"apigateway/server/controller"
    "github.com/gin-gonic/gin"
    "github.com/fvbock/endless"
    "log"
    "os"
)
func main() {
    router := gin.Default()
    web := router.Group("/Web")
    {
        guest :=web.Group("/Guest")
        {
            guest.POST("/login",controller.Login)
        }
        user :=web.Group("/User")
        {
            user.POST("/logout",controller.Logout)
            user.POST("/editPassword",controller.EditPassword)
            user.POST("/checkLogin",controller.CheckLogin)
        }
        group := web.Group("/Group")
        {
            group.POST("/addGroup",controller.AddGroup)
            group.POST("/editGroup",controller.EditGroup)
            group.POST("/deleteGroup",controller.DeleteGroup)
            group.POST("/getGroupList",controller.GetGroupList)
            group.POST("/getGroupName",controller.GetGroupName)
        }
        api := web.Group("/Api")
        {
            api.POST("/addApi",controller.AddApi)
            api.POST("/editApi",controller.EditApi)
            api.POST("/deleteApi",controller.DeleteApi)
            api.POST("/getApiList",controller.GetApiListOrderByName)
            api.POST("/getAllApiList",controller.GetAllApiListOrderByName)
            api.POST("/getApi",controller.GetApi)
            api.POST("/searchApi",controller.SearchApi)
            api.POST("/checkGatewayURLIsExist",controller.CheckGatewayURLIsExist)
        }
        backend := web.Group("/Backend")
        {
            backend.POST("/addBackend",controller.AddBackend)
            backend.POST("/editBackend",controller.EditBackend)
            backend.POST("/deleteBackend",controller.DeleteBackend)
            backend.POST("/getBackendList",controller.GetBackendList)
            backend.POST("/getBackend",controller.GetBackendInfo)
        }
        gateway := web.Group("/Gateway")
        {
            gateway.POST("/addGateway",controller.AddGateway)
            gateway.POST("/editGateway",controller.EditGateway)
            gateway.POST("/deleteGateway",controller.DeleteGateway)
            gateway.POST("/getGatewayList",controller.GetGatewayList)
            gateway.POST("/getGateway",controller.GetGatewayInfo)

        }
        ip := web.Group("/IP")
        {
            ip.POST("/editIPList",controller.EditIPList)
            ip.POST("/getIPInfo",controller.GetIPInfo)
            ip.POST("/checkIPListIsExist",controller.CheckIPListIsExist)
        }
        frequency := web.Group("/Frequency")
        {
            frequency.POST("/editFrequency",controller.EditFrequencyLimit)
            frequency.POST("/addFrequency",controller.AddFrequencyLimit)
            frequency.POST("/getFrequencyList",controller.GetFrequencyLimitList)
            frequency.POST("/checkFrequencyLimitIsExist",controller.CheckFrequencyLimitIsExist)
            frequency.POST("/deleteFrequency",controller.DeleteFrequencyLimit)
        }
        install := web.Group("/Install")
        {
            install.POST("/checkDBConnect",controller.CheckDBConnect)
            install.POST("/checkRedisConnect",controller.CheckRedisConnect)
            install.POST("/checkIsInstall",controller.CheckIsInstall)
            install.POST("/installConfigure",controller.InstallConfigure)
            install.POST("/install",controller.Install)
            install.POST("/test")
        }
    }
    log.Println(os.Getpid())
    err := endless.ListenAndServe(":8080",router)
    if err != nil {
		log.Println(err)
	}
	log.Println("Server on 8080 stopped")
	os.Exit(0)
}
