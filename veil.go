package veil

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/veil-services/veil-go/detectors"
)

// Exported errors for programmatic handling
var (
	// ErrEmptyInput is returned (or handled internally) when input is empty
	ErrEmptyInput = errors.New("veil: input cannot be empty")

	// ErrContextInvalid is returned when the restore context is nil or empty
	ErrContextInvalid = errors.New("veil: restore context is invalid or empty")
)

// RestoreContext stores the mapping required to restore original data.
// It must be JSON serializable.
type RestoreContext struct {
	// Data maps tokens to original values
	// e.g. "<<EMAIL_1>>" -> "john@example.com"
	Data map[string]string `json:"data"`
}

// Config defines the behavior of the Veil instance.
type Config struct {
	// Flags to enable standard detectors
	MaskEmail      bool
	MaskCPF        bool
	MaskCNPJ       bool
	MaskPhone      bool
	MaskCreditCard bool
	MaskIP         bool
	MaskUUID       bool

	// List of custom detectors registered by the user
	CustomDetectors []detectors.Detector

	// If true, the same value always gets the same token within the same request
	ConsistentTokenization bool
}

// Veil is the main engine.
type Veil struct {
	config    Config
	detectors []detectors.Detector
}

// New initializes a new Veil instance with the provided options.
// Example: veil.New(veil.WithEmail(), veil.WithCPF())
func New(opts ...Option) (*Veil, error) {
	// Default configuration
	cfg := Config{}

	// Apply options
	for _, opt := range opts {
		opt(&cfg)
	}

	v := &Veil{
		config:    cfg,
		detectors: make([]detectors.Detector, 0),
	}

	// Register standard detectors based on flags
	if cfg.MaskEmail {
		v.detectors = append(v.detectors, detectors.NewEmailDetector())
	}
	if cfg.MaskCPF {
		v.detectors = append(v.detectors, detectors.NewCPFDetector())
	}
	if cfg.MaskCNPJ {
		v.detectors = append(v.detectors, detectors.NewCNPJDetector())
	}
	if cfg.MaskCreditCard {
		v.detectors = append(v.detectors, detectors.NewCreditCardDetector())
	}
	if cfg.MaskIP {
		v.detectors = append(v.detectors, detectors.NewIPDetector())
	}
	if cfg.MaskPhone {
		v.detectors = append(v.detectors, detectors.NewPhoneDetector())
	}
	if cfg.MaskUUID {
		v.detectors = append(v.detectors, detectors.NewUUIDDetector())
	}

	// Register custom detectors
	v.detectors = append(v.detectors, cfg.CustomDetectors...)

	return v, nil
}

// Mask processes the text and returns the safe version + restoration context.
func (v *Veil) Mask(input string) (string, *RestoreContext, error) {
	if input == "" {
		return "", &RestoreContext{Data: make(map[string]string)}, nil
	}

	// 1. Scan: Collect all matches from all detectors
	var allMatches []detectors.Match
	for _, d := range v.detectors {
		matches := d.Scan(input)
		allMatches = append(allMatches, matches...)
	}

	// 2. Resolve conflicts (Greediest Match Wins)
	if len(allMatches) == 0 {
		return input, &RestoreContext{Data: make(map[string]string)}, nil
	}

	finalMatches := resolveOverlaps(allMatches)

	// 3. Tokenization and String Construction
	var sb strings.Builder
	// Pre-allocate builder size to avoid reallocations (heuristic: input size)
	sb.Grow(len(input))

	ctx := &RestoreContext{
		Data: make(map[string]string),
	}

	typeCounters := make(map[detectors.PIIType]int)
	valueCache := make(map[string]string)

	lastIndex := 0

	// Sort matches by start index for linear construction
	sort.Slice(finalMatches, func(i, j int) bool {
		return finalMatches[i].StartIndex < finalMatches[j].StartIndex
	})

	for _, m := range finalMatches {
		// Add non-masked text before the match
		if m.StartIndex > lastIndex {
			sb.WriteString(input[lastIndex:m.StartIndex])
		}

		// Determine token
		var token string
		if v.config.ConsistentTokenization {
			if existingToken, exists := valueCache[m.Value]; exists {
				token = existingToken
			}
		}

		if token == "" {
			typeCounters[m.Type]++
			count := typeCounters[m.Type]
			token = fmt.Sprintf("<<%s_%d>>", string(m.Type), count)

			if v.config.ConsistentTokenization {
				valueCache[m.Value] = token
			}
		}

		// Save to context and replace
		ctx.Data[token] = m.Value
		sb.WriteString(token)

		lastIndex = m.EndIndex
	}

	// Add remaining string
	if lastIndex < len(input) {
		sb.WriteString(input[lastIndex:])
	}

	return sb.String(), ctx, nil
}

// Restore takes the masked text and the original context to retrieve data.
func (v *Veil) Restore(maskedInput string, ctx *RestoreContext) (string, error) {
	if ctx == nil || len(ctx.Data) == 0 {
		return maskedInput, ErrContextInvalid
	}

	// Fast path: if no tokens marker exists, return immediately
	if !strings.Contains(maskedInput, "<<") {
		return maskedInput, nil
	}

	var sb strings.Builder
	sb.Grow(len(maskedInput)) // Optimistic allocation

	n := len(maskedInput)
	i := 0

	for i < n {
		// Find start of potential token
		if maskedInput[i] == '<' && i+1 < n && maskedInput[i+1] == '<' {
			// Found '<<', look for closing '>>'
			closing := strings.Index(maskedInput[i:], ">>")
			if closing != -1 {
				closingIndex := i + closing + 2 // include '>>' length
				tokenCandidate := maskedInput[i:closingIndex]

				// Check if this token exists in our context
				if originalValue, exists := ctx.Data[tokenCandidate]; exists {
					sb.WriteString(originalValue)
					i = closingIndex
					continue
				}
				// If not in context (or false positive like <<Shift>>), treat as normal text
			}
		}

		sb.WriteByte(maskedInput[i])
		i++
	}

	return sb.String(), nil
}

// Sanitize is a helper for logs that masks any input and returns the safe string.
// Unlike Mask(), it discards the restoration context (One-way mask).
func (v *Veil) Sanitize(input interface{}) interface{} {
	strInput := fmt.Sprintf("%v", input)
	masked, _, _ := v.Mask(strInput)
	return masked
}

// resolveOverlaps removes matches that are contained within larger matches or conflict.
// Strategy: Prioritize the largest match (Greediest). On size tie, prioritize Score.
func resolveOverlaps(matches []detectors.Match) []detectors.Match {
	if len(matches) <= 1 {
		return matches
	}

	// Sort by StartIndex asc, and then by Length desc (to pick the longest first)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].StartIndex != matches[j].StartIndex {
			return matches[i].StartIndex < matches[j].StartIndex
		}
		// If start at same position, longest comes first
		lenI := matches[i].EndIndex - matches[i].StartIndex
		lenJ := matches[j].EndIndex - matches[j].StartIndex
		return lenI > lenJ
	})

	var result []detectors.Match
	var lastMatch detectors.Match

	for i, m := range matches {
		if i == 0 {
			result = append(result, m)
			lastMatch = m
			continue
		}

		// Check overlap
		// If current match starts before previous match ends, there is a conflict.
		if m.StartIndex < lastMatch.EndIndex {
			// Conflict detected.
			// Since we sorted by length, 'lastMatch' is likely the "best" one.
			// We ignore 'm' as it is "inside" or conflicting with a prioritized match.
			continue
		}

		result = append(result, m)
		lastMatch = m
	}

	return result
}
