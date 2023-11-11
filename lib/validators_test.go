package lib

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestValidateURL(t *testing.T) {
	type test struct {
		input NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: ErrURLTooShort},
		{input: nullStringOfLength(0, false), want: ErrURLTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(2049, false), want: ErrURLTooLong},
		{input: nullStringOfLength(2048, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := validateURL(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateName(t *testing.T) {
	type test struct {
		input NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: ErrNameTooShort},
		{input: nullStringOfLength(0, false), want: ErrNameTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(2049, false), want: ErrNameTooLong},
		{input: nullStringOfLength(2048, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := validateName(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type test struct {
		input NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: nil},
		{input: nullStringOfLength(0, false), want: ErrDescriptionTooShort},
		{input: nullStringOfLength(1, false), want: nil},
		{input: nullStringOfLength(4097, false), want: ErrDescriptionTooLong},
		{input: nullStringOfLength(4096, false), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := validateDescription(tc.input)
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
		{input: []string{"x", "x"}, existing: []string{}, want: ErrDuplicateTag},
		{input: []string{"x"}, existing: []string{"x"}, want: ErrDuplicateTag},
		{input: []string{""}, existing: []string{}, want: ErrTagTooShort},
		{input: []string{"x"}, existing: []string{}, want: nil},
		{input: []string{stringOfLength("x", 129)}, existing: []string{}, want: ErrTagTooLong},
		{input: []string{stringOfLength("x", 128)}, existing: []string{}, want: nil},
		{input: []string{"?"}, existing: []string{}, want: ErrTagInvalidChar},
		{input: twentyFiveTags, existing: []string{}, want: ErrTooManyTags},
		{input: twentyFourTags, existing: []string{}, want: nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.input), func(t *testing.T) {
			got := validateTags(tc.input, tc.existing)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateFirst(t *testing.T) {
	type test struct {
		input NullInt64
		want  error
	}

	tests := []test{
		{input: nullInt64(1, false), want: nil},
		{input: nullInt64(0, true), want: nil},
		{input: nullInt64(0, false), want: ErrFirstTooSmall},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.input.Int64), func(t *testing.T) {
			got := validateFirst(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateOrder(t *testing.T) {
	type test struct {
		input Order
		want  error
	}

	tests := []test{
		{input: OrderName, want: nil},
		{input: OrderModified, want: nil},
		{input: "", want: ErrInvalidOrder},
		{input: "Description", want: ErrInvalidOrder},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got := validateOrder(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateDirection(t *testing.T) {
	type test struct {
		input Direction
		want  error
	}

	tests := []test{
		{input: DirectionAsc, want: nil},
		{input: DirectionDesc, want: nil},
		{input: "", want: ErrInvalidDirection},
		{input: "up", want: ErrInvalidDirection},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got := validateDirection(tc.input)
			validateValidator(t, tc.want, got)
		})
	}
}

func TestValidateQuery(t *testing.T) {
	type test struct {
		input NullString
		want  error
	}

	tests := []test{
		{input: nullStringOfLength(0, true), want: nil},
		{input: nullStringOfLength(3, false), want: nil},
		{input: nullStringOfLength(2, false), want: ErrQueryTooShort},
	}

	for _, tc := range tests {
		t.Run(tc.input.String, func(t *testing.T) {
			got := validateQuery(tc.input)
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
func nullStringOfLength(length int, isNull bool) NullString {
	str := stringOfLength("x", length)

	return NullString{
		NullString: sql.NullString{
			Valid:  !isNull,
			String: str,
		},
		Dirty: true,
	}
}

// nullInt64 generates a NullInt64 of a desired value.
func nullInt64(num int64, isNull bool) NullInt64 {
	return NullInt64{
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
