package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/glog"
)

// 发送端客户ID
var sendClientId string

// 启动MQTT发送服务
// 通过全局 MqttControlCh 获取发送请求
func StartSend() {
	sendMqtt := new(MQTT)

	defer func() {
		sendMqtt.Disconnect()
		fmt.Println("Send mqtt function is close.")
	}()

	sendClientId = "send-gorgons"
	sendMqtt.Connect(sendClientId, base.DefaultConfig.MqttUsername, base.DefaultConfig.MqttPassword, base.DefaultConfig.MqttServerAddress, sendOnConnect)

	for {
		pak := <-base.MqttControlCh

		var controlTopic = fmt.Sprintf("server/%d/%s/control_info", pak.DeviceType, pak.SerialNumber)
		sendMqtt.Publish(controlTopic, 2, pak.Payload)

		glog.Write(4, packageName, "send", fmt.Sprintf("PUBLISH Topic:%s, Payload: %s", controlTopic, pak.Payload))
	}
}

// 发送连接回调
var sendOnConnect paho.OnConnectHandler = func(client paho.Client) {
	glog.Write(3, packageName, "onConnect", fmt.Sprintf("%s connect to mqtt.", sendClientId))
}
