package swaggering

import "strings"

type Tree struct {
	Parent    *Tree
	Children  map[string]*Tree
	Endpoints map[string]Endpoint
}

func (t *Tree) child(dir string) *Tree {
	if t.Children == nil {
		t.Children = map[string]*Tree{}
	}

	child, ok := t.Children[dir]
	if ok {
		return child
	}

	child = &Tree{
		Parent: t,
	}
	t.Children[dir] = child
	return child
}

func (t *Tree) register(endpoint Endpoint) {
	path := endpoint.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	dirs := strings.Split(path, "/")

	tree := t
	for _, dir := range dirs {
		tree = tree.child(dir)
	}

	if tree.Endpoints == nil {
		tree.Endpoints = map[string]Endpoint{}
	}

	tree.Endpoints[endpoint.Method] = endpoint
}
