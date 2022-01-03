import time

import grpc
from loguru import logger
from passlib.handlers.pbkdf2 import pbkdf2_sha256

from user_srv.model.model import User
from user_srv.proto import user_pb2, user_pb2_grpc


class UserServices(user_pb2_grpc.UserServicer):
    @staticmethod
    def convert_user_to_rsp(user) -> user_pb2.UserInfoResponse:
        # 将user的model对象转换成message对象
        user_info_rsp = user_pb2.UserInfoResponse()
        user_info_rsp.id = user.id
        user_info_rsp.password = user.password
        user_info_rsp.mobile = user.mobile
        user_info_rsp.role = user.role

        if user.nick_name:
            user_info_rsp.nickName = user.nick_name
        if user.gender:
            user_info_rsp.gender = user.gender
        if user.birthday:
            user_info_rsp.birthDay = int(time.mktime(user.birthday.timetuple()))

        return user_info_rsp

    @logger.catch()
    def GetUserList(self, request: user_pb2.PageInfo, context):
        # 获取用户列表
        rsp = user_pb2.UserListResponse()
        users = User.select()
        rsp.total = users.count()

        start = 0
        per_page_nums = 10
        if request.pSize:
            per_page_nums = request.pSize

        if request.pn:
            start = per_page_nums * (request.pn - 1)

        users = users.limit(per_page_nums).offset(start)
        for user in users:
            rsp.data.append(self.convert_user_to_rsp(user))

        return rsp

    @logger.catch()
    def GetUserById(self, request, context):
        # 通过用户ID查询用户
        user_id = request.ID
        try:
            user = User.get_by_id(user_id)
            return self.convert_user_to_rsp(user)
        except User.DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在")
            return user_pb2.UserInfoResponse()

    @logger.catch()
    def GetUserByMobile(self, request, context):
        # 通过mobile查询用户
        mobile = request.mobile
        try:
            user = User.get(User.mobile == mobile)
            return self.convert_user_to_rsp(user)
        except User.DoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在")
            return user_pb2.UserInfoResponse()

    @logger.catch()
    def CreateUser(self, request, context):
        # 新建用户

        # 获取参数
        nick_name = request.nickName
        password = request.password
        mobile = request.mobile

        # 检验参数
        try:
            User.get(User.mobile == mobile)
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("手机号码已存在")
            return user_pb2.UserInfoResponse()
        except User.DoesNotExist:
            pass

        # 操作逻辑
        user = User()
        user.nick_name = nick_name
        user.mobile = mobile
        user.password = pbkdf2_sha256.hash(password)
        user.save()

        # 返回响应
        return self.convert_user_to_rsp(user)

    @logger.catch()
    def CheckPassword(self, request, context):
        # 检查密码

        if pbkdf2_sha256.verify(request.password, request.encryptedPassword):
            return user_pb2.CheckResponse(success=True)
        else:
            return user_pb2.CheckResponse(success=False)
