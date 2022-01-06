import consul

c = consul.Consul()

# rsp = c.agent.service.register(name="user-srv", service_id="user-srv",
#                                address="192.168.1.2", port=50051, tags=["mxshop"],
#                                check={
#                                    "GRPC": f"{'192.168.1.2'}:{50051}",
#                                    "GRPCUseTLS": False,
#                                    "Timeout": "5s",
#                                    "Interval": "5s",
#                                    "DeregisterCriticalServiceAfter": "15s"
#                                }
#                                )

# rsp = c.agent.service.deregister("user-srv")
rsp = c.agent.services()
print(rsp)
