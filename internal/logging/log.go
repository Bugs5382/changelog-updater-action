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
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"
	"unicode"

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
		return
	}

	// Default JSON output strips emojis from every log line.
	log.Logger = log.Output(&emojiStripper{w: os.Stderr})
}

// emojiStripper is an io.Writer that removes emoji (and another symbol)
// runes from the bytes written to it before delegating to the underlying
// writer. Intended to sit in front of zerolog's JSON output so structured
// logs stay ASCII-friendly.
type emojiStripper struct {
	w io.Writer
}

func (e *emojiStripper) Write(p []byte) (int, error) {
	cleaned := cleanJSONLine(p)
	if _, err := e.w.Write(cleaned); err != nil {
		return 0, err
	}
	// Report the original length so zerolog is happy.
	return len(p), nil
}

// cleanJSONLine strips emojis from every JSON string value in a single
// zerolog log line and left-trims whitespace from the "message" field.
// It preserves the original key order by streaming tokens instead of
// unmarshalling into a map. If the payload isn't valid JSON it falls back
// to a plain string strip.
func cleanJSONLine(p []byte) []byte {
	// Preserve a trailing newline if present so output stays line-delimited.
	trailing := ""
	trimmed := p
	if len(p) > 0 && p[len(p)-1] == '\n' {
		trailing = "\n"
		trimmed = p[:len(p)-1]
	}

	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.UseNumber()

	var out bytes.Buffer
	out.Grow(len(trimmed))

	// Track whether the next string value belongs to the "message" key so
	// we can additionally left-trim it.
	var (
		inObject     int
		expectKey    = true
		nextIsMsgVal bool
		needComma    bool
	)

	writeSep := func() {
		if needComma {
			out.WriteByte(',')
		}
	}

	for {
		tok, err := dec.Token()
		if err != nil {
			// EOF or malformed: fall back to rune-level strip on the whole payload.
			if out.Len() == 0 {
				return append([]byte(stripEmoji(string(trimmed))), trailing...)
			}
			break
		}

		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '{':
				writeSep()
				out.WriteByte('{')
				inObject++
				expectKey = true
				needComma = false
			case '}':
				out.WriteByte('}')
				inObject--
				needComma = true
			case '[':
				writeSep()
				out.WriteByte('[')
				expectKey = false
				needComma = false
			case ']':
				out.WriteByte(']')
				needComma = true
			}
		case string:
			writeSep()
			if inObject > 0 && expectKey {
				// This token is a key.
				keyBytes, _ := json.Marshal(v)
				out.Write(keyBytes)
				out.WriteByte(':')
				nextIsMsgVal = v == "message"
				expectKey = false
				needComma = false
			} else {
				// This token is a string value.
				cleaned := stripEmoji(v)
				if nextIsMsgVal {
					cleaned = strings.TrimLeft(cleaned, " \t")
					nextIsMsgVal = false
				}
				valBytes, _ := json.Marshal(cleaned)
				out.Write(valBytes)
				if inObject > 0 {
					expectKey = true
				}
				needComma = true
			}
		case json.Number:
			writeSep()
			out.WriteString(v.String())
			if inObject > 0 {
				expectKey = true
			}
			needComma = true
		case bool:
			writeSep()
			if v {
				out.WriteString("true")
			} else {
				out.WriteString("false")
			}
			if inObject > 0 {
				expectKey = true
			}
			needComma = true
		case nil:
			writeSep()
			out.WriteString("null")
			if inObject > 0 {
				expectKey = true
			}
			needComma = true
		}
	}

	return append(out.Bytes(), trailing...)
}

// stripEmoji removes runes that are part of Unicode "Symbol, Other" (So)
// or in the supplementary planes typically used by emoji. Safe to call on
// arbitrary UTF-8 strings — non-emoji content is returned unchanged.
func stripEmoji(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if isEmoji(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func isEmoji(r rune) bool {
	// Anything in the supplementary planes (>= U+1F000) is almost
	// exclusively emoji/pictographs in modern usage.
	if r >= 0x1F000 {
		return true
	}
	switch {
	case r >= 0x2100 && r <= 0x214F, // Letterlike Symbols (ℹ, ™, etc.)
		r >= 0x2190 && r <= 0x21FF, // Arrows
		r >= 0x2300 && r <= 0x23FF, // Misc Technical
		r >= 0x2500 && r <= 0x257F, // Box Drawing
		r >= 0x2580 && r <= 0x259F, // Block Elements
		r >= 0x25A0 && r <= 0x25FF, // Geometric Shapes
		r >= 0x2600 && r <= 0x27BF, // Misc Symbols + Dingbats
		r >= 0x2900 && r <= 0x297F, // Supplemental Arrows-B
		r >= 0x2B00 && r <= 0x2BFF, // Misc Symbols and Arrows
		r >= 0xFE00 && r <= 0xFE0F, // Variation selectors
		r == 0x200D,                // Zero-width joiner (emoji sequences)
		r == 0x20E3:                // Combining enclosing keycap
		return true
	}
	// Catch-all for category "Symbol, Other" — pictographic symbols.
	return unicode.Is(unicode.So, r)
}
