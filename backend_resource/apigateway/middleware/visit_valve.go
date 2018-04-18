package middleware

import (
	"apigateway/conf"
	"apigateway/dao"
	"apigateway/utils"
	"net/http"
	"strconv"

	"github.com/farseer810/yawf"
)

func GatewayValve(httpRequest *http.Request,context yawf.Context, info *utils.MappingInfo, httpResponse http.ResponseWriter) (bool, string) {
	var gatewayHashkey string
	gatewayHashkey = httpRequest.RequestURI[1:41]
	valveList := dao.GetGatewayValve(context,gatewayHashkey)
	minuteCount := dao.GetGatewayMinuteCount(context, info)
	var minuteValve = 0
	var secondValve = 0
	for _, valveInfo := range valveList {
		if valveInfo.IntervalType == 0 {
			secondValve = valveInfo.Count
		}else if valveInfo.IntervalType == 1{
			minuteValve = valveInfo.Count
		}
	}
	if minuteValve == 0{
		minuteValve, _ = strconv.Atoi(conf.Configure["minute_visit_limit"])
	}
	if minuteValve <= minuteCount {
		httpResponse.WriteHeader(403)
		return false, "minute visit limit exceeded"
	}
	secondCount := dao.GetGatewaySecondCount(context,info)
	if secondValve <= secondCount && secondValve != 0{
		httpResponse.WriteHeader(403)
		return false, "second visit limit exceeded"
	}
	return true,""
}
