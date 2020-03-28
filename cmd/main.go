package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/core"
	"github.com/ericebersohl/gobottas/discussion"
	"github.com/ericebersohl/gobottas/meme"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	DefaultChannelBuffer = 15
	DefaultDirPath       = "/store"
)

var (
	channelBuffer   int
	dirPath         string
	discussionQueue bool
	memeStash       bool
)

func init() {
	flag.IntVar(&channelBuffer, "buf", DefaultChannelBuffer, "Set the buffer size for the message channel [Default: 15] (must be greater than 0, 1 is an unbuffered channel)")
	flag.StringVar(&dirPath, "dir", DefaultDirPath, "Set the location on the local machine for gobottas to store files [Default: /store]")
	flag.BoolVar(&memeStash, "m", false, "Whether to include the MemeStash feature (default = false)")
	flag.BoolVar(&discussionQueue, "q", false, "Whether to include the Discussion Queue feature (default = false)")
}

// Returns a message handler for discord messages, a function is needed since we want the handler to have access to the channel
func messageHandler(c chan *gb.Message, r gb.Registry) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		// Ignore bot messages
		if m.Author.Bot {
			return
		}

		// Parse the message and send all messages go through the bot
		msg, err := r.Parse(m.Message)
		if err != nil {
			log.Printf("ignoring message (id = %s) due to error: %v", m.ID, err)
		}

		// send the parsed message through the channel
		c <- msg
	}
}

// function to be run in goroutine that handles parsed Messages coming out of the channel
func handleCommands(c chan *gb.Message, r gb.Registry, s *discordgo.Session) {

	// wait for messages come through, block until they do
	for msg := range c {
		err := r.Intercept(msg)
		if err != nil {
			log.Printf("CommandHandler: %v", err)
		}

		if msg.Response.Embed != nil {
			fmt.Printf("HANDLER: %v\n", msg.Response.Embed.Fields)
		}

		err = r.Execute(msg, s)
		if err != nil {
			log.Printf("CommandHandler: %v", err)
		}
	}
}

func getRegistryOpts() (opts []core.RegistryOpt) {
	// set the dir path
	opts = append(opts, core.WithPath(dirPath))

	// set discussion queue opts if applicable
	if discussionQueue {
		q := discussion.NewQueue()
		opts = append(opts, core.WithQueue(q))
		opts = append(opts, core.WithInterceptor(gb.Queue, discussion.Interceptor(q)))
	}

	// set the memeStash option
	if memeStash {
		s := meme.DefaultStash(dirPath)
		opts = append(opts, core.WithStash(s))
		opts = append(opts, core.WithInterceptor(gb.Meme, meme.Interceptor(&s)))
	}

	return opts
}

func main() {
	// parse flags
	flag.Parse()

	// create local file directory
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create local directory at %s: %v", dirPath, err)
	}

	// Load Environment Variables
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load variables from .env.\n%v", err)
	}

	// build a registry
	registry := core.NewRegistry(getRegistryOpts()...)

	// make a channel through which commands are sent and executed
	cmdChannel := make(chan *gb.Message, channelBuffer)
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
	log.Printf("Gobottas initialized with %d interceptors.", len(registry.Interceptors))

	// keep main open indefinitely
	<-make(chan interface{})
}
