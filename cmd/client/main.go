package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	// Build a connection to rabbitMQ
	connectStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectStr)
	if err != nil {
		log.Fatalf("Error setting up the Peril connection: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connecting Peril to RabbitMQ was successful!")

	// Get a user name
	user, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not welcome user: %s\n", err)
		return
	}

	// Bind to queue
	//channel, queue, err :=
	pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, user),
		routing.PauseKey,
		pubsub.QueueTransient,
	)

	//wait for interuption signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println("\nTerminating Peril connection and shutting down")
}
