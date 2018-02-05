package zone

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/mrmiguu/gpso/src/ex"
	"github.com/mrmiguu/sock"
)

var (
	hmax float64

	println = fmt.Println
	sprint  = fmt.Sprint
	newerr  = errors.New
	itoa    = strconv.Itoa
	must    = ex.Must
	abs     = ex.Abs
)

var (
	SantaClarita    = Node{Name: "santaclarita", Hwys: []int{5, 210}}
	SanFernando     = Node{Name: "sanfernando", Hwys: []int{5, 405}}
	LakeviewTerrace = Node{Name: "lakeviewterrace", Hwys: []int{210}}
	Sunland         = Node{Name: "sunland", Hwys: []int{210}}
	PanoramaCity    = Node{Name: "panoramacity", Hwys: []int{5, 170}}
	LaCanada        = Node{Name: "lacanada", Hwys: []int{2, 210}}
	Pasadena        = Node{Name: "pasadena", Hwys: []int{110}}
	VanNuys         = Node{Name: "vannuys", Hwys: []int{101, 405}}
	Burbank         = Node{Name: "burbank", Hwys: []int{5, 101}}
	StudioCity      = Node{Name: "studiocity", Hwys: []int{101, 170}}
	Glendale        = Node{Name: "glendale", Hwys: []int{2, 101}}
	Azusa           = Node{Name: "azusa", Hwys: []int{210, 605}}
	SanDimas        = Node{Name: "sandimas", Hwys: []int{57, 210}}
	Claremont       = Node{Name: "claremont", Hwys: []int{210}}
	HighlandPark    = Node{Name: "highlandpark", Hwys: []int{110}}
	LosAngeles      = Node{Name: "losangeles", Hwys: []int{5, 10, 110}}
	BeverlyHills    = Node{Name: "beverlyhills", Hwys: []int{2, 170}}
	Rosemead        = Node{Name: "rosemead", Hwys: []int{10}}
	WestCovina      = Node{Name: "westcovina", Hwys: []int{10}}
	Pomona          = Node{Name: "pomona", Hwys: []int{10}}
	ElMonte         = Node{Name: "elmonte", Hwys: []int{10, 605}}
	MontereyPark    = Node{Name: "montereypark", Hwys: []int{10, 710}}
	CalPoly         = Node{Name: "calpoly", Hwys: []int{57}}
	USC             = Node{Name: "usc", Hwys: []int{10, 60, 110}}
	Commerce        = Node{Name: "commerce", Hwys: []int{60, 710}}
	SouthElMonte    = Node{Name: "southelmonte", Hwys: []int{60, 605}}
	CulverCity      = Node{Name: "culvercity", Hwys: []int{10, 405}}
	ChinoHills      = Node{Name: "chinohills", Hwys: []int{60}}
	SantaMonica     = Node{Name: "santamonica", Hwys: []int{10}}
	DiamondBar      = Node{Name: "diamondbar", Hwys: []int{57}}
	CityOfIndustry  = Node{Name: "cityofindustry", Hwys: []int{60}}
	SantaFeSprings  = Node{Name: "santafesprings", Hwys: []int{5, 605}}
	LAX             = Node{Name: "lax", Hwys: []int{105, 405}}
	Watts           = Node{Name: "watts", Hwys: []int{105, 110}}
	Corona          = Node{Name: "corona", Hwys: []int{91}}
	Brea            = Node{Name: "brea", Hwys: []int{57}}
	Lynwood         = Node{Name: "lynwood", Hwys: []int{105, 710}}
	Compton         = Node{Name: "compton", Hwys: []int{91, 710}}
	Bellflower      = Node{Name: "bellflower", Hwys: []int{91, 605}}
	Hawthorne       = Node{Name: "hawthorne", Hwys: []int{91, 405}}
	Gardena         = Node{Name: "gardena", Hwys: []int{91, 110}}
	DominguezHills  = Node{Name: "dominguezhills", Hwys: []int{110}}
	NorthLongBeach  = Node{Name: "northlongbeach", Hwys: []int{710}}
	Torrance        = Node{Name: "torrance", Hwys: []int{405}}
	BuenaPark       = Node{Name: "buenapark", Hwys: []int{5, 91}}
	Carson          = Node{Name: "carson", Hwys: []int{110, 405}}
	Fullerton       = Node{Name: "fullerton", Hwys: []int{57, 91}}
	YorbaLinda      = Node{Name: "yorbalinda", Hwys: []int{55, 91}}
	Anaheim         = Node{Name: "anaheim", Hwys: []int{5}}
	Cypress         = Node{Name: "cypress", Hwys: []int{605}}
	Westminster     = Node{Name: "westminster", Hwys: []int{22, 405}}
	PalosVerdes     = Node{Name: "palosverdes", Hwys: []int{}}
	GardenGrove     = Node{Name: "gardengrove", Hwys: []int{5, 22, 57}}
	Orange          = Node{Name: "orange", Hwys: []int{22, 55}}
	LongBeach       = Node{Name: "longbeach", Hwys: []int{710}}
	SanPedro        = Node{Name: "sanpedro", Hwys: []int{110}}
	HuntingtonBeach = Node{Name: "huntingtonbeach", Hwys: []int{605}}
	SantaAna        = Node{Name: "santaana", Hwys: []int{5, 55, 405}}
	FountainValley  = Node{Name: "fountainvalley", Hwys: []int{405}}
)

var Nodes []Node

var name2node = func() map[string]Node {
	m := make(map[string]Node, len(Nodes))
	for _, n := range Nodes {
		m[n.Name] = n
	}
	return m
}()

func Aton(s string) (Node, error) {
	if n, found := name2node[s]; found {
		return n, nil
	}
	return Node{}, newerr("'" + s + "' not found in zone")
}

var srcDst = [][2]string{
	{Claremont.Name, SantaAna.Name},
	{SanFernando.Name, YorbaLinda.Name},
	{Carson.Name, Pomona.Name},
	{Azusa.Name, Torrance.Name},
	{Claremont.Name, Westminster.Name},
	{ChinoHills.Name, VanNuys.Name},
	{CulverCity.Name, Orange.Name},
	{SanPedro.Name, LaCanada.Name},
	{Burbank.Name, Corona.Name},
	{Corona.Name, PanoramaCity.Name},
	{FountainValley.Name, StudioCity.Name},
	{ElMonte.Name, SanPedro.Name},
}

func SrcDst() (src, dst Node) {
	sd := srcDst[rand.Intn(len(srcDst))]
	si := rand.Intn(len(sd))
	src, _ = Aton(sd[si])
	dst, _ = Aton(sd[(si+1)%len(sd)])
	return
}

func init() {
	SantaClarita.Near = []string{SanFernando.Name, LakeviewTerrace.Name}
	SanFernando.Near = []string{SantaClarita.Name, PanoramaCity.Name, VanNuys.Name}
	LakeviewTerrace.Near = []string{SantaClarita.Name, Sunland.Name}
	Sunland.Near = []string{LaCanada.Name, LakeviewTerrace.Name}
	PanoramaCity.Near = []string{SanFernando.Name, StudioCity.Name, Burbank.Name}
	LaCanada.Near = []string{Sunland.Name, Glendale.Name, HighlandPark.Name}
	Pasadena.Near = []string{HighlandPark.Name}
	VanNuys.Near = []string{SanFernando.Name, StudioCity.Name, CulverCity.Name}
	Burbank.Near = []string{PanoramaCity.Name, StudioCity.Name, Glendale.Name, LosAngeles.Name, BeverlyHills.Name}
	StudioCity.Near = []string{VanNuys.Name, PanoramaCity.Name, Burbank.Name, BeverlyHills.Name}
	Glendale.Near = []string{Azusa.Name, Pasadena.Name, LaCanada.Name, Burbank.Name, LosAngeles.Name, HighlandPark.Name, BeverlyHills.Name}
	Azusa.Near = []string{SanDimas.Name, HighlandPark.Name, ElMonte.Name}
	SanDimas.Near = []string{Azusa.Name, Claremont.Name, WestCovina.Name, Pomona.Name, CalPoly.Name}
	Claremont.Near = []string{SanDimas.Name}
	HighlandPark.Near = []string{Azusa.Name, LosAngeles.Name, Glendale.Name, LaCanada.Name}
	LosAngeles.Near = []string{Burbank.Name, BeverlyHills.Name, USC.Name, HighlandPark.Name, Glendale.Name, MontereyPark.Name, Commerce.Name, SantaFeSprings.Name, Lynwood.Name}
	BeverlyHills.Near = []string{USC.Name, LosAngeles.Name, StudioCity.Name, Burbank.Name, Glendale.Name, Commerce.Name, SantaFeSprings.Name, Lynwood.Name}
	Rosemead.Near = []string{MontereyPark.Name, ElMonte.Name}
	WestCovina.Near = []string{ElMonte.Name, SanDimas.Name, Pomona.Name, CalPoly.Name}
	Pomona.Near = []string{WestCovina.Name, CalPoly.Name, SanDimas.Name}
	ElMonte.Near = []string{Rosemead.Name, WestCovina.Name, Azusa.Name, SouthElMonte.Name}
	MontereyPark.Near = []string{LosAngeles.Name, Commerce.Name, Rosemead.Name}
	CalPoly.Near = []string{WestCovina.Name, DiamondBar.Name, Pomona.Name, SanDimas.Name}
	USC.Near = []string{Commerce.Name, LosAngeles.Name, BeverlyHills.Name, Watts.Name, CulverCity.Name, Lynwood.Name, SantaFeSprings.Name}
	Commerce.Near = []string{MontereyPark.Name, SouthElMonte.Name, USC.Name, Lynwood.Name, SantaFeSprings.Name, LosAngeles.Name, BeverlyHills.Name}
	SouthElMonte.Near = []string{Commerce.Name, SantaFeSprings.Name, CityOfIndustry.Name, ElMonte.Name}
	CulverCity.Near = []string{SantaMonica.Name, USC.Name, VanNuys.Name, LAX.Name}
	ChinoHills.Near = []string{DiamondBar.Name}
	SantaMonica.Near = []string{CulverCity.Name}
	DiamondBar.Near = []string{CalPoly.Name, ChinoHills.Name, CityOfIndustry.Name, Brea.Name}
	CityOfIndustry.Near = []string{SouthElMonte.Name, DiamondBar.Name}
	SantaFeSprings.Near = []string{SouthElMonte.Name, BuenaPark.Name, Bellflower.Name, Commerce.Name, LosAngeles.Name, BeverlyHills.Name, USC.Name, Lynwood.Name}
	LAX.Near = []string{CulverCity.Name, Watts.Name, Hawthorne.Name}
	Watts.Near = []string{USC.Name, LAX.Name, Lynwood.Name, Gardena.Name}
	Corona.Near = []string{YorbaLinda.Name}
	Brea.Near = []string{DiamondBar.Name, Fullerton.Name}
	Lynwood.Near = []string{Watts.Name, Compton.Name, Commerce.Name, SantaFeSprings.Name, USC.Name, BeverlyHills.Name, LosAngeles.Name}
	Compton.Near = []string{Lynwood.Name, Gardena.Name, NorthLongBeach.Name, Bellflower.Name}
	Bellflower.Near = []string{Compton.Name, SantaFeSprings.Name, BuenaPark.Name, Cypress.Name}
	Hawthorne.Near = []string{LAX.Name, Gardena.Name, Torrance.Name}
	Gardena.Near = []string{Watts.Name, Compton.Name, Hawthorne.Name, DominguezHills.Name}
	DominguezHills.Near = []string{Gardena.Name, Carson.Name}
	NorthLongBeach.Near = []string{Carson.Name, Compton.Name, Westminster.Name, LongBeach.Name}
	Torrance.Near = []string{Hawthorne.Name, Carson.Name, PalosVerdes.Name}
	BuenaPark.Near = []string{Bellflower.Name, SantaFeSprings.Name, Anaheim.Name, Fullerton.Name}
	Carson.Near = []string{Torrance.Name, DominguezHills.Name, SanPedro.Name, NorthLongBeach.Name, Westminster.Name, LongBeach.Name}
	Fullerton.Near = []string{BuenaPark.Name, Brea.Name, YorbaLinda.Name, GardenGrove.Name}
	YorbaLinda.Near = []string{Fullerton.Name, Corona.Name, Orange.Name}
	Anaheim.Near = []string{BuenaPark.Name, GardenGrove.Name}
	Cypress.Near = []string{GardenGrove.Name, Bellflower.Name, Westminster.Name, FountainValley.Name, HuntingtonBeach.Name}
	Westminster.Near = []string{Carson.Name, LongBeach.Name, NorthLongBeach.Name, Cypress.Name, HuntingtonBeach.Name, GardenGrove.Name, FountainValley.Name}
	PalosVerdes.Near = []string{Torrance.Name}
	GardenGrove.Near = []string{Cypress.Name, Westminster.Name, Orange.Name, SantaAna.Name, Fullerton.Name, Anaheim.Name, FountainValley.Name, HuntingtonBeach.Name}
	Orange.Near = []string{YorbaLinda.Name, GardenGrove.Name, SantaAna.Name}
	LongBeach.Near = []string{NorthLongBeach.Name, Westminster.Name, Carson.Name}
	SanPedro.Near = []string{Carson.Name}
	HuntingtonBeach.Near = []string{Westminster.Name, Cypress.Name, GardenGrove.Name, FountainValley.Name}
	SantaAna.Near = []string{Orange.Name, GardenGrove.Name, FountainValley.Name}
	FountainValley.Near = []string{Westminster.Name, SantaAna.Name, GardenGrove.Name, HuntingtonBeach.Name, Cypress.Name}

	nodes := []*Node{
		&SantaClarita, &SanFernando, &LakeviewTerrace, &Sunland, &PanoramaCity, &LaCanada, &Pasadena, &VanNuys,
		&Burbank, &StudioCity, &Glendale, &Azusa, &SanDimas, &Claremont, &HighlandPark, &LosAngeles,
		&BeverlyHills, &Rosemead, &WestCovina, &Pomona, &ElMonte, &MontereyPark, &CalPoly, &USC,
		&Commerce, &SouthElMonte, &CulverCity, &ChinoHills, &SantaMonica, &DiamondBar, &CityOfIndustry, &SantaFeSprings,
		&LAX, &Watts, &Corona, &Brea, &Lynwood, &Compton, &Bellflower, &Hawthorne,
		&Gardena, &DominguezHills, &NorthLongBeach, &Torrance, &BuenaPark, &Carson, &Fullerton, &YorbaLinda,
		&Anaheim, &Cypress, &Westminster, &PalosVerdes, &GardenGrove, &Orange, &LongBeach, &SanPedro,
		&HuntingtonBeach, &SantaAna, &FountainValley,
	}

	pts, diag, err := zoneData(59)
	must(err)
	hmax = diag
	must(zipNodePts(nodes, pts))

	Nodes = make([]Node, len(nodes))
	for i, node := range nodes {
		Nodes[i] = *node
	}
}

// zoneData by Jason Lin © 2014
func zoneData(n int) (pts [][2]int, diag float64, err error) {
	const dotW, dotH = 12, 13
	red := color.RGBA{R: 255, A: 255}

	f, err := os.Open(sock.Root + "/etc/map.png")
	if err != nil {
		return
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}
	size := img.Bounds().Size()
	diag = math.Sqrt(float64(size.X*size.X + size.Y*size.Y))

	for y := 0; y < size.Y; y++ { // top to bottom
		for x := 0; x < size.X; x++ { // left to right
			if img.At(x, y) != red {
				continue
			}

			dx, dy := x+1, y+dotH/2

			new := true
			for _, pt := range pts {
				if abs(pt[0]-dx) <= dotW && abs(pt[1]-dy) <= dotH {
					new = false
				}
			}
			if new {
				pts = append(pts, [2]int{dx, dy})
			}
		}
	}

	if n >= 0 && len(pts) != n {
		err = newerr(itoa(len(pts)) + " points counted (" + itoa(n) + " required)")
	}
	return
}
