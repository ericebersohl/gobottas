package core

import (
	"github.com/bwmarrin/discordgo"
	gb "github.com/ericebersohl/gobottas"
	"github.com/google/go-cmp/cmp"
	"testing"
)

/*
Test Cases:
- nil dgo msg
- nil author
- bad authorid, bad channelid
*/
func TestRegistry_Parse(t *testing.T) {
	r := NewRegistry()

	tests := []struct {
		name     string
		in       *discordgo.Message
		wantErr  bool
		wantType gb.Command
	}{
		{name: "nil dgo msg", in: nil, wantErr: true, wantType: gb.Error},
		{name: "nil auth", in: &discordgo.Message{Author: nil}, wantErr: true, wantType: gb.Error},
		{name: "bad authid", in: &discordgo.Message{Author: &discordgo.User{ID: "id"}}, wantErr: true, wantType: gb.Error},
		{name: "bad chanid", in: &discordgo.Message{Author: &discordgo.User{ID: "0"}, ChannelID: "id"}, wantErr: true, wantType: gb.Error},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := r.Parse(test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if err == nil && !test.wantErr {
				if out.Command != test.wantType {
					t.Errorf("out != want (out = %s, want = %s)", out.Command.String(), test.wantType.String())
				}
			}
		})
	}
}

/*
Test Cases:
- nil string
- start quotes, end quotes
- single quote with whitespace
*/
func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		out     []string
		wantErr bool
	}{
		{name: "nil", in: "", out: nil, wantErr: false},
		{name: "start-quotation", in: `"start quotes" and other args`, out: []string{"start quotes", "and", "other", "args"}, wantErr: false},
		{name: "end-quotation", in: `other args "end quotes"`, out: []string{"other", "args", "end quotes"}, wantErr: false},
		{name: "single-quote", in: `a " b`, out: []string{"a", "b"}, wantErr: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Tokenize(test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if !cmp.Equal(test.out, got) {
				t.Errorf("out != got (%s)", cmp.Diff(test.out, got))
			}
		})
	}
}
