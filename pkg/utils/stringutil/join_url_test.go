package stringutil

import (
	"reflect"
	"testing"
)

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name  string
		base  string
		paths []string
		want  string
	}{
		{
			name:  "join url",
			base:  "http://example.com",
			paths: []string{"foo", "bar"},
			want:  "http://example.com/foo/bar",
		},
		{
			name:  "join url with `/`",
			base:  "http://example.com/",
			paths: []string{"foo/", "/bar"},
			want:  "http://example.com/foo/bar",
		},
		{
			name:  "join url with `//`",
			base:  "http://example.com/foo//",
			paths: []string{"//bar"},
			want:  "http://example.com/foo/bar",
		},
		{
			name:  "join url with extension",
			base:  "http://example.com/foo",
			paths: []string{"bar.go"},
			want:  "http://example.com/foo/bar.go",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := JoinURL(test.base, test.paths...)

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}
