package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseGitStatus(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		expectedFiles int
	}{
		{name: "Empty status", status: "", expectedFiles: 0},
		{name: "Single file added", status: "A\tCLAUDE.md", expectedFiles: 1},
		{name: "Multiple files", status: "A\tCLAUDE.md\nM\tREADME.md\nD\tcmd/main_test.go", expectedFiles: 3},
		{name: "Files with spaces in names", status: "A\tfile with spaces.md\nM\tanother file.go", expectedFiles: 2},
		{name: "Renamed file", status: "R100\told.go\tnew.go", expectedFiles: 1},
		{name: "Copied file", status: "C064\toriginal.go\tcopy.go", expectedFiles: 1},
		{name: "Invalid line", status: "INVALID", expectedFiles: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := ParseGitStatus(tt.status)
			fileCount := CountFiles(tree)
			assert.Equal(t, tt.expectedFiles, fileCount)
		})
	}
}

func Test_parseGitStatus_TreeStructure(t *testing.T) {
	status := "A\tCLAUDE.md\nM\tinternal/app/cli/view.go\nD\tcmd/main_test.go"
	tree := ParseGitStatus(status)

	assert.NotNil(t, tree)
	assert.Equal(t, 3, CountFiles(tree))

	rendered := RenderTree(tree)
	assert.Contains(t, rendered, "CLAUDE.md")
	assert.Contains(t, rendered, "view.go")
	assert.Contains(t, rendered, "main_test.go")
}

func Test_countFiles(t *testing.T) {
	status := "A\tfile1.go\nM\tdir1/file2.go"
	tree := ParseGitStatus(status)

	count := CountFiles(tree)
	assert.Equal(t, 2, count)
}

func Test_renderTree(t *testing.T) {
	status := "A\tCLAUDE.md\nM\tREADME.md"
	tree := ParseGitStatus(status)

	result := RenderTree(tree)
	assert.Contains(t, result, "CLAUDE.md")
	assert.Contains(t, result, "README.md")
}

func Test_parseGitStatus_RenameWithScore(t *testing.T) {
	status := "R100\told.go\tnew.go"
	tree := ParseGitStatus(status)

	assert.Equal(t, 1, CountFiles(tree))

	rendered := RenderTree(tree)
	assert.Contains(t, rendered, "new.go", "Should show new filename for rename")
	assert.NotContains(t, rendered, "old.go", "Should not show old filename for rename")
	assert.Contains(t, rendered, "R ", "Should display R status")
}

func Test_parseGitStatus_CopyWithScore(t *testing.T) {
	status := "C064\toriginal.go\tcopy.go"
	tree := ParseGitStatus(status)

	assert.Equal(t, 1, CountFiles(tree))

	rendered := RenderTree(tree)
	assert.Contains(t, rendered, "copy.go", "Should show new filename for copy")
	assert.NotContains(t, rendered, "original.go", "Should not show original filename in tree")
	assert.Contains(t, rendered, "C ", "Should display C status")
}
