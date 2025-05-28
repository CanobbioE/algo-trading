package signals

type Operation string

const (
	Buy   Operation = "buy"
	Sell  Operation = "sell"
	NoOp  Operation = "noop"
	Setup Operation = "setup"
)
