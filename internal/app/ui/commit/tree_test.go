package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildFileTree(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		checkFn  func(*testing.T, []FileNode)
	}{
		{
			name:     "Success with empty input",
			input:    "",
			expected: 0,
			checkFn:  func(t *testing.T, nodes []FileNode) {},
		},
		{
			name:     "Success with whitespace only",
			input:    "   \n  \n  ",
			expected: 0,
			checkFn:  func(t *testing.T, nodes []FileNode) {},
		},
		{
			name:     "Success with single file added",
			input:    "A\tfile.txt",
			expected: 1,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "file.txt", nodes[0].Name)
				assert.Equal(t, "A", nodes[0].Status)
				assert.False(t, nodes[0].IsDir)
			},
		},
		{
			name:     "Success with single file modified",
			input:    "M\tfile.txt",
			expected: 1,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "file.txt", nodes[0].Name)
				assert.Equal(t, "M", nodes[0].Status)
			},
		},
		{
			name:     "Success with single file deleted",
			input:    "D\tfile.txt",
			expected: 1,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "D", nodes[0].Status)
			},
		},
		{
			name:     "Success with nested structure",
			input:    "A\tsrc/main.go\nM\tsrc/utils/helper.go",
			expected: 1,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "src", nodes[0].Name)
				assert.True(t, nodes[0].IsDir)
				assert.Equal(t, 2, len(nodes[0].Children))
			},
		},
		{
			name:     "Success with renamed file",
			input:    "R100\told.txt\tnew.txt",
			expected: 1,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "new.txt", nodes[0].Name)
				assert.Equal(t, "R100", nodes[0].Status)
			},
		},
		{
			name:     "Success with alphabetical ordering",
			input:    "A\tzebra.txt\nM\tapple.txt",
			expected: 2,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.Equal(t, "apple.txt", nodes[0].Name)
				assert.Equal(t, "zebra.txt", nodes[1].Name)
			},
		},
		{
			name:     "Success with directories first",
			input:    "A\tfile.txt\nA\tdir/file.txt",
			expected: 2,
			checkFn: func(t *testing.T, nodes []FileNode) {
				assert.True(t, nodes[0].IsDir)
				assert.False(t, nodes[1].IsDir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildFileTree(tt.input)
			assert.Equal(t, tt.expected, len(result))
			tt.checkFn(t, result)
		})
	}
}

func Test_SortChildren(t *testing.T) {
	tests := []struct {
		name    string
		node    *FileNode
		checkFn func(*testing.T, *FileNode)
	}{
		{
			name: "Success with directory first sort",
			node: &FileNode{
				Name:  "root",
				IsDir: true,
				Children: []FileNode{
					{Name: "file.txt", IsDir: false},
					{Name: "aaa", IsDir: true},
					{Name: "zzz.txt", IsDir: false},
					{Name: "bbb", IsDir: true},
				},
			},
			checkFn: func(t *testing.T, node *FileNode) {
				assert.Equal(t, 4, len(node.Children))
				assert.Equal(t, "aaa", node.Children[0].Name)
				assert.True(t, node.Children[0].IsDir)
				assert.Equal(t, "bbb", node.Children[1].Name)
				assert.True(t, node.Children[1].IsDir)
				assert.Equal(t, "file.txt", node.Children[2].Name)
				assert.False(t, node.Children[2].IsDir)
				assert.Equal(t, "zzz.txt", node.Children[3].Name)
				assert.False(t, node.Children[3].IsDir)
			},
		},
		{
			name: "Success with empty children",
			node: &FileNode{Name: "root", IsDir: true, Children: []FileNode{}},
			checkFn: func(t *testing.T, node *FileNode) {
				assert.Equal(t, 0, len(node.Children))
			},
		},
		{
			name: "Success with file node",
			node: &FileNode{Name: "file.txt", IsDir: false},
			checkFn: func(t *testing.T, node *FileNode) {
				assert.False(t, node.IsDir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortChildren(tt.node)
			tt.checkFn(t, tt.node)
		})
	}
}

func Test_FindChild(t *testing.T) {
	parent := &FileNode{
		Name:  "parent",
		IsDir: true,
		Children: []FileNode{
			{Name: "child1", IsDir: false},
			{Name: "child2", IsDir: false},
		},
	}

	tests := []struct {
		name      string
		childName string
		expected  bool
	}{
		{
			name:      "Success with existing child",
			childName: "child1",
			expected:  true,
		},
		{
			name:      "Failure with missing child",
			childName: "nonexistent",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findChild(parent, tt.childName)

			if tt.expected {
				assert.NotNil(t, result)
				assert.Equal(t, tt.childName, result.Name)
				result.Status = "M"
				assert.Equal(t, "M", parent.Children[0].Status)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
