package components

import "github.com/charmbracelet/lipgloss"

// Material Design 3 Color Palette
const (
	ColorPrimary = "#7D56F4"
	ColorSuccess = "#04B575"
	ColorError   = "#FF5252"
	ColorWarning = "#FFA726"
	ColorMuted   = "#9E9E9E"
	ColorText    = "#E0E0E0"
	ColorTextDim = "#666666"
	ColorBorder  = "8"
	ColorInfo    = "#5FA8D3"
	ColorTrace   = "#999999"
	ColorWhite   = "#FFFFFF"
)

// Material Design 3 Typography System

// Headline - High-emphasis text for major sections
var (
	HeadlineLarge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorPrimary)).MarginTop(1)
	HeadlineMedium = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorPrimary))
)

// Title - Medium-emphasis text for titles and subtitles
var (
	TitleLarge  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorSuccess))
	TitleMedium = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorSuccess))
)

// Body - Main content text
var (
	BodyLarge  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
	BodyMedium = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
)

// Label - Small text for labels, captions, hints
var (
	LabelLarge  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Italic(true).MarginTop(2)
	LabelMedium = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted))
	LabelSmall  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Faint(true)
)

// Semantic Styles - Map typography to use cases
var (
	SectionHeader = HeadlineLarge.Margin(1, 2)
	HelpText      = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Margin(1, 2)
	KeyHint       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorPrimary))
	CommandName   = TitleMedium
	ExampleCode   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorWarning))

	StatusSuccess = BodyMedium.Foreground(lipgloss.Color(ColorSuccess))
	StatusError   = BodyMedium.Foreground(lipgloss.Color(ColorError)).Bold(true)
	StatusWarning = BodyMedium.Foreground(lipgloss.Color(ColorWarning))

	SpinnerStyle  = lipgloss.NewStyle().Margin(1, 0, 0, 2).Foreground(lipgloss.Color(ColorPrimary))
	TextareaStyle = lipgloss.NewStyle().Margin(1, 0, 0, 0)
)

// Component-specific styles
var (
	ModeBadgeNormalStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite)).
				Background(lipgloss.Color(ColorPrimary)).
				Bold(true).
				Padding(0, 1)

	ModeBadgeInsertStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite)).
				Background(lipgloss.Color(ColorSuccess)).
				Bold(true).
				Padding(0, 1)
)

// Log level styles
var (
	LogDebugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextDim))
	LogInfoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo))
	LogWarnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarning))
	LogErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError))
	LogTraceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTrace))
)

// PanelStyle creates a bordered panel with consistent styling
func PanelStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorBorder)).
		Width(width-2).
		Height(height-2).
		Padding(0, 1).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top)
}
