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
	
	// Handle text with embedded quotes
	if ContainsQuotedText(sentence) {
		return ProcessTextWithQuotes(sentence, c.convertCasualToPoliteSegment)
	}
	
	return c.convertCasualToPoliteSegment(sentence)
}

// convertCasualToPoliteSegment converts a text segment (without quotes) from casual to polite form.
func (c *Converter) convertCasualToPoliteSegment(segment string) (string, error) {
	if strings.TrimSpace(segment) == "" {
		return segment, nil
	}
	
	morphemes, err := c.AnalyzeMorphemes(segment)
	if err != nil {
		return "", err
	}
	
	if len(morphemes) == 0 {
		return segment, nil
	}
	
	// Convert from the end of the sentence
	converted := c.convertVerbCasualToPolite(morphemes)
	converted = c.convertAdjectiveCasualToPolite(converted)
	converted = c.convertNounCasualToPolite(converted)
	converted = c.convertAuxiliaryCasualToPolite(converted)
	converted = c.convertConjunctionCasualToPolite(converted)
	
	result := c.reconstructSentence(converted)
	
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
		
		// Check if it's a verb in 基本形, 終止形 or 連体形
		if last.PartOfSpeech == "動詞" && 
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
			
			// Convert to 連用形 + ます
			renyoukei := c.getVerbRenyoukei(last)
			if renyoukei != "" {
				result[actualLastIdx].Surface = renyoukei
				result[actualLastIdx].InflectionForm = "連用形"
				
				// Add ます as a separate morpheme
				masuMorpheme := MorphemeInfo{
					Surface:        "ます",
					PartOfSpeech:   "助動詞",
					InflectionForm: "基本形",
					BaseForm:       "ます",
				}
				
				// Insert ます before punctuation
				newResult := make([]MorphemeInfo, len(result)+1)
				copy(newResult[:actualLastIdx+1], result[:actualLastIdx+1])
				newResult[actualLastIdx+1] = masuMorpheme
				copy(newResult[actualLastIdx+2:], result[actualLastIdx+1:])
				result = newResult
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
	case "五段・カ行イ音便", "五段・カ行促音便":
		// 書く → 書き, 行く → 行き
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
		
		
		// Check if it's an i-adjective in 基本形, 終止形 or 連体形
		if last.PartOfSpeech == "形容詞" && last.PartOfSpeechDetail1 == "自立" &&
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形" || last.InflectionForm == "連体形") {
			
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
		
		// Check for copula だ (基本形 or 終止形)
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "だ" && 
		   (last.InflectionForm == "基本形" || last.InflectionForm == "終止形") {
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
	result = c.handleNegativeCasualToPolite(result)
	
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
	
	// Skip punctuation at the end
	lastIdx := len(result) - 1
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		// Check if it's a past tense auxiliary た/だ
		if last.PartOfSpeech == "助動詞" && (last.Surface == "た" || last.Surface == "だ") &&
		   last.BaseForm == "た" && last.InflectionForm == "基本形" {
			
			// Find the verb before た/だ and convert to ました
			if actualLastIdx > 0 {
				prev := result[actualLastIdx-1]
				
				// Handle negative past tense: なかっ + た → ませんでした
				if prev.PartOfSpeech == "助動詞" && prev.Surface == "なかっ" && prev.BaseForm == "ない" {
					if actualLastIdx > 1 {
						verb := result[actualLastIdx-2]
						if verb.PartOfSpeech == "動詞" {
							// Convert to ませんでした
							renyoukei := c.getVerbRenyoukei(verb)
							if renyoukei != "" {
								result[actualLastIdx-2].Surface = renyoukei
								result[actualLastIdx-1].Surface = "ませんでし"
								result[actualLastIdx].Surface = "た"
								result[actualLastIdx].BaseForm = "た"
							}
						}
					}
				} else if prev.PartOfSpeech == "動詞" && (prev.InflectionForm == "連用タ接続" || prev.InflectionForm == "連用形") {
					// Convert verb to 連用形 and change た to ました
					renyoukei := c.getVerbRenyoukeiFromTaForm(prev)
					if renyoukei != "" {
						result[actualLastIdx-1].Surface = renyoukei
						result[actualLastIdx].Surface = "ました"
						result[actualLastIdx].BaseForm = "ます"
					}
				}
			}
		}
	}
	
	return result
}

// getVerbRenyoukeiFromTaForm converts a verb in タ接続 form to 連用形.
func (c *Converter) getVerbRenyoukeiFromTaForm(morpheme MorphemeInfo) string {
	baseForm := morpheme.BaseForm
	if baseForm == "" {
		baseForm = morpheme.Surface
	}
	
	// If it's already in a form that can be used with ます, use it directly
	if morpheme.InflectionForm == "連用タ接続" {
		// For タ接続 forms, we can often use them directly with ます
		surface := morpheme.Surface
		
		// Handle common patterns
		if strings.HasSuffix(surface, "っ") {
			// 言っ → 言い
			return strings.TrimSuffix(surface, "っ") + "い"
		}
		if strings.HasSuffix(surface, "ん") {
			// 読ん → 読み (for 読んだ)
			return strings.TrimSuffix(surface, "ん") + "み"
		}
		if strings.HasSuffix(surface, "い") {
			// 書い → 書き (for 書いた)
			return strings.TrimSuffix(surface, "い") + "き"
		}
		
		// Fallback: use the standard renyoukei conversion
		return c.getVerbRenyoukei(morpheme)
	}
	
	return c.getVerbRenyoukei(morpheme)
}

// handleNegativeCasualToPolite converts negative form from casual to polite.
// ～ない → ～ません
func (c *Converter) handleNegativeCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
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
		
		// Check if it's negative auxiliary ない
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "ない" {
			if actualLastIdx > 0 {
				prev := result[actualLastIdx-1]
				if prev.PartOfSpeech == "動詞" {
					// Convert verb to 連用形 and change ない to ません
					renyoukei := c.getVerbRenyoukei(prev)
					if renyoukei != "" {
						result[actualLastIdx-1].Surface = renyoukei
						result[actualLastIdx].Surface = "ません"
					}
				}
			}
		}
		
		// Handle i-adjective negative: ～くない → ～くありません or ～くないです
		if last.PartOfSpeech == "形容詞" && strings.HasSuffix(last.Surface, "くない") {
			base := strings.TrimSuffix(last.Surface, "くない")
			result[actualLastIdx].Surface = base + "くありません"
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
// convertConjunctionCasualToPolite converts conjunctions from casual to polite form.
// だから → ですから, だが → ですが, が → ですが (when used as conjunction)
func (c *Converter) convertConjunctionCasualToPolite(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Check all morphemes for conjunctions
	for i := 0; i < len(result); i++ {
		morpheme := result[i]
		
		// Convert specific conjunctions
		switch morpheme.Surface {
		case "だから":
			if morpheme.PartOfSpeech == "接続詞" {
				result[i].Surface = "ですから"
			}
		case "だが":
			if morpheme.PartOfSpeech == "接続詞" {
				result[i].Surface = "ですが"
			}
		case "が":
			// Convert "が" to "ですが" when it's used as a conjunction (接続助詞)
			// and preceded by "だ" (助動詞)
			if morpheme.PartOfSpeech == "助詞" && morpheme.PartOfSpeechDetail1 == "接続助詞" {
				// Check if the previous morpheme is "だ" (助動詞)
				if i > 0 && result[i-1].Surface == "だ" && result[i-1].PartOfSpeech == "助動詞" {
					// Remove "だ" and replace "が" with "ですが"
					result = append(result[:i-1], result[i:]...)
					result[i-1].Surface = "ですが"
					// Adjust index since we removed one element
					i--
				} else if i > 0 && result[i-1].Surface == "いる" && result[i-1].PartOfSpeech == "動詞" {
					// Convert "いる" to "います" and keep "が" as "が"
					result[i-1].Surface = "います"
					// Keep "が" as is
				} else {
					result[i].Surface = "ですが"
				}
			}
		}
	}
	
	return result
}

