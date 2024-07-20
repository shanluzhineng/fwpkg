// this file is copied from https://github.com/dnaeon/go-dependency-graph-algorithm

package depends

import (
	"errors"
	"fmt"

	"reflect"

	mapset "github.com/deckarep/golang-set"
	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/log"
)

// Node represents a single node in the graph with it's dependencies
type Node struct {
	// Index of the node
	index int

	// Name of the node
	name string

	// data of the node
	data *factory.MetaData

	// Dependencies of the node
	deps []*Node
}

// ErrCircularDependency report that circular dependency found
var ErrCircularDependency = errors.New("circular dependency found")

// NewNode creates a new node
func NewNode(index int, data interface{}, deps ...*Node) *Node {
	var md *factory.MetaData
	if reflect.TypeOf(data).Kind() == reflect.String {
		md = &factory.MetaData{Name: data.(string)}
	} else {
		md = data.(*factory.MetaData)
	}

	n := &Node{
		index: index,
		name:  md.Name,
		data:  md,
		deps:  deps,
	}

	return n
}

// Graph is the collection of node
type Graph []*Node

// Resolves the dependency graph
func resolveGraph(graph Graph) (Graph, error) {
	// A map containing the node names and the actual node object
	nodeNames := make(map[string]*Node)

	// A map containing the nodes and their dependencies
	nodeDependencies := make(map[string]mapset.Set)

	// Populate the maps
	for _, node := range graph {
		nodeNames[node.name] = node

		dependencySet := mapset.NewSet()
		for _, dep := range node.deps {
			dependencySet.Add(dep.name)
		}
		if nodeDependencies[node.name] != nil {
			log.Debugf("%v is already exist, overwrite!", node.name)
		}
		nodeDependencies[node.name] = dependencySet
	}

	//log.Debug(nodeDependencies)

	// Iteratively find and remove nodes from the graph which have no dependencies.
	// If at some point there are still nodes in the graph and we cannot find
	// nodes without dependencies, that means we have a circular dependency
	var resolved Graph
	for len(nodeDependencies) != 0 {
		// Get all nodes from the graph which have no dependencies
		readySet := mapset.NewSet()
		for name, deps := range nodeDependencies {
			if deps.Cardinality() == 0 {
				readySet.Add(name)
			}
		}

		// If there aren't any ready nodes, then we have a circular dependency
		if readySet.Cardinality() == 0 {
			var g Graph
			for name := range nodeDependencies {
				g = append(g, nodeNames[name])
			}

			return g, ErrCircularDependency
		}

		// Remove the ready nodes and add them to the resolved graph
		for name := range readySet.Iter() {
			delete(nodeDependencies, name.(string))
			resolved = append(resolved, nodeNames[name.(string)])
		}

		// Also make sure to remove the ready nodes from the
		// remaining node dependencies as well
		for name, deps := range nodeDependencies {
			diff := deps.Difference(readySet)
			nodeDependencies[name] = diff
		}
	}

	return resolved, nil
}

// Displays the dependency graph
func displayDependencyGraph(name string, graph Graph, logger func(v ...interface{})) {
	output := name + ":\n\nDependency tree:\n"
	for i, node := range graph {
		if len(node.deps) == 0 {
			output += fmt.Sprintf("\t%4d (%4d): %s ->\n", i, node.index, node.name)
		} else {
			for _, dep := range node.deps {
				output += fmt.Sprintf("\t%4d (%4d): %s -> %s\n", i, node.index, node.name, dep.name)
			}
		}
	}
	if reflect.TypeOf(logger).Kind() == reflect.Func {
		logger(output)
	}
}
