package util

import "github.com/fatih/color"

var titleColor = color.New(color.FgYellow).Add(color.Underline)

func PrintSectionTitle(title string) {
	titleColor.Println(title)
}
