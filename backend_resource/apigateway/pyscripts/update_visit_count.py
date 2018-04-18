#!/usr/bin/python
#coding: utf-8

import redis
import MySQLdb
import time,ConfigParser

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

def updateVisitCount():
    db = MySQLdb.connect(db=db_name, host=db_host, user=db_user, passwd=db_pass)
    r = redis.StrictRedis(host=redis_host, port=6379, db=redis_db, password=redis_pass)
    cursor = db.cursor()
    m = {}
    for key in r.keys('gatewayDayCount:*'):
        domain, date = key.split(':')[1:]
        count = r.get(key)
        cursor.execute('select gatewayID from eo_gateway where hashKey=%s;', (domain,))
        result = cursor.fetchall()
        if len(result) == 0:
            continue
        gatewayID = int(result[0][0])
        m[(gatewayID, date)] = int(count)

    for key in m:
        gatewayID, dateStr = key
        count = m[key]
        cursor.execute('select * from eo_gateway_count where gatewayID=%s and date=%s', (gatewayID, dateStr))
        if len(cursor.fetchall()) == 0:
            cursor.execute('insert into eo_gateway_count values(%s, %s, %s);', (gatewayID, count, dateStr))
        else:
            cursor.execute('update eo_gateway_count set visitCount=%s where gatewayID=%s and date=%s', (count, gatewayID, dateStr))
    db.commit()
    db.close()

while True:
    updateVisitCount()
    print 'Last update:', time.ctime()
    time.sleep(3600*2)
