package notebook

import (
	"fmt"
	"time"
)

func (n *Notebook) Allowed(t time.Time) bool {
	if n.WeekStart > 0 {
		_, week := t.ISOWeek()
		if week < n.WeekStart {
			return false
		}
	}

	if len(n.Weekdays) > 0 {
		var weekday = t.Weekday()
		var found bool
		for _, w := range n.Weekdays {
			if w == weekday {
				found = true
			}

		}
		if !found {
			return false
		}
	}

	return true
}

func (n *Notebook) AllowedDate(t time.Time) (time.Time, error) {
	if n.Allowed(t) {
		return t, nil
	}

	const day = time.Hour * 24

	// Go backwards to the beginning of the week
	for d := t; d.Weekday() >= 0; d = d.Add(-day) {
		if n.Allowed(d) {
			return d, nil
		}
	}

	// Try the next 7 days
	for i := 1; i <= 7; i++ {
		if d := t.Add(day * time.Duration(i)); n.Allowed(d) {
			return d, nil
		}
	}

	return t, fmt.Errorf("Could not find an allowed day")

}
