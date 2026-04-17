package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

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

type Options struct {
	Tag    string
	Notes  string
	Path   string
	Date   string
	DryRun bool
}

// Run Lets go!
func Run(opts Options) error {
	// tag check
	if opts.Tag == "" {
		log.Error().Msgf("%s Missing required --tag flag", emoji.Bomb.String())
		return errors.New("missing required --tag flag")
	}

	// notes check
	if len(opts.Notes) <= 0 {
		log.Error().Msgf("%s Notes are too short", emoji.Bomb.String())
		return errors.New("notes are too short")
	}

	opts.Notes = shiftHeaders(opts.Notes)

	targetFile := filepath.Join(opts.Path, "CHANGELOG.md")

	log.Debug().Msgf("%s Target File: %s", emoji.Construction.String(), targetFile)

	content, err := os.ReadFile(targetFile)
	if err != nil {
		log.Error().Msgf("%s could not read file %s: %s", emoji.Bomb.String(), targetFile, err)
		return fmt.Errorf("could not read file %s: %w", targetFile, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	replacing := false
	foundVersion := false
	targetHeader := "## " + opts.Tag // e.g., "## v0.1.0"

	for _, line := range lines {
		if replacing {
			if strings.HasPrefix(line, "## ") {
				replacing = false
				newLines = append(newLines, line) // Keep this next version header
			}
			continue // Skip the old notes
		}

		// Using HasPrefix allows us to match "## v0.1.0" even if it already has a date
		if strings.HasPrefix(line, targetHeader) {
			foundVersion = true
			replacing = true // Start skipping subsequent lines until the next ##

			// Rewrite the header with the specified date
			header := fmt.Sprintf("## %s - %s", opts.Tag, opts.Date)
			newLines = append(newLines, header)

			// Inject the new notes with proper markdown spacing
			newLines = append(newLines, "")
			newLines = append(newLines, strings.TrimSpace(opts.Notes))
			newLines = append(newLines, "")

			log.Debug().Msgf("%s Found target version header: %s", emoji.TestTube.String(), header)
			continue
		}

		newLines = append(newLines, line)
	}

	// Edge Case: What if this is a brand-new release and the header isn't in the file yet?
	// You might want to inject it right after the root "# Changelog" header here.
	if !foundVersion {
		log.Trace().Msgf("%s Version header not found in file. Ensure the file has a matching header: %s", emoji.TestTube.String(), targetHeader)
		// For now, returning an error, but you could implement an insert-at-top logic here!
		return fmt.Errorf("version %s not found in changelog", opts.Tag)
	}

	// Write back to the file
	output := strings.Join(newLines, "\n")

	// os.FileMode 0644 is standard for text files (read/write for owner, read for others)
	err = os.WriteFile(targetFile, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}

	return nil
}

// shiftHeaders
func shiftHeaders(text string) string {
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		hashCount := 0
		// Count leading '#' characters
		for _, ch := range line {
			if ch == '#' {
				hashCount++
			} else {
				break
			}
		}

		// If the line starts with '#' and is followed by a space, it's a header
		if hashCount > 0 && len(line) > hashCount && line[hashCount] == ' ' {
			lines[i] = "#" + line
		}
	}

	return strings.Join(lines, "\n")
}
