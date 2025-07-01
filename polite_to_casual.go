package kjconv

import (
	"strings"
)

// convertPoliteToCasual converts a sentence from polite form to casual form.
func (c *Converter) convertPoliteToCasual(sentence string) (string, error) {
	if IsQuotedText(sentence) {
		// Don't convert quoted text
		return sentence, nil
	}
	
	// Handle text with embedded quotes
	if ContainsQuotedText(sentence) {
		return ProcessTextWithQuotes(sentence, c.convertPoliteToCasualSegment)
	}
	
	return c.convertPoliteToCasualSegment(sentence)
}

// convertPoliteToCasualSegment converts a text segment (without quotes) from polite to casual form.
func (c *Converter) convertPoliteToCasualSegment(segment string) (string, error) {
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
	converted := c.convertVerbPoliteToCase(morphemes)
	converted = c.convertAdjectivePoliteToCase(converted)
	converted = c.convertNounPoliteToCase(converted)
	converted = c.convertAuxiliaryPoliteToCase(converted)
	converted = c.handleNegativePoliteToCase(converted)
	converted = c.convertConjunctionPoliteToCase(converted)
	
	result := c.reconstructSentence(converted)
	
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
	
	// Skip punctuation at the end
	lastIdx := len(result) - 1
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		
		// Handle ます forms
		if last.PartOfSpeech == "助動詞" && last.BaseForm == "ます" {
			switch last.InflectionForm {
			case "基本形": // ます → 終止形
				if actualLastIdx > 0 {
					prev := result[actualLastIdx-1]
					if prev.PartOfSpeech == "動詞" && prev.InflectionForm == "連用形" {
						// Convert to dictionary form
						baseForm := prev.BaseForm
						if baseForm != "" {
							result[actualLastIdx-1].Surface = baseForm
							result[actualLastIdx-1].InflectionForm = "基本形"
						}
						// Remove ます
						result = result[:actualLastIdx]
						// Add back any punctuation
						for i := actualLastIdx + 1; i <= lastIdx; i++ {
							result = append(result, morphemes[i])
						}
					}
				}
			case "過去": // ました → 過去形（タ形）
				if actualLastIdx > 0 {
					prev := result[actualLastIdx-1]
					if prev.PartOfSpeech == "動詞" {
						taForm := c.getVerbTaForm(prev)
						if taForm != "" {
							result[actualLastIdx-1].Surface = taForm
							result[actualLastIdx-1].InflectionForm = "終止形"
						}
						// Remove ました
						result = result[:actualLastIdx]
						// Add back any punctuation
						for i := actualLastIdx + 1; i <= lastIdx; i++ {
							result = append(result, morphemes[i])
						}
					}
				}
			case "否定": // ません → 否定形（ナイ形）
				if actualLastIdx > 0 {
					prev := result[actualLastIdx-1]
					if prev.PartOfSpeech == "動詞" {
						naiForm := c.getVerbNaiForm(prev)
						if naiForm != "" {
							result[actualLastIdx-1].Surface = naiForm
							result[actualLastIdx-1].InflectionForm = "終止形"
						}
						// Remove ません
						result = result[:actualLastIdx]
						// Add back any punctuation
						for i := actualLastIdx + 1; i <= lastIdx; i++ {
							result = append(result, morphemes[i])
						}
					}
				}
			}
		}
		
		// Handle ました pattern: まし + た
		if last.Surface == "た" && last.PartOfSpeech == "助動詞" && last.BaseForm == "た" &&
		   actualLastIdx > 0 {
			prev := result[actualLastIdx-1]
			if prev.Surface == "まし" && prev.PartOfSpeech == "助動詞" && prev.BaseForm == "ます" &&
			   actualLastIdx > 1 {
				verb := result[actualLastIdx-2]
				if verb.PartOfSpeech == "動詞" && verb.InflectionForm == "連用形" {
					// Convert 言い + まし + た → 言った
					taForm := c.getVerbTaForm(verb)
					if taForm != "" {
						result[actualLastIdx-2].Surface = taForm
						result[actualLastIdx-2].InflectionForm = "終止形"
						// Remove まし and た
						result = result[:actualLastIdx-1]
						// Add back any punctuation
						for i := actualLastIdx + 1; i <= lastIdx; i++ {
							result = append(result, morphemes[i])
						}
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
	
	// Skip punctuation at the end
	lastIdx := len(result) - 1
	actualLastIdx := lastIdx
	for actualLastIdx >= 0 && result[actualLastIdx].PartOfSpeech == "記号" {
		actualLastIdx--
	}
	
	if actualLastIdx >= 0 {
		last := result[actualLastIdx]
		
		// Handle いです → い
		if last.Surface == "です" && actualLastIdx > 0 {
			prev := result[actualLastIdx-1]
			if prev.PartOfSpeech == "形容詞" && strings.HasSuffix(prev.Surface, "い") {
				// Remove です
				result = result[:actualLastIdx]
				// Add back any punctuation
				for i := actualLastIdx + 1; i <= lastIdx; i++ {
					result = append(result, morphemes[i])
				}
			}
		}
		
		// Handle くありません → くない
		if strings.HasSuffix(last.Surface, "くありません") {
			base := strings.TrimSuffix(last.Surface, "くありません")
			result[actualLastIdx].Surface = base + "くない"
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
		
		switch last.Surface {
		case "です":
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
// convertConjunctionPoliteToCase converts conjunctions from polite to casual form.
// ですから → だから, ですが → だが, ですが → が (when used as conjunction)
func (c *Converter) convertConjunctionPoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
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
		case "ですから":
			if morpheme.PartOfSpeech == "接続詞" {
				result[i].Surface = "だから"
			}
		case "ですが":
			// Convert "ですが" to "だが" when it's used as a conjunction
			if morpheme.PartOfSpeech == "接続詞" || (morpheme.PartOfSpeech == "助詞" && morpheme.PartOfSpeechDetail1 == "接続助詞") {
				// Check if the previous morpheme is "います" (動詞)
				if i > 0 && result[i-1].Surface == "います" && result[i-1].PartOfSpeech == "助動詞" {
					// Convert "います" to "いる" and keep "ですが" as "が"
					result[i-1].Surface = "いる"
					result[i].Surface = "が"
					result[i].PartOfSpeech = "助詞"
					result[i].PartOfSpeechDetail1 = "接続助詞"
				} else {
					// Insert "だ" before "が"
					newMorpheme := MorphemeInfo{
						Surface:             "だ",
						PartOfSpeech:        "助動詞",
						PartOfSpeechDetail1: "*",
						PartOfSpeechDetail2: "*",
						PartOfSpeechDetail3: "*",
						InflectionType:      "特殊・ダ",
						InflectionForm:      "基本形",
						BaseForm:            "だ",
					}
					// Insert the new morpheme before current position
					result = append(result[:i], append([]MorphemeInfo{newMorpheme}, result[i:]...)...)
					// Update the current morpheme to "が"
					result[i+1].Surface = "が"
					result[i+1].PartOfSpeech = "助詞"
					result[i+1].PartOfSpeechDetail1 = "接続助詞"
					// Skip the next iteration since we added a morpheme
					i++
				}
			}
		case "が":
			// Convert "です" + "が" to "だが" when "が" is used as a conjunction (接続助詞)
			if morpheme.PartOfSpeech == "助詞" && morpheme.PartOfSpeechDetail1 == "接続助詞" {
				// Check if the previous morpheme is "ます" (助動詞) and before that is a verb
				if i > 1 && result[i-1].Surface == "ます" && result[i-1].PartOfSpeech == "助動詞" && result[i-2].PartOfSpeech == "動詞" {
					// Convert "い" + "ます" to "いる" and keep "が" as "が"
					// Get the base form of the verb before "ます"
					baseForm := result[i-2].BaseForm
					if baseForm == "いる" {
						// Remove "ます" and update the verb to its base form
						result = append(result[:i-1], result[i:]...)
						result[i-2].Surface = "いる"
						result[i-2].InflectionForm = "連体形"
						// Adjust index since we removed one element
						i--
					}
				} else if i > 0 && result[i-1].Surface == "です" && result[i-1].PartOfSpeech == "助動詞" {
					// Check if the morpheme before "です" is an adjective
					if i > 1 && result[i-2].PartOfSpeech == "形容詞" {
						// Remove "です" and keep "が" as "が"
						result = append(result[:i-1], result[i:]...)
						// Adjust index since we removed one element
						i--
					} else {
						// Replace "です" with "だ" and keep "が" as "が"
						result[i-1].Surface = "だ"
						// Keep "が" as is
					}
				}
			}
		}
	}
	
	return result
}
// handleNegativePoliteToCase converts negative forms from polite to casual.
// ～ません → ～ない
func (c *Converter) handleNegativePoliteToCase(morphemes []MorphemeInfo) []MorphemeInfo {
	if len(morphemes) == 0 {
		return morphemes
	}
	
	result := make([]MorphemeInfo, len(morphemes))
	copy(result, morphemes)
	
	// Look for ませ + ん pattern (ません)
	for i := 0; i < len(result)-1; i++ {
		if result[i].Surface == "ませ" && result[i].PartOfSpeech == "助動詞" &&
		   result[i+1].Surface == "ん" && result[i+1].PartOfSpeech == "助動詞" {
			
			// Find the verb before ませ
			if i > 0 && result[i-1].PartOfSpeech == "動詞" {
				verb := result[i-1]
				
				// Convert verb from 連用形 to 未然形 and replace ません with ない
				mizenkei := c.getVerbMizenkei(verb)
				if mizenkei != "" {
					result[i-1].Surface = mizenkei
					result[i].Surface = "ない"
					result[i].BaseForm = "ない"
					result[i].PartOfSpeech = "助動詞"
					
					// Remove the ん morpheme
					newResult := make([]MorphemeInfo, len(result)-1)
					copy(newResult[:i+1], result[:i+1])
					copy(newResult[i+1:], result[i+2:])
					result = newResult
					break
				}
			}
		}
	}
	
	return result
}

// getVerbMizenkei converts a verb to its 未然形 (irrealis form) for negative conjugation.
func (c *Converter) getVerbMizenkei(morpheme MorphemeInfo) string {
	baseForm := morpheme.BaseForm
	if baseForm == "" {
		baseForm = morpheme.Surface
	}
	
	// Handle common verb conjugation patterns
	switch morpheme.InflectionType {
	case "五段・カ行イ音便", "五段・カ行促音便":
		// 書く → 書か, 行く → 行か
		if strings.HasSuffix(baseForm, "く") {
			return strings.TrimSuffix(baseForm, "く") + "か"
		}
	case "五段・ガ行":
		// 泳ぐ → 泳が
		if strings.HasSuffix(baseForm, "ぐ") {
			return strings.TrimSuffix(baseForm, "ぐ") + "が"
		}
	case "五段・サ行":
		// 話す → 話さ
		if strings.HasSuffix(baseForm, "す") {
			return strings.TrimSuffix(baseForm, "す") + "さ"
		}
	case "五段・タ行":
		// 立つ → 立た
		if strings.HasSuffix(baseForm, "つ") {
			return strings.TrimSuffix(baseForm, "つ") + "た"
		}
	case "五段・ナ行":
		// 死ぬ → 死な
		if strings.HasSuffix(baseForm, "ぬ") {
			return strings.TrimSuffix(baseForm, "ぬ") + "な"
		}
	case "五段・バ行":
		// 呼ぶ → 呼ば
		if strings.HasSuffix(baseForm, "ぶ") {
			return strings.TrimSuffix(baseForm, "ぶ") + "ば"
		}
	case "五段・マ行":
		// 読む → 読ま
		if strings.HasSuffix(baseForm, "む") {
			return strings.TrimSuffix(baseForm, "む") + "ま"
		}
	case "五段・ラ行":
		// 作る → 作ら
		if strings.HasSuffix(baseForm, "る") {
			return strings.TrimSuffix(baseForm, "る") + "ら"
		}
	case "五段・ワ行促音便":
		// 言う → 言わ
		if strings.HasSuffix(baseForm, "う") {
			return strings.TrimSuffix(baseForm, "う") + "わ"
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
	
	// Fallback: return the surface form
	return morpheme.Surface
}
