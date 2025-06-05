package printer

// TODO: move to a separate pkg with colors.

import (
	"fmt"
	"os"
)

// Printer defines an interface to show output.
type Printer interface {
	// Printf prints a message formatted.
	Printf(format string, a ...any)
	// Println prints a message followed by a new line.
	Println(msg string)
	// PrintColored prints the message as colored output, if supported.
	PrintColored(c Color, format string, a ...any)
	// Reset the output, if supported.
	Reset()
	// CleanLine tells the printer to clean the line.
	CleanLine() Printer
}

// Standard prints to os.Stdout.
type Standard struct {
	cleanLine bool
}

// Printf the message to std output.
func (s *Standard) Printf(format string, a ...any) {
	if s.cleanLine {
		format += "\033[K"
	}
	_, _ = fmt.Fprintf(os.Stdout, format, a...)
}

// Println prints a message followed by a new line.
func (*Standard) Println(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "%s", msg+"\n")
}

// PrintColored prints the message in the given color.
func (s *Standard) PrintColored(color Color, format string, a ...any) {
	if c, ok := colorMap[color]; ok {
		s.Printf(c+format+colorReset, a...)
		return
	}
	s.Printf(format, a...)
}

// Reset the terminal, if supported.
func (s *Standard) Reset() {
	if s.cleanLine {
		s.Printf("\033[H")
		return
	}
	s.Printf("\033[H\033[2J")
}

// CleanLine returns the printer with the clean line option enabled.
func (*Standard) CleanLine() Printer {
	return &Standard{cleanLine: true}
}
