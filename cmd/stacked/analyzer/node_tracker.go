package analyzer

import "go/ast"

type nodeTracker int

func newNodeTracker() nodeTracker {
	return -1
}

func (n *nodeTracker) enter() {
	*n = 1
}

func (n *nodeTracker) step(node ast.Node) {
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

func (n *nodeTracker) isIn() bool {
	return *n > -1
}
