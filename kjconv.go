// Package kjconv provides Japanese text style conversion between casual (常体) and polite (敬体) forms.
package kjconv

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// ConversionMode represents the direction of text style conversion.
type ConversionMode int

const (
	// CasualToPolite converts from casual form (常体) to polite form (敬体)
	CasualToPolite ConversionMode = iota
	// PoliteToCasual converts from polite form (敬体) to casual form (常体)
	PoliteToCasual
)

// Converter handles Japanese text style conversion.
type Converter struct {
	tokenizer *tokenizer.Tokenizer
}

// NewConverter creates a new Converter instance with IPADIC dictionary.
func NewConverter() (*Converter, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	
	return &Converter{
		tokenizer: t,
	}, nil
}

// Convert converts the input text according to the specified mode.
func (c *Converter) Convert(text string, mode ConversionMode) (string, error) {
	// Split text into sentences
	sentences := SplitSentences(text)
	
	var convertedSentences []string
	
	for _, sentence := range sentences {
		var converted string
		var err error
		
		switch mode {
		case CasualToPolite:
			converted, err = c.convertCasualToPolite(sentence)
		case PoliteToCasual:
			converted, err = c.convertPoliteToCasual(sentence)
		default:
			return "", fmt.Errorf("unsupported conversion mode: %d", mode)
		}
		
		if err != nil {
			return "", err
		}
		
		convertedSentences = append(convertedSentences, converted)
	}
	
	// Join sentences back together
	return strings.Join(convertedSentences, ""), nil
}
