package main

import (
	"fmt"
	"log"
	// "os"
	// "os/signal"

	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var connectString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectString)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

    gamelogic.PrintServerHelp()
    for {
        input := gamelogic.GetInput()
        switch input[0] {
            case "pause":
                fmt.Println("Pause")
                publishMessage(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
                    IsPaused: true,
                })
            case "resume":
                fmt.Println("Resume")
                publishMessage(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
                    IsPaused: false,
                })
            case "exit":
                fmt.Println("Exit")
                break;
            default:
                fmt.Println("I dont have an idea")
                break;
        }
    }

	// fmt.Println("Starting Peril server...")
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}

func publishMessage[T any](ch *amqp.Channel, exchange string, key string, message T) error {
    err := pubsub.PublishJSON(ch, exchange, key, message)
    if err != nil {
        log.Fatal(err)
    }
    return err
}
