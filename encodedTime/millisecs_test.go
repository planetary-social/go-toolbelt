package encodedTime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestMillisecsUnmarshall(t *testing.T) {
	var v Millisecs

	err := json.Unmarshal([]byte(`1449808143436`), &v)
	if err != nil {
		t.Fatal(err)
	}

	if n := time.Time(v).Sub(time.UnixMilli(1449808143436)); n != 0 {
		t.Fatal(fmt.Errorf("times not equal: %v", n))
	}
}

func TestFloats(t *testing.T) {
	var v Millisecs

	err := json.Unmarshal([]byte(`1553708494043.0059`), &v)
	if err != nil {
		t.Fatal(err)
	}

	if n := time.Time(v).Sub(time.UnixMilli(1553708494043)); n != 0 {
		t.Fatal(fmt.Errorf("times not equal: %v", n))
	}
}

func TestMillisecsMarshal(t *testing.T) {
	v := struct {
		Date Millisecs
	}{Millisecs(time.Unix(12345, 0))}

	out, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(out, []byte(`{"Date":12345000}`)) != 0 {
		t.Fatal(fmt.Errorf("times not equal - got %q", out))
	}

}
