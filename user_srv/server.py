import argparse
import signal
import sys
from concurrent import futures

sys.path.insert(0, r"/Users/he/Documents/Python/mxshop_srvs")

import grpc
from loguru import logger

from user_srv.handler.user import UserServices
from user_srv.proto import user_pb2_grpc

logger.add("logs/user_srv{time}.log")


def on_exit(signo, frame):
    logger.info("进程中断")
    sys.exit(0)


def serve():
    parser = argparse.ArgumentParser()
    parser.add_argument('--ip',
                        nargs="?",
                        type=str,
                        default="127.0.0.1",
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
    user_pb2_grpc.add_UserServicer_to_server(UserServices(), server)

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

    print(f"启动服务: {args.ip}:{args.port}")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
