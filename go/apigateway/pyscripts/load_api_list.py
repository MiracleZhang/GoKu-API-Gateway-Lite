#!/usr/bin/python
#coding: utf-8

import redis
import MySQLdb,ConfigParser

cf = ConfigParser.ConfigParser()

cf.read("db_conf.conf")
db_host = cf.get("db","db_host")
db_user = cf.get("db","db_user")
db_pass = cf.get("db","db_pass")
db_name = cf.get("db","db_name")
redis_host = cf.get("redis","redis_host")
redis_pass = cf.get("redis","redis_pass")
redis_db = cf.get("redis","redis_db")
db = MySQLdb.connect(db=db_name, host=db_host, user=db_user, passwd=db_pass)


r = redis.StrictRedis(host=redis_host, port=6379, db=redis_db, password=redis_pass)
cursor = db.cursor()
cursor.execute('select g.hashKey, a.gatewayProtocol, a.gatewayRequestType, a.gatewayRequestURI ' + 
                'from eo_gateway as g, eo_gateway_api as a ' + 
                'where g.gatewayID=a.gatewayID')

for key in r.keys('*'):
    r.delete(key)

m = {}
for data in cursor.fetchall():
    hashKey, protocol, requestType, uri = data
    key = 'apiList:' + hashKey.lower()
    if m.get(key) == None:
        m[key] = []
    m[key].append(str(protocol) + ":" + str(requestType) + ":" + uri)

cursor.execute('select hashKey, token from eo_gateway;')
for data in cursor.fetchall():
    hashKey, token = data
    r.set('gatewayToken:' + hashKey, token)

for key in m:
    r.delete(key)
    r.rpush(key, *m[key])

