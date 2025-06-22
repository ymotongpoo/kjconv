package kjconv

import (
	"testing"
)

func TestNewConverter(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}
	if converter == nil {
		t.Fatal("NewConverter() returned nil")
	}
	if converter.tokenizer == nil {
		t.Fatal("NewConverter() did not initialize tokenizer")
	}
}

func TestConvert_CasualToPolite(t *testing.T) {
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
			name:     "名詞+だ → です",
			input:    "今日は晴れだ。",
			expected: "今日は晴れです。",
		},
		{
			name:     "動詞基本形 → ます形",
			input:    "本を読む。",
			expected: "本を読みます。",
		},
		{
			name:     "形容詞 → です",
			input:    "この花は美しい。",
			expected: "この花は美しいです。",
		},
		{
			name:     "五段動詞（カ行）",
			input:    "字を書く。",
			expected: "字を書きます。",
		},
		{
			name:     "五段動詞（サ行）",
			input:    "話す。",
			expected: "話します。",
		},
		{
			name:     "一段動詞",
			input:    "ご飯を食べる。",
			expected: "ご飯を食べます。",
		},
		{
			name:     "カ変動詞",
			input:    "学校に来る。",
			expected: "学校に来ます。",
		},
		{
			name:     "サ変動詞",
			input:    "勉強する。",
			expected: "勉強します。",
		},
		{
			name:     "複数文",
			input:    "今日は晴れだ。本を読む。",
			expected: "今日は晴れです。本を読みます。",
		},
		// TODO: 引用文の処理は後で実装
		// {
		// 	name:     "引用文は変換しない",
		// 	input:    "彼は「今日は晴れだ」と言った。",
		// 	expected: "彼は「今日は晴れだ」と言いました。",
		// },
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

func TestConvert_PoliteToCasual(t *testing.T) {
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
			name:     "です → だ",
			input:    "今日は晴れです。",
			expected: "今日は晴れだ。",
		},
		{
			name:     "ます → 基本形",
			input:    "本を読みます。",
			expected: "本を読む。",
		},
		{
			name:     "形容詞+です → 形容詞",
			input:    "この花は美しいです。",
			expected: "この花は美しい。",
		},
		{
			name:     "複数文",
			input:    "今日は晴れです。本を読みます。",
			expected: "今日は晴れだ。本を読む。",
		},
		// TODO: 引用文の処理は後で実装
		// {
		// 	name:     "引用文は変換しない",
		// 	input:    "彼は「今日は晴れです」と言った。",
		// 	expected: "彼は「今日は晴れです」と言った。",
		// },
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

func TestConvert_InvalidMode(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	_, err = converter.Convert("test", ConversionMode(999))
	if err == nil {
		t.Error("Convert() should fail with invalid mode")
	}
}

func TestConvert_EmptyString(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	result, err := converter.Convert("", CasualToPolite)
	if err != nil {
		t.Errorf("Convert() failed with empty string: %v", err)
	}
	if result != "" {
		t.Errorf("Convert() with empty string = %q, expected empty string", result)
	}
}
