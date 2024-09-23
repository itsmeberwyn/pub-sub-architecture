package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	durable   = 1
	transient = 2
)

func main() {
	user, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal(err)
	}

	var connectString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectString)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	gameState := gamelogic.NewGameState(user)
	pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+user, routing.PauseKey, transient, handlerPause(gameState))
	for {
		input := gamelogic.GetInput()
		switch input[0] {
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("I dont have an idea")
			break
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	defer fmt.Print("> ")
	return func(ps routing.PlayingState) {
		gs.HandlePause(routing.PlayingState{
			IsPaused: ps.IsPaused,
		})
	}
}

func publishMessage[T any](ch *amqp.Channel, exchange string, key string, message T) error {
	err := pubsub.PublishJSON(ch, exchange, key, message)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
