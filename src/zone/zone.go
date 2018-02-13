package zone

import (
	"encoding/json"
	"errors"
	"strings"

	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"math"
	"math/rand"
	"strconv"

	"github.com/mrmiguu/gpso/src/ex"
	"github.com/mrmiguu/gpso/src/node"
	"github.com/mrmiguu/jsutil"
	"github.com/mrmiguu/sock"
)

var (
	vtob       = json.Marshal
	join       = strings.Join
	open       = jsutil.Open
	newerr     = errors.New
	itoa       = strconv.Itoa
	must       = ex.Must
	abs        = ex.Abs
	zipNodePts = node.ZipNodePts

	Width  int
	Height int

	hmax float64
)

var (
	SantaClarita    = node.T{Name: "santaclarita", Hwys: []int{5, 210}}
	SanFernando     = node.T{Name: "sanfernando", Hwys: []int{5, 405}}
	LakeviewTerrace = node.T{Name: "lakeviewterrace", Hwys: []int{210}}
	Sunland         = node.T{Name: "sunland", Hwys: []int{210}}
	PanoramaCity    = node.T{Name: "panoramacity", Hwys: []int{5, 170}}
	LaCanada        = node.T{Name: "lacanada", Hwys: []int{2, 210}}
	Pasadena        = node.T{Name: "pasadena", Hwys: []int{110}}
	VanNuys         = node.T{Name: "vannuys", Hwys: []int{101, 405}}
	Burbank         = node.T{Name: "burbank", Hwys: []int{5, 101}}
	StudioCity      = node.T{Name: "studiocity", Hwys: []int{101, 170}}
	Glendale        = node.T{Name: "glendale", Hwys: []int{2, 101}}
	Azusa           = node.T{Name: "azusa", Hwys: []int{210, 605}}
	SanDimas        = node.T{Name: "sandimas", Hwys: []int{57, 210}}
	Claremont       = node.T{Name: "claremont", Hwys: []int{210}}
	HighlandPark    = node.T{Name: "highlandpark", Hwys: []int{110}}
	LosAngeles      = node.T{Name: "losangeles", Hwys: []int{5, 10, 110}}
	BeverlyHills    = node.T{Name: "beverlyhills", Hwys: []int{2, 170}}
	Rosemead        = node.T{Name: "rosemead", Hwys: []int{10}}
	WestCovina      = node.T{Name: "westcovina", Hwys: []int{10}}
	Pomona          = node.T{Name: "pomona", Hwys: []int{10}}
	ElMonte         = node.T{Name: "elmonte", Hwys: []int{10, 605}}
	MontereyPark    = node.T{Name: "montereypark", Hwys: []int{10, 710}}
	CalPoly         = node.T{Name: "calpoly", Hwys: []int{57}}
	USC             = node.T{Name: "usc", Hwys: []int{10, 60, 110}}
	Commerce        = node.T{Name: "commerce", Hwys: []int{60, 710}}
	SouthElMonte    = node.T{Name: "southelmonte", Hwys: []int{60, 605}}
	CulverCity      = node.T{Name: "culvercity", Hwys: []int{10, 405}}
	ChinoHills      = node.T{Name: "chinohills", Hwys: []int{60}}
	SantaMonica     = node.T{Name: "santamonica", Hwys: []int{10}}
	DiamondBar      = node.T{Name: "diamondbar", Hwys: []int{57}}
	CityOfIndustry  = node.T{Name: "cityofindustry", Hwys: []int{60}}
	SantaFeSprings  = node.T{Name: "santafesprings", Hwys: []int{5, 605}}
	LAX             = node.T{Name: "lax", Hwys: []int{105, 405}}
	Watts           = node.T{Name: "watts", Hwys: []int{105, 110}}
	Corona          = node.T{Name: "corona", Hwys: []int{91}}
	Brea            = node.T{Name: "brea", Hwys: []int{57}}
	Lynwood         = node.T{Name: "lynwood", Hwys: []int{105, 710}}
	Compton         = node.T{Name: "compton", Hwys: []int{91, 710}}
	Bellflower      = node.T{Name: "bellflower", Hwys: []int{91, 605}}
	Hawthorne       = node.T{Name: "hawthorne", Hwys: []int{91, 405}}
	Gardena         = node.T{Name: "gardena", Hwys: []int{91, 110}}
	DominguezHills  = node.T{Name: "dominguezhills", Hwys: []int{110}}
	NorthLongBeach  = node.T{Name: "northlongbeach", Hwys: []int{710}}
	Torrance        = node.T{Name: "torrance", Hwys: []int{405}}
	BuenaPark       = node.T{Name: "buenapark", Hwys: []int{5, 91}}
	Carson          = node.T{Name: "carson", Hwys: []int{110, 405}}
	Fullerton       = node.T{Name: "fullerton", Hwys: []int{57, 91}}
	YorbaLinda      = node.T{Name: "yorbalinda", Hwys: []int{55, 91}}
	Anaheim         = node.T{Name: "anaheim", Hwys: []int{5}}
	Cypress         = node.T{Name: "cypress", Hwys: []int{605}}
	Westminster     = node.T{Name: "westminster", Hwys: []int{22, 405}}
	PalosVerdes     = node.T{Name: "palosverdes", Hwys: []int{}}
	GardenGrove     = node.T{Name: "gardengrove", Hwys: []int{5, 22, 57}}
	Orange          = node.T{Name: "orange", Hwys: []int{22, 55}}
	LongBeach       = node.T{Name: "longbeach", Hwys: []int{710}}
	SanPedro        = node.T{Name: "sanpedro", Hwys: []int{110}}
	HuntingtonBeach = node.T{Name: "huntingtonbeach", Hwys: []int{605}}
	SantaAna        = node.T{Name: "santaana", Hwys: []int{5, 55, 405}}
	FountainValley  = node.T{Name: "fountainvalley", Hwys: []int{405}}
)

var Nodes []node.T

var name2node = map[string]node.T{}
var NodeMap []byte

func Aton(s string) (node.T, error) {
	if n, found := name2node[s]; found {
		return n, nil
	}
	return node.T{}, newerr("'" + s + "' not found in zone")
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

func SrcDst() (src, dst node.T) {
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

	nodes := []*node.T{
		&SantaClarita, &SanFernando, &LakeviewTerrace, &Sunland, &PanoramaCity, &LaCanada, &Pasadena, &VanNuys,
		&Burbank, &StudioCity, &Glendale, &Azusa, &SanDimas, &Claremont, &HighlandPark, &LosAngeles,
		&BeverlyHills, &Rosemead, &WestCovina, &Pomona, &ElMonte, &MontereyPark, &CalPoly, &USC,
		&Commerce, &SouthElMonte, &CulverCity, &ChinoHills, &SantaMonica, &DiamondBar, &CityOfIndustry, &SantaFeSprings,
		&LAX, &Watts, &Corona, &Brea, &Lynwood, &Compton, &Bellflower, &Hawthorne,
		&Gardena, &DominguezHills, &NorthLongBeach, &Torrance, &BuenaPark, &Carson, &Fullerton, &YorbaLinda,
		&Anaheim, &Cypress, &Westminster, &PalosVerdes, &GardenGrove, &Orange, &LongBeach, &SanPedro,
		&HuntingtonBeach, &SantaAna, &FountainValley,
	}

	pts, size, err := zoneData(59)
	must(err)
	Width, Height = size[0], size[1]
	hmax = math.Sqrt(float64(Width*Width + Height*Height))
	must(zipNodePts(nodes, pts))

	Nodes = make([]node.T, len(nodes))
	for i, node := range nodes {
		Nodes[i] = *node
		name2node[node.Name] = *node
	}

	NodeMap, err = vtob(name2node)
	must(err)

	// for _, node := range nodes {
	// 	if node.Pt[0] != 0 || node.Pt[1] != 0 {
	// 		return
	// 	}
	// }
	// must(newerr("all zone points are empty"))
}

// zoneData by Jason Lin © 2014
func zoneData(n int) (pts [][2]int, size [2]int, err error) {
	const dotW, dotH = 12, 13
	red := color.RGBA{R: 255, A: 255}

	r, err := open(sock.Root + "/etc/map.png")
	if err != nil {
		return
	}
	defer r.Close()
	img, _, err := image.Decode(r)
	if err != nil {
		return
	}
	rect := img.Bounds().Size()
	if rect.X == 0 || rect.Y == 0 {
		err = newerr("bad zone size")
		return
	}
	size[0], size[1] = rect.X, rect.Y

	for y := 0; y < rect.Y; y++ { // top to bottom
		for x := 0; x < rect.X; x++ { // left to right
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
	for _, pt := range pts {
		if pt[0] != 0 || pt[1] != 0 {
			return
		}
	}
	err = newerr("all zone points are empty")
	return
}
