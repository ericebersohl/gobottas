package command

import "testing"

func TestToSnowflake(t *testing.T) {

	tests := []struct {
		name    string
		in      string
		wantErr bool
		want    Snowflake
	}{
		{name: "nil-value", in: "", wantErr: true, want: 0},
		{name: "not-snowflake", in: "notasnowflake", wantErr: true, want: 0},
		{name: "snowflake", in: "123456789", wantErr: false, want: 123456789},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := ToSnowflake(test.in)
			if (err != nil) != test.wantErr {
				t.Logf("wantErr != err (want = %v, err = %v)", test.wantErr, err)
				t.FailNow()
			}

			if s != test.want {
				t.Errorf("got the wrong snowflake (want = %d, got = %d)", test.want, s)
			}
		})
	}
}

func TestStrToCommandType(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want CommandType
	}{
		{name: "none", in: "none", want: None},
		{name: "help", in: "help", want: Help},
		{name: "meme", in: "meme", want: Meme},
		{name: "err", in: "error", want: Error},
		{name: "unrec", in: "blargh", want: Unrecognized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ct := StrToCommandType(test.in)
			if ct != test.want {
				t.Errorf("unexpected CommandType (in = %s, out = %s)", test.in, ct.String())
			}
		})
	}
}

func TestCommandType_String(t *testing.T) {
	tests := []struct {
		name string
		in   CommandType
		want string
	}{
		{name: "none", in: None, want: "None"},
		{name: "help", in: Help, want: "Help"},
		{name: "unrec", in: Unrecognized, want: "Unrecognized"},
		{name: "meme", in: Meme, want: "Meme"},
		{name: "error", in: Error, want: "Error"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.in.String()
			if s != test.want {
				t.Errorf("wrong string (in = %s, out = %s)", test.in.String(), s)
			}
		})
	}
}
