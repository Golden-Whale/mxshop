from peewee import *

from user_srv.settings import settings


class BaseModel(Model):
    class Meta:
        database = settings.DB


class User(BaseModel):
    GENDER_CHOIES = (
        ("female", "女"),
        ("male", "男")
    )

    ROLE_CHOIES = (
        (1, "普通用户"),
        (2, "管理员")
    )

    mobile = CharField(max_length=11, index=True, verbose_name="手机号码", unique=True)
    password = CharField(max_length=100, verbose_name="密码")
    nick_name = CharField(max_length=20, null=True, verbose_name="昵称")
    head_url = CharField(max_length=200, null=True, verbose_name="头像")
    birthday = DateField(null=True, verbose_name="生日")
    address = CharField(max_length=200, null=True, verbose_name="地址")
    desc = TextField(null=True, verbose_name="简介")
    gender = CharField(max_length=6, choices=GENDER_CHOIES, verbose_name="性别")
    role = IntegerField(default=1, choices=ROLE_CHOIES, verbose_name="用户角色")


if __name__ == '__main__':
    settings.DB.create_tables([User])

    # for i in range(10):
    #     user = User()
    #     user.nick_name = f"bobby{i}"®
    #     user.mobile = f"1383838943{i}"
    #     user.password = pbkdf2_sha256.hash("admin123")
    #     user.save()
    users = User.select()
    users = users.offset(1).limit(2)
    for user in users:
        # if user.birthday:
        #     print(user.birthday)
        #     u_time = int(time.mktime(user.birthday.timetuple()))
        #     print(u_time)
        #     print(datetime.date.fromtimestamp(u_time))
        # print(pbkdf2_sha256.verify("admin123", user.password))
        print(user.id)