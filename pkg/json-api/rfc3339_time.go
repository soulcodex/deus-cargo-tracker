package xjsonapi

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// RFC3339Time serializes to/from RFC3339 strings.
type RFC3339Time time.Time

// Time returns the underlying time.Time.
func (t *RFC3339Time) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	return time.Time(*t)
}

// MarshalJSON implements json.Marshaler using RFC3339 (UTC).
func (t *RFC3339Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}

	tt := time.Time(*t)
	if tt.IsZero() {
		return []byte("null"), nil
	}

	utcTime := tt.Format(time.RFC3339)
	var buf bytes.Buffer
	buf.WriteByte('"')
	buf.WriteString(utcTime)
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler accepting RFC3339 or RFC3339Nano.
func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = RFC3339Time(time.Time{})
		return nil
	}

	unquotedTime := strings.Trim(string(data), `"`)
	if unquotedTime == "" {
		*t = RFC3339Time(time.Time{})
		return nil
	}

	tt, err := time.Parse(time.RFC3339, unquotedTime)
	if err != nil {
		tt, err = time.Parse(time.RFC3339Nano, unquotedTime)
		if err != nil {
			return fmt.Errorf("invalid RFC3339 time format provided: %w", err)
		}
	}

	*t = RFC3339Time(tt)
	return nil
}
