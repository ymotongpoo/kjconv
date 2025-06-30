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
		{
			name:     "接続詞",
			input:    "すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいが、この例の構成では便宜上「未指定」アドレスbbbを使用しています。",
			expected: []string{"すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいですが、この例の構成では便宜上「未指定」アドレスbbbを使用しています。"},
		},
		{
			name:     "接続詞2",
			input:    "プロセッサーはオプションだが、いくつかは推奨です。",
			expected: []string{"プロセッサーはオプションですが、いくつかは推奨です。"},
		},
		{
			name:     "接続詞3",
			input:    "コンポーネントが設定されているが、config節で定義されていない場合、そのコンポーネントは有効になりません。",
			expected: []string{"コンポーネントが設定されていますが、config節で定義されていない場合、そのコンポーネントは有効になりません。"},
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

func TestProcessTextWithQuotes(t *testing.T) {
	// Create a converter for testing
	converter, err := NewConverter()
	if err != nil {
		t.Fatalf("NewConverter() failed: %v", err)
	}

	// Processor that uses actual conversion logic
	processor := func(text string) (string, error) {
		return converter.convertCasualToPoliteSegment(text)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "引用文なし",
			input:    "今日は晴れだ。",
			expected: "今日は晴れです。",
		},
		{
			name:     "引用文あり（「」）",
			input:    "彼は「今日は晴れだ」と言った。",
			expected: "彼は「今日は晴れだ」と言いました。",
		},
		{
			name:     "引用文あり（『』）",
			input:    "彼は『今日は晴れだ』と言った。",
			expected: "彼は『今日は晴れだ』と言いました。",
		},
		{
			name:     "複数引用文",
			input:    "「今日は晴れだ」と「明日は雨だ」と言った。",
			expected: "「今日は晴れだ」と「明日は雨だ」と言いました。",
		},
		{
			name:     "引用文前後にテキスト",
			input:    "昨日彼は「今日は晴れだ」と言った。",
			expected: "昨日彼は「今日は晴れだ」と言いました。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessTextWithQuotes(tt.input, processor)
			if err != nil {
				t.Errorf("ProcessTextWithQuotes() failed: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("ProcessTextWithQuotes(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
