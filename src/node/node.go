package node

import (
	"errors"
	"strconv"
)

var (
	// sprint = fmt.Sprint
	newerr = errors.New
	itoa   = strconv.Itoa
)

// T holds a city's information.
type T struct {
	Name string
	Pt   [2]int
	Hwys []int
	Near []string
}

// Index finds the index of the node in nodes.
func (t T) Index(nodes []T) (int, error) {
	for i, node := range nodes {
		if node.Name == t.Name {
			return i, nil
		}
	}
	return -1, newerr("node not found in nodes")
}

// Contained checks if node is contained in nodes.
func (t T) Contained(nodes []T) bool {
	_, err := t.Index(nodes)
	return err == nil
}

// func (t T) String() string {
// 	return t.Name + sprint(t.Pt) + ":\n\tHwys=" + sprint(t.Hwys) + "\n\tNear=" + sprint(t.Near)
// }
