package types

import "time"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	t.Time, err = time.Parse("\"2006-01-02T15:04:05\"", string(b))
	return err
}
