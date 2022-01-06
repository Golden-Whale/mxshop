import requests

headers = {
    "contentType": "applications/json"
}


def register(name, id, address, port):
    url = "http://192.168.1.2:8500/v1/agent/service/register"
    rsp = requests.put(url, headers=headers, json={
        'Name': name,
        'ID': id,
        'Tags': ["mxshop", "web"],
        'Address': address,
        'Port': port,
        "Check": {
            "GRPC": f"{address}:{port}",
            "GRPCUseTLS": False,
            "Timeout": "5s",
            "Interval": "5s",
            "DeregisterCriticalServiceAfter": "15s"
        }
    })
    if rsp.status_code == 200:
        print("注册成功")
    else:
        print("注册失败：", rsp.status_code, rsp.text)


def deregister(id):
    url = "http://192.168.1.2:8500/v1/agent/service/deregister/" + id
    rsp = requests.put(url, headers=headers)
    if rsp.status_code == 200:
        print("注销成功")
    else:
        print("注销失败", rsp.status_code, rsp.text)


if __name__ == '__main__':
    register("mshop-web", "mshop-web", "192.168.1.2", 50051)
    # deregister('mshop-web')
