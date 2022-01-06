import argparse
import signal
import sys
from concurrent import futures

import grpc
from loguru import logger

from common.grpc_health.v1 import health_pb2_grpc
from common.grpc_health.v1.health import HealthServicer
from common.reggister import consul
from user_srv.handler.user import UserServicer
from user_srv.proto import user_pb2_grpc
from user_srv.settings import settings

logger.add("logs/user_srv{time}.log")
consulClient = consul.ConsulRegister(settings.CONSUL_HOST, settings.CONSUL_PORT)


def on_exit(signo, frame):
    logger.info("进程中断")
    consulClient.deregister(settings.SERVICE_NAME)
    sys.exit(0)


def serve():
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip',
                        nargs="?",
                        type=str,
                        default="192.168.1.2",
                        help="binding ip")

    parser.add_argument('--port',
                        nargs="?",
                        type=int,
                        default=50051,
                        help="the listening port")

    args = parser.parse_args()
    # 创建一个链接
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # 注册逻辑
    user_pb2_grpc.add_UserServicer_to_server(UserServicer(), server)

    health_pb2_grpc.add_HealthServicer_to_server(HealthServicer(), server)

    # 开启逻辑
    server.add_insecure_port(f"{args.ip}:{args.port}")

    # 主线程退出信号监听
    """
        windows下支持的信号是有限的：
            SIGINT ctrl+c 终端
            SIGTERM kill发出的软件终止
    """
    signal.signal(signal.SIGINT, on_exit)
    signal.signal(signal.SIGTERM, on_exit)

    logger.info(f"启动服务: {args.ip}:{args.port}")
    server.start()

    logger.info(f"服务注册开始")
    consul_register_status = consulClient.register(settings.SERVICE_NAME, settings.SERVICE_NAME, args.ip, args.port,
                                                   settings.SERVICE_TAGS)
    if not consul_register_status:
        logger.error("服务注册失败")
        sys.exit(0)
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
