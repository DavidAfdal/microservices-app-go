package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)


func declareEnxchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topics",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false, //
		nil,
	)
}