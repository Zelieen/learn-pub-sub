package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	QueueDurable SimpleQueueType = iota
	QueueTransient
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	load, err := json.Marshal(val)
	if err != nil {
		return errors.Join(errors.New("could not marshal to json: "), err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        load,
	})
	if err != nil {
		return errors.Join(errors.New("publish JSON failed: "), err)
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	//open a channel
	comChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel to Peril connection: %s\n", err)
		return comChan, amqp.Queue{}, err
	}

	// declare a queue
	comQueue, err := comChan.QueueDeclare(queueName,
		queueType == QueueDurable,
		queueType == QueueTransient,
		queueType == QueueTransient,
		false,
		nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s\n", err)
		return comChan, comQueue, err
	}

	// bind queue to channel
	err = comChan.QueueBind(comQueue.Name, key, exchange, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind exchange channel to its queue: %s\n", err)
		return comChan, comQueue, err
	}

	return comChan, comQueue, nil
}
