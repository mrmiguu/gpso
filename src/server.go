package src

import (
	"github.com/mrmiguu/jsutil"
	"github.com/mrmiguu/sock"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func init() {
	must(jsutil.CompileWithGzip("www/script.go"))

	sock.Addr = ":80"
	handshake := sock.Rstring()

	for msg := range handshake {
		println(msg)
	}
}
