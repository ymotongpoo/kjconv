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
		{
			name:     "文中の「が」接続詞変換",
			input:    "すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいが、この例の構成では便宜上「未指定」アドレスbbbを使用している。",
			expected: "すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいですが、この例の構成では便宜上「未指定」アドレスbbbを使用しています。",
		},
		{
			name:     "文中の「が」接続詞変換2",
			input:    "プロセッサーはオプションだが、いくつかは推奨だ。",
			expected: "プロセッサーはオプションですが、いくつかは推奨です。",
		},
		{
			name:     "文中の「が」接続詞変換3",
			input:    "コンポーネントが設定されているが、config節で定義されていない場合、そのコンポーネントは有効にならない。",
			expected: "コンポーネントが設定されていますが、config節で定義されていない場合、そのコンポーネントは有効になりません。",
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
		{
			name:     "文中の「ですが」接続詞変換",
			input:    "すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいですが、この例の構成では便宜上「未指定」アドレスbbbを使用しています。",
			expected: "すべてのクライアントがローカルの場合、エンドポイントをaaaにバインドするのが望ましいが、この例の構成では便宜上「未指定」アドレスbbbを使用している。",
		},
		{
			name:     "文中の「ですが」接続詞変換2",
			input:    "プロセッサーはオプションですが、いくつかは推奨です。",
			expected: "プロセッサーはオプションだが、いくつかは推奨だ。",
		},
		{
			name:     "文中の「ですが」接続詞変換3",
			input:    "コンポーネントが設定されていますが、config節で定義されていない場合、そのコンポーネントは有効になりません。",
			expected: "コンポーネントが設定されているが、config節で定義されていない場合、そのコンポーネントは有効にならない。",
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
