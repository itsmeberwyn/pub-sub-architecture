package main

import (
	"fmt"
	"log"

	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	durable   = 1
	transient = 2
)

func main() {
	var connectString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectString)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

    ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", durable)
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
}

func publishMessage[T any](ch *amqp.Channel, exchange string, key string, message T) error {
    err := pubsub.PublishJSON(ch, exchange, key, message)
    if err != nil {
        log.Fatal(err)
    }
    return err
}

