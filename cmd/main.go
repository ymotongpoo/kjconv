// Command kjconv provides a command-line tool for Japanese text style conversion.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ymotongpoo/kjconv"
)

func main() {
	var (
		mode = flag.String("mode", "casual-to-polite", "Conversion mode: 'casual-to-polite' or 'polite-to-casual'")
		text = flag.String("text", "", "Text to convert")
	)
	flag.Parse()

	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if *text == "" {
		slog.Error("text is required")
		flag.Usage()
		os.Exit(1)
	}

	converter, err := kjconv.NewConverter()
	if err != nil {
		slog.Error("failed to create converter", "error", err)
		os.Exit(1)
	}

	var convMode kjconv.ConversionMode
	switch *mode {
	case "casual-to-polite":
		convMode = kjconv.CasualToPolite
	case "polite-to-casual":
		convMode = kjconv.PoliteToCasual
	default:
		slog.Error("invalid mode", "mode", *mode)
		os.Exit(1)
	}

	result, err := converter.Convert(*text, convMode)
	if err != nil {
		slog.Error("conversion failed", "error", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
