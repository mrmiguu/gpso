package main

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/mrmiguu/jsutil"

	"github.com/mrmiguu/gpso/node"
	"github.com/mrmiguu/gpso/src/ex"

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

	vtob = json.Marshal
	btov = json.Unmarshal
	must = ex.Must
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
				jsutil.Panic(err)
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
			go func() {
				user := loginInput.Value()

				passhash := sha256.Sum256([]byte(user))
				userPass := [2]string{
					user,
					string(passhash[:]),
				}
				authb, err := vtob(userPass)
				must(err)
				println("[UserPass] sending")
				authc <- authb
				println("[UserPass] sent")

				// statsh := head(user, "stats")
				// statsc := sock.Rbytes()
				// defer sock.Close(statsh)
				nodesh := head(user, "nodes")
				nodesc := sock.Rbytes()
				defer sock.Close(nodesh)

				// go func() {
				// 	for statsb := range statsc {
				// 		stats := plr.Stats{}
				// 		btov(statsb, &stats)
				// 		expBar.Width(css.Pct(stats.Exp))
				// 	}
				// }()

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

		case <-mapScreen.Link:
			allScreens.Display(css.None)
			mapScreen.Display(css.Grid)
		}
	}
}
