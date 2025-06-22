package kjconv

import (
	"testing"
)

func TestConjunctionConversion_CasualToPolite(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "だから → ですから",
			input:    "だから今日は晴れだ。",
			expected: "ですから今日は晴れです。",
		},
		{
			name:     "だが → ですが",
			input:    "だが今日は晴れだ。",
			expected: "ですが今日は晴れです。",
		},
		{
			name:     "文中のだからは変換しない",
			input:    "今日は晴れだから嬉しい。",
			expected: "今日は晴れだから嬉しいです。",
		},
		{
			name:     "接続詞なしの文",
			input:    "今日は晴れだ。",
			expected: "今日は晴れです。",
		},
		{
			name:     "複数文での接続詞",
			input:    "だから今日は晴れだ。明日は雨だ。",
			expected: "ですから今日は晴れです。明日は雨です。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input, CasualToPolite)
			if err != nil {
				t.Errorf("Convert() failed: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Convert() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestConjunctionConversion_PoliteToCasual(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ですから → だから",
			input:    "ですから今日は晴れです。",
			expected: "だから今日は晴れだ。",
		},
		{
			name:     "ですが → だが",
			input:    "ですが今日は晴れです。",
			expected: "だが今日は晴れだ。",
		},
		{
			name:     "接続詞なしの文",
			input:    "今日は晴れです。",
			expected: "今日は晴れだ。",
		},
		{
			name:     "複数文での接続詞",
			input:    "ですから今日は晴れです。明日は雨です。",
			expected: "だから今日は晴れだ。明日は雨だ。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.Convert(tt.input, PoliteToCasual)
			if err != nil {
				t.Errorf("Convert() failed: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Convert() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
