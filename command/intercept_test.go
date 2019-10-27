package command

import "testing"

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
