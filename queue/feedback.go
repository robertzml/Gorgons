package queue

import (
	"encoding/json"
	"fmt"
	"github.com/robertzml/Gorgons/glog"
)

/**
 从Rabbit MQ 中获取反馈队列指令，并拼装TLV 协议
 */
func Feedback() {
	rbChannel, err := rmConnection.Channel()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Write(1, packageName, "Feedback", fmt.Sprintf("catch runtime panic in feedback: %v", r))
		}

		rbChannel.Close()
		glog.Write(1, packageName, "Feedback", "feedback queue is close.")
	}()

	queueName := "FeedbackQueue"

	queue, err := rbChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = rbChannel.Qos(1, 0, false)
	glog.Write(3, packageName, "Feedback", "declare feedback queue")

	msgs, err := rbChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msgs {

		pak := new(feedbackPacket)
		if err = json.Unmarshal(d.Body, pak); err != nil {
			glog.Write(2, packageName, "Feedback", "deserialize queue packet failed, "+err.Error())
			d.Ack(false)
			continue
		}

		glog.Write(4, packageName, "Feedback", fmt.Sprintf("receive queue tag: %d, packet: %+v", d.DeliveryTag, pak))

		if pak.DeviceType == 1 {
			waterHeaterFeedback(pak)
		} else {
			glog.Write(3, packageName, "Feedback", "unknown device.")
		}

		d.Ack(false)
	}
}

/**
拼接热水器反馈报文，并下发到emq
 */
func waterHeaterFeedback(qp *feedbackPacket) {

}