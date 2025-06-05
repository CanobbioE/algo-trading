package signals

// Operation instructs the receiver what to do.
type Operation string

const (
	// Buy instructs the caller to buy the product.
	Buy Operation = "buy"
	// Sell instructs the caller to sell the product.
	Sell Operation = "sell"
	// NoOp instructs the caller to hold the product.
	NoOp Operation = "noop"
	// Setup instructs the caller that a setup is in progress, look for an opportunity.
	Setup Operation = "setup"
)
