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

	ch, err := conn.Channel()

	gameState := gamelogic.NewGameState(user)
	pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+user, routing.PauseKey, transient,
		handlerPause(gameState))
	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+user, routing.ArmyMovesPrefix+".*",
		transient, handlerMove(ch, gameState))
	pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, "war", "war.*", durable, handlerWar(ch, gameState))
	for {
		input := gamelogic.GetInput()
		switch input[0] {
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
			err = publishMessage(ch, routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+user, move)
            fmt.Println("Move published successfully")
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

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) string {
	defer fmt.Print(">")
	return func(ps routing.PlayingState) string {
		gs.HandlePause(routing.PlayingState{
			IsPaused: ps.IsPaused,
		})
        return "Ack"
	}
}

func handlerMove(ch *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.ArmyMove) string {
	defer fmt.Print(">")
	return func(ga gamelogic.ArmyMove) string {
        outcome := gs.HandleMove(gamelogic.ArmyMove{
			Player:     ga.Player,
			Units:      ga.Units,
			ToLocation: ga.ToLocation,
		})
        if outcome == gamelogic.MoveOutComeSafe {
            return "Ack"
        }
        if outcome == gamelogic.MoveOutcomeMakeWar {
            err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, 
                routing.WarRecognitionsPrefix + "." + gs.GetUsername(),
                gamelogic.RecognitionOfWar{
                    Attacker: gs.Player,
                    Defender: ga.Player,
                })
            if err != nil {
                log.Fatal(err)
            }
            return "NackRequeue"
        }
        return "NackDiscard"
	}
}

func handlerWar(ch *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) string {
    defer fmt.Print(">")
    return func(rec gamelogic.RecognitionOfWar) string {
        outcome, _, _ := gs.HandleWar(rec)
        log.Println(outcome, " hello 123")
        if outcome == gamelogic.WarOutcomeNotInvolved {
            return "NackRequeue"
        } else if outcome == gamelogic.WarOutcomeNoUnits {
            return "NackDiscard"
        } 

        if outcome == gamelogic.WarOutcomeOpponentWon {
            log.Println("WarOutcomeOpponentWon")
            return "Ack"
        }
        if outcome == gamelogic.WarOutcomeYouWon {
            log.Println("WarOutcomeYouWon")
            return "Ack"
        } 
        if outcome == gamelogic.WarOutcomeDraw {
            log.Println("WarOutcomeDraw")
            return "Ack"
        }
        return "NackDiscard"
    }
}

func publishMessage[T any](ch *amqp.Channel, exchange string, key string, message T) error {
	err := pubsub.PublishJSON(ch, exchange, key, message)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
