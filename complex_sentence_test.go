package kjconv

import (
	"testing"
)

func TestComplexSentenceConversion_CasualToPolite(t *testing.T) {
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
			name:     "引用文と複合文1（現在の実装レベル）",
			input:    "彼は「明日雨だ」と言った。",
			expected: "彼は「明日雨だ」と言いました。",
		},
		{
			name:     "複合文2（現在の実装レベル）",
			input:    "お腹は空いている。",
			expected: "お腹は空いています。",
		},
		{
			name:     "複数の動詞を含む文（文末のみ変換）",
			input:    "朝起きて、朝食を食べて、学校に行く。",
			expected: "朝起きて、朝食を食べて、学校に行きます。",
		},
		{
			name:     "条件文",
			input:    "雨が降ったら、家にいる。",
			expected: "雨が降ったら、家にいます。",
		},
		{
			name:     "理由を表す文",
			input:    "疲れているから、早く寝る。",
			expected: "疲れているから、早く寝ます。",
		},
		{
			name:     "単純な否定文",
			input:    "本を読まない。",
			expected: "本を読みません。",
		},
		{
			name:     "単純な過去文",
			input:    "昨日映画を見た。",
			expected: "昨日映画を見ました。",
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

func TestComplexSentenceConversion_PoliteToCasual(t *testing.T) {
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
			name:     "引用文と複合文1（現在の実装レベル）",
			input:    "彼は「明日は雨だ」と言いました。",
			expected: "彼は「明日は雨だ」と言った。",
		},
		{
			name:     "複合文2（現在の実装レベル）",
			input:    "お腹が空いています。",
			expected: "お腹が空いている。",
		},
		{
			name:     "複数の動詞を含む文",
			input:    "朝起きて、朝食を食べて、学校に行きます。",
			expected: "朝起きて、朝食を食べて、学校に行く。",
		},
		{
			name:     "条件文",
			input:    "雨が降ったら、家にいます。",
			expected: "雨が降ったら、家にいる。",
		},
		{
			name:     "理由を表す文",
			input:    "疲れているから、早く寝ます。",
			expected: "疲れているから、早く寝る。",
		},
		{
			name:     "単純な否定文",
			input:    "本を読みません。",
			expected: "本を読まない。",
		},
		{
			name:     "単純な過去文",
			input:    "昨日映画を見ました。",
			expected: "昨日映画を見た。",
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
