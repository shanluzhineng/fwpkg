package depends

import (
	"fmt"
	"testing"

	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/stretchr/testify/assert"
)

func TestDepTree(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//
	// A working dependency graph
	//
	nodeA := NewNode(0, "A")
	nodeB := NewNode(1, "B")
	nodeC := NewNode(2, "C", nodeA)
	nodeD := NewNode(3, "D", nodeB)
	nodeE := NewNode(4, "E", nodeC, nodeD)
	nodeF := NewNode(5, "F", nodeA, nodeB)
	nodeG := NewNode(6, "G", nodeE, nodeF)
	nodeH := NewNode(7, "H", nodeG)
	nodeI := NewNode(8, "I", nodeA)
	nodeJ := NewNode(8, "J", nodeB)
	nodeK := NewNode(10, "K")

	var workingGraph Graph
	workingGraph = append(workingGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A working dependency graph\n")
	displayDependencyGraph("workingGraph", workingGraph, log.Debug)

	resolved, err := resolveGraph(workingGraph)
	assert.Equal(t, nil, err)
	if err != nil {
		log.Errorf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Debugf("The dependency graph resolved successfully")
	}
	displayDependencyGraph("resolved", resolved, log.Debug)

	//
	// A broken dependency graph with circular dependency
	//
	nodeA = NewNode(11, "A", nodeI)

	var brokenGraph Graph
	brokenGraph = append(brokenGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A broken dependency graph with circular dependency\n")
	displayDependencyGraph("brokenGraph", brokenGraph, log.Debug)

	resolved, err = resolveGraph(brokenGraph)
	assert.Equal(t, ErrCircularDependency, err)
	if err != nil {
		log.Errorf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Debugf("The dependency graph resolved successfully")
	}
	displayDependencyGraph("resolved", resolved, log.Debug)
}
