package ex

import (
	"bytes"
	"errors"
	"time"
)

var timeFmts = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

var (
	newerr = errors.New
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func Head(fields ...string) string {
	buf := bytes.Buffer{}
	for _, field := range fields {
		buf.WriteString(field + "\t")
	}
	return buf.String()
}

func Time(s string, fmts ...string) (time.Time, error) {
	fmts = append(fmts, timeFmts...)
	for _, fmt := range fmts {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, newerr("ex.Time: unable to parse time '" + s + "'")
}
