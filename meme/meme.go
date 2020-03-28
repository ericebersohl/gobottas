package meme

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discord"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Meme struct {
	Meme    string    `json:"meme"`
	Added   time.Time `json:"added"`
	AddedBy string    `json:"added-by"`
}

// create a msg.MessageEmbed to be sent to the discord channel
func (m *Meme) Embed() *discordgo.MessageEmbed {
	msg := discord.NewEmbed().
		EmbedColor(gb.MemeCol).
		EmbedTitle(m.Meme).
		EmbedFooter(fmt.Sprintf("Added by %s", m.AddedBy), "", "").
		EmbedTimestamp(m.Added)
	return msg.MessageEmbed
}

// create a new meme struct with the provided data
func NewMeme(meme, user string) *Meme {
	m := Meme{
		Meme:    meme,
		Added:   time.Now(),
		AddedBy: user,
	}
	return &m
}

// Slice of currently stored memes
type Stash struct {
	Memes     []*Meme `json:"memes"` // the memes in the stash
	LocalPath string  `json:"path"`  // path to local backup
}

// The default stash
func DefaultStash(localPath string) Stash {
	s := Stash{
		Memes: []*Meme{
			{Meme: "When did I do dangerous driving?", AddedBy: "Default Meme", Added: time.Now()},
			{Meme: "Stay out. IN! IN! IN! IN! IN! IN! IN!", AddedBy: "Default Meme", Added: time.Now()},
			{Meme: "Is his career over!?", AddedBy: "Default Meme", Added: time.Now()},
		},
		LocalPath: localPath,
	}

	return s
}

// Save the stash to a local folder
// Note that since this is a dockerized app, "local" means "inside the container"
// A volume is required for more permanent storage
func (s *Stash) Save(path string) error {
	// check for s != nil
	if s == nil {
		return fmt.Errorf("cannot save a nil stash")
	}

	// get bytes
	data, err := json.Marshal(s)
	if err != nil {
		log.Printf("Stash save error: %v", err)
		return err
	}

	// write to file
	err = ioutil.WriteFile(fmt.Sprintf("%s/meme.json", path), data, 0644)
	if err != nil {
		log.Printf("WriteFile error: %v", err)
		return err
	}

	return nil
}

// Load a stash from the meme.json file in the specified path
func (s *Stash) Load(path string) error {
	// get data
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/meme.json", path))
	if err != nil {
		log.Printf("Load error: %v", err)
		return err
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		log.Printf("Unmarshal err: %v", err)
		return err
	}

	return nil
}

// Enumeration of meme commands
type Command int

const (
	M = iota
	MAdd
	MRemove
	MList
	MError
)

func (c Command) String() string {
	return [...]string{"Meme", "Add", "Remove", "List"}[c]
}

func ArgToCommand(arg string) Command {
	switch arg {
	case "":
		return M
	case "add":
		return MAdd
	case "remove":
		return MRemove
	case "list":
		return MList
	default:
		return MError

	}
}

func Interceptor(s *Stash) gb.Interceptor {
	return func(msg *gb.Message) error {
		// skip if not a meme message
		if msg.Command != gb.Meme {
			return nil
		}

		// error if there isn't a stash
		if s == nil {
			return errors.New("cannot intercept without a stash")
		}

		// error if the meme stash is empty
		if len(s.Memes) == 0 {
			return errors.New("meme stash is empty")
		}

		// This command is returned to the same channel
		msg.Response.ChannelId = msg.Source.ChannelId

		// get the first arg (since args might be nil, have to check this way to avoid nil pointer deref)
		var arg string
		if len(msg.Args) > 0 {
			arg = msg.Args[0]
		}

		cmd := ArgToCommand(arg)
		switch cmd {
		case M:
			// select a meme at random
			rand.Seed(time.Now().UnixNano())
			fmt.Print("==", len(s.Memes), "==")
			meme := s.Memes[rand.Intn(len(s.Memes))]

			// set the embed, return nil
			msg.Response.Embed = meme.Embed()
			return nil

		case MAdd:
			// check args
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Args", "add requires a meme string\n").Embed()
				return nil
			}

			// create the meme
			meme := NewMeme(msg.Args[1], msg.Source.Username)

			// add it to the list
			s.Memes = append(s.Memes, meme)

			// save the list
			err := s.Save(s.LocalPath)
			if err != nil {
				msg.Response.Embed = discord.Error{
					Name: "Meme Save Error",
					Desc: err.Error(),
				}.Embed()
			}
			return nil

		case MRemove:
			// check args
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Args", "Remove requires a position argument of type int\n").Embed()
				return nil
			}

			// attempt to convert to index
			if idx, err := strconv.Atoi(msg.Args[1]); err == nil {
				if idx >= len(s.Memes) || idx < 0 {
					msg.Response.Embed = discord.NewError("Out of Bounds", "The provided index does not correspond to a meme\n").Embed()
					return nil
				}

				s.Memes = append(s.Memes[:idx], s.Memes[idx+1:]...)
			} else {
				msg.Response.Embed = discord.NewError("Invalid Index", "The provided meme index could not be converted to an integer\n").Embed()
				return nil
			}

			// save the list
			err := s.Save(s.LocalPath)

			if err != nil {
				msg.Response.Embed = discord.Error{
					Name: "Meme Save Error",
					Desc: err.Error(),
				}.Embed()
			}
			return nil

		case MList:
			var memes []string
			for i, m := range s.Memes {
				memes = append(memes, fmt.Sprintf("%d: %s", i, m.Meme))
			}

			// add backticks for discord
			memes = append([]string{"```"}, memes...)
			memes = append(memes, "```")

			e := discord.NewEmbed().
				EmbedColor(gb.MemeCol).
				EmbedTitle("Memes").
				EmbedTimestamp(time.Now()).
				EmbedDescription(strings.Join(memes, "\n"))

			msg.Response.Embed = e.MessageEmbed
			return nil

		case MError:
			msg.Response.Embed = discord.NewError("Unrecognized Command", "Gobottas did not recognize your command.").Embed()
			return nil
		}

		return errors.New("reached end of interceptor without returning from the switch")
	}
}
