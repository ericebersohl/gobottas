package main

import (
	"github.com/ericebersohl/gobottas/command"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	channelBuffer = 15
)

// Returns a message handler for discord messages, a function is needed since we want the handler to have access to the channel
func messageHandler(c chan *core.Message, r *core.Registry) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Parse the message and send all messages go through the bot
		msg, err := core.Parse(m.Message, r)
		if err != nil {
			log.Printf("ignoring message (id = %s) due to error: %v", m.ID, err)
		}

		// send the parsed message through the channel
		c <- msg
	}
}

// function to be run in goroutine that handles parsed Messages coming out of the channel
func handleCommands(c chan *core.Message, r *core.Registry, s *discordgo.Session) {

	// wait for messages come through, block until they do
	for msg := range c {
		err := r.Intercept(msg)
		if err != nil {
			log.Printf("CommandHandler: %v", err)
		}

		err = r.Execute(msg, s)
		if err != nil {
			log.Printf("CommandHandler: %v", err)
		}
	}
}

func main() {
	// Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load variables from .env.\n%v", err)
	}

	// build a registry
	// todo(ee): figure out how to persist functions
	registry := core.NewRegistry()

	// make a channel through which commands are sent and executed
	cmdChannel := make(chan *core.Message, channelBuffer)
	defer close(cmdChannel)

	// Get Connection to Server
	discord, err := discordgo.New("Bot " + os.Getenv("AUTH"))
	if err != nil {
		log.Fatalf("Failed to create discord client.\n%v", err)
	}

	// add a new message handler
	discord.AddHandler(messageHandler(cmdChannel, registry))

	// Open the connection
	if err := discord.Open(); err != nil {
		log.Fatalf("Failed to open connection.\n%v", err)
	}
	defer discord.Close()

	// spin up a goroutine to handle any commands that come through the channel
	go handleCommands(cmdChannel, registry, discord)

	// log that gobottas is running
	log.Printf("Gobottas initialized.")

	// keep main open indefinitely
	<-make(chan interface{})
}
