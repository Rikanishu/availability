package main

import (
	"fmt"
	"time"
)

type Availability struct {
	ObjectID   string
	AvailDates []Range
}

func (a *Availability) String() string {
	return fmt.Sprintf("Avail(%s) %v", a.ObjectID, a.AvailDates)
}

type Range struct {
	StartTS int64
	EndTS   int64
}

func (r Range) String() string {
	return fmt.Sprintf("(%v -> %v)", time.Unix(r.StartTS, 0), time.Unix(r.EndTS, 0))
}

func ConvertIntoAvailDates(rs []Range, gStartTS int64, gEndTS int64) (out []Range, err error) {
	if gEndTS <= gStartTS {
		return nil, fmt.Errorf("invalid range configuration")
	}
	out = make([]Range, 1)
	out[0].StartTS = gStartTS
	for _, r := range rs {
		startTS := r.StartTS
		endTS := r.EndTS
		if startTS < gStartTS {
			startTS = gStartTS
		}
		if endTS > gEndTS {
			endTS = gEndTS
		}
		if startTS == out[len(out)-1].StartTS {
			out[len(out)-1].StartTS = endTS + 1
			continue
		}
		out[len(out)-1].EndTS = startTS - 1
		out = append(out, Range{
			StartTS: endTS + 1,
		})
	}
	if out[len(out)-1].StartTS == gEndTS {
		out = out[:len(out)-1]
	} else {
		out[len(out)-1].EndTS = gEndTS - 1
	}
	return
}
