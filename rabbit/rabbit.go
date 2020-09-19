package rabbit

import (
	"github.com/streadway/amqp"
)

func Connect() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("control", "fanout", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
}