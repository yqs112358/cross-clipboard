package maputil

import (
	"reflect"
	"testing"

	"github.com/yqs112358/cross-clipboard/pkg/utils/pointerutil"
)

type TestStruct struct {
	TestString        string  `form:"string,limit=10"`
	TestPointerString *string `form:"pointerString,required"`
	TestInt           int     `form:"int"`
	TestPointerInt    *int    `form:"pointerInt"`
	TestDash          string  `form:"-"`
	TestEmpty         string  `form:""`
}

func TestToMapString(t *testing.T) {
	type args struct {
		in  interface{}
		tag string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "success struct",
			args: args{
				in: TestStruct{
					TestString:        "foo",
					TestPointerString: pointerutil.ToString("bar"),
					TestInt:           10,
					TestPointerInt:    pointerutil.ToInt(15),
					TestDash:          "foo",
					TestEmpty:         "bar",
				},
				tag: "form",
			},
			want: map[string]string{
				"string":        "foo",
				"pointerString": "bar",
				"int":           "10",
				"pointerInt":    "15",
			},
			wantErr: false,
		},
		{
			name: "success pointer struct",
			args: args{
				in: &TestStruct{
					TestString:        "foo",
					TestPointerString: pointerutil.ToString("bar"),
					TestInt:           10,
					TestPointerInt:    pointerutil.ToInt(15),
				},
				tag: "form",
			},
			want: map[string]string{
				"string":        "foo",
				"pointerString": "bar",
				"int":           "10",
				"pointerInt":    "15",
			},
			wantErr: false,
		},
		{
			name: "success without string",
			args: args{
				in: &TestStruct{
					TestPointerString: pointerutil.ToString("bar"),
					TestInt:           10,
					TestPointerInt:    pointerutil.ToInt(15),
				},
				tag: "form",
			},
			want: map[string]string{
				"string":        "",
				"pointerString": "bar",
				"int":           "10",
				"pointerInt":    "15",
			},
			wantErr: false,
		},
		{
			name: "success without pointer string",
			args: args{
				in: &TestStruct{
					TestString:        "foo",
					TestPointerString: nil,
					TestInt:           10,
					TestPointerInt:    pointerutil.ToInt(15),
				},
				tag: "form",
			},
			want: map[string]string{
				"string":     "foo",
				"int":        "10",
				"pointerInt": "15",
			},
			wantErr: false,
		},
		{
			name: "success without int",
			args: args{
				in: &TestStruct{
					TestString:        "foo",
					TestPointerString: pointerutil.ToString("bar"),
					TestPointerInt:    pointerutil.ToInt(15),
				},
				tag: "form",
			},
			want: map[string]string{
				"string":        "foo",
				"pointerString": "bar",
				"int":           "0",
				"pointerInt":    "15",
			},
			wantErr: false,
		},
		{
			name: "success without pointer int",
			args: args{
				in: &TestStruct{
					TestString:        "foo",
					TestPointerString: pointerutil.ToString("bar"),
					TestInt:           10,
					TestPointerInt:    nil,
				},
				tag: "form",
			},
			want: map[string]string{
				"string":        "foo",
				"pointerString": "bar",
				"int":           "10",
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ToMapString(test.args.in, test.args.tag)
			if (err != nil) != test.wantErr {
				t.Fatalf("ToMapString() error = %v, wantErr %v", err, test.wantErr)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ToMapString() got = %v, want %v", got, test.want)
			}
		})
	}
}
