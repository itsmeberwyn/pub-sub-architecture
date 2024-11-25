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
		exclusive = false
	}
	queue, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, amqp.Table{
        "x-dead-letter-exchange": "peril_dlx",
    })
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return ch, queue, err
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int,
	handler func(T) string,
) error {
    ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
    if err != nil {
        return err
    }

    msgs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)

    go func() {
        var body T
        for msg := range msgs {
            err := json.Unmarshal(msg.Body, &body)
            if err != nil {
                log.Fatal(err)
            }
            ackType := handler(body)
            if ackType == "Ack" {
                msg.Ack(false)
            } else if ackType == "NackRequeue" {
                msg.Nack(false, true)
            } else if ackType == "NackDiscard" {
                msg.Nack(false, false)
            }
        }
    }()


    return nil
}
