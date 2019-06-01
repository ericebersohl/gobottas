package main

import (
	"log"
	"os"

	dgo "github.com/bwmarrin/discordgo"
	gde "github.com/joho/godotenv"
)

var (
	commandPrefix string
	botID         string
)

func main() {
	err := gde.Load()
	errCheck("Failed to load from .env", err)

	discord, err := dgo.New("Bot " + os.Getenv("AUTH"))
	errCheck("Failed to create a discord client.", err)

	user, err := discord.User("@me")
	errCheck("Failed to retrieve user account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *dgo.Session, ready *dgo.Ready) {
		err = discord.UpdateStatus(0, "Gobottas initialized.")
		if err != nil {
			log.Println("Failed to set status.")
		}

		servers := discord.State.Guilds
		log.Printf("Gobottas is running on %d servers\n", len(servers))
	})

	err = discord.Open()
	errCheck("Failed to connect to discord.", err)
	defer discord.Close()

	commandPrefix = "&"

	<-make(chan struct{})
}

func errCheck(msg string, err error) {
	if err != nil {
		log.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(discord *dgo.Session, message *dgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		return // bot is talking, does not react to itself
	}

	// content := message.Content
	log.Printf("Message: %+v || From %s\n", message.Message, message.Author)
}
