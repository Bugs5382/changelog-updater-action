package process

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
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

const (
	diffFromLabel = "CHANGELOG.md (before)"
	diffToLabel   = "CHANGELOG.md (after)"
)

// logDiff emits a unified diff between the original and updated
// CHANGELOG.md contents via log.Info. Powered by hexops/gotextdiff, which
// is the same Myers-based diff implementation gopls uses internally, so
// the output is a standard unified diff familiar to any developer.
func logDiff(original, updated string) {
	if original == updated {
		log.Info().Msgf("%s Diff: no changes detected", emoji.Memo.String())
		return
	}

	edits := myers.ComputeEdits(span.URIFromPath("CHANGELOG.md"), original, updated)
	unified := fmt.Sprint(gotextdiff.ToUnified(diffFromLabel, diffToLabel, original, edits))

	log.Info().Msgf("%s Diff:", emoji.Memo.String())

	// Emit each diff line as its own log entry so it flows nicely through
	// zerolog's line-oriented output (console or JSON).
	scanner := bufio.NewScanner(strings.NewReader(unified))
	for scanner.Scan() {
		log.Info().Msg(scanner.Text())
	}
}
