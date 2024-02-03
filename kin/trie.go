package kin

import (
	"strings"
)

//   "/" as root node
// pattern serve as a sign of an exist path
type node struct {
	pattern  string  // complete router path to match
	part     string  // segment of router path at current node
	children []*node // child nodes
	isWild   bool    // contain parameter(:id) or wildcard (*)
}

func (n *node) mathChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// insert  router path
func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern // save in the last node
		return	// finished
	}

	// find a child match current level path
	part := parts[height]
	child := n.mathChild(part)
	if child == nil {	// not exist
		flag := part[0] == ':' || part[0] == '*'

		child = &node{
			part: parts[height],
			isWild: flag,
		}
		n.children = append(n.children, child)
	}

	// match next level
	child.insert(pattern, parts, height + 1)
}


// find all match node
// e.g.:  /go/12  matches /go/12  and /:lang/12  and /:lang/:page etc.
func (n *node) matchChildren(part string) []*node {
	ret := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			ret = append(ret, child)
		}
	}
	return ret
}

// query  router path
func(n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*"){
		if n.pattern == "" { // no sign -- not a tail node
			return nil
		}
		return n 
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height + 1)
		if result != nil {
			return result
		}
	}

	return nil
}