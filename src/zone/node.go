package zone

// Node holds a city's information.
type Node struct {
	Name string
	Pt   [2]int
	Hwys []int
	Near []string
}

// Index finds the index of the node in nodes.
func (n Node) Index(nodes []Node) (int, error) {
	for i, node := range nodes {
		if node.Name == n.Name {
			return i, nil
		}
	}
	return -1, newerr("node not found in nodes")
}

// Contained checks if node is contained in nodes.
func (n Node) Contained(nodes []Node) bool {
	_, err := n.Index(nodes)
	return err == nil
}

func (n Node) String() string {
	return n.Name + sprint(n.Pt) + ":\n\tHwys=" + sprint(n.Hwys) + "\n\tNear=" + sprint(n.Near)
}

func zipNodePts(nodes []*Node, pts [][2]int) error {
	nL, pL := len(nodes), len(pts)
	if nL != pL {
		return newerr(itoa(nL) + "-nodes vs " + itoa(pL) + "-pts")
	}
	for i, node := range nodes {
		node.Pt = pts[i]
	}
	return nil
}
