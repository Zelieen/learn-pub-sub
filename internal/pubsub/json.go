package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	// check queue and exchange
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}
	// new channel
	newchan, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	// unmarshal in an unnamed function to dodge having to hand over parameter: T any
	go func() {
		defer ch.Close()
		var msgBody T
		for msg := range newchan {
			err := json.Unmarshal(msg.Body, &msgBody)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			acknowl := handler(msgBody)
			switch acknowl {
			case Ack:
				msg.Ack(false)
				log.Printf("Ack: message processed")
			case NackRequeue:
				msg.Nack(false, true)
				log.Printf("Nack: message requeue")
			case NackDiscard:
				msg.Nack(false, false)
				log.Printf("Nack: message discarded")
			default:
				msg.Nack(false, false)
				log.Printf("Nack: message is problematic and discarded")
			}
		}
	}()

	return nil
}
