package printer

import (
	"fmt"
	"strings"
)

// String prints to a given output string.
type String struct {
	output *strings.Builder
}

// NewStringsPrinter creates a new Printer that stores its output into a string builder.
func NewStringsPrinter(out *strings.Builder) *String {
	return &String{
		output: out,
	}
}

// Printf prints a message formatted.
func (b *String) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(b.output, format, a...)
}

// Println prints a message followed by a new line.
func (b *String) Println(msg string) {
	b.output.WriteString(msg + "\n")
}

// PrintColored prints the message as colored output, if supported.
func (b *String) PrintColored(_ Color, format string, a ...any) {
	_, _ = fmt.Fprintf(b.output, format, a...)
}

// Reset the output, if supported.
func (*String) Reset() {}

// CleanLine tells the printer to clean the line.
func (b *String) CleanLine() Printer {
	return b
}

// CleanOutput removes all color-related formatting from the output.
func CleanOutput(out *strings.Builder) string {
	s := out.String()
	for _, c := range colorMap {
		s = strings.ReplaceAll(s, c, "")
	}
	return strings.ReplaceAll(s, colorReset, "")
}
