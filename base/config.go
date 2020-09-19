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

	// MQTT 状态订阅消息 channel
	MqttStatusCh  chan *ReceivePacket
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

	// MQTT主题订阅频道
	MqttChannel int

	// 日志级别
	LogLevel	int

	// 输出日志到控制台
	LogToConsole bool
}

// 初始化默认配置
func InitConfig(channel int) {
	DefaultConfig.MqttServerAddress = "tcp://192.168.1.120:1883"
	DefaultConfig.MqttServerHttp = "http://192.168.1.120:18083"
	DefaultConfig.MqttChannel = channel
	DefaultConfig.MqttUsername = "glaucus"
	DefaultConfig.MqttPassword = "123456"
	DefaultConfig.LogLevel = 3
	DefaultConfig.LogToConsole = true
}

// 载入配置文件
func LoadConfig(channel int)  {
	file, err := os.Open("./conf.json")
	if err != nil {
		fmt.Printf("cannot open the config file.\n")
		InitConfig(channel)
		return
	}

	defer func() {
		_  = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&DefaultConfig)
	if err != nil {
		fmt.Printf("cannot parse the config file.\n")
		InitConfig(channel)
		return
	}

	DefaultConfig.MqttChannel = channel
}

// 初始化全局 channel
func InitChannel() {
	MqttControlCh = make(chan *SendPacket)
	MqttStatusCh = make(chan *ReceivePacket, 10)
}
