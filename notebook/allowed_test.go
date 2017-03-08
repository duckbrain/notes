package notebook

import (
	"testing"
	"time"
)

func TestAllowed(t *testing.T) {
	n := Notebook{
		WeekStart: 0,
		Weekdays:  []time.Weekday{0},
	}

	if !n.Allowed(time.Date(2017, 2, 26, 0, 0, 0, 0, time.Local)) {
		t.Error("Should be allowed")
	}
}
