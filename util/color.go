package util

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	Purple      = lipgloss.Color("99")
	PurpleStyle = lipgloss.NewStyle().Foreground(Purple)
	Gray        = lipgloss.Color("245")
	LightGray   = lipgloss.Color("241")
	ErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	TitleColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#fdc981")).Underline(true)
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fdc981")).
			Bold(true).
			Align(lipgloss.Center)

	cellStyle = lipgloss.NewStyle().Padding(0, 1)
)

func StyledTable() *table.Table {
	return table.New().
		Headers("username", "UID", "GID", "shell").
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(Gray))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		})
}
