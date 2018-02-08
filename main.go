package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/mrmiguu/gpso/src/ipn"
	"github.com/mrmiguu/gpso/src/plr"

	"github.com/mrmiguu/gpso/src/ex"

	"github.com/mrmiguu/jsutil"
	"github.com/mrmiguu/sock"
)

var (
	compile = jsutil.CompileWithGzip
	println = fmt.Println
	must    = ex.Must
	head    = ex.Head
	newerr  = errors.New
	vtob    = json.Marshal
	btov    = json.Unmarshal
	itoa    = strconv.Itoa
	atoi    = strconv.Atoi
)

func main() {
	must(compile(sock.Root + "/script.go"))

	accts := plr.Accounts{}

	ipn.Accounts = &accts
	http.HandleFunc("/gpso_ipn", ipn.Handler)

	sock.Addr = ":80"
	plrc := sock.Wbytes()
	authc := sock.Rbytes()
	errc := sock.Werror()

	go saveAccounts(&accts)

	for authb := range authc {
		var userPass [2]string
		btov(authb, &userPass)
		user := userPass[0]
		plr, err := accts.Get(user, []byte(userPass[1]))
		if err != nil {
			errc <- err
			continue
		}
		go play(plrc, user, plr)
	}
}

func saveAccounts(accts *plr.Accounts) {
	for range time.Tick(1 * time.Minute) {
		accts.Save()
	}
}

func play(plrc chan<- []byte, user string, plr *plr.Player) {
	plr.AddExp(25)

	println("[" + user + "] sending")
	plrc <- plr.Bytes()
	println("[" + user + "] sent")

	jumph := head(user, "jump")
	jumpc := sock.Rbool(jumph)
	defer sock.Close(jumph)

	sideh := head(user, "side")
	sidec := sock.Wint(sideh)
	defer sock.Close(sideh)

	for range jumpc {
		sidec <- rand.Intn(6)
		time.Sleep(250)

		plr.Move()

		println("[" + user + "] sending")
		plrc <- plr.Bytes()
		println("[" + user + "] sent")
	}

	select {} // keep the player's sockets alive (for now)
}
