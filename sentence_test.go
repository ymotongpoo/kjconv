package kjconv

import (
	"reflect"
	"testing"
)

func TestSplitSentences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "単一文（句点）",
			input:    "今日は晴れです。",
			expected: []string{"今日は晴れです。"},
		},
		{
			name:     "複数文（句点）",
			input:    "今日は晴れです。明日は雨です。",
			expected: []string{"今日は晴れです。", "明日は雨です。"},
		},
		{
			name:     "疑問符",
			input:    "元気ですか？はい、元気です。",
			expected: []string{"元気ですか？", "はい、元気です。"},
		},
		{
			name:     "感嘆符",
			input:    "すごい！本当にすごいです。",
			expected: []string{"すごい！", "本当にすごいです。"},
		},
		{
			name:     "混合",
			input:    "今日は晴れです。元気ですか？はい、元気です！",
			expected: []string{"今日は晴れです。", "元気ですか？", "はい、元気です！"},
		},
		{
			name:     "句読点なし",
			input:    "今日は晴れです",
			expected: []string{"今日は晴れです"},
		},
		{
			name:     "空文字列",
			input:    "",
			expected: []string{},
		},
		{
			name:     "空白のみ",
			input:    "   ",
			expected: []string{},
		},
		{
			name:     "句読点後の空白",
			input:    "今日は晴れです。  明日は雨です。",
			expected: []string{"今日は晴れです。", "明日は雨です。"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitSentences(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitSentences(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsQuotedText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "「」で囲まれた文",
			input:    "「今日は晴れです」",
			expected: true,
		},
		{
			name:     "『』で囲まれた文",
			input:    "『今日は晴れです』",
			expected: true,
		},
		{
			name:     "引用符なし",
			input:    "今日は晴れです",
			expected: false,
		},
		{
			name:     "開始引用符のみ",
			input:    "「今日は晴れです",
			expected: false,
		},
		{
			name:     "終了引用符のみ",
			input:    "今日は晴れです」",
			expected: false,
		},
		{
			name:     "混合引用符",
			input:    "「今日は晴れです』",
			expected: false,
		},
		{
			name:     "空文字列",
			input:    "",
			expected: false,
		},
		{
			name:     "引用符のみ",
			input:    "「」",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsQuotedText(tt.input)
			if result != tt.expected {
				t.Errorf("IsQuotedText(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
