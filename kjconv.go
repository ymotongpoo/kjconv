// Package kjconv provides Japanese text style conversion between casual (常体) and polite (敬体) forms.
package kjconv

import (
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
	// TODO: Implement conversion logic
	return text, nil
}
