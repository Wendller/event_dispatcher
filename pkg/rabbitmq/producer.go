package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func Publish(channel *amqp.Channel, body string) error {
	err := channel.Publish(ExchangeType, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		return err
	}

	return nil
}
