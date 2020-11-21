package base

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	// 默认配置
	DefaultConfig Config

	// MQTT 发送控制指令 channel
	MqttControlCh  chan *SendPacket
)

// 配置
type Config struct {
	// MQTT 服务器地址
	MqttServerAddress string

	// MQTT HTTP服务地址
	MqttServerHttp string

	// MQTT 用户名
	MqttUsername string

	// MQTT 密码
	MqttPassword string

	// Rabbit MQ 连接字符串
	RabbitMQAddress string

	// 日志级别
	LogLevel	int

	// 输出日志到控制台
	LogToConsole bool
}

// 初始化默认配置
func InitConfig() {
	DefaultConfig.MqttServerAddress = "tcp://192.168.1.120:1883"
	DefaultConfig.MqttServerHttp = "http://192.168.1.120:18083"
	DefaultConfig.MqttUsername = "glaucus"
	DefaultConfig.MqttPassword = "123456"
	DefaultConfig.RabbitMQAddress = "amqp://guest:guest@localhost:5672/"
	DefaultConfig.LogLevel = 3
	DefaultConfig.LogToConsole = true
}

// 载入配置文件
func LoadConfig()  {
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Printf("cannot open the config file.\n")
		InitConfig()
		return
	}

	defer func() {
		_  = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&DefaultConfig)
	if err != nil {
		fmt.Printf("cannot parse the config file.\n")
		InitConfig()
		return
	}
}

// 初始化全局 channel
func InitChannel() {
	MqttControlCh = make(chan *SendPacket)
}
