package discussion

import (
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discord"
	"github.com/ericebersohl/gobottas/mock"
	"testing"
)

/*
Test Cases:
- Not a Queue Message
- Nil Queue
- Bad Command
- Add: too few args, name only, name and description, duplicate
- Remove: too few args, not found (dErr)
- Next: empty queue (dErr), normal
- Bump: too few args, not found (dErr), normal
- Skip: too few args, not found (dErr), normal
- Attach: too few args, not found (dErr), normal
- Detach: too few args, bad Atoi, Index Oob (dErr), normal
*/
func TestInterceptor(t *testing.T) {
	q := NewQueue()
	eq := NewQueue()

	tests := []struct {
		name        string
		queue       *Queue
		in          *gb.Message
		wantErr     bool
		wantDiscErr bool
		wantEmbed   bool
	}{
		// Error cases
		{name: "not-queue-command", queue: q, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "nil-queue", queue: nil, in: mock.NewMessage(gb.Queue), wantErr: true, wantDiscErr: false, wantEmbed: false},
		{name: "bad-command", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("not", "valid", "args")), wantErr: false, wantDiscErr: false, wantEmbed: true},

		// Add
		{name: "add-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("add")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "add-name-only", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("add", "testName")), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "add-both", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("add", "testName2", "testDesc")), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "add-dup", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("add", "testName2")), wantErr: true, wantDiscErr: true, wantEmbed: true},

		// Remove
		{name: "rem-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("remove")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "rem-not-found", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("remove", "not-topic")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "rem-normal", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("remove", "testName")), wantErr: false, wantDiscErr: false, wantEmbed: false},

		// Next
		{name: "next-normal", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("next")), wantErr: false, wantDiscErr: false, wantEmbed: true},
		{name: "next-empty", queue: eq, in: mock.NewMessage(gb.Queue, mock.WithArgs("next")), wantErr: true, wantDiscErr: true, wantEmbed: true},

		// Bump
		{name: "bump-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("bump")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "bump-not-found", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("bump", "not-found")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "bump-normal", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("bump", "testName2")), wantErr: false, wantDiscErr: false, wantEmbed: false},

		// Skip
		{name: "skip-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("skip")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "skip-not-found", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("skip", "not-found")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "skip-normal", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("skip", "testName2")), wantErr: false, wantDiscErr: false, wantEmbed: false},

		// Attach
		{name: "attach-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("attach")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "attach-not-found", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("attach", "not-found", "https://google.com")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "attach-normal", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("attach", "testName2", "https://google.com")), wantErr: true, wantDiscErr: true, wantEmbed: true},

		// Detach
		{name: "det-too-few", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("detach")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "det-bad-atoi", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("detach", "testName2", "zer0")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "det-oob", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("detach", "testName2", "5")), wantErr: true, wantDiscErr: true, wantEmbed: true},
		{name: "det-norm", queue: q, in: mock.NewMessage(gb.Queue, mock.WithArgs("detach", "testName2", "0")), wantErr: false, wantDiscErr: false, wantEmbed: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := Interceptor(test.queue)
			err := i(test.in)
			if err != nil {
				if !test.wantErr {
					t.Errorf("err != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
				}

				// check for discord error
				_, ok := err.(discord.Error)

				if test.wantErr && ok != test.wantDiscErr {
					t.Errorf("wanted discord error, got not discord error")
				}
			}

			if err == nil && !test.wantErr {
				if (test.in.Response.Embed != nil) != test.wantEmbed {
					t.Errorf("Embed is nil when it should not be (embed == nil: %t)", test.in.Response.Embed == nil)
				}
			}
		})
	}
}
