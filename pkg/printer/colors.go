package printer

// Color defines an output color.
type Color int

const (
	// None sets the output to the default color.
	None Color = iota - 1
	// Red sets the output to red.
	Red
	// Green sets the output to green.
	Green
	// Yellow sets the output to yellow.
	Yellow
	// Blue sets the output to blue.
	Blue
	// White sets the output to white.
	White
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorWhite  = "\033[37m"
)

var colorMap = map[Color]string{
	Red:    colorRed,
	Green:  colorGreen,
	Yellow: colorYellow,
	Blue:   colorBlue,
	White:  colorWhite,
}

// WrapInColor adds color to a string without printing it.
func WrapInColor(s string, color Color) string {
	if c, ok := colorMap[color]; ok {
		return c + s + colorReset
	}
	return s
}
