package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

var QueueName string = "mainqueue"
var ExchangeType string = "amq.direct"

func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return channel, nil
}
