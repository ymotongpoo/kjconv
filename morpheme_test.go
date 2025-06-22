package kjconv

import (
	"testing"
)

func TestAnalyzeMorphemes(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected []MorphemeInfo
	}{
		{
			name:  "簡単な文",
			input: "今日は晴れだ。",
			expected: []MorphemeInfo{
				{Surface: "今日", PartOfSpeech: "名詞", BaseForm: "今日"},
				{Surface: "は", PartOfSpeech: "助詞", BaseForm: "は"},
				{Surface: "晴れ", PartOfSpeech: "名詞", BaseForm: "晴れ"},
				{Surface: "だ", PartOfSpeech: "助動詞", BaseForm: "だ", InflectionForm: "基本形"},
				{Surface: "。", PartOfSpeech: "記号", BaseForm: "。"},
			},
		},
		{
			name:  "動詞を含む文",
			input: "本を読む。",
			expected: []MorphemeInfo{
				{Surface: "本", PartOfSpeech: "名詞", BaseForm: "本"},
				{Surface: "を", PartOfSpeech: "助詞", BaseForm: "を"},
				{Surface: "読む", PartOfSpeech: "動詞", BaseForm: "読む", InflectionForm: "基本形"},
				{Surface: "。", PartOfSpeech: "記号", BaseForm: "。"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.AnalyzeMorphemes(tt.input)
			if err != nil {
				t.Errorf("AnalyzeMorphemes() failed: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("AnalyzeMorphemes() returned %d morphemes, expected %d", len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing morpheme at index %d", i)
					continue
				}

				actual := result[i]
				if actual.Surface != expected.Surface {
					t.Errorf("Morpheme %d: Surface = %q, expected %q", i, actual.Surface, expected.Surface)
				}
				if actual.PartOfSpeech != expected.PartOfSpeech {
					t.Errorf("Morpheme %d: PartOfSpeech = %q, expected %q", i, actual.PartOfSpeech, expected.PartOfSpeech)
				}
				if actual.BaseForm != expected.BaseForm {
					t.Errorf("Morpheme %d: BaseForm = %q, expected %q", i, actual.BaseForm, expected.BaseForm)
				}
				if expected.InflectionForm != "" && actual.InflectionForm != expected.InflectionForm {
					t.Errorf("Morpheme %d: InflectionForm = %q, expected %q", i, actual.InflectionForm, expected.InflectionForm)
				}
			}
		})
	}
}

func TestAnalyzeMorphemes_EmptyString(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	result, err := converter.AnalyzeMorphemes("")
	if err != nil {
		t.Errorf("AnalyzeMorphemes() failed with empty string: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("AnalyzeMorphemes() with empty string returned %d morphemes, expected 0", len(result))
	}
}
