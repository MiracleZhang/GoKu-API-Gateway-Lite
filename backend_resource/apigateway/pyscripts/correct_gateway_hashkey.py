#!/usr/bin/python
# coding: utf-8
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

cursor.execute('select cacheID, apiJson, gatewayHashKey from eo_gateway_api_cache')
# 将hashKey全部变为小写保存在缓存表中
m = {}
for data in cursor.fetchall():
    cacheID, apiJson, hashKey = data
    apiJson = json.loads(apiJson)
    hashKey = hashKey.lower()
    apiJson['gatewayHashKey'] = hashKey
    apiJson = json.dumps(apiJson)
    m[cacheID] = [apiJson, hashKey]

for id in m:
    apiJson, hashKey = m[id]
    print cursor.execute('update eo_gateway_api_cache set apiJson=%s, gatewayHashKey=%s where cacheID=%s', (apiJson, hashKey, id))

db.commit()