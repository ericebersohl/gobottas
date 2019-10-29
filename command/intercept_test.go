package command

import (
	"github.com/ericebersohl/gobottas/discussion"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestHelpInterceptor(t *testing.T) {

	noCommand := Message{CommandType: None}
	unrecognized := Message{CommandType: Unrecognized}
	help := Message{CommandType: Help}
	meme := Message{CommandType: Meme}
	errorMsg := Message{CommandType: Error}

	tests := []struct {
		name      string
		in        *Message
		out       CommandType
		outString string
		wantNil   bool
	}{
		{name: "no-command", in: &noCommand, out: None, wantNil: true},
		{name: "unrec", in: &unrecognized, out: Unrecognized, wantNil: false},
		{name: "help", in: &help, out: Help, wantNil: false},
		{name: "meme", in: &meme, out: Meme, wantNil: false},
		{name: "error", in: &errorMsg, out: Error, wantNil: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := HelpInterceptor(test.in)
			if err != nil {
				t.Errorf("help interceptor error: %v", err)
			}
			if (test.in.HelpData == nil) != test.wantNil {
				t.Errorf("helpData not initialized properly.")
			}

			if test.in.CommandType != test.out {
				t.Errorf("HelpData has wrong command type (want = %s, got = %s)", test.out.String(), test.in.CommandType.String())
			}
		})
	}
}

func TestMemeInterceptor(t *testing.T) {

	none := Message{CommandType: None}
	unrecognized := Message{CommandType: Unrecognized}
	help := Message{CommandType: Help}
	meme := Message{CommandType: Meme}
	errorMsg := Message{CommandType: Error}

	tests := []struct {
		name    string
		in      *Message
		out     CommandType
		wantNil bool
	}{
		{name: "none", in: &none, out: None, wantNil: true},
		{name: "unrec", in: &unrecognized, out: Unrecognized, wantNil: true},
		{name: "help", in: &help, out: Help, wantNil: true},
		{name: "meme", in: &meme, out: Meme, wantNil: false},
		{name: "error", in: &errorMsg, out: Error, wantNil: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := MemeInterceptor(test.in)
			if err != nil {
				t.Errorf("meme interceptor error: %v", err)
			}

			if (test.in.MemeData == nil) != test.wantNil {
				t.Errorf("meme struct incorrectly initialized")
			}

			if test.in.CommandType != test.out {
				t.Errorf("in and out CommandTypes don't match (want = %s, got = %s)", test.out.String(), test.in.CommandType.String())
			}
		})
	}
}

func TestQueueInterceptor(t *testing.T) {
	none := Message{CommandType: Unrecognized, Args: []string{"add", "topicName", "description"}}
	err := Message{CommandType: Queue, Args: []string{"not", "relevant", "to", "queue"}}
	add := Message{CommandType: Queue, Args: []string{"add", "topicName", "description"}}
	remove := Message{CommandType: Queue, Args: []string{"remove", "topicName"}}
	next := Message{CommandType: Queue, Args: []string{"next"}}
	bump := Message{CommandType: Queue, Args: []string{"bump", "topicName"}}
	skip := Message{CommandType: Queue, Args: []string{"skip", "topicName"}}
	attach := Message{CommandType: Queue, Args: []string{"attach", "topicName", "source"}}
	detach := Message{CommandType: Queue, Args: []string{"detach", "topicName", "source"}}

	tests := []struct {
		name     string
		in       *Message
		out      discussion.QueueCommand
		wantArgs []string
		wantNil  bool
	}{
		{name: "error", in: &err, out: discussion.QError, wantArgs: []string{"not", "relevant", "to", "queue"}, wantNil: false},
		{name: "add", in: &add, out: discussion.QAdd, wantArgs: []string{"topicName", "description"}, wantNil: false},
		{name: "remove", in: &remove, out: discussion.QRemove, wantArgs: []string{"topicName"}, wantNil: false},
		{name: "next", in: &next, out: discussion.QNext, wantArgs: []string{}, wantNil: false},
		{name: "bump", in: &bump, out: discussion.QBump, wantArgs: []string{"topicName"}, wantNil: false},
		{name: "skip", in: &skip, out: discussion.QSkip, wantArgs: []string{"topicName"}, wantNil: false},
		{name: "attach", in: &attach, out: discussion.QAttach, wantArgs: []string{"topicName", "source"}, wantNil: false},
		{name: "detach", in: &detach, out: discussion.QDetach, wantArgs: []string{"topicName", "source"}, wantNil: false},
		{name: "none", in: &none, out: discussion.QError, wantArgs: []string{"add", "topicName", "description"}, wantNil: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := QueueInterceptor(test.in)
			if err != nil {
				t.Errorf("queue interceptor got an error: %v", err)
			}

			if (test.in.QueueData == nil) != test.wantNil {
				t.Errorf("queue data is nil when it shouldn't be")
			}

			if !test.wantNil {
				if test.in.QueueData.Command != test.out {
					t.Errorf("commands don't match (want = %s, got = %s)", test.out.String(), test.in.QueueData.Command.String())
				}

				if !cmp.Equal(test.wantArgs, test.in.Args) {
					t.Errorf("args != wantArgs (args = %v, wantArgs = %v)", test.in.Args, test.wantArgs)
				}
			}
		})
	}
}
