package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

    // err = ch.ExchangeDeclare(routing.ExchangePerilDirect, "direct", true, false, false, false, nil)
    // if err != nil {
    //     log.Fatal(err)
    // }

    // _, err = ch.QueueDeclare("pause_test", true, false, false, false, nil)
    // if err != nil {
    //     log.Fatal(err)
    // }

    // err = ch.QueueBind("pause_test", routing.PauseKey, routing.ExchangePerilDirect, false, nil)
    // if err != nil {
    //     log.Fatal(err)
    // }

	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
    if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Starting Peril server...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
