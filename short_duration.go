package main

import (
	"fmt"
	"time"
)

// we use the arbitrary boundaries:
// 360d, 48 hrs, 90m
const (
	day_boundary    = 360
	hour_boundary   = 48
	minute_boundary = 90
	second_boundary = 90
)

// The default Duration.String() implementation can be wide (w/ respect to terminal widths)
// ex: 5.004587543s
// ex: 2540400h10m10.000000000s
// This method attemps to shorten that:
// ex: 12.3y
// ex: 5.51d
// ex: 23.5h
// ex: 3.53h
// ex: 5.42s
// ex: 542ms
// ex: 542µs
// ex: 542ns
func ShortString(d time.Duration) string {
	u := uint64(d)
	neg := d < 0
	if neg {
		u = uint64(-1 * d)
	}
	f := float64(u)
	tmp := "x"

	if u < uint64(time.Second) {
		switch {
		case u > 100*uint64(time.Millisecond):
			tmp = fmt.Sprintf("%dms", u/uint64(time.Millisecond))
		case u > 10*uint64(time.Millisecond):
			tmp = fmt.Sprintf("%.01fms", f/float64(time.Millisecond))
		case u > uint64(time.Millisecond):
			tmp = fmt.Sprintf("%.02fms", f/float64(time.Millisecond))
		case u > 100*uint64(time.Microsecond):
			tmp = fmt.Sprintf("%dµs", u/uint64(time.Microsecond))
		case u > 10*uint64(time.Microsecond):
			tmp = fmt.Sprintf("%.01fµs", f/float64(time.Microsecond))
		case u > uint64(time.Microsecond):
			tmp = fmt.Sprintf("%.02fµs", f/float64(time.Microsecond))
		case u == 0:
			return "0"
		default:
			tmp = fmt.Sprintf("%dns", u/uint64(time.Nanosecond))
		}
	} else {
		switch {
		case u < uint64(10*time.Second):
			tmp = fmt.Sprintf("%0.02fs", f/float64(time.Second))
		case u < uint64(second_boundary*time.Second):
			tmp = fmt.Sprintf("%0.01fs", f/float64(time.Second))
		case u < uint64(minute_boundary*time.Minute):
			tmp = fmt.Sprintf("%dm%ds", int(time.Duration(u).Minutes()), (u%(60*uint64(time.Second)))/uint64(time.Second))
		case u < uint64(hour_boundary*time.Hour):
			tmp = fmt.Sprintf("%dh%dm", int(time.Duration(u).Hours()), (u%(60*uint64(time.Minute)))/uint64(time.Minute))
		case u < uint64(day_boundary*time.Hour*24):
			tmp = fmt.Sprintf("%0.01fd", float64(time.Duration(u).Hours())/24)
		default:
			tmp = time.Duration(u).String()
		}
	}

	if neg {
		tmp = "-" + tmp
	}

	return tmp
}
