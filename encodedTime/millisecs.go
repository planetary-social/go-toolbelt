package encodedTime

import (
	"encoding/json"
	"strconv"
	"time"
)

// Millisecs is used to get a time from a JSON number that represents a timestamp in milliseconds.
type Millisecs time.Time

func NewMillisecs(secs int64) Millisecs {
	return Millisecs(time.Unix(secs, 0))
}

func (t *Millisecs) UnmarshalJSON(in []byte) error {
	var milliseconds float64
	err := json.Unmarshal(in, &milliseconds)
	if err != nil {
		return err
	}
	*t = Millisecs(time.UnixMilli(int64(milliseconds)))
	return nil
}

func (t Millisecs) MarshalJSON() ([]byte, error) {
	milliseconds := time.Time(t).UnixMilli()
	return []byte(strconv.FormatInt(milliseconds, 10)), nil
}
