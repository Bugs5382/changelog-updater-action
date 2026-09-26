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
	"strings"

	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

// baselineWarningMarker is the text release-drafter v7 puts in the warning it
// appends to a release body when it finds no previous published release to
// compare against (the first release of a repository, or after a history
// rewrite). See #42.
const baselineWarningMarker = "Release Drafter could not find a previous"

// noChangesLine is release-drafter's default no-changes-template, which v7
// also emits when it has no comparison baseline.
const noChangesLine = "* No changes"

// stripBaselineWarning removes the release-drafter v7 missing-baseline
// warning from drafted notes. v7 wraps it in two "---" rules:
//
//	---
//	> [!WARNING]
//	> Release Drafter could not find a previous **published release** ...
//	...
//	---
//
// The block is advice for whoever publishes the draft, not release notes, and
// the "---" directly after the notes would turn the line above it into a
// setext heading in CHANGELOG.md. Anything around the block (a configured
// footer, for example) is kept. When the block is not in the expected shape
// the notes are returned untouched. The bool reports whether it stripped.
func stripBaselineWarning(notes string) (string, bool) {
	if !strings.Contains(notes, baselineWarningMarker) {
		log.Trace().Msg("Notes carry no release-drafter baseline warning")
		return notes, false
	}

	lines := strings.Split(strings.ReplaceAll(notes, "\r\n", "\n"), "\n")

	marker := -1
	for i, line := range lines {
		if strings.Contains(line, baselineWarningMarker) {
			marker = i
			break
		}
	}

	// Walk back over the quoted warning lines to the opening rule.
	start := -1
	for i := marker - 1; i >= 0; i-- {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "---" {
			start = i
			break
		}
		if trimmed != "" && !strings.HasPrefix(trimmed, ">") {
			break
		}
	}

	// The closing rule is the next "---" after the marker.
	end := -1
	for i := marker + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}

	if start < 0 || end < 0 {
		log.Warn().Msgf("%s Notes mention a missing release-drafter baseline but not in the expected --- block; leaving them unchanged", emoji.Warning.String())
		log.Debug().Msgf("Baseline warning shape: marker line %d, opening rule %d, closing rule %d", marker, start, end)
		return notes, false
	}

	log.Debug().Msgf("%s Removing release-drafter baseline warning (lines %d-%d of the notes)", emoji.Scissors.String(), start+1, end+1)

	before := strings.TrimRight(strings.Join(lines[:start], "\n"), "\n")
	after := strings.Trim(strings.Join(lines[end+1:], "\n"), "\n")

	out := before
	if after != "" {
		out += "\n\n" + after
	}
	return out + "\n", true
}

// prepareDraftedNotes applies the release-drafter specific clean-up to the
// drafted notes before they are nested under the version header, and logs
// what it found so a maintainer can tell why a changelog entry is thin.
func prepareDraftedNotes(tag, notes string) string {
	cleaned, stripped := stripBaselineWarning(notes)
	if stripped {
		log.Warn().Msgf("%s release-drafter found no previous published release to compare against, so its notes for %s list no changes. Set changelog-body to write the notes for this release by hand.", emoji.Warning.String(), tag)
	}

	for _, line := range strings.Split(cleaned, "\n") {
		if strings.TrimSpace(line) == noChangesLine {
			log.Warn().Msgf("%s The drafted notes for %s say %q", emoji.Warning.String(), tag, noChangesLine)
			break
		}
	}

	return cleaned
}
