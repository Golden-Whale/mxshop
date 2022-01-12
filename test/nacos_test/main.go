package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"mxshop-api/test/config"
)

func main() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.1.2",
			Port:   8848,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         "65b94fdb-5ecc-48a9-950e-df24f07e6e22", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//LogDir:              "tmp/nacos/log",
		//CacheDir:            "tmp/nacos/cache",
		Username:   "nacos",
		Password:   "nacos",
		RotateTime: "1h",
		MaxAge:     3,
		LogLevel:   "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group:  "dev"})
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	serverConfig := config.ServerConfig{}
	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: "user-web.json",
	//	Group:  "dev",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("配置文件产生变化，" + "group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
	//time.Sleep(3000 * time.Second)
}
