package printer

// Composite prints to all underlying printers.
type Composite struct {
	printers []Printer
}

// NewCompositePrinter creates a Printer that collects multiple printers.
func NewCompositePrinter(printers ...Printer) *Composite {
	return &Composite{
		printers: printers,
	}
}

// Printf prints a message formatted.
func (c *Composite) Printf(format string, a ...any) {
	for _, p := range c.printers {
		p.Printf(format, a...)
	}
}

// Println prints a message followed by a new line.
func (c *Composite) Println(msg string) {
	for _, p := range c.printers {
		p.Println(msg)
	}
}

// PrintColored prints the message as colored output, if supported.
func (c *Composite) PrintColored(clr Color, format string, a ...any) {
	for _, p := range c.printers {
		p.PrintColored(clr, format, a...)
	}
}

// Reset the output, if supported.
func (c *Composite) Reset() {
	for _, p := range c.printers {
		p.Reset()
	}
}

// CleanLine tells the printer to clean the line.
func (c *Composite) CleanLine() Printer {
	for _, p := range c.printers {
		p.CleanLine()
	}
	return c
}
