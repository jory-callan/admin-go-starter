package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Load 加载配置文件并返回强类型 AppConfig
// 流程: DefaultConfig() → ReadInConfig() → Unmarshal()
// 配置文件中未指定的字段自动使用 DefaultConfig() 中的默认值
func Load(configFile string) *AppConfig {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigType("yml")

	// 环境变量配置，默认前缀为 APP，下划线替换为点号
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 2. 查找并读取配置文件，优先级 1. 命令行参数 2. 环境变量
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// 默认查找顺序, 优先级 1. 当前目录 2. ./config 3. ./conf
		for _, p := range []string{".", "./config", "./conf"} {
			v.AddConfigPath(p)
		}
		// 默认配置文件名为 config
		v.SetConfigName("config")
	}

	// 3. 读取并合并配置
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Sprintf("read config: %v", err))
		}
		log.Println("config file not found, using default configuration")
	}

	// 3. Unmarshal 到强类型
	var cfg AppConfig
	// 支持「人类可读的时间格式」
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
	)
	if err := v.Unmarshal(&cfg, viper.DecodeHook(decodeHook)); err != nil {
		panic(fmt.Sprintf("unmarshal config: %v", err))
	}

	return &cfg
}
