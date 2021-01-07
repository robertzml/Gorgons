package queue

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/equipment"
	"github.com/robertzml/Gorgons/glog"
	"github.com/robertzml/Gorgons/send"
)

/**
 从Rabbit MQ 中获取特殊队列指令，并拼装TLV 协议
 */
func Special() {
	rbChannel, err := rmConnection.Channel()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "Special", fmt.Sprintf("catch runtime panic in special: %v", r))
		}

		rbChannel.Close()
		glog.Write(1, packageName, "Special", "special queue is close.")
	}()

	queueName := "SpecialQueue"

	queue, err := rbChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = rbChannel.Qos(1, 0, false)
	glog.Write(3, packageName, "Special", "declare special queue")

	msgs, err := rbChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs {

		pak := new(specialPacket)
		if err = json.Unmarshal(d.Body, pak); err != nil {
			glog.Write(2, packageName, "Special", "deserialize queue packet failed, "+err.Error())
			d.Ack(false)
			continue
		}

		glog.Write(4, packageName, "Special", fmt.Sprintf("receive queue tag: %d, packet: %+v", d.DeliveryTag, pak))

		if pak.DeviceType == 1 {
			waterHeaterSpecial(pak)
		} else {
			glog.Write(3, packageName, "Special", "unknown device.")
		}

		d.Ack(false)
	}

	glog.Write(3, packageName, "Special", "special queue reach end")
}

/**
拼接热水器特殊指令报文，并下发到emq
*/
func waterHeaterSpecial(qp *specialPacket) {
	waterHeater := equipment.NewWaterHeaterContext(snapshot)

	if mainboardNumber, exist := waterHeater.GetMainboardNumber(qp.SerialNumber); exist {
		specialMsg := send.NewWHSpecialMessage(qp.SerialNumber, mainboardNumber)

		sendPak := new(base.SendPacket)
		sendPak.SerialNumber = qp.SerialNumber
		sendPak.DeviceType = qp.DeviceType

		switch qp.ControlType {
		case 1:
			sendPak.Payload = specialMsg.SoftFunction(qp.Option)
		case 2:
			sendPak.Payload = specialMsg.Special(qp.Option)
		case 3:
			sendPak.Payload = specialMsg.Manual(qp.Option)
		case 4:
			sendPak.Payload = specialMsg.Duplicate(qp.Option)
		default:
			glog.Write(3, packageName, "special", "wrong control type.")
			return
		}

	} else {
		glog.Write(2, packageName, "special", fmt.Sprintf("sn: %s. equipment cannot found mainboard number.", qp.SerialNumber))
	}
}