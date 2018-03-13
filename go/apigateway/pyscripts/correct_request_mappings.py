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
cursor.execute('select cacheID, apiJson from eo_gateway_api_cache;')
# 删除部分字段
m = {}
for data in cursor.fetchall():
    id = data[0]
    data = json.loads(data[1])
    for param in data['requestParams']:
        del param['checkbox']
        del param['$$hashKey']
        param['gwParamPosition'] = param['gwParamPostion']
        del param['gwParamPostion']
    for param in data['constantParams']:
        del param['$$hashKey']
    data = json.dumps(data)
    m[id] = data

for id in m:
    print cursor.execute('update eo_gateway_api_cache set apiJson=%s where cacheID=%s', (m[id], id))

db.commit()
