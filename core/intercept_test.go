package core

import (
	"github.com/ericebersohl/gobottas/discussion"
	"github.com/ericebersohl/gobottas/model"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestHelpInterceptor(t *testing.T) {

	noCommand := model.Message{CommandType: model.None}
	unrecognized := model.Message{CommandType: model.Unrecognized}
	help := model.Message{CommandType: model.Help}
	meme := model.Message{CommandType: model.Meme}
	errorMsg := model.Message{CommandType: model.Error}
	args := model.Message{CommandType: model.Meme, Args: []string{"meme"}}

	tests := []struct {
		name       string
		in         *model.Message
		out        model.CommandType
		outString  string
		wantNil    bool
		wantSubMsg bool
	}{
		{name: "no-command", in: &noCommand, out: model.None, wantNil: true},
		{name: "unrec", in: &unrecognized, out: model.Unrecognized, wantNil: false},
		{name: "help", in: &help, out: model.Help, wantNil: false},
		{name: "meme", in: &meme, out: model.Meme, wantNil: false},
		{name: "error", in: &errorMsg, out: model.Error, wantNil: false},
		{name: "args", in: &args, out: model.Meme, wantNil: false, wantSubMsg: true},
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

			if test.wantSubMsg {
				if test.in.HelpData.SubMsg == "" {
					t.Errorf("HelpData has empty submsg (want not empty)")
				}
			}
		})
	}
}

func TestMemeInterceptor(t *testing.T) {

	none := model.Message{CommandType: model.None}
	unrecognized := model.Message{CommandType: model.Unrecognized}
	help := model.Message{CommandType: model.Help}
	meme := model.Message{CommandType: model.Meme}
	errorMsg := model.Message{CommandType: model.Error}

	tests := []struct {
		name    string
		in      *model.Message
		out     model.CommandType
		wantNil bool
	}{
		{name: "none", in: &none, out: model.None, wantNil: true},
		{name: "unrec", in: &unrecognized, out: model.Unrecognized, wantNil: true},
		{name: "help", in: &help, out: model.Help, wantNil: true},
		{name: "meme", in: &meme, out: model.Meme, wantNil: false},
		{name: "error", in: &errorMsg, out: model.Error, wantNil: true},
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
	none := model.Message{CommandType: model.Unrecognized, Args: []string{"add", "topicName", "description"}}
	err := model.Message{CommandType: model.Queue, Args: []string{"not", "relevant", "to", "queue"}}
	add := model.Message{CommandType: model.Queue, Args: []string{"add", "topicName", "description"}}
	remove := model.Message{CommandType: model.Queue, Args: []string{"remove", "topicName"}}
	next := model.Message{CommandType: model.Queue, Args: []string{"next"}}
	bump := model.Message{CommandType: model.Queue, Args: []string{"bump", "topicName"}}
	skip := model.Message{CommandType: model.Queue, Args: []string{"skip", "topicName"}}
	attach := model.Message{CommandType: model.Queue, Args: []string{"attach", "topicName", "source"}}
	detach := model.Message{CommandType: model.Queue, Args: []string{"detach", "topicName", "source"}}

	tests := []struct {
		name     string
		in       *model.Message
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
