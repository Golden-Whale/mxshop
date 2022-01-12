package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
)

func GetEnvInfo(env string) int {
	viper.AutomaticEnv()
	return viper.GetInt(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pron.yaml", configFilePrefix)
	if debug == 1 {
		configFileName = fmt.Sprintf("user-web/%s-debug-work.yaml", configFilePrefix)
	} else if debug == 2 {
		configFileName = fmt.Sprintf("user-web/%s-debug-home.yaml", configFilePrefix)
	}

	v := viper.New()
	// 文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	// viper的功能- 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化: %s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息: %v", global.NacosConfig)
	})

	// nacos 配置中心获取配置
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		Username:            global.NacosConfig.User,
		Password:            global.NacosConfig.Password,
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.S().Fatal("连接配置中心失败: %v", err.Error())
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Fatal("从配置中心获取配置失败: %v", err.Error())
	}
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Fatal("配置中心解析json失败: %v", err.Error())
	}
}
