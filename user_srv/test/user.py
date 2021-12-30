import grpc

from user_srv.proto import user_pb2, user_pb2_grpc


class UserTest:
    def __init__(self):
        _channel = grpc.insecure_channel('127.0.0.1:50051')
        self.stub = user_pb2_grpc.UserStub(_channel)

    def user_list(self):
        rsp = self.stub.GetUserList(user_pb2.PageInfo(pn=2, pSize=2))
        print(rsp.total)
        for _user in rsp.data:
            print(_user.mobile, _user.birthDay)

    def user_id(self, id):
        rsp = self.stub.GetUserById(user_pb2.IdReqeust(ID=id))
        print(rsp.mobile, rsp.birthDay)

    def create_uesr(self, nick_name, password, mobile):
        rsp: user_pb2.UserInfoResponse = self.stub.CreateUser(user_pb2.CreateUserInfo(
            nickName=nick_name,
            password=password,
            mobile=mobile
        ))
        print(rsp.id)


if __name__ == '__main__':
    user = UserTest()
    # user.user_list()
    # user.user_id(11)
    user.create_uesr("booby", "18787878787", "admin1234")