package components

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FileEntry represents a file or directory in the tree
type FileEntry struct {
	name     string
	status   string
	isDir    bool
	children []*FileEntry
}

// FileTree builds a hierarchical tree from git status output
type FileTree struct {
	root *FileEntry
}

// ParseGitStatus parses git status output into a file tree
func ParseGitStatus(status string) *FileTree {
	tree := &FileTree{
		root: &FileEntry{
			name:     "",
			isDir:    true,
			children: make([]*FileEntry, 0),
		},
	}

	if status == "" {
		return tree
	}

	lines := strings.Split(strings.TrimSpace(status), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		statusCode := strings.TrimSpace(parts[0])
		path := parts[1]

		if len(parts) >= 3 && (strings.HasPrefix(statusCode, "R") || strings.HasPrefix(statusCode, "C")) {
			path = parts[2]
		}

		displayStatus := statusCode
		if len(statusCode) > 1 && (statusCode[0] == 'R' || statusCode[0] == 'C') {
			displayStatus = string(statusCode[0])
		}

		tree.addFile(path, displayStatus)
	}

	tree.sortTree(tree.root)
	return tree
}

// addFile adds a file to the tree, creating directories as needed
func (t *FileTree) addFile(path string, status string) {
	parts := strings.Split(path, "/")
	current := t.root

	for i, part := range parts {
		isFile := i == len(parts)-1

		var child *FileEntry
		for _, c := range current.children {
			if c.name == part {
				child = c
				break
			}
		}

		if child == nil {
			child = &FileEntry{
				name:     part,
				isDir:    !isFile,
				children: make([]*FileEntry, 0),
			}
			if isFile {
				child.status = status
			}
			current.children = append(current.children, child)
		}

		current = child
	}
}

// sortTree sorts the tree: directories first, then files, both alphabetically
func (t *FileTree) sortTree(node *FileEntry) {
	sort.Slice(node.children, func(i, j int) bool {
		if node.children[i].isDir != node.children[j].isDir {
			return node.children[i].isDir
		}
		return node.children[i].name < node.children[j].name
	})

	for _, child := range node.children {
		if child.isDir {
			t.sortTree(child)
		}
	}
}

// RenderTree renders the file tree as a string
func RenderTree(tree *FileTree) string {
	if tree == nil || tree.root == nil {
		return ""
	}

	var b strings.Builder
	renderNode(tree.root, "", true, &b)
	return b.String()
}

// renderEntryContent renders a file or directory entry
func renderEntryContent(node *FileEntry, prefix string, treeChar string, b *strings.Builder) {
	statusStyle := lipgloss.NewStyle().Bold(true)
	nameStyle := lipgloss.NewStyle()

	statusStr := ""
	if !node.isDir {
		switch node.status {
		case "A":
			statusStyle = statusStyle.Foreground(lipgloss.Color("2"))
			statusStr = "A"
		case "M":
			statusStyle = statusStyle.Foreground(lipgloss.Color("3"))
			statusStr = "M"
		case "D":
			statusStyle = statusStyle.Foreground(lipgloss.Color("1"))
			statusStr = "D"
		case "R":
			statusStyle = statusStyle.Foreground(lipgloss.Color("6"))
			statusStr = "R"
		case "C":
			statusStyle = statusStyle.Foreground(lipgloss.Color("5"))
			statusStr = "C"
		default:
			statusStr = node.status
		}

		b.WriteString(prefix)
		b.WriteString(treeChar + " ")
		b.WriteString(statusStyle.Render(statusStr))
		b.WriteString(" ")
		b.WriteString(nameStyle.Render(node.name))
		b.WriteString("\n")
	} else {
		b.WriteString(prefix)
		b.WriteString(treeChar + " ")
		b.WriteString(nameStyle.Render(node.name))
		b.WriteString("/\n")
	}
}

// renderNode recursively renders a node and its children
func renderNode(node *FileEntry, prefix string, isRoot bool, b *strings.Builder) {
	if !isRoot {
		renderEntryContent(node, prefix, "├─", b)
	}

	childPrefix := prefix
	if !isRoot {
		childPrefix = prefix + "│  "
	}

	for i, child := range node.children {
		if i == len(node.children)-1 {
			lastPrefix := strings.TrimSuffix(childPrefix, "│  ")
			if !isRoot {
				lastPrefix += "   "
			}
			renderLastNode(child, lastPrefix, b)
		} else {
			renderNode(child, childPrefix, false, b)
		}
	}
}

// renderLastNode renders the last child with proper tree characters
func renderLastNode(node *FileEntry, prefix string, b *strings.Builder) {
	renderEntryContent(node, prefix, "└─", b)

	childPrefix := prefix + "   "
	for i, child := range node.children {
		if i == len(node.children)-1 {
			renderLastNode(child, childPrefix, b)
		} else {
			renderNode(child, childPrefix+"│  ", false, b)
		}
	}
}

// CountFiles counts total files in the tree
func CountFiles(tree *FileTree) int {
	if tree == nil || tree.root == nil {
		return 0
	}
	return countFilesInNode(tree.root)
}

func countFilesInNode(node *FileEntry) int {
	count := 0
	for _, child := range node.children {
		if child.isDir {
			count += countFilesInNode(child)
		} else {
			count++
		}
	}
	return count
}
