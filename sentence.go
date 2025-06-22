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
		// No sentence-ending punctuation found, return the whole text
		return []string{strings.TrimSpace(text)}
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
