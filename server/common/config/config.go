/*
 * @Author: JimZhang
 * @Date: 2025-07-24 10:50:38
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-24 23:39:53
 * @FilePath: /server/common/config/config.go
 * @Description:
 *
 */
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Server        server        `yaml:"server"`
	DB            db            `yaml:"db"`
	Redis         redis         `yaml:"redis"`
	RabbitMQ      rabbitmq      `yaml:"rabbitmq"`
	ImageSettings imageSettings `yaml:"imageSettings"`
	Log           log           `yaml:"log"`
	Jwt           JWT           `yaml:"jwt"`
}

type server struct {
	Port  int    `yaml:"port"`
	Host  string `yaml:"host"`
	Model string `yaml:"model"`
}

// 数据库配置
type db struct {
	Dialects           string       `yaml:"dialects"`
	Host               string       `yaml:"host"`
	Port               int          `yaml:"port"`
	DBName             string       `yaml:"db"`
	Username           string       `yaml:"username"`
	Password           string       `yaml:"password"`
	Charset            string       `yaml:"charset"`
	MaxIdleConns       int          `yaml:"maxIdleConns"`
	MaxOpenConns       int          `yaml:"maxOpenConns"`
	SetConnMaxLifetime int          `yaml:"setConnMaxLifetime"`
	ConnMaxIdleTime    int          `yaml:"connMaxIdleTime"`
	SlowThreshold      int          `yaml:"slowThreshold"`
	LogLevel           int          `yaml:"logLevel"`
	PrepareStmt        bool         `yaml:"prepareStmt"`
	ReadReplicas       readReplicas `yaml:"readReplicas"`
}

// 读写分离配置
type readReplicas struct {
	Enabled bool     `yaml:"enabled"`
	Hosts   []string `yaml:"hosts"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type rabbitmq struct {
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
	Username             string `yaml:"username"`
	Password             string `yaml:"password"`
	Vhost                string `yaml:"vhost"`
	MaxConnections       int    `yaml:"maxConnections"`
	MaxChannels          int    `yaml:"maxChannels"`
	ReconnectDelay       string `yaml:"reconnectDelay"`
	MaxReconnectAttempts int    `yaml:"maxReconnectAttempts"`
}

type log struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Model      string `yaml:"model"`
	MaxAge     int    `yaml:"maxAge"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
}

// imageSettings图片上传配置
type imageSettings struct {
	UploadDir string `yaml:"uploadDir"`
	ImageHost string `yaml:"imageHost"`
}

type JWT struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

var Config *config

func init() {
	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		panic(err)
	}
}
