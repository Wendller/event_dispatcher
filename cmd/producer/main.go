package main

import "github.com/Wendller/goexpert/eventdispatcher/pkg/rabbitmq"

func main() {
	channel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	rabbitmq.Publish(channel, "Hello from RabbitMQ")
}
