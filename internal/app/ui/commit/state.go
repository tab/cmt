package commit

// State holds the data for the commit TUI
type State struct {
	Files         []FileNode
	CommitMessage string
	Prefix        string
	Diff          string
	Accepted      bool
	Error         error
}

// FileNode represents a file or directory in the tree
type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []FileNode
	Status   string // A=added, M=modified, D=deleted
}
