package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// Conf 定义一个全局配置变量
var Conf = new(AppConfig)

// AppConfig  网站配置
type AppConfig struct {
	PodIpAddr       string `mapstructure:"pod_ip_addr"` //服务ip地址
	SnowflakeId     int64  `mapstructure:"snowflake_id"`
	*GinConfig      `mapstructure:"gin"`
	*MysqlConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*AliyunConfig   `mapstructure:"aliyun"`
	*ServiceConfig  `mapstructure:"service"`
	*EtcdConfig     `mapstructure:"etcd"`
	*RabbitMQConfig `mapstructure:"rabbitmq"`
	*QiniuConfig    `mapstructure:"qiniuyv"`
	*LogConfig      `mapstructure:"log"`
	*TracerConfig   `mapstructure:"tracer"`
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

// EtcdConfig etcd配置
type EtcdConfig struct {
	EtcdHost string `mapstructure:"host"`
	EtcdPort int    `mapstructure:"port"`
}

// RabbitMQConfig 消息队列配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	Bucket      string `mapstructure:"bucket"`
	QiniuServer string `mapstructure:"qiniuServer"`
}

// LogConfig 日志配置
type LogConfig struct {
	FileName   string `mapstructure:"file_name"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// TracerConfig 追踪配置
type TracerConfig struct {
	TracingEndPoint string  `mapstructure:"tracing_end_point"` //追踪数据要发送到的服务器地址
	OtelState       string  `mapstructure:"otel_state"`
	OtelSampler     float64 `mapstructure:"otel_sampler"`
}

func init() {
	//设置读取配置文件路径
	viper.SetConfigFile("C:\\Users\\浅梦\\Desktop\\star\\app\\constant\\settings\\config.yaml")
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		//utils.Logger.Error("读取配置文件失败", zap.Error(err))
		log.Println("viper.ReadInConfig failed, err:", err)
		return
	}
	//将读取配置信息反序列化入全局变量
	if err := viper.Unmarshal(Conf); err != nil {
		//utils.Logger.Error("配置文件反序列化入全局变量失败", zap.Error(err))
		log.Println("viper.Unmarshal failed, err:", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//utils.Logger.Info("配置文件修改了")
		//将更改的配置文件信息反序列化入全局变量
		if err := viper.Unmarshal(Conf); err != nil {
			//utils.Logger.Error("配置文件反序列化入全局变量失败", zap.Error(err))
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
			return
		}
	})
}
