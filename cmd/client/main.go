package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+user, routing.PauseKey, durable)

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
