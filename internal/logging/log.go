package logging

/*
ISC License

Copyright (c) 2026 Shane

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init sets up logging configuration
func Init(verbose bool) {
	// If we are in a test, stop everything immediately.
	if isTest() {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return
	}

	setLogFormat()

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr != "" {
		if parsedLevel, err := zerolog.ParseLevel(strings.ToLower(levelStr)); err == nil {
			zerolog.SetGlobalLevel(parsedLevel)
			return
		}
	}

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func isTest() bool {
	return strings.HasSuffix(os.Args[0], ".test") ||
		strings.Contains(strings.Join(os.Args, " "), "-test.")
}

func setLogFormat() {
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))

	if format == "text" {
		// Otherwise, default to ConsoleWriter (Text)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	} else {
		return
	}
}

func setLogLevel() {
	levelStr := os.Getenv("LOG_LEVEL")

	if levelStr == "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		return
	}

	parsedLevel, err := zerolog.ParseLevel(levelStr)
	if err == nil {
		zerolog.SetGlobalLevel(parsedLevel)
		return
	}

	if levelInt, err := strconv.Atoi(levelStr); err == nil {
		zerolog.SetGlobalLevel(zerolog.Level(levelInt))
		return
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msgf("Invalid log level '%s', defaulting to info", levelStr)
}
