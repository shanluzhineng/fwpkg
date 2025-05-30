// Package depends provides dependency resolver for factory
package depends

import (
	"reflect"

	"github.com/shanluzhineng/fwpkg/system/factory"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/fwpkg/utils/str"
)

// depResolver sort by the configuration dependency which specified by tag depends
type depResolver []*factory.MetaData

func (s depResolver) Resolve() (resolved Graph, err error) {

	var workingGraph Graph
	var node *Node
	for i, item := range s {
		// find the index of its dependency
		dep, ok := s.findDependencies(item)
		// for debug only
		if ok {
			node = NewNode(i, item, dep...)
		} else {
			node = NewNode(i, item)
		}
		workingGraph = append(workingGraph, node)
	}
	resolved, err = resolveGraph(workingGraph)
	//displayDependencyGraph("working graph", workingGraph, log.Debug)
	return
}

func (s depResolver) findDependencyIndex(depName string) int {
	for i, item := range s {
		// find type name
		if item.TypeName == depName {
			return i
		}

		// find item name or package name
		if item.ShortName == depName || item.Name == depName || item.PkgName == depName {
			return i
		}

		// else find method name
		if item.Kind == factory.Method {
			method := item.MetaObject.(reflect.Method)
			methodName := str.ToLowerCamel(method.Name)
			if depName == methodName || depName == item.PkgName+"."+methodName {
				return i
			}
		}
	}
	return -1
}

func (s depResolver) findDependencies(item *factory.MetaData) (dep []*Node, ok bool) {
	// iterate dependencies
	if len(item.DepNames) > 0 {
		for _, dp := range item.DepNames {
			depIdx := s.findDependencyIndex(dp)
			if depIdx >= 0 {
				depMetaData := s[depIdx]
				item.DepMetaData = append(item.DepMetaData, depMetaData)
				dep = append(dep, NewNode(depIdx, depMetaData))
			} else {
				// found external dependency
				extData := &factory.MetaData{Name: dp}
				item.DepMetaData = append(item.DepMetaData, extData)
				dep = append(dep, NewNode(depIdx, extData))
				log.Warnf("dependency %v is not found", dp)
			}
		}
		ok = true
	}

	return
}

// Resolve resolve dependencies
func Resolve(data []*factory.MetaData) (result []*factory.MetaData, err error) {
	if len(data) != 0 {

		dep := depResolver(data)
		var resolved Graph
		resolved, err = dep.Resolve()

		if err != nil {
			//log.Errorf("Failed to resolve dependencies: %s", err)
			displayDependencyGraph("missing dependency or circular dependency graph", resolved, log.Error)
		} else {
			log.Debugf("The dependency graph resolved successfully")
			displayDependencyGraph("resolved dependency graph", resolved, log.Debug)
			for _, item := range resolved {
				if item.index >= 0 {
					result = append(result, data[item.index])
				}
			}
		}
	}

	return
}
