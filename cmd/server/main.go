package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	//Build a connection to rabbitMQ
	connectStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectStr)
	if err != nil {
		log.Fatalf("Error setting up the Peril connection: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connecting Peril to RabbitMQ was successful!")

	//wait for interuption signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("\nTerminating Peril connection and shutting down")
}
