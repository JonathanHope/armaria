package validate

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/internal/null"
)

func TestValidateURL(t *testing.T) {
	type test struct {
		input null.NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: armaria.ErrURLTooShort},
		{input: nullStringOfLength(0, false), want: armaria.ErrURLTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(2049, false), want: armaria.ErrURLTooLong},
		{input: nullStringOfLength(2048, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := URL(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateName(t *testing.T) {
	type test struct {
		input null.NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: armaria.ErrNameTooShort},
		{input: nullStringOfLength(0, false), want: armaria.ErrNameTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(2049, false), want: armaria.ErrNameTooLong},
		{input: nullStringOfLength(2048, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := Name(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type test struct {
		input null.NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: nil},
		{input: nullStringOfLength(0, false), want: armaria.ErrDescriptionTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(4097, false), want: armaria.ErrDescriptionTooLong},
		{input: nullStringOfLength(4096, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := Description(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateTags(t *testing.T) {
	type test struct {
		input    []string
		existing []string
		want     error
	}

	twentyFourTags := make([]string, 0)
	for i := 0; i < 24; i++ {
		twentyFourTags = append(twentyFourTags, uuid.New().String())
	}

	twentyFiveTags := make([]string, 0)
	for i := 0; i < 25; i++ {
		twentyFiveTags = append(twentyFiveTags, uuid.New().String())
	}

	tests := []test{
		{input: []string{"x", "x"}, existing: []string{}, want: armaria.ErrDuplicateTag},
		{input: []string{"x"}, existing: []string{"x"}, want: armaria.ErrDuplicateTag},
		{input: []string{""}, existing: []string{}, want: armaria.ErrTagTooShort},
		{input: []string{"x"}, existing: []string{}, want: nil},
		{input: []string{stringOfLength("x", 129)}, existing: []string{}, want: armaria.ErrTagTooLong},
		{input: []string{stringOfLength("x", 128)}, existing: []string{}, want: nil},
		{input: []string{"?"}, existing: []string{}, want: armaria.ErrTagInvalidChar},
		{input: twentyFiveTags, existing: []string{}, want: armaria.ErrTooManyTags},
		{input: twentyFourTags, existing: []string{}, want: nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := Tags(tc.input, tc.existing)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateFirst(t *testing.T) {
	type test struct {
		input null.NullInt64
		want  error
	}

	tests := []test{
		{input: nullInt64(1, false), want: nil},
		{input: nullInt64(0, true), want: nil},
		{input: nullInt64(0, false), want: armaria.ErrFirstTooSmall},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.input.Int64), func(t *testing.T) {
			got := First(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateOrder(t *testing.T) {
	type test struct {
		input armaria.Order
		want  error
	}

	tests := []test{
		{input: armaria.OrderName, want: nil},
		{input: armaria.OrderModified, want: nil},
		{input: "", want: armaria.ErrInvalidOrder},
		{input: "Description", want: armaria.ErrInvalidOrder},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got := Order(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateDirection(t *testing.T) {
	type test struct {
		input armaria.Direction
		want  error
	}

	tests := []test{
		{input: armaria.DirectionAsc, want: nil},
		{input: armaria.DirectionDesc, want: nil},
		{input: "", want: armaria.ErrInvalidDirection},
		{input: "up", want: armaria.ErrInvalidDirection},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got := Direction(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateQuery(t *testing.T) {
	type test struct {
		input null.NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: nil},
		{input: nullStringOfLength(3, false), want: nil},
		{input: nullStringOfLength(2, false), want: armaria.ErrQueryTooShort},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := Query(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

// nullStringOfLength generates a string of a desired length
func stringOfLength(substr string, length int) string {
	var str = ""
	for i := 0; i < length; i++ {
		str = str + substr
	}

	return str
}

// nullStringOfLength generates a NullString of a desired length
func nullStringOfLength(length int, isNull bool) null.NullString {
	str := stringOfLength("x", length)

	return null.NullString{
		NullString: sql.NullString{
			Valid:  !isNull,
			String: str,
		},
		Dirty: true,
	}
}

// nullInt64 generates a NullInt64 of a desired value.
func nullInt64(num int64, isNull bool) null.NullInt64 {
	return null.NullInt64{
		NullInt64: sql.NullInt64{
			Valid: !isNull,
			Int64: num,
		},
		Dirty: true,
	}
}

// validateValidator validates the result of running a validator function.
func validateValidator(t *testing.T, want error, got error) {
	if want == nil && got == nil {
		return
	} else if errors.Is(got, want) {
		return
	} else {
		t.Errorf("got %+v; want %+v", got, want)
	}
}
