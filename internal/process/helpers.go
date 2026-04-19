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

import "strings"

// normalizeSpacing rewrites the changelog so every markdown heading has
// exactly one blank line before it (except when it's the first line of
// the file), and collapses any runs of 2+ consecutive blank lines into a
// single blank line. It also guarantees the file ends with exactly one
// trailing newline.
func normalizeSpacing(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))

	isHeading := func(s string) bool {
		trimmed := strings.TrimLeft(s, " \t")
		return strings.HasPrefix(trimmed, "# ") ||
			strings.HasPrefix(trimmed, "## ") ||
			strings.HasPrefix(trimmed, "### ") ||
			strings.HasPrefix(trimmed, "#### ") ||
			strings.HasPrefix(trimmed, "##### ") ||
			strings.HasPrefix(trimmed, "###### ")
	}

	for _, line := range lines {
		if isHeading(line) {
			// Ensure exactly one blank line before the heading (unless at start).
			// Trim any trailing blanks we already emitted.
			for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
				out = out[:len(out)-1]
			}
			if len(out) > 0 {
				out = append(out, "")
			}
			out = append(out, line)
			continue
		}

		// Collapse multiple consecutive blank lines into a single one.
		if strings.TrimSpace(line) == "" {
			if len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
				continue
			}
		}
		out = append(out, line)
	}

	// Trim trailing blank lines, then guarantee one final newline.
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
		out = out[:len(out)-1]
	}
	return strings.Join(out, "\n") + "\n"
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
			lines[i] = "##" + line
		}
	}

	return strings.Join(lines, "\n")
}
