package command

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

// todo(ee): better testing, lacking a good way to get a test session to use
func TestHelpExecutor(t *testing.T) {
	s := discordgo.Session{}

	nilHD := Message{}

	tests := []struct{
		name string
		in *Message
		wantErr bool
	}{
		{name: "nilHD", in: &nilHD, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := HelpExecutor(&s, test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}
		})
	}
}

// todo(ee): better testing; need better session solution (see above)
func TestMemeExecutor(t *testing.T) {
	s := discordgo.Session{}
	nilMD := Message{}

	tests := []struct{
		name string
		in *Message
		wantErr bool
	}{
		{name: "nilMD", in: &nilMD, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := MemeExecutor(&s, test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}
		})
	}
}