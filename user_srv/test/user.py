import grpc

from user_srv.proto import user_pb2, user_pb2_grpc


class UserTest:
    def __init__(self):
        _channel = grpc.insecure_channel('localhost:50051')
        self.stub = user_pb2_grpc.UserStub(_channel)

    def user_list(self):
        rsp = self.stub.GetUserList(user_pb2.PageInfo(pn=2, pSize=2))
        print(rsp.total)
        for user in rsp.data:
            print(user.mobile, user.birthDay)


if __name__ == '__main__':
    user = UserTest()
    user.user_list()
