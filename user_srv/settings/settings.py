import json
import uuid

import nacos
from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin
from loguru import logger

# 使用ReconnectMixin 来防止出现连接断开，来阻止查询失败
from common.reggister import consul


class ReconnectMysqlDatabase(PooledMySQLDatabase, ReconnectMixin):
    # python的mro
    def sequence_exists(self, seq):
        pass


NACOS = {
    "Host": "192.168.1.2",
    "Port": 8848,
    "NameSpace": "65b94fdb-5ecc-48a9-950e-df24f07e6e22",
    "User": "nacos",
    "Password": "nacos",
    "DataId": "user-srv.json",
    "Group": "dev"
}

client = nacos.NacosClient(f"{NACOS['Host']}:{NACOS['Port']}", namespace=NACOS['NameSpace'], username=NACOS['User'],
                           password=NACOS['Password'])
data = client.get_config(NACOS['DataId'], NACOS['Group'])
data = json.loads(data)

logger.info(f"配置信息:{data}")

# consul的配置
CONSUL_CLIENT = consul.ConsulRegister(data['consul']['host'], data['consul']['port'])

# 服务相关的配置
SERVICE_NAME = data['name']
SERVICE_TAGS = data['tags']
SERVICE_ID = str(uuid.uuid1())

DB = ReconnectMysqlDatabase(database=data['mysql']['db'], host=data['mysql']['host'], port=data['mysql']['port'],
                            user=data['mysql']['user'], password=data['mysql']['password'])


def update_cfg(args):
    global DB, CONSUL_CLIENT, SERVICE_NAME, SERVICE_TAGS
    new_data = json.loads(args['content'])

    logger.info(f"配置产生变化:{new_data}")

    DB = ReconnectMysqlDatabase(database=new_data['mysql']['db'], host=new_data['mysql']['host'],
                                port=new_data['mysql']['port'], user=new_data['mysql']['user'],
                                password=new_data['mysql']['password'])
    SERVICE_NAME = new_data['name']
    SERVICE_TAGS = new_data['tags']
    CONSUL_CLIENT = consul.ConsulRegister(new_data['consul']['host'], new_data['consul']['port'])
