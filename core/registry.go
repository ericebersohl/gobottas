package core

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discussion"
	"log"
	"regexp"
	"strconv"
)

const (
	DefaultCommandPrefix = '&'
	DefaultDirPath       = "../build"
)

// Contains Gobottas functions and data
type Registry struct {
	Interceptors    map[gb.Command]gb.Interceptor // all built-in interceptors
	DirPath         string                        // path to local data
	CommandPrefix   uint8                         // character that precedes all Gobottas commands
	DiscussionQueue *discussion.Queue             // the data structure that holds discussion queue data
}

type RegistryOpt func(*Registry)

func NewRegistry(opts ...RegistryOpt) *Registry {
	r := Registry{
		Interceptors:  make(map[gb.Command]gb.Interceptor),
		DirPath:       DefaultDirPath,
		CommandPrefix: DefaultCommandPrefix,
	}

	for _, o := range opts {
		o(&r)
	}

	return &r
}

// Opt Functions
func WithQueue(q *discussion.Queue) RegistryOpt {
	return func(r *Registry) {
		r.DiscussionQueue = q
	}
}

func WithPrefix(p uint8) RegistryOpt {
	return func(r *Registry) {
		r.CommandPrefix = p
	}
}

// set the interceptor for the given command type; will overwrite any existing interceptor
func WithInterceptor(c gb.Command, i gb.Interceptor) RegistryOpt {
	return func(r *Registry) {
		r.Interceptors[c] = i
	}
}

func WithPath(s string) RegistryOpt {
	return func(r *Registry) {
		r.DirPath = s
	}
}

// Function to parse incoming messages
func (r *Registry) Parse(dMsg *discordgo.Message) (cmd *gb.Message, err error) {
	// Default to command none
	cmd = &gb.Message{
		Command:  gb.None,
		Response: &gb.Response{},
	}

	// check for nil in msg
	if dMsg == nil {
		cmd.Command = gb.Error
		return cmd, errors.New("discord message is nil")
	}

	// Check for nil author
	if dMsg.Author == nil {
		cmd.Command = gb.Error
		return cmd, errors.New("discord message has empty author")
	}

	// Create source
	src := gb.Source{
		Content: dMsg.Content,
	}

	// Convert AuthorId and ChannelId
	src.AuthorId, err = gb.ToSnowflake(dMsg.Author.ID)
	if err != nil {
		log.Printf("Failed to parse author from discord: %v", err)
		return cmd, err
	}

	src.ChannelId, err = gb.ToSnowflake(dMsg.ChannelID)
	if err != nil {
		log.Printf("Failed to parse channel from discord: %v", err)
		return cmd, err
	}

	// get username
	src.Username = dMsg.Author.Username

	// attach src to msg
	cmd.Source = &src

	// get the command
	args, err := Tokenize(cmd.Source.Content)
	if err != nil {
		log.Printf("Failed to tokenize content: %v", err)
		return cmd, err
	}

	// check for prefix
	if args[0][0] != r.CommandPrefix {
		return cmd, nil
	}

	// set the command
	cmd.Command = gb.StrToCommand(args[0][1:])
	if cmd.Command == gb.Unrecognized {
		return cmd, nil
	}

	cmd.Args = args[1:]
	return cmd, nil
}

// Function to call all Interceptors on a message
func (r *Registry) Intercept(msg *gb.Message) error {
	fmt.Printf("Intercept: %d\n", len(r.Interceptors))
	for _, i := range r.Interceptors {
		err := i(msg)
		if err != nil {
			log.Printf("Registry.Intercept: %v", err)
			return err
		}
	}
	return nil
}

// Calls the Executor to which the Registry points for the Message CommandType
func (r *Registry) Execute(msg *gb.Message, s gb.Session) error {

	// prefer embeds, then messages, then not found
	if msg.Response.Embed != nil {
		fmt.Printf("%v\n", msg.Response.Embed)
		_, err := s.ChannelMessageSendEmbed(msg.Response.ChannelId.String(), msg.Response.Embed)
		if err != nil {
			log.Printf("Error on Execute: %v", err)
			return err
		}

		// exit successfully
		return nil
	}

	// only executes if Embed is nil
	if msg.Response.Text != "" {
		_, err := s.ChannelMessageSend(msg.Response.ChannelId.String(), msg.Response.Text)
		if err != nil {
			log.Printf("Error on Execute: %v", err)
			return err
		}

		// exit sucessfully
		return nil
	}

	// Following unix norm that no response indicates success
	return nil
}

// Split the message content by spaces while leaving segments in quotes intact
// Uses regex
func Tokenize(s string) (tok []string, err error) {

	// parse the string using the dark art of regular expressions
	tok = regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`).FindAllString(s, -1)

	// remove quotes
	for i := range tok {

		// only call on strings with quotes
		if string(tok[i][0]) == "\"" {
			tok[i], err = strconv.Unquote(tok[i])
			if err != nil {
				log.Printf("Tokenize failed to Unquote string: %v", err)
				return nil, err
			}
		}
	}

	return tok, nil
}
