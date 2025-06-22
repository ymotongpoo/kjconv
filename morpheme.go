package kjconv

import (
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// MorphemeInfo represents morphological analysis information for a token.
type MorphemeInfo struct {
	Surface    string // 表層形
	PartOfSpeech string // 品詞
	PartOfSpeechDetail1 string // 品詞細分類1
	PartOfSpeechDetail2 string // 品詞細分類2
	PartOfSpeechDetail3 string // 品詞細分類3
	InflectionType string // 活用型
	InflectionForm string // 活用形
	BaseForm string // 原形
}

// AnalyzeMorphemes performs morphological analysis on the input text.
func (c *Converter) AnalyzeMorphemes(text string) ([]MorphemeInfo, error) {
	tokens := c.tokenizer.Tokenize(text)
	
	var morphemes []MorphemeInfo
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		
		features := token.Features()
		
		morpheme := MorphemeInfo{
			Surface: token.Surface,
		}
		
		// Extract features according to IPADIC format
		if len(features) > 0 {
			morpheme.PartOfSpeech = features[0]
		}
		if len(features) > 1 {
			morpheme.PartOfSpeechDetail1 = features[1]
		}
		if len(features) > 2 {
			morpheme.PartOfSpeechDetail2 = features[2]
		}
		if len(features) > 3 {
			morpheme.PartOfSpeechDetail3 = features[3]
		}
		if len(features) > 4 {
			morpheme.InflectionType = features[4]
		}
		if len(features) > 5 {
			morpheme.InflectionForm = features[5]
		}
		if len(features) > 6 {
			morpheme.BaseForm = features[6]
		}
		
		morphemes = append(morphemes, morpheme)
	}
	
	return morphemes, nil
}
