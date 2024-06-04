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
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
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
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
