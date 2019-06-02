package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	commandPrefix string
	botID         string
)

func main() {
	// Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load variables from .env.\n%v", err)
	}

	commandPrefix = "`"

	// Get Connection to Server
	discord, err := discordgo.New("Bot " + os.Getenv("AUTH"))
	if err != nil {
		log.Fatalf("Failed to create discord client.\n%v", err)
	}

	// Open the connection
	if err := discord.Open(); err != nil {
		log.Fatalf("Failed to open connection.\n%v", err)
	}
	defer discord.Close()

	// keep main open indefinitely
	<-make(chan interface{})
}
