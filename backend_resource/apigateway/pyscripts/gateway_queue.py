#-*- coding: utf-8 -*-

import redis
import MySQLdb
import json
import datetime
import time,ConfigParser

cf = ConfigParser.ConfigParser()
cf.read("db_conf.conf")

redis_host = cf.get("redis","redis_host")
redis_pass = cf.get("redis","redis_pass")
redis_db = cf.get("redis","redis_db")

db_host = cf.get("db","db_host")
db_user = cf.get("db","db_user")
db_pass = cf.get("db","db_pass")
db_name = cf.get("db","db_name")

r = redis.StrictRedis(host=redis_host, port=6379, db=redis_db, password=redis_pass)
sql = ('select gatewayProtocol, gatewayRequestType, gatewayRequestURI ' + 
    'from eo_gateway_api ' + 
    'where gatewayID=%s')

def mysql_conn():
    conn=MySQLdb.connect(db=db_name,host=db_host,user=db_user,passwd=db_pass)
    return conn


def passAddGateway(hash_key, token):
    r.set('gatewayToken:' + hash_key, token)


def passDeleteGateway(hash_key):
    r.delete('apiList:' + hash_key)
    r.delete('gatewayToken:' + hash_key)

def deleteApiInfo(gateway_id,hash_key):
    # 删除apiInfo
    keys = r.keys("apiInfo:" + hash_key + "*")
    if len(keys)>0:
        r.delete(*keys)

def loadAPIList(gateway_id, hash_key):
    conn = mysql_conn()
    cursor = conn.cursor()
    cursor.execute(sql, (gateway_id, ))
    apis = []
    for api in cursor.fetchall():
        print api
        protocol, requestType, uri = (str(api[0]), str(api[1]), api[2])
        apis.append(protocol + ':' + requestType + ':' + uri)
    print 'apis:', apis
    listname = 'apiList:' + hash_key
    cursor.close()
    conn.close()
    with r.pipeline() as pipe:
        while True:
            try:
                pipe.watch(listname)
                original_len = pipe.llen(listname)
                print 'original_len', original_len
                pipe.multi()
                print 'api', '*' * 10
                for api in apis:
                    print api
                    pipe.rpush(listname, api)
                for i in xrange(original_len):
                    pipe.lpop(listname)
                pipe.execute()
                break
            except redis.WatchError:
                continue
    

last = datetime.datetime.now()
print last

while True:
    task = r.blpop('gatewayQueue', 1)

    # 检查任务队列是否为空
    if task == None:
        # 是否退出程序
        if r.get('gatewayQueueCloseSignal'):
            break
        # 打印时间
        now = datetime.datetime.now()
        if now > last + datetime.timedelta(minutes=1):
            last = now
            print last
        continue

    # 处理任务
    task = json.loads(task[1])
    print datetime.datetime.now(), ':', task
    if task['type'] == 'gateway':
        data = task['data']
        if task['operation'] == 'add':
            gateway_id, hash_key, token, gateway_area = data['gatewayID'], data['gatewayHashKey'], data['token'], data['gatewayArea']
            passAddGateway(hash_key, token)
        elif task['operation'] == 'delete':
            hash_key = data['gatewayHashKey']
            passDeleteGateway(hash_key)
    elif task['type'] == 'api':
        data = task['data']
        gateway_id, hash_key = data['gatewayID'], data['gatewayHashKey']
        time.sleep(1)
        loadAPIList(gateway_id, hash_key)
    elif task['type'] == "backend":
        data = task['data']
        gateway_id, hash_key = data['gatewayID'], data['gatewayHashKey']
        deleteApiInfo(gateway_id, hash_key)
    else:
        r.rpush('gatewayQueue', json.dumps(task))
        time.sleep(1)


