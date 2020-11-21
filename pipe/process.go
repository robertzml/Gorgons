package pipe

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/glog"
	"github.com/streadway/amqp"
)

const (
	// 当前包名称
	packageName = "pipe"
)

/*
 从Rabbit MQ 中获取指令，并拼装TLV 协议
 */
func Process() {
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

		pak := new(packet)
		if err = json.Unmarshal(d.Body, pak); err != nil {
			glog.Write(2, packageName, "process", "deserialize queue packet failed, " + err.Error())
			d.Ack(false)
			continue
		}

		glog.Write(4, packageName, "process", fmt.Sprintf("receive queue tag: %d, packet: %+v", d.DeliveryTag, pak))


		d.Ack(false)
	}
}

func waterHeaterControl(pak packet) {

}