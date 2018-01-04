package main

import (
	"github.com/mrmiguu/page"
	"github.com/mrmiguu/sock"
)

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	splash, err := page.ID("splash")
	must(err)

	done, err := splash.Anim("splash-anim")
	must(err)
	<-done

	println("showing splash fully!")

	sock.Addr = "goplaysmile.com"
	handshake := sock.Wstring()
	handshake <- "*a handshake*"
}
