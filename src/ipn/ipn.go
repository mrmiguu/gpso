package ipn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/mrmiguu/gpso/src/ex"
	"github.com/mrmiguu/gpso/src/plr"
)

var price2gpsos = map[float64]int{
	2.00:  10,
	4.00:  20,
	9.00:  50,
	15.99: 100,
	25.95: 200,
	49.99: 500,
	69.95: 1000,
}

var (
	// Accounts references player accounts.
	Accounts *plr.Accounts

	ipn IPN

	println = fmt.Println
	sprint  = fmt.Sprint
	itoa    = strconv.Itoa
	stof    = strconv.ParseFloat
	ftos    = strconv.FormatFloat
	must    = ex.Must
	time    = ex.Time
)

// Handler handles instant payment notifications, processing account changes.
// This uses the default instant payment notification instance.
func Handler(w http.ResponseWriter, r *http.Request) {
	ipn.Handler(w, r)
}

// IPN instant payment notification.
type IPN struct {
	once     sync.Once
	Accounts *plr.Accounts
}

func (ipn *IPN) init() {
	if ipn.Accounts == nil {
		ipn.Accounts = Accounts
	}
}

// Handler handles instant payment notifications, processing account changes.
func (ipn *IPN) Handler(w http.ResponseWriter, r *http.Request) {
	ipn.once.Do(ipn.init)

	println("ipn: /gpso_ipn")
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("No route for %v", r.Method), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println("ipn: " + err.Error())
		return
	}
	defer r.Body.Close()

	body = append([]byte(`cmd=_notify-validate&`), body...)

	resp, err := http.Post(
		"https://www.paypal.com/cgi-bin/webscr",
		r.Header.Get("Content-Type"),
		bytes.NewBuffer(body),
	)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("ipn: " + err.Error())
		return
	}
	defer resp.Body.Close()
	verifyStatus := string(b)

	println("ipn: " + verifyStatus)
	if verifyStatus != "VERIFIED" {
		return
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		println("ipn: " + err.Error())
		return
	}

	// official parsing
	// official parsing
	// official parsing

	users := vals["custom"]
	if len(users) == 0 {
		println("ipn: empty custom")
		return
	}
	user := users[0]

	grossStrs := vals["mc_gross"]
	if len(grossStrs) == 0 {
		println("ipn: empty mc_gross")
		return
	}
	grossStr := grossStrs[0]
	gross, err := stof(grossStr, 64)
	if err != nil {
		println("ipn: could not parse gross from '" + grossStr + "'")
		return
	}

	currencies := vals["mc_currency"]
	if len(currencies) == 0 {
		println("ipn: empty mc_currency")
		return
	}
	currency := currencies[0]

	dateStrs := vals["payment_date"]
	if len(dateStrs) == 0 {
		println("ipn: empty payment_date")
		return
	}
	dateStr := dateStrs[0]
	date, err := time(dateStr, `15:04:05 Jan 02, 2006 MST`)
	if err != nil {
		println("ipn: " + err.Error())
		return
	}

	println("ipn: payment of " + ftos(gross, 'f', 2, 64) + " " + currency + " received " + date.String())

	price := gross // TODO: conversion goes here

	gpsos, found := price2gpsos[price]
	if !found {
		println("ipn: price of " + ftos(price, 'f', 2, 64) + " " + currency + " not found")
		return
	}

	if err := ipn.Accounts.AddGpsos(user, gpsos); err != nil {
		println("ipn: " + err.Error())
		return
	}

	println("ipn: " + user + " bought " + itoa(gpsos) + " gpsos")
}
