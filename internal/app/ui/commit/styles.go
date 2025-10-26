package commit

import "github.com/charmbracelet/lipgloss"

const (
	ColorPrimary  = lipgloss.Color("#7D56F4") // Purple - primary/focus color
	ColorBorder   = lipgloss.Color("8")       // Gray - borders and help text
	ColorMuted    = lipgloss.Color("7")       // Light gray - muted elements
	ColorAdded    = lipgloss.Color("10")      // Green - added files
	ColorModified = lipgloss.Color("11")      // Yellow - modified files
	ColorDeleted  = lipgloss.Color("9")       // Red - deleted files
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Padding(1, 2, 0, 2)

	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
		Foreground(ColorBorder).
		Padding(0, 2)
)
