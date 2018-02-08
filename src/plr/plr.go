package plr

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"github.com/mrmiguu/gpso/src/zone"
)

var (
	newerr = errors.New
	vtob   = json.Marshal
	aton   = zone.Aton
	itoa   = strconv.Itoa
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
func (a *Accounts) Get(user string, pass []byte) (*Player, error) {
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
			acct = &Account{
				Passhash: string(pass),
				Player: Player{
					Name:  user,
					Level: 1,
				},
			}
			a.m[user] = acct
			if err := acct.unsafeSave(user); err != nil {
				return nil, err
			}
		} else if err := json.NewDecoder(f).Decode(&acct); err != nil { // file found, but undecodable
			return nil, err
		}
		acct.mux = &a.mux
		acct.Online = false
	}
	// TODO: add something that changes Online=true to false
	//       modify 'sock' to spring this event
	//
	// if acct.Online {
	// 	return nil, newerr("single user per account only")
	// }
	city, goal := zone.SrcDst()
	acct.Online = true
	acct.City = city.Name
	acct.Goal = goal.Name
	println("acct.City=" + acct.City)
	println("acct.Goal=" + acct.Goal)
	if acct.Passhash != string(pass) {
		return nil, newerr("bad password")
	}
	println("[LoggedIn] " + user)
	return &acct.Player, nil
}

// AddGpsos adds gpsos into the user's account if existant.
func (a *Accounts) AddGpsos(user string, gpsos int) error {
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
	acct.Gpsos += gpsos
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

func (a *Accounts) Atos() map[string]*Player {
	a.mux.Lock()
	defer a.mux.Unlock()
	stats := make(map[string]*Player, len(a.m))
	for user, acct := range a.m {
		stats[user] = &acct.Player
	}
	return stats
}

type Account struct {
	Passhash string
	Online   bool
	Player
}

type Player struct {
	mux *sync.Mutex

	Name string
	Exp,
	Level,
	Gpsos int
	Color Color
	City,
	Goal string
}

func (p *Player) AddExp(exp int) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.Exp += exp
}

func (p *Player) unsafeScramble() {
	city, goal := zone.SrcDst()
	p.City, p.Goal = city.Name, goal.Name
}

func (p *Player) Move() {
	p.mux.Lock()
	defer p.mux.Unlock()
	city, err := aton(p.City)
	if err != nil {
		println("plr.Move: " + err.Error())
		return
	}
	goal, err := aton(p.Goal)
	if err != nil {
		println("plr.Move: " + err.Error())
		return
	}
	if p.City == p.Goal {
		p.unsafeScramble()
		return
	}
	path, err := zone.Find(city, goal)
	if err != nil {
		println("plr.Move: " + err.Error())
		return
	}
	for i, n := range path {
		println("path-" + itoa(i) + ": " + n.Name)
	}
	if len(path) < 2 {
		println("plr.Move: no more steps in path")
		return
	}
	p.City = path[1].Name
}

func (p Player) Bytes() []byte {
	p.mux.Lock()
	defer p.mux.Unlock()
	b, _ := vtob(p)
	return b
}
