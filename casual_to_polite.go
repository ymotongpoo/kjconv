package kjconv

import (
	"log/slog"
	"strings"
)

// convertCasualToPolite converts a sentence from casual form to polite form.
func (c *Converter) convertCasualToPolite(sentence string) (string, error) {
	slog.Debug("converting casual to polite", "sentence", sentence)
	
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
	converted := c.convertVerbCasualToPolite(morphemes)
	converted = c.convertAdjectiveCasualToPolite(converted)
	converted = c.convertNounCasualToPolite(converted)
	converted = c.convertAuxiliaryCasualToPolite(converted)
	
	result := c.reconstructSentence(converted)
	slog.Debug("conversion result", "original", sentence, "converted", result)
	
	return result, nil
}

// convertVerbCasualToPolite converts verbs from casual to polite form.
// 動詞の終止形・連体形 → 連用形 + 「ます」
func (c *Converter) convertVerbCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
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
		
		slog.Debug("checking verb conversion", "surface", last.Surface, "pos", last.PartOfSpeech, "inflection_form", last.InflectionForm, "base_form", last.BaseForm)
		
		// Check if it's a verb in 基本形, 終止形 or 連体形
		if last.PartOfSpeech == "動詞" && 
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
			
			// Convert to 連用形 + ます
			renyoukei := c.getVerbRenyoukei(last)
			if renyoukei != "" {
				slog.Debug("converting verb", "original", last.Surface, "renyoukei", renyoukei)
				result[actualLastIdx].Surface = renyoukei + "ます"
				result[actualLastIdx].InflectionForm = "連用形"
			}
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
	case "サ変・スル":
		// する → し
		if baseForm == "する" {
			return "し"
		}
	case "サ変":
		// する → し
		if baseForm == "する" {
			return "し"
		}
	}
	
	// Special handling for する verb
	if baseForm == "する" {
		return "し"
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
	
	// Skip punctuation at the end
	lastIdx := len(result) - 1
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		slog.Debug("checking adjective conversion", "surface", last.Surface, "pos", last.PartOfSpeech, "detail1", last.PartOfSpeechDetail1, "inflection_form", last.InflectionForm)
		
		// Check if it's an i-adjective in 基本形, 終止形 or 連体形
		if last.PartOfSpeech == "形容詞" && last.PartOfSpeechDetail1 == "自立" &&
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
			
			slog.Debug("converting adjective", "original", last.Surface)
			result[actualLastIdx].Surface = last.Surface + "です"
		}
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
	
	// Skip punctuation at the end
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		slog.Debug("checking noun conversion", "surface", last.Surface, "pos", last.PartOfSpeech, "base_form", last.BaseForm, "inflection_form", last.InflectionForm)
		
		// Check for copula だ (基本形 or 終止形)
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "だ" && 
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形") {
			slog.Debug("converting だ to です")
			result[actualLastIdx].Surface = "です"
			result[actualLastIdx].BaseForm = "です"
		}
		
		// Check for である (might be split into multiple morphemes)
		if actualLastIdx > 0 {
			secondLast := result[actualLastIdx-1]
			if secondLast.Surface == "で" && last.Surface == "ある" {
				// Replace である with です
				result = result[:actualLastIdx-1] // Remove both で and ある
				result = append(result, MorphemeInfo{
					Surface: "です",
					PartOfSpeech: "助動詞",
					BaseForm: "です",
				})
				// Add back any punctuation
				for i := actualLastIdx + 1; i <= lastIdx; i++ {
					result = append(result, morphemes[i])
				}
			}
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
	
	// Handle complex patterns first (multi-morpheme expressions)
	result = c.handleComplexCasualToPolite(result)
	
	// Handle single morpheme patterns
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
		
		// Handle ようだ → ようです
		if morpheme.Surface == "ようだ" {
			result[i].Surface = "ようです"
		}
		
		// Handle わけだ → わけです
		if morpheme.Surface == "わけだ" {
			result[i].Surface = "わけです"
		}
		
		// Handle はずだ → はずです
		if morpheme.Surface == "はずだ" {
			result[i].Surface = "はずです"
		}
	}
	
	return result
}

// handleComplexCasualToPolite handles multi-morpheme expressions.
func (c *Converter) handleComplexCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) < 2 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Handle のだ/んだ → のです/んです
	for i := len(result) - 2; i >= 0; i-- {
		if i+1 < len(result) {
			current := result[i]
			next := result[i+1]
			
			// Handle のだ → のです
			if current.Surface == "の" && (next.Surface == "だ" || next.BaseForm == "だ") {
				result[i+1].Surface = "です"
				result[i+1].BaseForm = "です"
			}
			
			// Handle んだ → んです  
			if current.Surface == "ん" && (next.Surface == "だ" || next.BaseForm == "だ") {
				result[i+1].Surface = "です"
				result[i+1].BaseForm = "です"
			}
		}
	}
	
	// Handle past tense た → ました
	result = c.handlePastTenseCasualToPolite(result)
	
	// Handle negative ない → ません
	result = c.handleNegativeCasualToPolite(result)
	
	return result
}

// handlePastTenseCasualToPolite converts past tense from casual to polite.
// ～た → ～ました
func (c *Converter) handlePastTenseCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	if lastIdx >= 0 {
		last := result[lastIdx]
		
		// Check if it's a past tense auxiliary た/だ
		if last.PartOfSpeech == "助動詞" && (last.Surface == "た" || last.Surface == "だ") &&
		   last.InflectionForm == "終止形" {
			
			// Find the verb before た/だ and convert to ました
			if lastIdx > 0 {
				prev := result[lastIdx-1]
				if prev.PartOfSpeech == "動詞" {
					// Convert verb to 連用形 and change た to ました
					renyoukei := c.getVerbRenyoukei(prev)
					if renyoukei != "" {
						result[lastIdx-1].Surface = renyoukei
						result[lastIdx].Surface = "ました"
					}
				}
			}
		}
	}
	
	return result
}

// handleNegativeCasualToPolite converts negative form from casual to polite.
// ～ない → ～ません
func (c *Converter) handleNegativeCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	lastIdx := len(result) - 1
	if lastIdx >= 0 {
		last := result[lastIdx]
		
		// Check if it's negative auxiliary ない
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "ない" {
			if lastIdx > 0 {
				prev := result[lastIdx-1]
				if prev.PartOfSpeech == "動詞" {
					// Convert verb to 連用形 and change ない to ません
					renyoukei := c.getVerbRenyoukei(prev)
					if renyoukei != "" {
						result[lastIdx-1].Surface = renyoukei
						result[lastIdx].Surface = "ません"
					}
				}
			}
		}
		
		// Handle i-adjective negative: ～くない → ～くありません or ～くないです
		if last.PartOfSpeech == "形容詞" && strings.HasSuffix(last.Surface, "くない") {
			base := strings.TrimSuffix(last.Surface, "くない")
			result[lastIdx].Surface = base + "くありません"
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
