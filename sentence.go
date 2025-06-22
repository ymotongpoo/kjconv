package kjconv

import (
	"regexp"
	"strings"
)

// SplitSentences splits text into sentences based on sentence-ending punctuation.
// It splits on 句点（。）, 疑問符（？）, and 感嘆符（！）.
func SplitSentences(text string) []string {
	// Regular expression to split on sentence-ending punctuation
	// Matches 。, ？, ！ followed by optional whitespace
	re := regexp.MustCompile(`[。？！]\s*`)
	
	// Find all matches to preserve the punctuation
	matches := re.FindAllStringIndex(text, -1)
	
	if len(matches) == 0 {
		// No sentence-ending punctuation found
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return []string{}
		}
		return []string{trimmed}
	}
	
	var sentences []string
	start := 0
	
	for _, match := range matches {
		end := match[1] // Include the punctuation and any following whitespace
		sentence := strings.TrimSpace(text[start:end])
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
		start = end
	}
	
	// Add any remaining text after the last punctuation
	if start < len(text) {
		remaining := strings.TrimSpace(text[start:])
		if remaining != "" {
			sentences = append(sentences, remaining)
		}
	}
	
	return sentences
}

// IsQuotedText checks if the text is within quotation marks (「」or『』).
// Quoted text should not be converted according to the specification.
func IsQuotedText(text string) bool {
	// Check for 「」quotation marks
	if strings.HasPrefix(text, "「") && strings.HasSuffix(text, "」") {
		return true
	}
	
	// Check for 『』quotation marks
	if strings.HasPrefix(text, "『") && strings.HasSuffix(text, "』") {
		return true
	}
	
	return false
}

// ContainsQuotedText checks if the text contains quoted sections that should be preserved.
func ContainsQuotedText(text string) bool {
	return strings.Contains(text, "「") || strings.Contains(text, "」") ||
		   strings.Contains(text, "『") || strings.Contains(text, "』")
}

// ProcessTextWithQuotes processes text while preserving quoted sections.
func ProcessTextWithQuotes(text string, processor func(string) (string, error)) (string, error) {
	if !ContainsQuotedText(text) {
		return processor(text)
	}
	
	// Simple approach: split by quotes and process non-quoted parts
	var result strings.Builder
	var inQuote bool
	var quoteChar rune
	
	for _, r := range text {
		if !inQuote && (r == '「' || r == '『') {
			inQuote = true
			quoteChar = r
			result.WriteRune(r)
		} else if inQuote && ((r == '」' && quoteChar == '「') || (r == '』' && quoteChar == '『')) {
			inQuote = false
			result.WriteRune(r)
		} else if inQuote {
			result.WriteRune(r)
		} else {
			// Process non-quoted character
			processed, err := processor(string(r))
			if err != nil {
				return "", err
			}
			result.WriteString(processed)
		}
	}
	
	return result.String(), nil
}
