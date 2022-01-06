from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin


# 使用ReconnectMixin 来防止出现连接断开，来阻止查询失败
class ReconnectMysqlDatabase(PooledMySQLDatabase, ReconnectMixin):
    # python的mro
    def sequence_exists(self, seq):
        pass


MYSQL_DB = "mxshop_user_srv"
MYSQL_HOST = "127.0.0.1"
MYSQL_PORT = 3306
MYSQL_USER = "root"
MYSQL_PASSWORD = "123456"

# consul的配置
CONSUL_HOST = "192.168.1.2"
CONSUL_PORT = "8500"

# 服务相关的配置
SERVICE_NAME = "user-srv"
SERVICE_TAGS = ["imooc", "bobby", "python", "srv"]

DB = ReconnectMysqlDatabase(database=MYSQL_DB, host=MYSQL_HOST, port=MYSQL_PORT, user=MYSQL_USER,
                            password=MYSQL_PASSWORD)
