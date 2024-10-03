package stringutil

import (
	"reflect"
	"testing"
)

func TestLimitStringLen(t *testing.T) {
	tests := []struct {
		name string
		str  string
		len  int
		want string
	}{
		{
			name: "string length > limit",
			str:  "aaaaaaaaaaaaaa",
			len:  5,
			want: "aaaaa",
		},
		{
			name: "string length < limit",
			str:  "aaaaa",
			len:  10,
			want: "aaaaa",
		},
		{
			name: "string length == limit",
			str:  "aaaaa",
			len:  5,
			want: "aaaaa",
		},
		{
			name: "empty string",
			str:  "",
			len:  5,
			want: "",
		},
		{
			name: "limit 0",
			str:  "aaaaa",
			len:  0,
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := LimitStringLen(test.str, test.len)

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
