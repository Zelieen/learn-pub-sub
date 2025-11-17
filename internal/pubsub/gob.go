package pubsub

import (
	"bytes"
	"encoding/gob"
	"time"
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
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func decode(data []byte) (GameLog, error) {
	var gl GameLog
	r := bytes.NewReader(data)
	err := gob.NewDecoder(r).Decode(&gl)
	if err != nil {
		return GameLog{}, err
	}
	return gl, nil
}

func PublishGob() {}
