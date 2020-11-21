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
	"time"
)

const (
	// 当前包名称
	packageName = "pipe"

	// 队列名称
	queueName = "ControlQueue"
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
			waterHeaterControl(pak)
		} else {
			glog.Write(3, packageName, "process", "unknown device.")
		}

		d.Ack(false)
	}
}

/*
 拼接热水器控制报文，并下发到mqtt
*/
func waterHeaterControl(qp *queuePacket) {
	waterHeater := equipment.NewWaterHeaterContext(snapshot)

	if mainboardNumber, exist := waterHeater.GetMainboardNumber(qp.SerialNumber); exist {
		controlMsg := send.NewWHControlMessage(qp.SerialNumber, mainboardNumber)

		sendPak := new(base.SendPacket)
		sendPak.SerialNumber = qp.SerialNumber
		sendPak.DeviceType = 1

		// 获取已保存的设置信息
		_, set := waterHeater.LoadSetting(qp.SerialNumber)
		set.SerialNumber = qp.SerialNumber

		switch qp.ControlType {
		case 1: //开关机
			sendPak.Payload = controlMsg.Power(qp.Option)
		case 2: //激活,非激活
			sendPak.Payload = controlMsg.Activate(qp.Option)
			set.Activate = int8(qp.Option)
			if qp.Option == 1 {
				set.SetActivateTime = time.Now().Unix() * 1000
			}
		case 3: // 加锁
			sendPak.Payload = controlMsg.Lock()
			set.Unlock = 0
		case 4: // 解锁
			sendPak.Payload = controlMsg.Unlock(qp.Option, qp.Deadline)
			set.Unlock = 1
			if qp.Option == 1 {
				set.DeadlineTime = qp.Deadline
			}
		case 5: // 设置允许使用时间
			sendPak.Payload = controlMsg.SetDeadline(qp.Deadline)
			set.DeadlineTime = qp.Deadline
		case 6: // 设定温度
			sendPak.Payload = controlMsg.SetTemp(qp.Option)
		case 7:
			sendPak.Payload = controlMsg.Clean(qp.Option)
		case 8:
			sendPak.Payload = controlMsg.Clean(qp.Option)
		default:
			glog.Write(3, packageName, "control", "wrong control type.")
			return
		}

		// 保存设置信息
		if qp.ControlType >= 2 && qp.ControlType <= 5 {
			waterHeater.SaveSetting(set)
		}

		base.MqttControlCh <- sendPak

	} else {	// 未找到设备直接保存设置值
		_, set := waterHeater.LoadSetting(qp.SerialNumber)

		set.SerialNumber = qp.SerialNumber

		switch qp.ControlType {
		case 2:
			set.Activate = int8(qp.Option)
			if qp.Option == 1 {
				set.SetActivateTime = time.Now().Unix() * 1000
			}
		case 3:
			set.Unlock = 0
		case 4:
			set.Unlock = 1
			if qp.Option == 1 {
				set.DeadlineTime = qp.Deadline
			}
		case 5:
			set.DeadlineTime = qp.Deadline
		default:
			glog.Write(3, packageName, "control", "wrong control type.")
			return
		}

		waterHeater.SaveSetting(set)

		glog.Write(3, packageName, "control", fmt.Sprintf("sn: %s. save setting for not found equipment.", qp.SerialNumber))
	}

	return
}
