package kjconv

import (
	"testing"
)

func TestGetVerbRenyoukei(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		morpheme MorphemeInfo
		expected string
	}{
		{
			name: "五段・カ行イ音便（書く）",
			morpheme: MorphemeInfo{
				Surface:        "書く",
				BaseForm:       "書く",
				InflectionType: "五段・カ行イ音便",
			},
			expected: "書き",
		},
		{
			name: "五段・サ行（話す）",
			morpheme: MorphemeInfo{
				Surface:        "話す",
				BaseForm:       "話す",
				InflectionType: "五段・サ行",
			},
			expected: "話し",
		},
		{
			name: "五段・マ行（読む）",
			morpheme: MorphemeInfo{
				Surface:        "読む",
				BaseForm:       "読む",
				InflectionType: "五段・マ行",
			},
			expected: "読み",
		},
		{
			name: "一段（食べる）",
			morpheme: MorphemeInfo{
				Surface:        "食べる",
				BaseForm:       "食べる",
				InflectionType: "一段",
			},
			expected: "食べ",
		},
		{
			name: "カ変（来る）",
			morpheme: MorphemeInfo{
				Surface:        "来る",
				BaseForm:       "来る",
				InflectionType: "カ変",
			},
			expected: "来",
		},
		{
			name: "サ変（する）",
			morpheme: MorphemeInfo{
				Surface:        "する",
				BaseForm:       "する",
				InflectionType: "サ変",
			},
			expected: "し",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.getVerbRenyoukei(tt.morpheme)
			if result != tt.expected {
				t.Errorf("getVerbRenyoukei() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestGetVerbTaForm(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		morpheme MorphemeInfo
		expected string
	}{
		{
			name: "五段・カ行イ音便（書く）",
			morpheme: MorphemeInfo{
				Surface:        "書く",
				BaseForm:       "書く",
				InflectionType: "五段・カ行イ音便",
			},
			expected: "書いた",
		},
		{
			name: "五段・サ行（話す）",
			morpheme: MorphemeInfo{
				Surface:        "話す",
				BaseForm:       "話す",
				InflectionType: "五段・サ行",
			},
			expected: "話した",
		},
		{
			name: "五段・マ行（読む）",
			morpheme: MorphemeInfo{
				Surface:        "読む",
				BaseForm:       "読む",
				InflectionType: "五段・マ行",
			},
			expected: "読んだ",
		},
		{
			name: "一段（食べる）",
			morpheme: MorphemeInfo{
				Surface:        "食べる",
				BaseForm:       "食べる",
				InflectionType: "一段",
			},
			expected: "食べた",
		},
		{
			name: "カ変（来る）",
			morpheme: MorphemeInfo{
				Surface:        "来る",
				BaseForm:       "来る",
				InflectionType: "カ変",
			},
			expected: "来た",
		},
		{
			name: "サ変（する）",
			morpheme: MorphemeInfo{
				Surface:        "する",
				BaseForm:       "する",
				InflectionType: "サ変",
			},
			expected: "した",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.getVerbTaForm(tt.morpheme)
			if result != tt.expected {
				t.Errorf("getVerbTaForm() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestGetVerbNaiForm(t *testing.T) {
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	tests := []struct {
		name     string
		morpheme MorphemeInfo
		expected string
	}{
		{
			name: "五段・カ行イ音便（書く）",
			morpheme: MorphemeInfo{
				Surface:        "書く",
				BaseForm:       "書く",
				InflectionType: "五段・カ行イ音便",
			},
			expected: "書かない",
		},
		{
			name: "五段・サ行（話す）",
			morpheme: MorphemeInfo{
				Surface:        "話す",
				BaseForm:       "話す",
				InflectionType: "五段・サ行",
			},
			expected: "話さない",
		},
		{
			name: "五段・マ行（読む）",
			morpheme: MorphemeInfo{
				Surface:        "読む",
				BaseForm:       "読む",
				InflectionType: "五段・マ行",
			},
			expected: "読まない",
		},
		{
			name: "一段（食べる）",
			morpheme: MorphemeInfo{
				Surface:        "食べる",
				BaseForm:       "食べる",
				InflectionType: "一段",
			},
			expected: "食べない",
		},
		{
			name: "カ変（来る）",
			morpheme: MorphemeInfo{
				Surface:        "来る",
				BaseForm:       "来る",
				InflectionType: "カ変",
			},
			expected: "来ない",
		},
		{
			name: "サ変（する）",
			morpheme: MorphemeInfo{
				Surface:        "する",
				BaseForm:       "する",
				InflectionType: "サ変",
			},
			expected: "しない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.getVerbNaiForm(tt.morpheme)
			if result != tt.expected {
				t.Errorf("getVerbNaiForm() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
