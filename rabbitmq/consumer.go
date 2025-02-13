package rabbitmq

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume() string {
	conn, err := amqp.Dial("")
	if err != nil {
		slog.Debug("Failed to connect to RabbitMQ: ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		slog.Debug("Failed to open a channel: ", err)
	}
	channel.Close()

	err = channel.ExchangeDeclare(
		"users",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Debug("Failed to declare an exchange", err)
	}

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		slog.Debug("Failed to declare a queue", err)
	}

	err = channel.QueueBind(
		queue.Name,
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		slog.Debug("Failed to bind a queue")
	}

	message, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Debug("Failed to register a consumer", err)
	}

	var forever chan struct{}
	var body string = ""

	go func() {
		for i := range message {
			body += string(i.Body)
		}
	}()
	<-forever

	return body
}
