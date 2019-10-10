package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func getDiscordMessage() discordgo.Message {
	u := discordgo.User{
		ID:            "3",
		Bot:           false,
	}

	return discordgo.Message{
		ID:              "1",
		ChannelID:       "2",
		Content:         "&help do some things",
		Author:          &u,
	}
}

func getCommandMesage() Message {
	s := Source{
		AuthorId:  3,
		ChannelId: 2,
		Content:   "&help do some things",
	}

	return Message{
		CommandType: Help,
		Source:      &s,
		Args:        []string{"do", "some", "things"},
	}
}

func TestParse(t *testing.T) {
	r := NewRegistry()
	// normal: normal, bot, not prefixed
	// nil Message
	errMessage := Message{
		CommandType: Error,
	}

	// nil user
	nilUserIn := getDiscordMessage()
	nilUserIn.Author = nil

	// bad authorid
	badAuthorIn := getDiscordMessage()
	badAuthorIn.Author.ID = "not-a-snowflake"

	// bad channel id
	badChannelIn := getDiscordMessage()
	badChannelIn.ChannelID = "not-a-channel"

	// empty content string
	emptyContentIn := getDiscordMessage()
	emptyContentIn.Content = ""

	// no error, no prefix
	noPrefixIn := getDiscordMessage()
	noPrefixIn.Content = "help do some things"
	noPrefixOut := getCommandMesage()
	noPrefixOut.CommandType = None
	noPrefixOut.Source.Content = "help do some things"
	noPrefixOut.Args = nil

	// no error, bot
	botIn := getDiscordMessage()
	botIn.Author.Bot = true
	botOut := getCommandMesage()
	botOut.CommandType = None
	botOut.Args = nil

	// no error, unrecognized command
	unRecIn := getDiscordMessage()
	unRecIn.Content = "&unrecognized do some things"
	unRecOut := getCommandMesage()
	unRecOut.Source.Content = "&unrecognized do some things"
	unRecOut.CommandType = Unrecognized
	unRecOut.Args = nil

	// normal
	in := getDiscordMessage()
	out := getCommandMesage()

	tests := []struct{
		name string
		in *discordgo.Message
		want *Message
		wantErr bool
	}{
		{name: "nil-user", in: &nilUserIn, want: &errMessage, wantErr: true},
		{name: "bad-author-id", in: &badAuthorIn, want: &errMessage, wantErr: true},
		{name: "bad-channel-id", in: &badChannelIn, want: &errMessage, wantErr: true},
		{name: "empty-content-string", in: &emptyContentIn, want: &errMessage, wantErr: true},
		{name: "no-prefix", in: &noPrefixIn, want: &noPrefixOut, wantErr: false},
		{name: "bot", in: &botIn, want: &botOut, wantErr: false},
		{name: "unrecognized", in: &unRecIn, want: &unRecOut, wantErr: false},
		{name: "normal-case", in: &in, want: &out, wantErr: false},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := Parse(test.in, r)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if s := cmp.Diff(out, test.want); s != "" {
				t.Logf(s)
				t.Fail()
			}
		})
	}
}
