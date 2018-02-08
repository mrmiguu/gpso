package main

import (
	"crypto/sha256"
	"encoding/json"
	"strconv"

	"github.com/mrmiguu/gpso/src/zone"

	"github.com/mrmiguu/gpso/src/plr"

	"github.com/mrmiguu/jsutil"

	"github.com/mrmiguu/gpso/src/ex"

	"github.com/mrmiguu/page"
	"github.com/mrmiguu/page/css"
	"github.com/mrmiguu/sock"
)

var (
	pct   = css.Pct
	atoi  = zone.Aton
	redir = jsutil.Redirect
	panic = jsutil.Panic
	itoa  = strconv.Itoa
	vtob  = json.Marshal
	btov  = json.Unmarshal
	head  = ex.Head

	allScreens   = page.Class("screen")
	splashScreen = page.Class("splash")
	loginScreen  = page.Class("login")
	loginInput   = page.Class("logininput")
	gameScreen   = page.Class("game")
	expBar       = page.Class("exp")
	dieElem      = page.Class("diebody")
	mapScreen    = page.Class("map")
	gpsosScreen  = page.Class("gpsos")
	g10          = page.Class("g10")
	g20          = page.Class("g20")
	g50          = page.Class("g50")
	g100         = page.Class("g100")
	g200         = page.Class("g200")
	g500         = page.Class("g500")
	g1000        = page.Class("g1000")
	tabElem      = page.Class("tab")

	carElems = map[plr.Color]page.Elem{
		plr.Red:    page.Class("rcar"),
		plr.Orange: page.Class("ocar"),
		plr.Yellow: page.Class("ycar"),
		plr.Green:  page.Class("gcar"),
		plr.Blue:   page.Class("bcar"),
		plr.Indigo: page.Class("icar"),
		plr.Violet: page.Class("vcar"),
	}

	user  string
	jumpc chan<- bool
	sidec <-chan int
	expc  <-chan int
	lvlc  <-chan int
)

func main() {
	sock.Addr = "goplaysmile.com"

	go syncPlrs(sock.Rbytes())
	authc := sock.Wbytes()
	errc := sock.Rerror()

	var cycling bool

	for {
		println(`reselecting`)
		select {
		case err := <-errc:
			if err != nil {
				panic(err)
			}

		case <-splashScreen.Link:
			allScreens.Display(css.None)
			splashScreen.Display(css.Grid)
			splashScreen.Animation("showsplash")

		case <-loginScreen.Link:
			allScreens.Display(css.None)
			loginScreen.Display(css.Grid)
			<-loginInput.Hit
			loginInput.Value("")

		case <-gameScreen.Link:
			if gpsosScreen.Display() == css.Grid {
				gpsosScreen.Animation("gpsosup")
			}
			allScreens.Display(css.None)
			gameScreen.Display(css.Grid)

			if cycling {
				break
			}
			cycling = true
			user = loginInput.Value()
			go func() {
				passhash := sha256.Sum256([]byte(user))
				userPass := [2]string{
					user,
					string(passhash[:]),
				}
				authb, err := vtob(userPass)
				must(err)

				jumpc = sock.Wbool(head(user, "jump"))
				sidec = sock.Rint(head(user, "side"))
				expc = sock.Rint(head(user, "exp"))
				lvlc = sock.Rint(head(user, "lvl"))

				println("[UserPass] sending")
				authc <- authb
				println("[UserPass] sent")
			}()

		case <-dieElem.Hit:
			dieElem.Animation("diejump")
			go dieElem.Animation("none")
			jumpc <- true
		case side := <-sidec:
			dieElem.BgImg("dice/white-" + itoa(side+1) + ".png")

		case <-gpsosScreen.Link:
			gpsosScreen.Display(css.Grid)
			gpsosScreen.Animation("gpsosdown")

		case <-g10.Hit:
			payPalRedir("3AGKVQVLS9WF2")
		case <-g20.Hit:
			payPalRedir("QAZLX4G4DRYXY")
		case <-g50.Hit:
			payPalRedir("T3MQB3N6G3ZPL")
		case <-g100.Hit:
			payPalRedir("8EN6JGPG65ELU")
		case <-g200.Hit:
			payPalRedir("KDXU5V9TUXUMU")
		case <-g500.Hit:
			payPalRedir("V79FQ6H9M2TW2")
		case <-g1000.Hit:
			payPalRedir("VG8SE9DGAGAFC")

		case <-mapScreen.Link:
			allScreens.Display(css.None)
			mapScreen.Display(css.Grid)
		}
	}
}

func syncPlrs(plrc <-chan []byte) {
	for plrb := range plrc {
		var plr plr.Player
		if err := btov(plrb, &plr); err != nil {
			println(err)
			continue
		}

		if plr.Name == user { // our user interface
			expBar.Width(pct(plr.Exp))
		}

		car := carElems[plr.Color]
		city, err := atoi(plr.City)
		if err != nil {
			println(err)
			continue
		}

		x, y := city.Pt[0]/zone.Width, city.Pt[1]/zone.Height
		go car.Move(pct(x*100), pct(y*100))
	}
}

func payPalRedir(id string) {
	if len(user) == 0 {
		return
	}
	redir("https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=" + id + "&custom=" + user)
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
