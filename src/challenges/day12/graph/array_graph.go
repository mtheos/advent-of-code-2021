package graph

import "math"

type ArrayGraph struct {
	nodes []string
	edges [][]bool
}

func (g *ArrayGraph) Neighbours(node int) []int {
	var neighbours []int
	for to, edge := range g.edges[node] {
		if edge {
			neighbours = append(neighbours, to)
		}
	}
	return neighbours
}

func (g *ArrayGraph) ChannelIterator() <-chan ThisCouldHaveBeenAvoidedIfGoLetYouImplementRangeOverCustomTypes {
	ch := make(chan ThisCouldHaveBeenAvoidedIfGoLetYouImplementRangeOverCustomTypes)
	go func() {
		for i, node := range g.nodes {
			ch <- ThisCouldHaveBeenAvoidedIfGoLetYouImplementRangeOverCustomTypes{i, node}
		}
		close(ch)
	}()
	return ch
}

func (g *ArrayGraph) Size() int {
	return len(g.nodes)
}

func (g *ArrayGraph) exists(name string) bool {
	for _, node := range g.nodes {
		if node == name {
			return true
		}
	}
	return false
}

func (g *ArrayGraph) MaybeCreate(name string) {
	if !g.exists(name) {
		g.create(name)
	}
}

func (g *ArrayGraph) create(name string) {
	g.nodes = append(g.nodes, name)
	g.edges = append(g.edges, []bool{})
}

func (g *ArrayGraph) Name(idx int) string {
	return g.nodes[idx]
}

func (g *ArrayGraph) Connect(from string, to string, both bool) {
	fIdx, tIdx := g.Idx(from), g.Idx(to)
	max := int(math.Max(float64(fIdx), float64(tIdx)))
	for i := 0; i < len(g.edges); i++ {
		for j := len(g.edges[i]); j <= max; j++ {
			g.edges[i] = append(g.edges[i], false)
		}
	}
	g.edges[fIdx][tIdx] = true
	if both {
		g.edges[tIdx][fIdx] = true
	}
}

func (g *ArrayGraph) Idx(name string) int {
	for i, node := range g.nodes {
		if node == name {
			return i
		}
	}
	return -1
}
