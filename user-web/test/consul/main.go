package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, id string, tags []string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	registration := new(api.AgentServiceRegistration)
	registration.ID = id
	registration.Name = name
	registration.Address = address
	registration.Port = port
	registration.Tags = tags
	registration.Check = &api.AgentServiceCheck{
		Interval:                       "0.5s",
		Timeout:                        "0.5s",
		HTTP:                           "http://192.168.1.2:8081/health",
		DeregisterCriticalServiceAfter: "1.5s",
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func AllService() error {
	cfg := api.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	data, err := client.Agent().Services()
	if err != nil {
		return err
	}
	for key, _ := range data {
		fmt.Println(key)
	}
	return nil
}

func FilterService() error {
	cfg := api.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		return err
	}
	for key, _ := range data {
		fmt.Println(key)
	}
	return nil
}

func main() {
	//_ = Register("192.168.1.2", 8081, "user-web", "user-web", []string{"mxshop"})
	FilterService()
}
