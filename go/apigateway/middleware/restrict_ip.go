package middleware

import (
	"apigateway/conf"
	"apigateway/dao"
	"apigateway/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/farseer810/yawf"
)

func IPValve(httpRequest *http.Request, context yawf.Context,
	httpResponse http.ResponseWriter, headers yawf.Headers) (bool, string) {
	var info *utils.IPListInfo
	var gatewayHashkey string
	gatewayHashkey = httpRequest.RequestURI[1:41]

	remoteAddr := httpRequest.RemoteAddr
	remoteIP := interceptIP(remoteAddr, ":")

	minuteCount := dao.GetIPMinuteCount(context, remoteIP)
	fmt.Println(minuteCount)
	minuteCountLimit, _ := strconv.Atoi(conf.Configure["ip_minute_visit_limit"])
	if minuteCountLimit < minuteCount {
		httpResponse.WriteHeader(403)
		dao.UpdateBlackList(context, remoteIP)
		return false, "ip visit limit exceeded"
	}
	info = dao.GetIPList(context, gatewayHashkey)
	if info == nil{
		return true,""
	}
	chooseType := info.ChooseType

	if chooseType == 1 {
		for _, ipList := range info.IPList {
			if ipList == remoteIP {
				fmt.Println(remoteIP)
				httpResponse.WriteHeader(403)
				return false, "illegal IP"
			}
		}
	} else if chooseType == 2 {
		for _, ipList := range info.IPList {
			if ipList == remoteIP {
				return true, ""
			}
		}
		fmt.Println(remoteIP)
		httpResponse.WriteHeader(403)
		return false, "illegal IP"
	}

	return true, ""
}

func interceptIP(str, substr string) string {
	result := strings.Index(str, substr)
	var rs string
	if result > 7 {
		rs = str[:result]
	}
	return rs
}
