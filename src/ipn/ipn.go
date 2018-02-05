package ipn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/mrmiguu/gpso/src/ex"
	"github.com/mrmiguu/gpso/src/plr"
)

var (
	// Accounts references player accounts.
	Accounts *plr.Accounts

	ipn IPN

	must   = ex.Must
	sprint = fmt.Sprint
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

	println("ipn: /gpso_ipn [!]")
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("No route for %v", r.Method), http.StatusNotFound)
		return
	}

	println("ipn: Write Status 200")
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

	if verifyStatus != "VERIFIED" {
		println("ipn: " + verifyStatus)
		return
	}
	println("ipn: " + verifyStatus)

	vals, err := url.ParseQuery(string(body))
	for k, v := range vals {
		println("ipn: " + k + "=" + sprint(v))
		// for _, user := range v {
		// 	ipn.Accounts.AddGpsos(user)
		// }
	}
}
