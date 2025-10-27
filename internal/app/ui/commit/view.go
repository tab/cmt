package commit

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the commit UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing…"
	}

	if m.stateMachine.WorkflowMode() == Editing {
		return m.renderEditMode()
	}

	return m.renderNormalMode()
}

// renderNormalMode renders the view in normal (viewing) mode
func (m Model) renderNormalMode() string {
	var sections []string

	titleText := ">_ commit message"

	if m.stateMachine.ViewPane() == AppLogsPane {
		titleText = ">_ logs"
	}

	if m.stateMachine.IsGenerating() {
		titleText = m.spinner.View() + " loading…"
	}

	title := titleStyle.Render(titleText)
	sections = append(sections, title)
	sections = append(sections, "")

	if m.stateMachine.IsGenerating() {
		return strings.Join(sections, "\n")
	}

	if m.stateMachine.ViewPane() == AppLogsPane {
		logsPanel := panelStyle.
			Width(m.width - 2).
			Height(m.height - 10).
			Render(m.viewport.View())
		sections = append(sections, logsPanel)
	} else {
		treeWidth := m.width / 3
		messageWidth := m.width - treeWidth - 4

		treePanelStyle := panelStyle
		messagePanelStyle := panelStyle

		if m.focusPane == TreeFocus {
			treePanelStyle = treePanelStyle.BorderForeground(ColorPrimary)
		} else {
			messagePanelStyle = messagePanelStyle.BorderForeground(ColorPrimary)
		}

		treePanel := treePanelStyle.
			Width(treeWidth).
			Height(m.height - 10).
			Render(m.treeViewport.View())

		messagePanel := messagePanelStyle.
			Width(messageWidth).
			Height(m.height - 10).
			Render(m.viewport.View())

		panels := lipgloss.JoinHorizontal(lipgloss.Top, treePanel, messagePanel)
		sections = append(sections, panels)
	}

	sections = append(sections, "")

	helpView := m.help.View(m.keys)
	sections = append(sections, helpStyle.Render(helpView))

	return strings.Join(sections, "\n")
}

// renderEditMode renders the view in edit mode
func (m Model) renderEditMode() string {
	var sections []string

	title := titleStyle.Render(">_ edit")
	sections = append(sections, title)
	sections = append(sections, "")

	sections = append(sections, m.textarea.View())

	return strings.Join(sections, "\n")
}

// renderFileTree renders the file tree as a string
func (m Model) renderFileTree() string {
	if len(m.state.Files) == 0 {
		return "No files staged"
	}

	var sb strings.Builder
	sb.WriteString(m.renderTreeNodes(m.state.Files, "", true))

	return sb.String()
}

// renderTreeNodes recursively renders tree nodes
func (m Model) renderTreeNodes(nodes []FileNode, indent string, isRoot bool) string {
	var sb strings.Builder

	for i, node := range nodes {
		isLast := i == len(nodes)-1

		if !isRoot {
			if isLast {
				sb.WriteString(indent + "└── ")
			} else {
				sb.WriteString(indent + "├── ")
			}
		}

		statusIcon := ""
		statusColor := ColorMuted
		switch node.Status {
		case "A":
			statusIcon = "[+] "
			statusColor = ColorAdded
		case "M":
			statusIcon = "[~] "
			statusColor = ColorModified
		case "D":
			statusIcon = "[-] "
			statusColor = ColorDeleted
		}

		if statusIcon != "" {
			sb.WriteString(lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon))
		}

		name := node.Name
		if node.IsDir {
			name += "/"
		}
		sb.WriteString(name)
		sb.WriteString("\n")

		if node.IsDir && len(node.Children) > 0 {
			childIndent := indent
			if !isRoot {
				if isLast {
					childIndent += "  "
				} else {
					childIndent += "│ "
				}
			}
			sb.WriteString(m.renderTreeNodes(node.Children, childIndent, false))
		}
	}

	return sb.String()
}
