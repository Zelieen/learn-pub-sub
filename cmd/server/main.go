package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	//open a channel
	comChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel to Peril connection: %s\n", err)
		return
	}

	// print help for start up
	gamelogic.PrintServerHelp()

	// start the game loop
	for {
		words := gamelogic.GetInput()
		if len(words) < 1 {
			continue
		}
		switch word := words[0]; word {
		case "pause":
			log.Println("Sending pause signal to game")
			err = publishPause(true, comChan)
			if err != nil {
				log.Printf("Could not publish pause signal: %s\n", err)
			}
		case "resume":
			log.Println("Sending resume signal to game")
			err = publishPause(false, comChan)
			if err != nil {
				log.Printf("Could not publish resume signal: %s\n", err)
			}
		case "quit":
			log.Println("Exiting the game. Stay safe.")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			log.Println("Unknown command. Please try something else.")
		}
	}
}
