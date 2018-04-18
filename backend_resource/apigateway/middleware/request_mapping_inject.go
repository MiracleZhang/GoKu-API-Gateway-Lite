package middleware

import (
	"apigateway/conf"
	"apigateway/dao"
	"net/http"

	"fmt"

	"github.com/farseer810/yawf"
)

var (
	methodIndicator = map[string]string{"POST": "0", "GET": "1", "PUT": "2", "DELETE": "3", "HEAD": "4",
		"OPTIONS": "5", "PATCH": "6"}
)

func isURIMatched(context yawf.Context, incomingURI, testURI string) bool {
	isMatched := incomingURI == testURI
	return isMatched
}

//注入请求映射
func InjectRequestMapping(httpRequest *http.Request, context yawf.Context,
	httpResponse http.ResponseWriter, headers yawf.Headers) (bool, string) {
	var domain, method, scheme, gatewayHashkey, requestURL string

	// TODO: 0 for http, 1 for https
	scheme = "0"
	fmt.Println(httpRequest.RemoteAddr)
	method = methodIndicator[httpRequest.Method]
	if method == "" {
		httpResponse.WriteHeader(404)
		return false, ""
	}

	domain = httpRequest.Host
	requestURL = httpRequest.RequestURI
	fmt.Println(len(requestURL))
	fmt.Println(requestURL)
	fmt.Println(domain)
	if conf.Configure["is_debug"] != "true" && len(requestURL) <41 {
		httpResponse.WriteHeader(404)
		return false, ""
	}

	gatewayHashkey = requestURL[1:41]
	token := dao.GetGatewayToken(context, gatewayHashkey)
	
	if token == "" {
		httpResponse.WriteHeader(404)
		return false, ""
	}

	if token != headers["Eo-Gateway-Token"] {
		httpResponse.WriteHeader(401)
		return false, ""
	}
	fmt.Println(gatewayHashkey)
	paths := dao.GetAllAPIPaths(context, gatewayHashkey)
	fmt.Println(paths)
	var matchedURI string
	for _, uri := range paths {
		if uri[0:4] != scheme+":"+method+":" {
			continue
		}
		if isURIMatched(context, httpRequest.URL.Path[41:], uri[4:]) {
			matchedURI = uri
		}
	}
	if matchedURI == "" {
		httpResponse.WriteHeader(404)
		return false, ""
	}
	info := dao.GetMapping(context, gatewayHashkey, matchedURI)
	context.Map(info)

	return true, ""
}
