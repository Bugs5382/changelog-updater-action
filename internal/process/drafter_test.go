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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// The testdata/*.md fixtures are release-drafter v7 bodies (the `body`
// output). They live in Markdown files because the canonical category titles
// carry emoji, which the repo does not allow in Go source.
//
//   - release-drafter-v7.md: captured from a release-drafter v7.7.0 dry run
//     against this repository with the canonical .github/release-drafter.yml.
//   - release-drafter-v7-no-baseline.md: captured from the same dry run with
//     every published release filtered out, so v7 found no comparison
//     baseline. v7 then lists no changes and appends its warning block.
//   - release-drafter-v7-collapsed.md: built with the v7.7.0 changelog layout
//     for a category over its collapse-after limit (a <details> block whose
//     summary reads "N changes").
//
// testdata/golden/*.md holds the CHANGELOG.md each fixture must produce.

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return string(data)
}

func TestReleaseDrafterV7(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Parallel()

	tests := []struct {
		name    string
		fixture string
		tag     string
		date    string
	}{
		{
			name:    "body with a comparison baseline",
			fixture: "release-drafter-v7.md",
			tag:     "v0.4.1",
			date:    "2026-09-26",
		},
		{
			name:    "first release without a comparison baseline",
			fixture: "release-drafter-v7-no-baseline.md",
			tag:     "v1.0.0",
			date:    "2026-09-26",
		},
		{
			name:    "collapsed category",
			fixture: "release-drafter-v7-collapsed.md",
			tag:     "v1.1.0",
			date:    "2026-10-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := writeChangelog(t, baseContent)

			err := Run(Options{
				Tag:   tt.tag,
				Notes: readFixture(t, tt.fixture),
				Date:  tt.date,
				Path:  dir,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := readChangelog(t, dir)
			want := readFixture(t, filepath.Join("golden", tt.fixture))
			if got != want {
				t.Errorf("CHANGELOG.md mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}

func TestReleaseDrafterV7ReplacesExistingEntry(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Parallel()

	// A second run for the same version (a later push before the release is
	// published) replaces the entry rather than adding a second one.
	dir := writeChangelog(t, baseContent)
	for i := 0; i < 2; i++ {
		err := Run(Options{
			Tag:   "v0.4.1",
			Notes: readFixture(t, "release-drafter-v7.md"),
			Date:  "2026-09-26",
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("run %d: unexpected error: %v", i+1, err)
		}
	}

	got := readChangelog(t, dir)
	want := readFixture(t, filepath.Join("golden", "release-drafter-v7.md"))
	if got != want {
		t.Errorf("CHANGELOG.md mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestChangelogBodyKeepsBaselineWarning(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Parallel()

	// changelog-body is written verbatim, so nothing is stripped from it.
	body := readFixture(t, "release-drafter-v7-no-baseline.md")
	dir := writeChangelog(t, baseContent)

	err := Run(Options{
		Tag:           "v1.0.0",
		ChangelogBody: body,
		Date:          "2026-09-26",
		Path:          dir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readChangelog(t, dir)
	if !strings.Contains(got, baselineWarningMarker) {
		t.Errorf("expected changelog-body to be written verbatim, got:\n%s", got)
	}
}

func TestNotesOnlyBaselineWarning(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Parallel()

	// Notes that are nothing but the warning block leave no body to write.
	dir := writeChangelog(t, baseContent)
	err := Run(Options{
		Tag:   "v1.0.0",
		Notes: "---\n> [!WARNING]\n> " + baselineWarningMarker + " x.\n\n---\n",
		Date:  "2026-09-26",
		Path:  dir,
	})
	if err == nil {
		t.Fatal("expected an error for notes that only hold the baseline warning")
	}
	if got := readChangelog(t, dir); got != baseContent {
		t.Errorf("expected CHANGELOG.md to be untouched, got:\n%s", got)
	}
}

func TestStripBaselineWarning(t *testing.T) {
	t.Parallel()

	noBaseline := readFixture(t, "release-drafter-v7-no-baseline.md")

	tests := []struct {
		name         string
		in           string
		wantStripped bool
		want         string
	}{
		{
			name:         "v7 warning block is removed",
			in:           noBaseline,
			wantStripped: true,
			want: "# What Changed\n\n* No changes\n\n# Extra\n\n" +
				"**Full Changelog**: https://github.com/Bugs5382/changelog-updater-action/compare/...v0.0.1\n",
		},
		{
			name:         "body without a warning is unchanged",
			in:           "# What Changed\n\n- fix: a thing (#1)\n",
			wantStripped: false,
			want:         "# What Changed\n\n- fix: a thing (#1)\n",
		},
		{
			name:         "warning without the closing rule is left alone",
			in:           "notes\n\n---\n> [!WARNING]\n> " + baselineWarningMarker + " x.\n",
			wantStripped: false,
			want:         "notes\n\n---\n> [!WARNING]\n> " + baselineWarningMarker + " x.\n",
		},
		{
			name:         "text after the block is kept",
			in:           "notes\n\n---\n> [!WARNING]\n> " + baselineWarningMarker + " x.\n\n---\n\nfooter\n",
			wantStripped: true,
			want:         "notes\n\nfooter\n",
		},
		{
			name:         "CRLF line endings",
			in:           "notes\r\n\r\n---\r\n> [!WARNING]\r\n> " + baselineWarningMarker + " x.\r\n\r\n---\r\n",
			wantStripped: true,
			want:         "notes\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, stripped := stripBaselineWarning(tt.in)
			if stripped != tt.wantStripped {
				t.Errorf("stripped = %v, want %v", stripped, tt.wantStripped)
			}
			// The fixture's title carries an emoji; compare without it.
			got = strings.ReplaceAll(got, " \U0001F440", "")
			if got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
