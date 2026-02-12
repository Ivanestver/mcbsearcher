package types

type Path []*Edge

func (path *Path) Last() *Edge {
	return (*path)[len(*path)-1]
}
