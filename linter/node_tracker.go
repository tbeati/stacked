package linter

import "go/ast"

type nodeTracker int

func newNodeTracker() nodeTracker {
	return -1
}

func (n *nodeTracker) enterNode() {
	*n = 1
}

func (n *nodeTracker) depthFirstSearchStep(node ast.Node) {
	if *n == 0 {
		*n = -1
	} else if *n > -1 {
		if node == nil {
			*n--
		} else {
			*n++
		}
	}
}

func (n *nodeTracker) isInNode() bool {
	return *n > -1
}
