package controller

import (
	"apigateway/dao"
	"apigateway/utils"
	"net/http"
	"strings"

	"github.com/farseer810/requests"
	"github.com/farseer810/yawf"
)

func CreateRequest(httpRequest *http.Request, info *utils.MappingInfo, queryParams yawf.QueryParams,
	formParams yawf.FormParams, httpResponse http.ResponseWriter, context yawf.Context) (int, []byte) {
	var backendDomain string
	if info.BackendProtocol == "0" {
		backendDomain = "http://" + info.BackendURI + info.BackendPath
	} else {
		backendDomain = "https://" + info.BackendURI + info.BackendPath
	}

	session := requests.NewSession()
	request, err := session.Method(info.BackendRequestType, backendDomain)
	if err != nil {
		panic(err)
	}

	var backendHeaders map[string][]string = make(map[string][]string)
	var backendQueryParams map[string][]string = make(map[string][]string)
	var backendFormParams map[string][]string = make(map[string][]string)

	for _, reqParam := range info.RequestParams {
		var param []string

		switch reqParam.ParamPosition {
		case "header":
			param = httpRequest.Header[reqParam.ParamKey]
		case "body":
			if httpRequest.Method == "POST" || httpRequest.Method == "PUT" || httpRequest.Method == "PATCH" {
				param = formParams[reqParam.ParamKey]
			} else {
				continue
			}
		case "query":
			param = queryParams[reqParam.ParamKey]
		}
		if param == nil {
			if reqParam.IsNotNull {
				// missing required parameters
				return 400, []byte("")
			} else {
				continue
			}
		}
		switch reqParam.BackendParamPosition {
		case "header":
			backendHeaders[reqParam.BackendParamKey] = param
		case "body":
			if info.BackendRequestType == "POST" || info.BackendRequestType == "PUT" || info.BackendRequestType == "PATCH" {
				backendFormParams[reqParam.BackendParamKey] = param
			}
		case "query":
			backendQueryParams[reqParam.BackendParamKey] = param
		}
	}

	for _, constParam := range info.ConstantParams {
		switch constParam.ParamPosition {
		case "header":
			backendHeaders[constParam.BackendParamKey] = []string{constParam.ParamValue}
		case "body":
			if info.BackendRequestType == "POST" || info.BackendRequestType == "PUT" || info.BackendRequestType == "PATCH" {
				backendFormParams[constParam.BackendParamKey] = []string{constParam.ParamValue}
			} else {
				backendQueryParams[constParam.BackendParamKey] = []string{constParam.ParamValue}
			}
		}
	}

	for key, values := range backendHeaders {
		request.SetHeader(key, values...)
	}
	for key, values := range backendQueryParams {
		request.SetQueryParam(key, values...)
	}
	for key, values := range backendFormParams {
		request.SetFormParam(key, values...)
	}

	cookies := make(map[string]string)
	for _, cookie := range httpRequest.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	session.SetCookies(cookies)

	res, err := request.Send()
	if err != nil {
		panic(err)
	}

	httpResponseHeader := httpResponse.Header()
	for key, _ := range httpResponseHeader {
		httpResponseHeader[key] = nil
	}
	for key, values := range res.Headers() {
		httpResponseHeader[key] = values
	}
	index := strings.Index(httpRequest.RemoteAddr, ":")
	remoteIP := httpRequest.RemoteAddr[:index]
	// go dao.UpdateIPVisitCount(context, remoteIP, res)
	go dao.UpdateVisitCount(context, info, res,remoteIP)

	return res.StatusCode(), res.Body()
}
