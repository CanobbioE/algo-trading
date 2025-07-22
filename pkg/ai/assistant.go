package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

//nolint:lll
const (
	// todo: consider making these configurable
	premise      = `You are a expert financial advisor. Your job is to analyze the technical data and recent market context of a stock or ETF, and recommend a single, clear action: BUY, HOLD, or SELL.`
	instructions = `You must:
1. Analyze the provided technical signals.
2. Search for and incorporate up-to-date **news, events, and sentiment** about the asset from the web. Consider macroeconomic news, earnings, geopolitical events, or sector-wide changes.
3. Weigh your recommendation according to the provided RISK TOLERANCE:
   - LOW: Avoid risky trades; prioritize capital preservation.
   - MEDIUM: Take moderate risk for reasonable returns.
   - HIGH: Accept volatility in pursuit of higher gains.`
	outputRules = `Your response must:
- Be **no more than 2 concise sentences**.
- Provide a direct recommendation: BUY, HOLD, or SELL.
- Justify the recommendation based on both technicals and recent events.
- In case your recommendation is to:
	- BUY: include a suggested purchase price (e.g. near current BID, near daily lowest, market value, etc)
	- SELL: include a suggested sell price (e.g. near current ASK, near daily highest, market value, etc)
- Do **not** repeat the raw input.
- Think through your reasoning silently; return only the final result.`
)

// Config for the AI assistant.
type Config struct {
	APIKey string `json:"api_key"`
}

// Assistant uses Gemini API to assist with stock picking.
type Assistant struct {
	cli *genai.Client
}

// NewAssistant creates a new AI Assistant using Gemini, with the given config.
func NewAssistant(ctx context.Context, cfg *Config) (*Assistant, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &Assistant{cli: client}, nil
}

// Analyse calls the underlying Gemini client and asks it to run a technical analysis on the given stock,
// based on the user maximum risk tolerance and the ticker analysis.
func (a *Assistant) Analyse(ctx context.Context, input, risk, ticker string) (string, error) {
	modelCfg := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: premise},
				{Text: instructions},
				{Text: outputRules},
			},
		},
		Temperature:      genai.Ptr[float32](0.2),
		CandidateCount:   1,
		ResponseMIMEType: "text/plain",
	}

	parts := []*genai.Part{
		{Text: "Assume the maximum risk the user is willing to take is " + risk},
		{Text: "The ticker for the asset is: " + ticker},
		{Text: "Input:\n" + input},
	}

	// TODO: consider caching
	resp, err := a.cli.Models.GenerateContent(ctx, "gemini-2.5-pro", []*genai.Content{{Parts: parts}}, modelCfg)
	if err != nil {
		return "", fmt.Errorf("ai assistant: failed to generate content: %w", err)
	}

	return resp.Text(), nil
}
