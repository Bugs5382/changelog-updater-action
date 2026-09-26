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
	"testing"

	"github.com/rs/zerolog"
)

// Regression tests for #43: the version header must match the tag exactly,
// not by prefix, or a tag like v1.0.1 overwrites the v1.0.10 entry.
func TestRunMatchesVersionHeaderExactly(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Parallel()

	tests := []struct {
		name    string
		initial string
		tag     string
		want    string
	}{
		{
			name: "v1.0.1 does not match v1.0.10",
			initial: "# Changelog\n\n## v1.0.10 - 2026-01-10\n\nten notes\n\n" +
				"## v1.0.1 - 2026-01-01\n\none notes\n",
			tag: "v1.0.1",
			want: "# Changelog\n\n## v1.0.10 - 2026-01-10\n\nten notes\n\n" +
				"## v1.0.1 - 2026-02-02\n\nnew notes\n",
		},
		{
			name: "v1.0.0 does not match v1.0.0-rc.1",
			initial: "# Changelog\n\n## v1.0.0-rc.1 - 2026-01-10\n\nrc notes\n\n" +
				"## v0.9.0 - 2026-01-01\n\nold notes\n",
			tag: "v1.0.0",
			want: "# Changelog\n\n## v1.0.0 - 2026-02-02\n\nnew notes\n\n" +
				"## v1.0.0-rc.1 - 2026-01-10\n\nrc notes\n\n" +
				"## v0.9.0 - 2026-01-01\n\nold notes\n",
		},
		{
			name:    "bare header still matches",
			initial: "# Changelog\n\n## v1.0.1\n\n## v1.0.0 - 2026-01-01\n\nold notes\n",
			tag:     "v1.0.1",
			want:    "# Changelog\n\n## v1.0.1 - 2026-02-02\n\nnew notes\n\n## v1.0.0 - 2026-01-01\n\nold notes\n",
		},
		{
			name:    "bracketed header still matches",
			initial: "# Changelog\n\n## [v1.0.1] - 2026-01-01\n\nold notes\n",
			tag:     "v1.0.1",
			want:    "# Changelog\n\n## v1.0.1 - 2026-02-02\n\nnew notes\n",
		},
		{
			name:    "dated header is replaced in place",
			initial: "# Changelog\n\n## v1.0.1 - 2026-01-01\n\nold notes\n",
			tag:     "v1.0.1",
			want:    "# Changelog\n\n## v1.0.1 - 2026-02-02\n\nnew notes\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := writeChangelog(t, tt.initial)

			err := Run(Options{
				Tag:   tt.tag,
				Notes: "new notes",
				Date:  "2026-02-02",
				Path:  dir,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := readChangelog(t, dir); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestIsVersionHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line string
		tag  string
		want bool
	}{
		{"## v1.0.1", "v1.0.1", true},
		{"## v1.0.1 - 2026-01-01", "v1.0.1", true},
		{"## [v1.0.1] - 2026-01-01", "v1.0.1", true},
		{"## (v1.0.1)", "v1.0.1", true},
		{"## v1.0.10 - 2026-01-01", "v1.0.1", false},
		{"## v1.0.1-rc.1", "v1.0.1", false},
		{"### v1.0.1", "v1.0.1", false},
		{"## ", "v1.0.1", false},
		{"v1.0.1", "v1.0.1", false},
	}

	for _, tt := range tests {
		if got := isVersionHeader(tt.line, tt.tag); got != tt.want {
			t.Errorf("isVersionHeader(%q, %q) = %v, want %v", tt.line, tt.tag, got, tt.want)
		}
	}
}
