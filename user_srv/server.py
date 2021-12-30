from concurrent import futures

import grpc

from user_srv.handler.user import UserServices
from user_srv.proto import user_pb2_grpc


def serve():
    # 创建一个链接
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # 注册逻辑
    user_pb2_grpc.add_UserServicer_to_server(UserServices(), server)

    # 开启逻辑
    server.add_insecure_port("0.0.0.0:50051")
    print("启动服务: 127.0.0.1:50051")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
