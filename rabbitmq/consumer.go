package rabbitmq

import (
	"context"
	"database/sql"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(db *sql.DB) string {
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
		"",
		"logs",
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

	result := []string{}

	for j := 0; j != len(body); j += 1 {
		elem := ""
		if string(body[j]) != "," {
			elem += string(body[j])
		} else {
			result = append(result, elem)
		}
	}

	db.ExecContext(context.Background(), "INSERT INTO users (id, name, pass, about_me) VALUES ($1, $2, $3, $4)", result[0], result[1], result[2], result[3])

	result = nil

	return body
}
