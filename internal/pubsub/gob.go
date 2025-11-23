package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GameLog struct {
	CurrentTime time.Time
	Message     string
	Username    string
}

func encode(gl GameLog) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(gl)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(data []byte) (GameLog, error) {
	var gl GameLog
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(&gl)
	if err != nil {
		return GameLog{}, err
	}
	return gl, nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, log GameLog) error {

	load, err := encode(log)
	if err != nil {
		return errors.Join(errors.New("could not encode to gob: "), err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        load,
	})
	if err != nil {
		return errors.Join(errors.New("publish gob failed: "), err)
	}

	return nil
}

func PublishGameLog(ch *amqp.Channel, user string, msg string) error {
	return nil
}
