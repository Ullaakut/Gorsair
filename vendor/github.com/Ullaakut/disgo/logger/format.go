package logger

import (
	"github.com/fatih/color"
)

var (
	// Success colors a message in bold green to represent success.
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()

	// Failure colors a message in bold red to represent failure.
	Failure = color.New(color.FgRed, color.Bold).SprintFunc()

	// Trace colors a message in faint white (usually rendered in gray)
	// to represent a log of low importance for the user.
	Trace = color.New(color.FgHiWhite, color.Faint).SprintFunc()

	// Important colors a message in bold to represent an important
	// information.
	Important = color.New(color.Bold).SprintFunc()

	// Link colors a message in underlined blue to represent a clickable link.
	Link = color.New(color.FgBlue, color.Underline).SprintFunc()
)
