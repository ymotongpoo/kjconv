package kjconv

import (
	"log/slog"
	"strings"
)

// convertPoliteToCasual converts a sentence from polite form to casual form.
func (c *Converter) convertPoliteToCasual(sentence string) (string, error) {
	slog.Debug("converting polite to casual", "sentence", sentence)
	
	if IsQuotedText(sentence) {
		// Don't convert quoted text
		slog.Debug("skipping quoted text", "sentence", sentence)
		return sentence, nil
	}
	
	morphemes, err := c.AnalyzeMorphemes(sentence)
	if err != nil {
		return "", err
	}
	
	slog.Debug("morphemes analyzed", "count", len(morphemes))
	for i, m := range morphemes {
		slog.Debug("morpheme", "index", i, "surface", m.Surface, "pos", m.PartOfSpeech, "inflection_form", m.InflectionForm, "base_form", m.BaseForm)
	}
	
	if len(morphemes) == 0 {
		return sentence, nil
	}
	
	// Convert from the end of the sentence
	converted := c.convertVerbPoliteToCase(morphemes)
	converted = c.convertAdjectivePoliteToCase(converted)
	converted = c.convertNounPoliteToCase(converted)
	converted = c.convertAuxiliaryPoliteToCase(converted)
	
	result := c.reconstructSentence(converted)
	slog.Debug("conversion result", "original", sentence, "converted", result)
	
	return result, nil
}

// convertVerbPoliteToCase converts verbs from polite to casual form.
// ～ます → 終止形, ～ました → 過去形, ～ません → 否定形
func (c *Converter) convertVerbPoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	if lastIdx >= 0 {
		last := result[lastIdx]
		
		// Handle ます forms
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "ます" {
			switch last.InflectionForm {
			case "終止形": // ます → 終止形
				if lastIdx > 0 {
					prev := result[lastIdx-1]
					if prev.PartOfSpeech == "動詞" {
						// Convert to dictionary form
						baseForm := prev.BaseForm
						if baseForm != "" {
							result[lastIdx-1].Surface = baseForm
							result[lastIdx-1].InflectionForm = "終止形"
						}
						// Remove ます
						result = result[:lastIdx]
					}
				}
			case "過去": // ました → 過去形（タ形）
				if lastIdx > 0 {
					prev := result[lastIdx-1]
					if prev.PartOfSpeech == "動詞" {
						taForm := c.getVerbTaForm(prev)
						if taForm != "" {
							result[lastIdx-1].Surface = taForm
							result[lastIdx-1].InflectionForm = "終止形"
						}
						// Remove ました
						result = result[:lastIdx]
					}
				}
			case "否定": // ません → 否定形（ナイ形）
				if lastIdx > 0 {
					prev := result[lastIdx-1]
					if prev.PartOfSpeech == "動詞" {
						naiForm := c.getVerbNaiForm(prev)
						if naiForm != "" {
							result[lastIdx-1].Surface = naiForm
							result[lastIdx-1].InflectionForm = "終止形"
						}
						// Remove ません
						result = result[:lastIdx]
					}
				}
			}
		}
	}
	
	return result
}

// getVerbTaForm converts a verb to its past form (タ形).
func (c *Converter) getVerbTaForm(morpheme MorphemeInfo) string {
	baseForm := morpheme.BaseForm
	if baseForm == "" {
		baseForm = morpheme.Surface
	}
	
	switch morpheme.InflectionType {
	case "五段・カ行イ音便":
		if strings.HasSuffix(baseForm, "く") {
			return strings.TrimSuffix(baseForm, "く") + "いた"
		}
	case "五段・ガ行":
		if strings.HasSuffix(baseForm, "ぐ") {
			return strings.TrimSuffix(baseForm, "ぐ") + "いだ"
		}
	case "五段・サ行":
		if strings.HasSuffix(baseForm, "す") {
			return strings.TrimSuffix(baseForm, "す") + "した"
		}
	case "五段・タ行":
		if strings.HasSuffix(baseForm, "つ") {
			return strings.TrimSuffix(baseForm, "つ") + "った"
		}
	case "五段・ナ行":
		if strings.HasSuffix(baseForm, "ぬ") {
			return strings.TrimSuffix(baseForm, "ぬ") + "んだ"
		}
	case "五段・バ行":
		if strings.HasSuffix(baseForm, "ぶ") {
			return strings.TrimSuffix(baseForm, "ぶ") + "んだ"
		}
	case "五段・マ行":
		if strings.HasSuffix(baseForm, "む") {
			return strings.TrimSuffix(baseForm, "む") + "んだ"
		}
	case "五段・ラ行":
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "った"
		}
	case "五段・ワ行促音便":
		if strings.HasSuffix(baseForm, "う") {
			return strings.TrimSuffix(baseForm, "う") + "った"
		}
	case "一段":
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "た"
		}
	case "カ変":
		if baseForm == "来る" {
			return "来た"
		}
	case "サ変":
		if baseForm == "する" {
			return "した"
		}
	}
	
	return baseForm + "た" // fallback
}

// getVerbNaiForm converts a verb to its negative form (ナイ形).
func (c *Converter) getVerbNaiForm(morpheme MorphemeInfo) string {
	baseForm := morpheme.BaseForm
	if baseForm == "" {
		baseForm = morpheme.Surface
	}
	
	switch morpheme.InflectionType {
	case "五段・カ行イ音便":
		if strings.HasSuffix(baseForm, "く") {
			return strings.TrimSuffix(baseForm, "く") + "かない"
		}
	case "五段・ガ行":
		if strings.HasSuffix(baseForm, "ぐ") {
			return strings.TrimSuffix(baseForm, "ぐ") + "がない"
		}
	case "五段・サ行":
		if strings.HasSuffix(baseForm, "す") {
			return strings.TrimSuffix(baseForm, "す") + "さない"
		}
	case "五段・タ行":
		if strings.HasSuffix(baseForm, "つ") {
			return strings.TrimSuffix(baseForm, "つ") + "たない"
		}
	case "五段・ナ行":
		if strings.HasSuffix(baseForm, "ぬ") {
			return strings.TrimSuffix(baseForm, "ぬ") + "なない"
		}
	case "五段・バ行":
		if strings.HasSuffix(baseForm, "ぶ") {
			return strings.TrimSuffix(baseForm, "ぶ") + "ばない"
		}
	case "五段・マ行":
		if strings.HasSuffix(baseForm, "む") {
			return strings.TrimSuffix(baseForm, "む") + "まない"
		}
	case "五段・ラ行":
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "らない"
		}
	case "五段・ワ行促音便":
		if strings.HasSuffix(baseForm, "う") {
			return strings.TrimSuffix(baseForm, "う") + "わない"
		}
	case "一段":
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "ない"
		}
	case "カ変":
		if baseForm == "来る" {
			return "来ない"
		}
	case "サ変":
		if baseForm == "する" {
			return "しない"
		}
	}
	
	return baseForm + "ない" // fallback
}

// convertAdjectivePoliteToCase converts adjectives from polite to casual form.
// ～いです → ～い, ～くありません → ～くない
func (c *Converter) convertAdjectivePoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	if lastIdx >= 0 {
		last := result[lastIdx]
		
		// Handle いです → い
		if last.Surface == "です" && lastIdx > 0 {
			prev := result[lastIdx-1]
			if prev.PartOfSpeech == "形容詞" && strings.HasSuffix(prev.Surface, "い") {
				// Remove です
				result = result[:lastIdx]
			}
		}
		
		// Handle くありません → くない
		if strings.HasSuffix(last.Surface, "くありません") {
			base := strings.TrimSuffix(last.Surface, "くありません")
			result[lastIdx].Surface = base + "くない"
		}
	}
	
	return result
}

// convertNounPoliteToCase converts nouns with polite copula to casual form.
// です → だ, でした → だった, ではありません → ではない
func (c *Converter) convertNounPoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Skip punctuation at the end
	lastIdx := len(result) - 1
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		slog.Debug("checking noun polite conversion", "surface", last.Surface, "pos", last.PartOfSpeech, "base_form", last.BaseForm)
		
		switch last.Surface {
		case "です":
			slog.Debug("converting です to だ")
			result[actualLastIdx].Surface = "だ"
			result[actualLastIdx].BaseForm = "だ"
		case "でした":
			result[actualLastIdx].Surface = "だった"
		case "ではありません":
			result[actualLastIdx].Surface = "ではない"
		case "でしょう":
			result[actualLastIdx].Surface = "だろう"
		}
	}
	
	return result
}

// convertAuxiliaryPoliteToCase converts auxiliary expressions from polite to casual.
func (c *Converter) convertAuxiliaryPoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	for i := len(result) - 1; i >= 0; i-- {
		morpheme := result[i]
		
		// Handle various polite expressions
		switch morpheme.Surface {
		case "かもしれません":
			result[i].Surface = "かもしれない"
		case "ようです":
			result[i].Surface = "ようだ"
		case "わけです":
			result[i].Surface = "わけだ"
		case "はずです":
			result[i].Surface = "はずだ"
		}
	}
	
	// Handle のです/んです → のだ/んだ
	for i := len(result) - 2; i >= 0; i-- {
		if i+1 < len(result) {
			current := result[i]
			next := result[i+1]
			
			if (current.Surface == "の" || current.Surface == "ん") && next.Surface == "です" {
				result[i+1].Surface = "だ"
				result[i+1].BaseForm = "だ"
			}
		}
	}
	
	return result
}
