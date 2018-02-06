package main

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/mrmiguu/jsutil"

	"github.com/mrmiguu/gpso/src/ex"
	"github.com/mrmiguu/gpso/src/node"

	"github.com/mrmiguu/page"
	"github.com/mrmiguu/page/css"
	"github.com/mrmiguu/sock"
)

var (
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

	carElems = []page.Elem{
		page.Class("rcar"),
		page.Class("ocar"),
		page.Class("ycar"),
		page.Class("gcar"),
		page.Class("bcar"),
		page.Class("icar"),
		page.Class("vcar"),
	}

	user string

	vtob  = json.Marshal
	btov  = json.Unmarshal
	panic = jsutil.Panic
	// must  = jsutil.Must
	head = ex.Head
)

func main() {
	sock.Addr = "goplaysmile.com"
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

				// statsh := head(user, "stats")
				// statsc := sock.Rbytes(statsh)
				// defer sock.Close(statsh)
				nodesh := head(user, "nodes")
				nodesc := sock.Rbytes(nodesh)
				defer sock.Close(nodesh)

				println("[UserPass] sending")
				authc <- authb
				println("[UserPass] sent")

				println("[Nodes] receiving")
				nodesb := <-nodesc
				nodes := []node.T{}
				must(btov(nodesb, &nodes))
				println("[Nodes] " + string(nodesb))

				nodei, cari := 0, 0
				for {
					car := carElems[cari]
					x, y := float64(nodes[nodei].Pt[0])/1454, float64(nodes[nodei].Pt[1])/1210

					go car.Move(css.Pct(x*100), css.Pct(y*100))
					time.Sleep(1 * time.Second)

					// cycle through the pts
					nodei = (nodei + 1) % len(nodes)
					// cycle through the cars
					cari = (cari + 1) % len(carElems)
				}
			}()

		case <-dieElem.Hit:
			dieElem.Animation("diejump")
			go dieElem.Animation("none")

		case <-gpsosScreen.Link:
			gpsosScreen.Display(css.Grid)
			gpsosScreen.Animation("gpsosdown")

		case <-g10.Hit:
			payPalRedirect("3AGKVQVLS9WF2")
		case <-g20.Hit:
			payPalRedirect("QAZLX4G4DRYXY")
		case <-g50.Hit:
			payPalRedirect("T3MQB3N6G3ZPL")
		case <-g100.Hit:
			payPalRedirect("8EN6JGPG65ELU")
		case <-g200.Hit:
			payPalRedirect("KDXU5V9TUXUMU")
		case <-g500.Hit:
			payPalRedirect("V79FQ6H9M2TW2")
		case <-g1000.Hit:
			payPalRedirect("VG8SE9DGAGAFC")

		case <-mapScreen.Link:
			allScreens.Display(css.None)
			mapScreen.Display(css.Grid)
		}
	}
}

func payPalRedirect(id string) {
	jsutil.Redirect("https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=" + id + "&custom=" + user)
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
