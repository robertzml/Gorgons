package queue

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/equipment"
	"github.com/robertzml/Gorgons/glog"
	"github.com/robertzml/Gorgons/send"
	"time"
)

/*
 从Rabbit MQ 中获取控制队列指令，并拼装TLV 协议
*/
func Control() {
	rbChannel, err := rmConnection.Channel()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "Control", fmt.Sprintf("catch runtime panic in control: %v", r))
		}

		rbChannel.Close()
		glog.Write(1, packageName, "Control", "control queue is close.")
	}()

	// 队列名称
	queueName := "ControlQueue"

	queue, err := rbChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = rbChannel.Qos(1, 0, false)

	glog.Write(3, packageName, "Control", "declare control queue")

	msgs, err := rbChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs {

		pak := new(controlPacket)
		if err = json.Unmarshal(d.Body, pak); err != nil {
			glog.Write(2, packageName, "Control", "deserialize queue packet failed, "+err.Error())
			_ = d.Ack(false)
			continue
		}

		glog.Write(4, packageName, "Control", fmt.Sprintf("receive queue tag: %d, packet: %+v", d.DeliveryTag, pak))

		if pak.DeviceType == 1 {
			waterHeaterControl(pak)
		} else {
			glog.Write(3, packageName, "Control", "unknown device.")
		}

		_ = d.Ack(false)
	}
}

/*
 拼接热水器控制报文，并下发到mqtt
*/
func waterHeaterControl(qp *controlPacket) {
	waterHeater := equipment.NewWaterHeaterContext(snapshot)

	if mainboardNumber, exist := waterHeater.GetMainboardNumber(qp.SerialNumber); exist {
		controlMsg := send.NewWHControlMessage(qp.SerialNumber, mainboardNumber)

		sendPak := new(base.SendPacket)
		sendPak.SerialNumber = qp.SerialNumber
		sendPak.DeviceType = qp.DeviceType

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
		case 7: // 手动清洗
			sendPak.Payload = controlMsg.Clean(qp.Option)
		case 8: // 数据清零
			sendPak.Payload = controlMsg.Clear(int8(qp.Option))
		case 10:
			sendPak.Payload = controlMsg.SetHeatTime(qp.Option)
		case 11:
			sendPak.Payload = controlMsg.SetHotWater(qp.Option)
		case 12:
			sendPak.Payload = controlMsg.SetWorkTime(qp.Option)
		case 13:
			sendPak.Payload = controlMsg.SetUsedPower(qp.Option)
		case 14:
			sendPak.Payload = controlMsg.SetSavePower(qp.Option)
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
