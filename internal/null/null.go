package null

import (
	"database/sql"
	"encoding/json"
)

// Armaria deals with databases and JSON a fair bit.
// Both of these things make use of NULL.
// The abstractions in this file allow us to support NULL robustly for both of them.
// The abstractions also keep track of whether a field was set at all which is useful for optional params.

// NullString is a string that can be NULL.
type NullString struct {
	sql.NullString      // allows use with sql/database and grants null support (null if Valid = false)
	Dirty          bool // tracks if the argument has been provided at all (provided if Dirty = true)
}

// MarshalJSON marshalls a NullString to JSON.
func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

// UnmarshalJSON unmarshalls a NullString from JSON.
func (s *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == `null` {
		s.Valid = false
		return nil
	}

	s.Valid = true
	return json.Unmarshal(data, &s.String)
}

// NullStringFrom converts a string to a NullString.
// The returned NullString will have Dirty set to true.
func NullStringFrom(str string) NullString {
	return NullString{
		NullString: sql.NullString{
			Valid:  true,
			String: str,
		},
		Dirty: true,
	}
}

// NullStringFromPtr converts a *string to a NullString.
// If the string is nil it will be treated as null.
// The returned NullString will have Dirty set to true.
func NullStringFromPtr(str *string) NullString {
	if str == nil {
		return NullString{
			NullString: sql.NullString{
				Valid:  false,
				String: "",
			},
			Dirty: true,
		}
	}

	return NullString{
		NullString: sql.NullString{
			Valid:  true,
			String: *str,
		},
		Dirty: true,
	}
}

// PtrFromNullString converts a NullString to a *string.
// If the NullString is not Valid nil is returned.
func PtrFromNullString(str NullString) *string {
	if str.Valid {
		return &str.String
	}

	return nil
}

// NullInt64 is an int64 that can be NULL.
type NullInt64 struct {
	sql.NullInt64      // allows use with sql/database and grants null support (null if Valid = false)
	Dirty         bool // tracks if the argument has been provided at all (provided if Dirty = true)
}

// MarshalJSON marshalls a NullInt64 to JSON.
func (s NullInt64) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Int64)
	}
	return []byte(`null`), nil
}

// UnmarshalJSON unmarshalls a NullInt64 from JSON.
func (s *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == `null` {
		s.Valid = false
		return nil
	}

	s.Valid = true
	return json.Unmarshal(data, &s.Int64)
}

// NullInt64From converts an int to a NullInt.
// If the int is zero it will be treated as null.
// The returned NullInt64 will have Dirty set to true.
func NullInt64From(num int64) NullInt64 {
	return NullInt64{
		NullInt64: sql.NullInt64{
			Valid: true,
			Int64: num,
		},
		Dirty: true,
	}
}

// NullInt64FromPtr converts a *int64 to a NullInt64.
// If the int64 is nil it will be treated as null.
// The returned NullInt64 will have Dirty set to true.
func NullInt64FromPtr(num *int64) NullInt64 {
	if num == nil {
		return NullInt64{
			NullInt64: sql.NullInt64{
				Valid: false,
				Int64: 0,
			},
			Dirty: true,
		}
	}

	return NullInt64{
		NullInt64: sql.NullInt64{
			Valid: true,
			Int64: *num,
		},
		Dirty: true,
	}
}

// PtrFromNullInt64 converts a NullInt64 to a *int64.
// If the NullInt64 is not Valid nil is returned.
func PtrFromNullInt64(num NullInt64) *int64 {
	if num.Valid {
		return &num.Int64
	}

	return nil
}
