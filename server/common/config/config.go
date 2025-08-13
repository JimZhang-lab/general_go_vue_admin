/*
 * @Author: JimZhang
 * @Date: 2025-07-24 10:50:38
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-08-11 00:00:00
 * @FilePath: /server/common/config/config.go
 * @Description: 使用 viper 读取配置
 */
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Server        server        `mapstructure:"server" yaml:"server"`
	DB            db            `mapstructure:"db" yaml:"db"`
	Redis         redis         `mapstructure:"redis" yaml:"redis"`
	RabbitMQ      rabbitmq      `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	ImageSettings imageSettings `mapstructure:"imageSettings" yaml:"imageSettings"`
	Log           log           `mapstructure:"log" yaml:"log"`
	Jwt           JWT           `mapstructure:"jwt" yaml:"jwt"`
	Seed          seed          `mapstructure:"seed" yaml:"seed"`
}

type server struct {
	Port  int    `mapstructure:"port" yaml:"port"`
	Host  string `mapstructure:"host" yaml:"host"`
	Model string `mapstructure:"model" yaml:"model"`
}

// 数据库配置
type db struct {
	Dialects           string       `mapstructure:"dialects" yaml:"dialects"`
	Host               string       `mapstructure:"host" yaml:"host"`
	Port               int          `mapstructure:"port" yaml:"port"`
	DBName             string       `mapstructure:"db" yaml:"db"`
	Username           string       `mapstructure:"username" yaml:"username"`
	Password           string       `mapstructure:"password" yaml:"password"`
	Charset            string       `mapstructure:"charset" yaml:"charset"`
	MaxIdleConns       int          `mapstructure:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns       int          `mapstructure:"maxOpenConns" yaml:"maxOpenConns"`
	SetConnMaxLifetime int          `mapstructure:"setConnMaxLifetime" yaml:"setConnMaxLifetime"`
	ConnMaxIdleTime    int          `mapstructure:"connMaxIdleTime" yaml:"connMaxIdleTime"`
	SlowThreshold      int          `mapstructure:"slowThreshold" yaml:"slowThreshold"`
	LogLevel           int          `mapstructure:"logLevel" yaml:"logLevel"`
	PrepareStmt        bool         `mapstructure:"prepareStmt" yaml:"prepareStmt"`
	ReadReplicas       readReplicas `mapstructure:"readReplicas" yaml:"readReplicas"`
}

// 读写分离配置
type readReplicas struct {
	Enabled bool     `mapstructure:"enabled" yaml:"enabled"`
	Hosts   []string `mapstructure:"hosts" yaml:"hosts"`
}

type redis struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
}

type rabbitmq struct {
	Host                 string `mapstructure:"host" yaml:"host"`
	Port                 int    `mapstructure:"port" yaml:"port"`
	Username             string `mapstructure:"username" yaml:"username"`
	Password             string `mapstructure:"password" yaml:"password"`
	Vhost                string `mapstructure:"vhost" yaml:"vhost"`
	MaxConnections       int    `mapstructure:"maxConnections" yaml:"maxConnections"`
	MaxChannels          int    `mapstructure:"maxChannels" yaml:"maxChannels"`
	ReconnectDelay       string `mapstructure:"reconnectDelay" yaml:"reconnectDelay"`
	MaxReconnectAttempts int    `mapstructure:"maxReconnectAttempts" yaml:"maxReconnectAttempts"`
}

type log struct {
	Path       string `mapstructure:"path" yaml:"path"`
	Name       string `mapstructure:"name" yaml:"name"`
	Model      string `mapstructure:"model" yaml:"model"`
	MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge"`
	MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
}

// imageSettings图片上传配置
type imageSettings struct {
	UploadDir string `mapstructure:"uploadDir" yaml:"uploadDir"`
	ImageHost string `mapstructure:"imageHost" yaml:"imageHost"`
}

type JWT struct {
	Secret string `mapstructure:"secret" yaml:"secret"`
	Expire int    `mapstructure:"expire" yaml:"expire"`
}

// 首次初始化种子配置
// 允许通过配置关闭种子写入，或指定默认管理员账号
// 示例：
// seed:
//
//	enable: true
//	admin:
//	  username: admin
//	  password: admin123
//	  nickname: 系统管理员
//	  email: admin@example.com
//	  phone: 13800138000
type seed struct {
	Enable bool     `mapstructure:"enable" yaml:"enable"`
	Admin  seedUser `mapstructure:"admin" yaml:"admin"`
}

type seedUser struct {
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Nickname string `mapstructure:"nickname" yaml:"nickname"`
	Email    string `mapstructure:"email" yaml:"email"`
	Phone    string `mapstructure:"phone" yaml:"phone"`
}

var Config *config

func init() {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	v.SetConfigType("yaml")

	// 支持环境变量覆盖（如 SERVER_PORT=8080 覆盖 server.port）
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置失败: %w", err))
	}

	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("解析配置失败: %w", err))
	}
	Config = &cfg
}
