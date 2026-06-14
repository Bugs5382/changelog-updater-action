package process

/*
ISC License

Copyright (c) 2026 Shane & Contributors

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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

type Options struct {
	Tag     string
	Notes   string
	Diff    bool
	DryRun  bool
	Path    string
	Date    string
	Cleanup bool

	// ChangelogBody is a manual override for the version's changelog body.
	// When non-empty it is written verbatim and the drafted Notes are
	// ignored, providing an escape hatch for when automatic generation
	// (e.g. Release Drafter) yields "No changes" — for instance after a
	// history rewrite orphans the PR/commit associations. Unlike Notes, the
	// body is not passed through shiftHeaders: the maintainer controls the
	// exact markdown that appears under the version header.
	ChangelogBody string
}

// ErrCleanupExclusive is returned when --cleanup is combined with any of
// the version-editing flags (--tag, --notes, --date). Exported so callers
// and tests can assert against it via errors.Is.
var ErrCleanupExclusive = errors.New("--cleanup cannot be combined with --tag, --notes, --changelog-body, or --date")

// validateCleanupExclusivity enforces that --cleanup is not combined with
// any of the version-editing flags. It returns a descriptive error naming
// every offending flag so the user sees them all at once instead of
// discovering them one failed run at a time.
func validateCleanupExclusivity(opts Options) error {
	if !opts.Cleanup {
		return nil
	}
	var offenders []string
	if opts.Tag != "" {
		offenders = append(offenders, "--tag")
	}
	if opts.Notes != "" {
		offenders = append(offenders, "--notes")
	}
	if opts.ChangelogBody != "" {
		offenders = append(offenders, "--changelog-body")
	}
	if opts.Date != "" {
		offenders = append(offenders, "--date")
	}
	if len(offenders) == 0 {
		return nil
	}
	return fmt.Errorf("%w: got %s", ErrCleanupExclusive, strings.Join(offenders, ", "))
}

// Run Let's go!
func Run(opts Options) error {
	// --cleanup is a standalone mode: it must not be combined with --tag,
	// --notes, or --date. Validate before branching so the error surfaces
	// regardless of which code path would have run next.
	if err := validateCleanupExclusivity(opts); err != nil {
		log.Error().Msgf("%s %s", emoji.Bomb.String(), err)
		return err
	}

	if opts.Cleanup {
		return runCleanup(opts)
	}

	// Outside of cleanup mode, --date defaults to today when the caller
	// didn't supply one. main.go now passes an empty string when the user
	// didn't set --date, so we fill the default in here.
	if opts.Date == "" {
		opts.Date = time.Now().Format("2006-01-02")
	}

	// tag check
	if opts.Tag == "" {
		return errors.New("missing required --tag flag")
	}

	// Resolve the body to write under the version header. A manual
	// --changelog-body takes precedence over the drafted --notes: when it's
	// set we write it verbatim and skip both the notes check and the header
	// shift, letting a maintainer repair a release whose generated notes are
	// empty or wrong without another tag/push cycle. When it's absent we fall
	// back to the existing behavior: require notes and demote their headers so
	// they nest cleanly under the version header.
	if opts.ChangelogBody != "" {
		opts.Notes = opts.ChangelogBody
	} else {
		if len(opts.Notes) <= 0 {
			return errors.New("notes are too short")
		}
		opts.Notes = shiftHeaders(opts.Notes)
	}

	targetFile := filepath.Join(opts.Path, "CHANGELOG.md")

	log.Debug().Msgf("%s Target File: %s", emoji.Construction.String(), targetFile)

	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", targetFile, err)
	}

	original := string(content)

	lines := strings.Split(original, "\n")
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

	// Brand-new release: the version header isn't in the file yet. Insert it
	// at the top of the changelog, right after the first top-level "# " title
	// (if present), otherwise at the very top.
	if !foundVersion {
		log.Debug().Msgf("%s Version %s not found; inserting new entry at top", emoji.TestTube.String(), opts.Tag)

		header := fmt.Sprintf("## %s - %s", opts.Tag, opts.Date)
		// Leading blank line guarantees separation from the title above.
		entry := []string{"", header, "", strings.TrimSpace(opts.Notes), ""}

		insertAt := 0
		for i, line := range newLines {
			if strings.HasPrefix(line, "# ") {
				// Insert immediately after the title line.
				insertAt = i + 1
				// If the title is already followed by a blank line, skip it so
				// we don't end up with two consecutive blanks.
				if insertAt < len(newLines) && strings.TrimSpace(newLines[insertAt]) == "" {
					// Drop our entry's own leading blank since one already exists.
					entry = entry[1:]
				}
				break
			}
		}

		// Splice entry into newLines at insertAt.
		merged := make([]string, 0, len(newLines)+len(entry))
		merged = append(merged, newLines[:insertAt]...)
		merged = append(merged, entry...)
		merged = append(merged, newLines[insertAt:]...)
		newLines = merged
	}

	// Write back to the file
	output := strings.Join(newLines, "\n")
	output = normalizeSpacing(output)

	// When --diff is set, log a human-readable diff between the original file
	// contents and the freshly generated output before we do anything else.
	if opts.Diff {
		logDiff(original, output)
	}

	// Respect --dry: parse, diff, and log without touching the file.
	if opts.DryRun {
		log.Info().Msgf("%s Dry run enabled; %s left unchanged", emoji.TestTube.String(), targetFile)
		return nil
	}

	// os.FileMode 0644 is standard for text files (read/write for owner, read for others)
	err = os.WriteFile(targetFile, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}

	return nil
}
