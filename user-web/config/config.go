package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type RLYSmsConfig struct {
	APIAccount string `mapstructure:"account"`
	APIToken   string `mapstructure:"token"`
	APPID      string `mapstructure:"app_id"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstructure:"expire"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	RLYSmsInfo  RLYSmsConfig  `mapstructure:"rly_sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
}
