package doctor

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

var (
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow
	failStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
	dimStyle  = lipgloss.NewStyle().Faint(true)
)

// Print writes the report as one line per check, with the fix under every check that is not passing.
func Print(w io.Writer, r Report) {
	width := 0
	for _, c := range r.Checks {
		width = max(width, len(c.Name))
	}
	for _, c := range r.Checks {
		mark := passStyle.Render("✓")
		switch c.Status {
		case StatusWarn:
			mark = warnStyle.Render("!")
		case StatusFail:
			mark = failStyle.Render("✗")
		}
		line := fmt.Sprintf("%s %-*s  %s", mark, width, c.Name, c.Found)
		if c.Required != "" && c.Status != StatusPass {
			line += dimStyle.Render(" (needs " + c.Required + ")")
		}
		fmt.Fprintln(w, line)
		if c.Fix != "" {
			fmt.Fprintf(w, "  %*s  %s\n", width, "", dimStyle.Render(c.Fix))
		}
	}
}
