package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var (
	GlobalConfig *Config
	viper_       *viper.Viper
	env          = flag.String("env", "dev", "运行环境 (dev|test|prod)")
	configFile   = flag.String("config", "", "配置文件路径")
)

type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	Redis    *RedisConfig    `mapstructure:"redis"`
}

// 初始化配置
func Init() error {
	flag.Parse()
	v := viper.New()
	// 设置配置文件
	// 默认配置路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile("./config/config-common.yaml")
	*configFile = GetConfigFileByEnv(*env)

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	v.SetConfigFile(*configFile)
	v.MergeInConfig()

	viper_ = v

	// 解析配置到结构体
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	GlobalConfig = cfg
	// 设置配置文件监听
	setupConfigWatch(viper_, cfg)
	log.Printf("Config loaded successfully from: %s", v.ConfigFileUsed())
	return nil
}

// 设置配置监听
func setupConfigWatch(v *viper.Viper, cfg *Config) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		// 重新解析配置
		if err := v.Unmarshal(cfg); err != nil {
			log.Printf("Failed to reload config: %v", err)
			return
		}
		log.Printf("Config reloaded successfully")
		// 可以在这里触发配置更新事件
		triggerConfigUpdate(cfg)
	})
}

// 触发配置更新事件
func triggerConfigUpdate(cfg *Config) {
	// 这里可以添加配置更新时的回调逻辑
	// 例如：重新连接数据库、更新日志级别等
	log.Printf("Configuration updated - App: %s, Debug: %t",
		cfg.Server.Name, cfg.Server.Debug)
}

// 触发配置更新事件
func addWatchConfig(key []string, fn func(config Config)) {

}

// 获取配置实例
func GetConfig() *Config {
	return GlobalConfig
}

// 重新加载配置（手动触发）
func Reload() error {
	if GlobalConfig == nil || viper_ == nil {
		return fmt.Errorf("config not initialized")
	}

	if err := viper_.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(GlobalConfig)
}

// 根据环境获取配置文件名
func GetConfigFileByEnv(env string) string {
	switch env {
	case "dev", "development":
		return "./config/config-dev.yaml"
	case "prod", "production":
		return "./config/config-prod.yaml"
	case "test":
		return "./config/config-test.yaml"
	default:
		return "./config/config.yaml"
	}
}
