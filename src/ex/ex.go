package ex

import (
	"bytes"
	"fmt"
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

func Head(fields ...interface{}) string {
	buf := bytes.Buffer{}
	for _, field := range fields {
		buf.WriteString(fmt.Sprintln(field) + "\t")
	}
	return buf.String()
}
