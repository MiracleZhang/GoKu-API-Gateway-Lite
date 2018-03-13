#!/usr/bin/python
#coding: utf-8

import MySQLdb
import json,ConfigParser

cf = ConfigParser.ConfigParser()

cf.read("db_conf.conf")
db_host = cf.get("db","db_host")
db_user = cf.get("db","db_user")
db_pass = cf.get("db","db_pass")
db_name = cf.get("db","db_name")

db = MySQLdb.connect(db=db_name, host=db_host, user=db_user, passwd=db_pass)
cursor = db.cursor()


new_json_str = """{
  "requestParams": [
    {
      "gwParamKey": "Accept-Language",
      "paramType": "0",
      "backendParamPosition": "1",
      "isNotNull": "0",
      "gwParamPosition": "0",
      "backendParamKey": "accept_language_from_header"
    },
    {
      "gwParamKey": "name",
      "paramType": "0",
      "backendParamPosition": "1",
      "isNotNull": "0",
      "gwParamPosition": "2",
      "backendParamKey": "name"
    },
    {
      "gwParamKey": "name",
      "paramType": "0",
      "backendParamPosition": "2",
      "isNotNull": "0",
      "gwParamPosition": "2",
      "backendParamKey": "name"
    }
  ],
  "constantParams": [
    {
      "paramValue": "\u8868\u5355\u5e38\u91cf\u53c2\u6570\u503c",
      "backendParamKey": "const_body_param",
      "paramName": "\u8868\u5355\u5e38\u91cf\u53c2\u6570\u8bf4\u660e",
      "paramPosition": "1"
    }
  ],
  "backendRequestType": "0",
  "backendProtocol": "0",
  "backendURI": "webtest.farseer810.cn/eotest/test1",
  "gatewayHashKey": "ugcal5jf48028b5438802a3b15ce0d6ea21dd4c6ff6b697",
  "groupID": "7",
  "gatewayID": 8,
  "gatewayRequestType": "1",
  "gatewayProtocol": "0",
  "isRequestBody": "0"
}"""
cursor.execute('update eo_gateway_api_cache set apiJson=%s where apiID=32;', (json.dumps(json.loads(new_json_str)),))
db.commit()
cursor.execute('select apiJson from eo_gateway_api_cache where apiID=32;')
info = json.loads(cursor.fetchall()[0][0])
print info


