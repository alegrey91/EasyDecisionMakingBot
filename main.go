package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	apiToken string
	verbose  bool
	version  string = "v0.0.1"
)

func main() {
	flag.StringVar(&apiToken, "token", "", "Telegram API Token")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbosity")
	v := flag.Bool("version", false, "Print software version")
	flag.Parse()

	// Print version and exit.
	if *v {
		fmt.Println(version)
		return
	}
	// check if API Token has been provided by flags.
	if apiToken == "" {
		fmt.Println(fmt.Errorf("Please provide token.\n"))
	}

	// create new bot instance.
	b, err := tb.NewBot(tb.Settings{
		Token:  apiToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// setup a new seed.
	rand.Seed(time.Now().UnixNano())

	// reply to help message sending instructions on how to use bot's commands.
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Chat, fmt.Sprintf("Welcome to the EasyDecisionMakingBot page\nPlease take a look here:\n%s", help()))
		if verbose {
			log.Printf("Replied to /starticommand\n")
		}
	})

	// reply to help message sending instructions on how to use bot's commands.
	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Chat, help())
		if verbose {
			log.Printf("Replied to /help command\n")
		}
	})

	// flip the coin and reply with "head" or "tail".
	b.Handle("/coin", func(m *tb.Message) {
		coinFace, err := getCoinFace(rand.Intn(2))
		if err != nil {
			log.Fatal(err)
		}

		b.Send(m.Chat, coinFace)
		if verbose {
			log.Printf("Replied to /coin command with %s\n", coinFace)
		}
	})

	// throw the dice and reply with a number between 1 and 6.
	b.Handle("/dice", func(m *tb.Message) {
		diceFace := rand.Intn(6) + 1
		b.Send(m.Chat, fmt.Sprintf("🎲 %d", diceFace))
		if verbose {
			log.Printf("Replied to /dice command with %d\n", diceFace)
		}
	})

	// reply in case of invalid command.
	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Chat, fmt.Sprintf("Invalid command provided"))
		if verbose {
			log.Printf("Invalid command provided\n")
		}
	})

	b.Start()
}

// Return help string.
func help() string {
	return fmt.Sprintf("/help  print this help message\n/coin  flip a coin\n/dice  throw the dice\n")
}

// Return "Head" or "Tail" from int value.
func getCoinFace(value int) (string, error) {
	switch value {
	case 0:
		return "Head", nil
	case 1:
		return "Tail", nil
	default:
		return "", errors.New("Invalid coin value")
	}
}
