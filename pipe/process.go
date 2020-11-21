package pipe

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/db"
	"github.com/robertzml/Gorgons/equipment"
	"github.com/robertzml/Gorgons/glog"
	"github.com/robertzml/Gorgons/send"
	"github.com/streadway/amqp"
)

const (
	// 当前包名称
	packageName = "pipe"
)

var (
	// 用于注入实时数据访问类
	snapshot db.Snapshot
)

/*
 从Rabbit MQ 中获取指令，并拼装TLV 协议
*/
func Process(snap db.Snapshot) {
	snapshot = snap

	rmConnection, err := amqp.Dial(base.DefaultConfig.RabbitMQAddress)
	if err != nil {
		panic(err)
	}

	rbChannel, err := rmConnection.Channel()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "process", fmt.Sprintf("catch runtime panic in process: %v", r))
		}

		rmConnection.Close()
		rbChannel.Close()
		glog.Write(3, packageName, "process", "process service is close.")
	}()

	queueName := "ControlQueue"
	queue, err := rbChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = rbChannel.Qos(1, 0, false)

	msgs, err := rbChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs {

		pak := new(queuePacket)
		if err = json.Unmarshal(d.Body, pak); err != nil {
			glog.Write(2, packageName, "process", "deserialize queue packet failed, "+err.Error())
			d.Ack(false)
			continue
		}

		glog.Write(4, packageName, "process", fmt.Sprintf("receive queue tag: %d, packet: %+v", d.DeliveryTag, pak))

		if pak.DeviceType == 1 {
			_ = waterHeaterControl(pak)
		} else {
			glog.Write(3, packageName, "process", "unknown device.")
		}

		d.Ack(false)
	}
}

/*
 拼接热水器控制报文，并下发到mqtt
*/
func waterHeaterControl(qp *queuePacket) error {
	waterHeater := equipment.NewWaterHeaterContext(snapshot)

	if mainboardNumber, exist := waterHeater.GetMainboardNumber(qp.SerialNumber); exist {
		controlMsg := send.NewWHControlMessage(qp.SerialNumber, mainboardNumber)

		sendPak := new(base.SendPacket)
		sendPak.SerialNumber = qp.SerialNumber
		sendPak.DeviceType = 1

		switch qp.ControlType {
		case 1:
			sendPak.Payload = controlMsg.Power(qp.Option)
		}

		base.MqttControlCh <- sendPak
	} else {
		glog.Write(3, packageName, "control", "cannot find mainboard number")
	}

	return nil
}
