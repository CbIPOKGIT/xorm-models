package xormmodels

import (
	"database/sql"
	"encoding/json"
)

type SqlNullstring struct {
	sql.NullString
}

type SqlNullDate struct {
	sql.NullTime
}

type SqlNullFloat struct {
	sql.NullFloat64
}

type SqlNullInt struct {
	sql.NullInt32
}

type SqlNullBigInt struct {
	sql.NullInt64
}

func (s SqlNullstring) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

func (s *SqlNullstring) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != "" {
		s.Valid = true
		s.String = str
	} else {
		s.Valid = false
	}

	return nil
}

func (s SqlNullFloat) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Float64)
	}
	return []byte(`null`), nil
}

func (s *SqlNullFloat) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != 0 {
		s.Valid = true
		s.Float64 = f
	} else {
		s.Valid = false
	}

	return nil
}

func (s SqlNullInt) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Int32)
	}
	return []byte(`null`), nil
}

func (s *SqlNullInt) UnmarshalJSON(data []byte) error {
	var f int
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != 0 {
		s.Valid = true
		s.Int32 = int32(f)
	} else {
		s.Valid = false
	}

	return nil
}

func (s SqlNullBigInt) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Int64)
	}
	return []byte(`null`), nil
}

func (s *SqlNullBigInt) UnmarshalJSON(data []byte) error {
	var f int
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != 0 {
		s.Valid = true
		s.Int64 = int64(f)
	} else {
		s.Valid = false
	}

	return nil
}

func (s SqlNullDate) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Time)
	}
	return []byte(`null`), nil
}

func (nt *SqlNullDate) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		nt.Valid = false
		return err
	}

	nt.Valid = !nt.Time.IsZero()
	return nil
}
