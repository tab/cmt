package commit

import (
	"sort"
	"strings"
)

// BuildFileTree constructs a hierarchical tree from a flat list of file paths with status
func BuildFileTree(statusOutput string) []FileNode {
	if strings.TrimSpace(statusOutput) == "" {
		return []FileNode{}
	}

	root := make(map[string]*FileNode)
	lines := strings.Split(strings.TrimSpace(statusOutput), "\n")

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		path := parts[1]

		if strings.HasPrefix(status, "R") && len(parts) >= 3 {
			path = parts[2]
		}

		segments := strings.Split(path, "/")

		var parent *FileNode
		currentPath := ""

		for i, segment := range segments {
			if segment == "" {
				continue
			}

			if currentPath == "" {
				currentPath = segment
			} else {
				currentPath = currentPath + "/" + segment
			}

			isDir := i < len(segments)-1

			var node *FileNode
			if parent == nil {
				node = root[segment]
				if node == nil {
					node = &FileNode{
						Name:  segment,
						Path:  currentPath,
						IsDir: isDir,
					}
					root[segment] = node
				}
			} else {
				node = findChild(parent, segment)
				if node == nil {
					parent.Children = append(parent.Children, FileNode{
						Name:  segment,
						Path:  currentPath,
						IsDir: isDir,
					})
					node = &parent.Children[len(parent.Children)-1]
				}
			}

			if isDir {
				node.IsDir = true
				node.Status = ""
			} else {
				node.IsDir = false
				node.Status = status
			}

			parent = node
		}
	}

	rootNodes := make([]*FileNode, 0, len(root))
	for _, node := range root {
		rootNodes = append(rootNodes, node)
	}

	sort.Slice(rootNodes, func(i, j int) bool {
		if rootNodes[i].IsDir && !rootNodes[j].IsDir {
			return true
		}
		if !rootNodes[i].IsDir && rootNodes[j].IsDir {
			return false
		}
		return rootNodes[i].Name < rootNodes[j].Name
	})

	result := make([]FileNode, 0, len(rootNodes))
	for _, node := range rootNodes {
		sortChildren(node)
		result = append(result, *node)
	}

	return result
}

// sortChildren recursively sorts children of a node
func sortChildren(node *FileNode) {
	if !node.IsDir || len(node.Children) == 0 {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir && !node.Children[j].IsDir {
			return true
		}
		if !node.Children[i].IsDir && node.Children[j].IsDir {
			return false
		}
		return node.Children[i].Name < node.Children[j].Name
	})

	for i := range node.Children {
		sortChildren(&node.Children[i])
	}
}

// findChild returns a pointer to a child node with the given name, if present
func findChild(parent *FileNode, name string) *FileNode {
	for i := range parent.Children {
		if parent.Children[i].Name == name {
			return &parent.Children[i]
		}
	}
	return nil
}
