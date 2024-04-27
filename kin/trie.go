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
	} else {
		/////////////////////////////
		if child.part[0] == ':' {
			height -= 1
		}
		/////////////////////////////////
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
func(n *node) search(parts []string, height int, flag bool) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*"){
		if n.pattern == "" { // no sign -- not a tail node
			return nil
		}
		return n 
	}

	part := parts[height]
	// flag := false
	children := n.matchChildren(part)

	for _, child := range children {
		if child.part[0] == ':' {	// e.g.  : 
			flag = true	// match a dynamic parameter path
			height -= 1
		}

		result := child.search(parts, height + 1, flag)
		if result != nil {
			return result
		} else {
			if flag {
				return child.search(parts, height + 2, flag)	// TODO
			}
		}
	}

	// if flag {	// children = nil or no match child (in children)
	// 	return n // return current router (a dynamic parameter router)
	// }
	return nil
}