package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 定义一个全局配置变量
var Conf = new(AppConfig)

// AppConfig  网站配置
type AppConfig struct {
	*GinConfig     `mapstructure:"gin"`
	*MysqlConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
	*AliyunConfig  `mapstructure:"aliyun"`
	*ServiceConfig `mapstructure:"service"`
	*EtcdConfig    `mapstructure:"etcd"`
}

type GinConfig struct {
	HttpHost string `mapstructure:"http_host"`
	HttpPort int    `mapstructure:"http_port"`
}

// MysqlConfig  mysql配置
type MysqlConfig struct {
	MysqlPort     int    `mapstructure:"port"`
	MysqlHost     string `mapstructure:"host"`
	MysqlUser     string `mapstructure:"user"`
	MysqlPassword string `mapstructure:"password"`
	MysqlDbname   string `mapstructure:"dbname"`
}

// RedisConfig redis配置
type RedisConfig struct {
	RedisHost     string `mapstructure:"host"`
	RedisPassword string `mapstructure:"password"`
	RedisPort     int    `mapstructure:"port"`
	RedisDb       int    `mapstructure:"db"`
	RedisPoolSize int    `mapstructure:"pool_size"`
}

// AliyunConfig  阿里云短信配置
type AliyunConfig struct {
	AccessKeyId        string `mapstructure:"access_key_id"`
	AccessKeySecret    string `mapstructure:"access_key_secret"`
	SignName           string `mapstructure:"sign_name"`
	SignupTemplateCode string `mapstructure:"signup_template_code"`
	LoginTemplateCode  string `mapstructure:"login_template_code"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	UserServiceAddress    string `mapstructure:"user_service_address"`
	SendSmsServiceAddress string `mapstructure:"send_sms_service_address"`
}

type EtcdConfig struct {
	EtcdHost string `mapstructure:"host"`
	EtcdPort int    `mapstructure:"port"`
}

func Init() (err error) {
	//设置读取配置文件路径
	viper.SetConfigFile("../../../settings/config.yaml")
	//读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("viper ReadInConfig failed, err:%v\n", err)
		return
	}
	//将读取配置信息反序列化入全局变量
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改了")
		//将更改的配置文件信息反序列化入全局变量
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed, err:%v\n", err)
			return
		}
	})
	return nil
}
