package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
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

	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.ServerConfig)

	// viper的功能- 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化: %s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息: %v", global.ServerConfig)
	})
}