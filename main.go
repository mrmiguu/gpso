package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/mrmiguu/gpso/src/ipn"
	"github.com/mrmiguu/gpso/src/plr"

	"github.com/mrmiguu/gpso/src/ex"

	"github.com/mrmiguu/gpso/src/zone"

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
	authc := sock.Rbytes()
	errc := sock.Werror()

	go saveAccounts(&accts)

	for authb := range authc {
		var userPass [2]string
		btov(authb, &userPass)
		user := userPass[0]
		stats, err := accts.Get(user, []byte(userPass[1]))
		if err != nil {
			errc <- err
			continue
		}
		go play(user, stats)
	}
}

func saveAccounts(accts *plr.Accounts) {
	for range time.Tick(1 * time.Minute) {
		accts.Save()
	}
}

func play(user string, stats *plr.Stats) {
	statsh := head(user, "stats")
	statsc := sock.Wbytes(statsh)
	defer sock.Close(statsh)
	nodesh := head(user, "nodes")
	nodesc := sock.Wbytes(nodesh)
	defer sock.Close(nodesh)

	println("[Stats] sending")
	statsc <- stats.Bytes()
	println("[Stats] sent")

	nodesb, err := vtob(zone.Nodes)
	if err != nil {
		println(err)
		return
	}
	println("[Nodes] sending")
	nodesc <- nodesb
	println("[Nodes] sent")
}
