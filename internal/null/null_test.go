package null

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestNullStringMarshallJSON(t *testing.T) {
	type test struct {
		input NullString
		want  string
	}

	tests := []test{
		{input: NullStringFrom("string"), want: `"string"`},
		{input: NullStringFromPtr(nil), want: `null`},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got, err := tc.input.MarshalJSON()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %+v; want %+v", string(got), tc.want)
			}
		})
	}
}

func TestNullStringUnmarshallJSON(t *testing.T) {
	type test struct {
		input string
		want  NullString
	}

	tests := []test{
		{input: `null`, want: NullStringFromPtr(nil)},
		{input: `"string"`, want: NullStringFrom("string")},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := NullStringFromPtr(nil)
			err := got.UnmarshalJSON([]byte(tc.input))
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			equal := reflect.DeepEqual(got, tc.want)
			if !equal {
				t.Errorf("got %+v; want %+v", got, tc.want)
			}
		})
	}
}

func TestNullStringFrom(t *testing.T) {
	got := NullStringFrom("test")
	want := NullString{
		NullString: sql.NullString{
			Valid:  true,
			String: "test",
		},
		Dirty: true,
	}

	equal := reflect.DeepEqual(got, want)
	if !equal {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

func TestNullStringFromPtr(t *testing.T) {
	type test struct {
		input *string
		want  NullString
	}

	input := "string"
	tests := []test{
		{input: nil, want: NullString{
			NullString: sql.NullString{
				Valid:  false,
				String: "",
			},
			Dirty: true,
		}},
		{input: &input, want: NullString{
			NullString: sql.NullString{
				Valid:  true,
				String: input,
			},
			Dirty: true,
		}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := NullStringFromPtr(tc.input)
			equal := reflect.DeepEqual(got, tc.want)
			if !equal {
				t.Errorf("got %+v; want %+v", got, tc.want)
			}
		})
	}
}

func TestPtrFromNullString(t *testing.T) {
	type test struct {
		input NullString
		want  *string
	}

	input := "string"
	tests := []test{
		{input: NullStringFrom("string"), want: &input},
		{input: NullStringFromPtr(nil), want: nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := PtrFromNullString(tc.input)
			if (got == nil || tc.want == nil) && (got != nil || tc.want != nil) {
				t.Errorf("got %+v; want %+v", got, tc.want)
			} else if got != nil && tc.want != nil && *got != *tc.want {
				t.Errorf("got %+v; want %+v", *got, *tc.want)
			}
		})
	}
}

func TestNullInt64MarshallJSON(t *testing.T) {
	type test struct {
		input NullInt64
		want  string
	}

	tests := []test{
		{input: NullInt64From(1), want: `1`},
		{input: NullInt64FromPtr(nil), want: `null`},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got, err := tc.input.MarshalJSON()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %+v; want %+v", string(got), tc.want)
			}
		})
	}
}

func TestNullInt64UnmarshallJSON(t *testing.T) {
	type test struct {
		input string
		want  NullInt64
	}

	tests := []test{
		{input: `null`, want: NullInt64FromPtr(nil)},
		{input: `1`, want: NullInt64From(1)},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := NullInt64FromPtr(nil)
			err := got.UnmarshalJSON([]byte(tc.input))
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			equal := reflect.DeepEqual(got, tc.want)
			if !equal {
				t.Errorf("got %+v; want %+v", got, tc.want)
			}
		})
	}
}

func TestNullInt64From(t *testing.T) {
	got := NullInt64From(1)
	want := NullInt64{
		NullInt64: sql.NullInt64{
			Valid: true,
			Int64: 1,
		},
		Dirty: true,
	}

	equal := reflect.DeepEqual(got, want)
	if !equal {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

func TestNullInt64FromPtr(t *testing.T) {
	type test struct {
		input *int64
		want  NullInt64
	}

	var input int64 = 1
	tests := []test{
		{input: nil, want: NullInt64{
			NullInt64: sql.NullInt64{
				Valid: false,
				Int64: 0,
			},
			Dirty: true,
		}},
		{input: &input, want: NullInt64{
			NullInt64: sql.NullInt64{
				Valid: true,
				Int64: input,
			},
			Dirty: true,
		}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := NullInt64FromPtr(tc.input)
			equal := reflect.DeepEqual(got, tc.want)
			if !equal {
				t.Errorf("got %+v; want %+v", got, tc.want)
			}
		})
	}
}

func TestPtrFromNullInt64(t *testing.T) {
	type test struct {
		input NullInt64
		want  *int64
	}

	var input int64 = 1
	tests := []test{
		{input: NullInt64From(1), want: &input},
		{input: NullInt64FromPtr(nil), want: nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := PtrFromNullInt64(tc.input)
			if (got == nil || tc.want == nil) && (got != nil || tc.want != nil) {
				t.Errorf("got %+v; want %+v", got, tc.want)
			} else if got != nil && tc.want != nil && *got != *tc.want {
				t.Errorf("got %+v; want %+v", *got, *tc.want)
			}
		})
	}
}
