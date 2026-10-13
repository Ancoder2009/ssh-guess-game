package main

import "charm.land/lipgloss/v2"

var TextBoxStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Green).Foreground(lipgloss.Black).BorderStyle(lipgloss.DoubleBorder()).
	Padding(0).Align(lipgloss.Center)

var GreenTextStyle = lipgloss.NewStyle().
	Bold(true).Foreground(lipgloss.Green)
