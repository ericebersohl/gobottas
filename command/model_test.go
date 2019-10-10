package command

import "testing"

func TestToSnowflake(t *testing.T) {

	tests := []struct{
		name string
		in string
		wantErr bool
		want Snowflake
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
	// nil
	// not a snowflake
	// snowflake
}
