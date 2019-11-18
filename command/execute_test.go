package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericebersohl/gobottas/discussion"
	"testing"
)

/*
Cases:
- nil helpdata struct
*/
func TestHelpExecutor(t *testing.T) {
	s := MockSession{}
	r := NewRegistry()

	nilHD := Message{}

	tests := []struct {
		name    string
		in      *Message
		wantErr bool
	}{
		{name: "nilHD", in: &nilHD, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := HelpExecutor(&s, r, test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}
		})
	}
}

/*
Cases:
- nil memedata struct
*/
func TestMemeExecutor(t *testing.T) {
	s := discordgo.Session{}
	r := NewRegistry()
	nilMD := Message{}

	tests := []struct {
		name    string
		in      *Message
		wantErr bool
	}{
		{name: "nilMD", in: &nilMD, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := MemeExecutor(&s, r, test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}
		})
	}
}

/*
Cases:
- Nil QueueData
- Add: normal, with description, too few args
- Remove: normal, too few args
- Next: normal
- BUmp: normal, too few args
- Skip: normal, too few args
- Attach: normal, too few args
- Detach: normal, too few args, bad str to num conversion
- List: normal
- Error: error state passed to executor
*/
func TestQueueExecutor(t *testing.T) {
	s := MockSession{}
	r := NewRegistry()

	nilqd := Message{CommandType: Queue, QueueData: nil}
	addname := Message{CommandType: Queue, Args: []string{"tname"}, QueueData: &discussion.QueueData{Command: discussion.QAdd}, Source: &Source{ChannelId: 0}}
	adddesc := Message{CommandType: Queue, Args: []string{"tname2", "tdesc"}, QueueData: &discussion.QueueData{Command: discussion.QAdd}, Source: &Source{ChannelId: 0}}
	addargs := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QAdd}, Source: &Source{ChannelId: 0}}
	remove := Message{CommandType: Queue, Args: []string{"tname"}, QueueData: &discussion.QueueData{Command: discussion.QRemove}, Source: &Source{ChannelId: 0}}
	removeargs := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QRemove}, Source: &Source{ChannelId: 0}}
	next := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QNext}, Source: &Source{ChannelId: 0}}
	bump := Message{CommandType: Queue, Args: []string{"tname2"}, QueueData: &discussion.QueueData{Command: discussion.QBump}, Source: &Source{ChannelId: 0}}
	bumpargs := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QBump}, Source: &Source{ChannelId: 0}}
	skip := Message{CommandType: Queue, Args: []string{"tname2"}, QueueData: &discussion.QueueData{Command: discussion.QSkip}, Source: &Source{ChannelId: 0}}
	skipargs := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QSkip}, Source: &Source{ChannelId: 0}}
	attach := Message{CommandType: Queue, Args: []string{"tname2", "google.com"}, QueueData: &discussion.QueueData{Command: discussion.QAttach}, Source: &Source{ChannelId: 0}}
	attachargs := Message{CommandType: Queue, Args: []string{"tname2"}, QueueData: &discussion.QueueData{Command: discussion.QAttach}, Source: &Source{ChannelId: 0}}
	detach := Message{CommandType: Queue, Args: []string{"tname2", "0"}, QueueData: &discussion.QueueData{Command: discussion.QDetach}, Source: &Source{ChannelId: 0}}
	detachnum := Message{CommandType: Queue, Args: []string{"tname2", "xylophone"}, QueueData: &discussion.QueueData{Command: discussion.QDetach}, Source: &Source{ChannelId: 0}}
	detachargs := Message{CommandType: Queue, Args: []string{"tname2"}, QueueData: &discussion.QueueData{Command: discussion.QDetach}, Source: &Source{ChannelId: 0}}
	list := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QList}, Source: &Source{ChannelId: 0}}
	errorCT := Message{CommandType: Queue, Args: []string{}, QueueData: &discussion.QueueData{Command: discussion.QError}, Source: &Source{ChannelId: 0}}

	tests := []struct {
		name         string
		in           *Message
		wantErr      bool
		wantQCommand discussion.QueueCommand
		wantCommand  CommandType
	}{
		{name: "nil-qd", in: &nilqd, wantErr: true, wantCommand: Queue},
		{name: "add-name", in: &addname, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QAdd},
		{name: "add-desc", in: &adddesc, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QAdd},
		{name: "add-args", in: &addargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "remove-normal", in: &remove, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QRemove},
		{name: "remove-args", in: &removeargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "next-normal", in: &next, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QNext},
		{name: "bump-normal", in: &bump, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QBump},
		{name: "bump-args", in: &bumpargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "skip-normal", in: &skip, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QSkip},
		{name: "skip-args", in: &skipargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "attach-normal", in: &attach, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QAttach},
		{name: "attach-args", in: &attachargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "detach-normal", in: &detach, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QDetach},
		{name: "detach-bad-num", in: &detachnum, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "detach-args", in: &detachargs, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QError},
		{name: "list-normal", in: &list, wantErr: false, wantCommand: Queue, wantQCommand: discussion.QList},
		{name: "error", in: &errorCT, wantErr: false, wantCommand: Error, wantQCommand: discussion.QError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := QueueExecutor(s, r, test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if !test.wantErr {
				if test.in.CommandType != test.wantCommand {
					t.Errorf("CT != wantCT (ct = %v, want = %v)", test.in.CommandType, test.wantCommand)
				}

				if test.in.QueueData.Command != test.wantQCommand {
					t.Errorf("QueuCommand != wantQCommand (QC = %v, wantQC = %v)", test.in.QueueData.Command, test.wantQCommand)
				}
			}
		})
	}
}

// mock for Session interface
type MockSession struct{}

func (m MockSession) ChannelMessageSend(string, string) (*discordgo.Message, error) {
	return nil, nil
}

func (m MockSession) ChannelMessageSendEmbed(string, *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return nil, nil
}