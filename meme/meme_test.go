package meme

import (
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discord"
	"github.com/ericebersohl/gobottas/mock"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

/*
Test Cases:
- not meme msg
- nil stash
- bad command
- Meme: normal
- Add: too few args, normal
- Remove: too few args, bad atoi, oob, normal
- List: normal
*/
func TestInterceptor(t *testing.T) {
	var s Stash
	var ds = DefaultStash("meme_stash")

	tests := []struct {
		name        string
		stash       *Stash
		in          *gb.Message
		wantErr     bool
		wantDiscErr bool
		wantEmbed   bool
	}{
		// General errors
		{name: "not-meme", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "nil-stash", stash: &s, in: mock.NewMessage(gb.Meme), wantErr: true, wantDiscErr: false, wantEmbed: false},
		{name: "bad-arg", stash: &ds, in: mock.NewMessage(gb.Meme, mock.WithArgs("not", "valid", "args")), wantErr: true, wantDiscErr: false, wantEmbed: true}, // todo(ee): this test passes whatever wantErr val is

		// Meme
		{name: "normal", stash: &ds, in: mock.NewMessage(gb.Meme), wantErr: false, wantDiscErr: false, wantEmbed: true},

		// Add
		{name: "too-few-args", stash: &ds, in: mock.NewMessage(gb.Meme, mock.WithArgs("add")), wantErr: true, wantDiscErr: false, wantEmbed: true},
		{name: "", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},

		// Remove
		{name: "", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},
		{name: "", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},

		// List
		{name: "", stash: &ds, in: mock.NewMessage(gb.None), wantErr: false, wantDiscErr: false, wantEmbed: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := Interceptor(test.stash)
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

/*
Test Cases:
- Save: normal, bad path, nil stash
*/
func TestSave(t *testing.T) {
	s := DefaultStash("meme_test")

	tests := []struct {
		name    string
		path    string
		stash   *Stash
		wantErr bool
	}{
		{name: "normal", path: "meme_test", stash: &s, wantErr: false},
		{name: "bad-path", path: "not-a-dir", stash: &s, wantErr: true},
		{name: "nil-stash", path: "meme_test", stash: nil, wantErr: true},
	}

	for _, test := range tests {
		_ = os.Mkdir("meme_test", 0755)
		defer os.RemoveAll("meme_test")

		t.Run(test.name, func(t *testing.T) {
			err := test.stash.Save(test.path)
			if (err != nil) != test.wantErr {
				t.Errorf("(err != nil) != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}
		})
	}
}

/*
Test Cases:
- Load: normal, bad path, nil stash
*/
func TestLoad(t *testing.T) {
	s := DefaultStash("meme_test")

	tests := []struct {
		name      string
		path      string
		stash     *Stash
		wantStash *Stash
		wantErr   bool
	}{
		{name: "normal", path: "meme_test", stash: &Stash{}, wantStash: &s, wantErr: false},
		{name: "bad-path", path: "not-a-dir", stash: &Stash{}, wantStash: nil, wantErr: true},
		{name: "nil-stash", path: "meme_test", stash: nil, wantStash: nil, wantErr: true},
	}

	for _, test := range tests {
		// create test dir
		_ = os.Mkdir("meme_test", 0755)
		defer os.RemoveAll("meme_test")

		// save the stash to a file
		_ = s.Save("meme_test")

		// run test
		t.Run(test.name, func(t *testing.T) {
			err := test.stash.Load(test.path)
			if (err != nil) != test.wantErr {
				t.Errorf("(err != nil) != wantErr (err = %v, wantErr = %v)", err, test.wantErr)
			}

			if err == nil && !test.wantErr {
				if !cmp.Equal(*test.stash, *test.wantStash) {
					t.Errorf("incorrect stash after load:\n%s", cmp.Diff(*test.stash, *test.wantStash))
				}
			}
		})
	}
}
