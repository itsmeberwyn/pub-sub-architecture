package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
	message, err := json.Marshal(val)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		return err
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	durable := true
	autoDelete := false
	exclusive := false
	if simpleQueueType == 2 { // 2 mean transient
		durable = false
		autoDelete = true
		exclusive = true
	}
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, queue, err
}
