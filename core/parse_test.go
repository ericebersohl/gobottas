package core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func getDiscordMessage() discordgo.Message {
	u := discordgo.User{
		ID:  "3",
		Bot: false,
	}

	return discordgo.Message{
		ID:        "1",
		ChannelID: "2",
		Content:   "&help do some things",
		Author:    &u,
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

	tests := []struct {
		name    string
		in      *discordgo.Message
		want    *Message
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

/*
Cases:
- string without quotes
- normal case
- single quote
- empty quotes
- quotes at the beginning
- quotes at the end
- quotes around a single word
- quotes around whitespace
- multiple quoted sections
*/
func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []string
		wantErr bool
	}{
		{name: "no-quotes", in: "some space delimited words", want: []string{"some", "space", "delimited", "words"}, wantErr: false},
		{name: "quotes", in: "some \"non space delimited\" words", want: []string{"some", "non space delimited", "words"}, wantErr: false},
		{name: "single", in: "random single \"quote", want: []string{"random", "single", "quote"}, wantErr: false},
		{name: "empty", in: "empty \"\" quotes", want: []string{"empty", "", "quotes"}, wantErr: false},
		{name: "beginning", in: "\"quotes at\" the beginning", want: []string{"quotes at", "the", "beginning"}, wantErr: false},
		{name: "end", in: "quotes at \"the end\"", want: []string{"quotes", "at", "the end"}, wantErr: false},
		{name: "ineffectual", in: "non \"effective\" quotes", want: []string{"non", "effective", "quotes"}, wantErr: false},
		{name: "quoted-space", in: "\"  \" test", want: []string{"  ", "test"}, wantErr: false},
		{name: "multiple", in: "a \"b c\" d \"e f\" g \"h i j\" k l", want: []string{"a", "b c", "d", "e f", "g", "h i j", "k", "l"}, wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tok, err := Tokenize(test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if !test.wantErr {
				if !reflect.DeepEqual(test.want, tok) {
					t.Errorf("Slices don't agree:\n%q\n%q", test.want, tok)
				}
			}
		})
	}
}
