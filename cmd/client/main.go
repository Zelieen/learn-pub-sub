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

	gamestate := gamelogic.NewGameState(user)

	// listen for server PAUSE
	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", user),
		routing.PauseKey,
		pubsub.QueueTransient,
		handlerPause(gamestate))
	if err != nil {
		log.Fatalf("Could not subscribe to pause: %v", err)
	}
	// listen to any army_moves
	err = pubsub.SubscribeJSON(conn,
		string(routing.ExchangePerilTopic),
		fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, user),
		fmt.Sprintf("%s.*", routing.ArmyMovesPrefix),
		pubsub.QueueTransient,
		handlerMove(gamelogic.ArmyMove{}))
	if err != nil {
		log.Fatalf("Could not subscribe to moves: %v", err)
	}

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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
