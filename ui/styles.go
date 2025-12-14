package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// SplitFlapStyles contains all the styling for the split-flap display
type SplitFlapStyles struct {
	Background   lipgloss.Style
	Text         lipgloss.Style
	Header       lipgloss.Style
	StatusLight  func(color string) lipgloss.Style
	AirportLabel lipgloss.Style
	PageInfo     lipgloss.Style
	Error        lipgloss.Style
}

// NewSplitFlapStyles creates a new set of split-flap styles
func NewSplitFlapStyles() *SplitFlapStyles {
	// Retro split-flap color scheme
	bgColor := lipgloss.Color("#1a1a1a") // Dark gray/black background
	textColor := lipgloss.Color("#f0f0f0") // High contrast white text
	headerColor := lipgloss.Color("#ffffff") // White headers
	errorColor := lipgloss.Color("#ff0000") // Red for errors

	return &SplitFlapStyles{
		Background: lipgloss.NewStyle().
			Background(bgColor).
			Foreground(textColor).
			Padding(1, 2),

		Text: lipgloss.NewStyle().
			Foreground(textColor),

		Header: lipgloss.NewStyle().
			Foreground(headerColor).
			Bold(true).
			Underline(true),

		StatusLight: func(color string) lipgloss.Style {
			var statusColor lipgloss.Color
			switch color {
			case "green":
				statusColor = lipgloss.Color("#00ff00")
			case "yellow":
				statusColor = lipgloss.Color("#ffff00")
			case "orange":
				statusColor = lipgloss.Color("#ff8800")
			case "red":
				statusColor = lipgloss.Color("#ff0000")
			default:
				statusColor = lipgloss.Color("#ffffff")
			}
			return lipgloss.NewStyle().
				Foreground(statusColor).
				Bold(true)
		},

		AirportLabel: lipgloss.NewStyle().
			Foreground(headerColor).
			Bold(true).
			MarginBottom(1),

		PageInfo: lipgloss.NewStyle().
			Foreground(textColor).
			MarginTop(1),

		Error: lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true),
	}
}

