package config

import (
	"github.com/spf13/viper"
)

// Load 加载配置文件
func Load(configFile string) *viper.Viper {
	conf := viper.New()
	
	// 查找配置文件
	if configFile != "" {
		conf.SetConfigFile(configFile)
	} else {
		// 默认查找顺序：./config.yml > ./config/config.yml > ./conf/config.yml
		conf.SetConfigName("config")
		conf.SetConfigType("yaml")
		conf.AddConfigPath(".")
		conf.AddConfigPath("./config")
		conf.AddConfigPath("./conf")
	}
	
	// 读取配置
	if err := conf.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}
	
	return conf
}
