package kjconv

import (
	"strings"
)

// convertCasualToPolite converts a sentence from casual form to polite form.
func (c *Converter) convertCasualToPolite(sentence string) (string, error) {
	if IsQuotedText(sentence) {
		// Don't convert quoted text
		return sentence, nil
	}
	
	morphemes, err := c.AnalyzeMorphemes(sentence)
	if err != nil {
		return "", err
	}
	
	if len(morphemes) == 0 {
		return sentence, nil
	}
	
	// Convert from the end of the sentence
	converted := c.convertVerbCasualToPolite(morphemes)
	converted = c.convertAdjectiveCasualToPolite(converted)
	converted = c.convertNounCasualToPolite(converted)
	converted = c.convertAuxiliaryCasualToPolite(converted)
	
	return c.reconstructSentence(converted), nil
}

// convertVerbCasualToPolite converts verbs from casual to polite form.
// 動詞の終止形・連体形 → 連用形 + 「ます」
func (c *Converter) convertVerbCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Check the last morpheme for verb conversion
	lastIdx := len(result) - 1
	last := result[lastIdx]
	
	// Check if it's a verb in 終止形 or 連体形
	if last.PartOfSpeech == "動詞" && 
	   (last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
		
		// Convert to 連用形 + ます
		renyoukei := c.getVerbRenyoukei(last)
		if renyoukei != "" {
			result[lastIdx].Surface = renyoukei + "ます"
			result[lastIdx].InflectionForm = "連用形"
		}
	}
	
	return result
}

// getVerbRenyoukei converts a verb to its 連用形 (continuative form).
func (c *Converter) getVerbRenyoukei(morpheme MorphemeInfo) string {
	baseForm := morpheme.BaseForm
	if baseForm == "" {
		baseForm = morpheme.Surface
	}
	
	// Handle common verb conjugation patterns
	switch morpheme.InflectionType {
	case "五段・カ行イ音便":
		// 書く → 書き
		if strings.HasSuffix(baseForm, "く") {
			return strings.TrimSuffix(baseForm, "く") + "き"
		}
	case "五段・ガ行":
		// 泳ぐ → 泳ぎ
		if strings.HasSuffix(baseForm, "ぐ") {
			return strings.TrimSuffix(baseForm, "ぐ") + "ぎ"
		}
	case "五段・サ行":
		// 話す → 話し
		if strings.HasSuffix(baseForm, "す") {
			return strings.TrimSuffix(baseForm, "す") + "し"
		}
	case "五段・タ行":
		// 立つ → 立ち
		if strings.HasSuffix(baseForm, "つ") {
			return strings.TrimSuffix(baseForm, "つ") + "ち"
		}
	case "五段・ナ行":
		// 死ぬ → 死に
		if strings.HasSuffix(baseForm, "ぬ") {
			return strings.TrimSuffix(baseForm, "ぬ") + "に"
		}
	case "五段・バ行":
		// 呼ぶ → 呼び
		if strings.HasSuffix(baseForm, "ぶ") {
			return strings.TrimSuffix(baseForm, "ぶ") + "び"
		}
	case "五段・マ行":
		// 読む → 読み
		if strings.HasSuffix(baseForm, "む") {
			return strings.TrimSuffix(baseForm, "む") + "み"
		}
	case "五段・ラ行":
		// 作る → 作り
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "り"
		}
	case "五段・ワ行促音便":
		// 言う → 言い
		if strings.HasSuffix(baseForm, "う") {
			return strings.TrimSuffix(baseForm, "う") + "い"
		}
	case "一段":
		// 食べる → 食べ
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る")
		}
	case "カ変":
		// 来る → 来
		if baseForm == "来る" {
			return "来"
		}
	case "サ変":
		// する → し
		if baseForm == "する" {
			return "し"
		}
	}
	
	// Fallback: try to use the surface form if it looks like 連用形
	if morpheme.InflectionForm == "連用形" {
		return morpheme.Surface
	}
	
	// Default fallback
	return strings.TrimSuffix(baseForm, "る")
}

// convertAdjectiveCasualToPolite converts adjectives from casual to polite form.
// 形容詞の終止形・連体形 → + 「です」
func (c *Converter) convertAdjectiveCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	last := result[lastIdx]
	
	// Check if it's an i-adjective in 終止形 or 連体形
	if last.PartOfSpeech == "形容詞" && last.PartOfSpeechDetail1 == "自立" &&
	   (last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
		
		result[lastIdx].Surface = last.Surface + "です"
	}
	
	return result
}

// convertNounCasualToPolite converts nouns with copula from casual to polite form.
// 「だ」→「です」、「である」→「です」
func (c *Converter) convertNounCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	last := result[lastIdx]
	
	// Check for copula だ
	if last.PartOfSpeech == "助動詞" && last.BaseForm == "だ" {
		result[lastIdx].Surface = "です"
		result[lastIdx].BaseForm = "です"
	}
	
	// Check for である (might be split into multiple morphemes)
	if lastIdx > 0 {
		secondLast := result[lastIdx-1]
		if secondLast.Surface == "で" && last.Surface == "ある" {
			// Replace である with です
			result = result[:lastIdx-1] // Remove both で and ある
			result = append(result, MorphemeInfo{
				Surface: "です",
				PartOfSpeech: "助動詞",
				BaseForm: "です",
			})
		}
	}
	
	return result
}

// convertAuxiliaryCasualToPolite converts auxiliary verbs and expressions.
func (c *Converter) convertAuxiliaryCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Handle various auxiliary expressions
	// This is a simplified implementation - more patterns can be added
	
	for i := len(result) - 1; i >= 0; i-- {
		morpheme := result[i]
		
		// Handle だろう → でしょう
		if morpheme.Surface == "だろう" {
			result[i].Surface = "でしょう"
			result[i].BaseForm = "でしょう"
		}
		
		// Handle かもしれない → かもしれません
		if morpheme.Surface == "かもしれない" {
			result[i].Surface = "かもしれません"
		}
	}
	
	return result
}

// reconstructSentence reconstructs a sentence from morphemes.
func (c *Converter) reconstructSentence(morphemes []MorphemeInfo) string {
	var parts []string
	for _, morpheme := range morphemes {
		parts = append(parts, morpheme.Surface)
	}
	return strings.Join(parts, "")
}
