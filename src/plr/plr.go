package plr

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	newerr = errors.New
	vtob   = json.Marshal
)

type Color int

const (
	Red Color = iota
	Orange
	Yellow
	Green
	Blue
	Indigo
	Violet
)

type Accounts struct {
	mux sync.Mutex
	m   map[string]*Account
}

// Get gets the stats of the player using a username and password.
func (a *Accounts) Get(user string, pass []byte) (*Stats, error) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if a.m == nil {
		a.m = make(map[string]*Account)
	}
	acct, found := a.m[user]
	if !found { // not found in memory, check file
		f, err := os.Open(userPath(user))
		defer f.Close()
		if err != nil { // file not found; set new account password
			println("[NewAccount] " + user)
			acct = &Account{Passhash: string(pass)}
		} else if err := json.NewDecoder(f).Decode(&acct); err != nil { // file found, but undecodable
			return nil, err
		}
		acct.mux = &a.mux
		acct.Online = false
		a.m[user] = acct
	}
	if acct.Online {
		return nil, newerr("single user per account only")
	}
	if acct.Passhash != string(pass) {
		return nil, newerr("bad password")
	}
	println("[LoggedIn] " + user)
	return &acct.Stats, nil
}

// AddGpsos adds gpsos into the user's account if existant.
func (a *Accounts) AddGpsos(user string) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	if a.m == nil {
		a.m = make(map[string]*Account)
	}
	acct, found := a.m[user]
	if !found { // not found in memory, check file
		f, err := os.Open(userPath(user))
		if err != nil {
			return newerr("account not found")
		} else if err := json.NewDecoder(f).Decode(&acct); err != nil { // file found, but undecodable
			f.Close()
			return err
		}
		f.Close()
		a.m[user] = acct
	}
	acct.Gpsos++
	acct.unsafeSave(user)
	if !found { // they weren't logged in
		delete(a.m, user)
	}
	return nil
}

func userPath(user string) string {
	return "accts/" + user + ".json"
}

func (a *Accounts) Save() {
	a.mux.Lock()
	defer a.mux.Unlock()
	for user, acct := range a.m {
		if err := acct.unsafeSave(user); err != nil {
			println(err)
		}
	}
}
func (a Account) Save(user string) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.unsafeSave(user)
}
func (a Account) unsafeSave(user string) error {
	println("[" + user + "] saving")
	f, err := os.Create(userPath(user))
	if err != nil {
		return err
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(a); err != nil {
		return err
	}
	println("[" + user + "] saved")
	return nil
}

func (a *Accounts) Atos() map[string]*Stats {
	a.mux.Lock()
	defer a.mux.Unlock()
	stats := make(map[string]*Stats, len(a.m))
	for user, acct := range a.m {
		stats[user] = &acct.Stats
	}
	return stats
}

type Account struct {
	Passhash string
	Online   bool
	Stats
}

type Stats struct {
	mux *sync.Mutex

	Exp   int
	Level int
	Gpsos int
	Color Color
}

func (s Stats) Bytes() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()
	b, _ := vtob(s)
	return b
}
