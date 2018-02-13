package zone

import (
	"fmt"
	"math"

	"github.com/mrmiguu/gpso/src/node"
)

// Find finds the shortest path from the source to the destination.
func Find(src, dst node.T) ([]node.T, error) {
	return find2(src, dst, hmax)
}

func String(path []node.T) string {
	return "path: " + Join(path, " -> ")
}

func Join(path []node.T, sep string) string {
	cities := make([]string, len(path))
	for i, city := range path {
		cities[i] = city.Name
	}
	return join(cities, sep)
}

type Set map[string]node.T

func NewSet(n ...node.T) Set {
	s := Set{}
	for _, N := range n {
		s[N.Name] = N
	}
	return s
}
func (s Set) Add(n node.T) {
	s[n.Name] = n
}
func (s Set) Remove(n node.T) {
	delete(s, n.Name)
}
func (s Set) Contains(n node.T) bool {
	_, found := s[n.Name]
	return found
}

type Map map[string]float64

func (m Map) Get(n node.T) float64 {
	if f, found := m[n.Name]; found {
		return f
	}
	return math.Inf(0)
}
func (m Map) Put(n node.T, f float64) {
	m[n.Name] = f
}

func find2(start, goal node.T, hmax float64) ([]node.T, error) {
	println("find2 begin")
	defer println("find2 end")
	closed := map[string]node.T{}
	open := map[string]node.T{start.Name: start}
	from := map[string]node.T{}
	g := Map{start.Name: 0} // new entires inf
	h, err := heur(start, goal, hmax)
	if err != nil {
		return nil, err
	}
	f := Map{start.Name: h} // new entires inf

	for len(open) > 0 {
		var cur node.T
		score := math.Inf(0)
		for _, city := range open {
			if h := f.Get(city); h <= score {
				cur, score = city, h
			}
		}

		if cur.Name == goal.Name {
			return redraw(from, cur), nil
		}

		delete(open, cur.Name)
		if _, found := open[cur.Name]; found {
			return nil, newerr("could not remove '" + cur.Name + "' from " + fmt.Sprint(open))
		}
		closed[cur.Name] = cur
		if _, found := closed[cur.Name]; !found {
			return nil, newerr("could not remove '" + cur.Name + "' from " + fmt.Sprint(closed))
		}

		for _, nearStr := range cur.Near {
			near, err := Aton(nearStr)
			if err != nil {
				return nil, err
			}

			if _, found := closed[near.Name]; found {
				continue
			}
			open[near.Name] = near

			tempg := g.Get(cur) + dist(cur.Pt, near.Pt)
			if tempg >= g.Get(near) {
				continue
			}

			h, err := heur(near, goal, hmax)
			if err != nil {
				return nil, err
			}

			from[near.Name] = cur
			g.Put(near, tempg)
			f.Put(near, g.Get(near)+h)
		}
	}

	return nil, newerr("could not find path")
}

func heur(start, goal node.T, max float64) (float64, error) {
	Nstart, err := Aton(start.Name)
	if err != nil {
		return 0, err
	}
	Ngoal, err := Aton(goal.Name)
	if err != nil {
		return 0, err
	}
	if max != 0 {
		return dist(Nstart.Pt, Ngoal.Pt) / max, nil
	}
	return 0, newerr("divide by zero")
}

func dist(a, b [2]int) float64 {
	x, y := b[0]-a[0], b[1]-a[1]
	return math.Sqrt(float64(x*x + y*y))
}

func redraw(from map[string]node.T, cur node.T) []node.T {
	path := []node.T{cur}
	for found := false; ; {
		if cur, found = from[cur.Name]; !found {
			break
		}
		path = append([]node.T{cur}, path...)
	}
	return path
}
