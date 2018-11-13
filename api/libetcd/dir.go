package libetcd

import (
	"strings"

	etcd "go.etcd.io/etcd/client"
)

// lastKeyIs indicates if the last key in a key path is the specified value
func lastKeyIs(fullPath, key string) bool {
	parts := strings.Split(fullPath, "/")
	return parts[len(parts)-1] == key
}

// TraverseDir gets all values with a specified key from a
// directory recursively
func TraverseDir(node *etcd.Node, key string) []string {
	// If node not directory
	if !node.Dir && lastKeyIs(node.Key, key) {
		return []string{node.Value}
	}

	// Otherwise
	values := []string{}

	for _, n := range node.Nodes {
		// If directory call recursively
		if n.Dir {
			for _, subN := range TraverseDir(n, key) {
				values = append(values, subN)
			}
		} else {
			// Otherwise add value
			if lastKeyIs(n.Key, key) {
				values = append(values, n.Value)
			}
		}
	}

	return values
}
