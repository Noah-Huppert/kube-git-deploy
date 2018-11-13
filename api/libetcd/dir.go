package libetcd

import (
	etcd "go.etcd.io/etcd/client"
)

// TraverseDir gets all values from a directory recursively
func TraverseDir(node *etcd.Node) []string {
	// If node not directory
	if !node.Dir {
		return []string{node.Value}
	}

	// Otherwise
	values := []string{}

	for _, n := range node.Nodes {
		// If directory call recursively
		if n.Dir {
			for _, subN := range TraverseDir(n) {
				values = append(values, subN)
			}
		} else {
			// Otherwise add value
			values = append(values, n.Value)
		}
	}

	return values
}
