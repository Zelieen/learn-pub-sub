package main

import (
	"fmt"
	"log"

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

	gamestate := gamelogic.NewGameState(user)

	// start REPL
	for {
		words := gamelogic.GetInput()

		switch command := words[0]; command {
		case "spawn":
			err = gamestate.CommandSpawn(words)
			if err != nil {
				log.Printf("Could not spawn unit: %s\n", err)
			}
		case "move":
			_, err := gamestate.CommandMove(words)
			if err != nil {
				log.Printf("Could not move unit: %s\n", err)
			}
		case "status":
			gamestate.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		case "help":
			gamelogic.PrintClientHelp()
		default:
			fmt.Println("Unknown command. Please try something else.")
		}
	}
}
