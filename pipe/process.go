package pipe

import (
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/streadway/amqp"
	"log"
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
		rmConnection.Close()
		rbChannel.Close()
		fmt.Println("send service is close.")
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

	// forever := make(chan bool)


	for d := range msgs {

		log.Printf("Received a tag: %d, message: %s", d.DeliveryTag, d.Body)

		d.Ack(false)
	}


	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <- forever
}