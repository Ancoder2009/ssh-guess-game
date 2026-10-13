package main

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const rawBanner = `
███████╗███████╗██╗  ██╗     ██████╗ ██╗   ██╗███████╗███████╗███████╗     ██████╗  █████╗ ███╗   ███╗███████╗
██╔════╝██╔════╝██║  ██║    ██╔════╝ ██║   ██║██╔════╝██╔════╝██╔════╝    ██╔════╝ ██╔══██╗████╗ ████║██╔════╝
███████╗███████╗███████║    ██║  ███╗██║   ██║█████╗  ███████╗███████╗    ██║  ███╗███████║██╔████╔██║█████╗  
╚════██║╚════██║██╔══██║    ██║   ██║██║   ██║██╔══╝  ╚════██║╚════██║    ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝  
███████║███████║██║  ██║    ╚██████╔╝╚██████╔╝███████╗███████║███████║    ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗
╚══════╝╚══════╝╚═╝  ╚═╝     ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝     ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝
`

var BANNER = gradientBanner(rawBanner)

func gradientBanner(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = gradientLine(line)
	}
	return strings.Join(lines, "\n")
}

func gradientLine(s string) string {
	visibleRunes := 0
	for _, r := range s {
		if r != ' ' && r != '\r' {
			visibleRunes++
		}
	}
	if visibleRunes == 0 {
		return s
	}

	colors := lipgloss.Blend1D(visibleRunes,
		lipgloss.Color("#a779ed"),
		lipgloss.Color("#89e9f4"),
	)

	var b strings.Builder
	colorIndex := 0
	for _, r := range s {
		switch r {
		case ' ', '\r':
			b.WriteRune(r)
			continue
		default:
			b.WriteString(lipgloss.NewStyle().Foreground(colors[colorIndex]).Render(string(r)))
			colorIndex++
		}
	}

	return b.String()
}
