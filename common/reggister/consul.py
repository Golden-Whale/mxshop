import consul
import requests

from common.reggister.base import Register


class ConsulRegister(Register):

    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.c = consul.Consul(host, port)

    def register(self, name, id, address, port, tags, check=None) -> bool:
        if check is None:
            check = {
                "GRPC": f"{address}:{port}",
                "GRPCUseTLS": False,
                "Timeout": "5s",
                "Interval": "5s",
                "DeregisterCriticalServiceAfter": "15s"
            }
        return self.c.agent.service.register(name, id, address, port, tags, check)

    def deregister(self, service_id) -> bool:
        return self.c.agent.service.deregister(service_id)

    def get_all_service(self) -> dict:
        return self.c.agent.services()

    def filter_service(self, filter) -> dict:
        url = f"http://{self.host}:{self.port}/v1/agent/services"
        params = {
            "filter": filter
        }
        return requests.get(url, params=params).json()
