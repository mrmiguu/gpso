package zone

import (
	"math"
	"sort"

	"github.com/mrmiguu/gpso/node"
)

// Find finds the shortest path from the source to the destination.
func Find(src, dst node.T) ([]node.T, error) {
	return find(src, dst, hmax)
}

func find(src, dst node.T, hmax float64) ([]node.T, error) {
	closed := []node.T{}
	open := []node.T{src}
	from := map[string]node.T{}
	g := map[string]float64{src.Name: 0}

	h, err := heur(src, dst, hmax)
	if err != nil {
		return nil, err
	}
	f := map[string]float64{src.Name: h}

	for {
		if len(open) == 0 {
			break
		}

		cur := open[0]
		open = open[1:]

		if cur.Name == dst.Name {
			return redraw(from, cur), nil
		}

		closed = append(closed, cur)

		for _, name := range cur.Near {
			near, err := Aton(name)
			if err != nil {
				return nil, err
			}

			if near.Contained(closed) {
				continue
			}

			tempg := g[cur.Name] + 1

			if near.Contained(open) {
				open = append(open, near)
				sort.Slice(open, func(i, j int) bool { return f[open[i].Name] < f[open[j].Name] })
			} else if tempg >= g[name] {
				continue
			}

			from[name] = cur
			g[name] = tempg

			h, err := heur(near, dst, hmax)
			if err != nil {
				return nil, err
			}
			f[name] = g[name] + h
		}
	}

	return []node.T{src}, nil
}

func heur(src, dst node.T, max float64) (float64, error) {
	Nsrc, err := Aton(src.Name)
	if err != nil {
		return 0, err
	}
	Ndst, err := Aton(dst.Name)
	if err != nil {
		return 0, err
	}
	if max != 0 {
		return dist(Nsrc.Pt, Ndst.Pt) / max, nil
	}
	return 0, newerr("divide by zero")
}

func dist(a, b [2]int) float64 {
	x, y := b[0]-a[0], b[1]-a[1]
	return math.Sqrt(float64(x*x + y*y))
}

func redraw(from map[string]node.T, cur node.T) []node.T {
	path := []node.T{cur}
	for {
		if _, found := from[cur.Name]; !found {
			break
		}
		cur = from[cur.Name]
		path = append([]node.T{cur}, path...)
	}
	return path
}
