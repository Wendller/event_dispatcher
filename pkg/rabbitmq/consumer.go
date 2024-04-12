package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func Consume(channel *amqp.Channel, out chan<- amqp.Delivery) error {
	messages, err := channel.Consume(QueueName, "go-consumer", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for message := range messages {
		out <- message
	}

	return nil
}
