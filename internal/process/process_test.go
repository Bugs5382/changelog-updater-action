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

// writeChangelog CHANGELOG.md inside a fresh temp directory, returning the directory path.
func writeChangelog(t *testing.T, content string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "changelog-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp CHANGELOG.md: %v", err)
	}
	return dir
}

// readChangelog is a test helper that reads and returns CHANGELOG.md from dir.
func readChangelog(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("failed to read temp CHANGELOG.md: %v", err)
	}
	return string(data)
}

func TestRun(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	tests := []struct {
		name    string
		tag     string
		notes   string
		wantErr bool
	}{
		{
			name:    "empty tag",
			tag:     "",
			notes:   "Some notes here",
			wantErr: true,
		},
		{
			name:    "notes too short",
			tag:     "v1.1.0",
			notes:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{}

			opts.Tag = tt.tag
			opts.Notes = tt.notes

			err := Run(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	t.Run("step 1 -- add first version", func(t *testing.T) {
		// A brand-new changelog with just the top title — no version header yet.
		dir := writeChangelog(t, baseContent)

		err := Run(Options{
			Tag:   "v1.0.0",
			Notes: notesV100,
			Date:  "2026-01-01",
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := readChangelog(t, dir)

		if !strings.Contains(got, "## v1.0.0 - 2026-01-01") {
			t.Errorf("expected versioned header in output, got:\n%s", got)
		}
		if !strings.Contains(got, "### What Changed 👀") {
			t.Errorf("expected top-level section heading in output, got:\n%s", got)
		}
		if !strings.Contains(got, "#### 🐛 Bug Fixes") {
			t.Errorf("expected bug fixes section heading in output, got:\n%s", got)
		}
		if !strings.Contains(got, "fix: nil pointer on empty tag input") {
			t.Errorf("expected bug fix note in output, got:\n%s", got)
		}
	})

	t.Run("step 2 -- add second version", func(t *testing.T) {
		// v1.0.0 is already stamped; a bare v1.1.0 header is ready to be filled.
		initial := baseContent + "\n## v1.1.0\n\n## v1.0.0 - 2026-01-01\n\n" + notesV100 + "\n"
		dir := writeChangelog(t, initial)

		err := Run(Options{
			Tag:   "v1.1.0",
			Notes: notesV110,
			Date:  "2026-02-01",
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := readChangelog(t, dir)

		if !strings.Contains(got, "## v1.1.0 - 2026-02-01") {
			t.Errorf("expected v1.1.0 header in output, got:\n%s", got)
		}
		if !strings.Contains(got, "feat: add --dry flag to skip writing changes") {
			t.Errorf("expected v1.1.0 feature note in output, got:\n%s", got)
		}
		if !strings.Contains(got, "fix: incorrect header shift when notes contain nested headers") {
			t.Errorf("expected v1.1.0 bug fix note in output, got:\n%s", got)
		}
		// v1.0.0 must be preserved and ordered after v1.1.0.
		if !strings.Contains(got, "## v1.0.0 - 2026-01-01") {
			t.Errorf("expected v1.0.0 header still present, got:\n%s", got)
		}
		if !strings.Contains(got, "fix: nil pointer on empty tag input") {
			t.Errorf("expected v1.0.0 notes still present, got:\n%s", got)
		}
		v110Pos := strings.Index(got, "## v1.1.0")
		v100Pos := strings.Index(got, "## v1.0.0")
		if v110Pos >= v100Pos {
			t.Errorf("expected v1.1.0 to appear before v1.0.0 in output, got:\n%s", got)
		}
	})

	t.Run("step 3 -- modify latest version same date", func(t *testing.T) {
		initial := baseContent +
			"\n## v1.1.0 - 2026-02-01\n\n" + shiftHeaders(notesV110) + "\n\n" +
			"## v1.0.0 - 2026-01-01\n\n" + shiftHeaders(notesV100) + "\n"
		dir := writeChangelog(t, initial)

		err := Run(Options{
			Tag:   "v1.1.0",
			Notes: notesV110SameDatePatch,
			Date:  "2026-02-01", // same date
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := readChangelog(t, dir)

		if !strings.Contains(got, "## v1.1.0 - 2026-02-01") {
			t.Errorf("expected v1.1.0 header preserved with same date, got:\n%s", got)
		}
		if !strings.Contains(got, "fix: preserve trailing newline on changelog write") {
			t.Errorf("expected new patch fix note in output, got:\n%s", got)
		}
		// v1.0.0 must remain intact.
		if !strings.Contains(got, "## v1.0.0 - 2026-01-01") {
			t.Errorf("expected v1.0.0 header preserved, got:\n%s", got)
		}
	})

	t.Run("step 4 -- change date and notes of latest version", func(t *testing.T) {
		initial := baseContent +
			"\n## v1.1.0 - 2026-02-01\n\n" + shiftHeaders(notesV110SameDatePatch) + "\n\n" +
			"## v1.0.0 - 2026-01-01\n\n" + shiftHeaders(notesV100) + "\n"
		dir := writeChangelog(t, initial)

		err := Run(Options{
			Tag:   "v1.1.0",
			Notes: notesV110NewDate,
			Date:  "2026-03-15", // new date
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got := readChangelog(t, dir)

		// Header must carry the new date.
		if !strings.Contains(got, "## v1.1.0 - 2026-03-15") {
			t.Errorf("expected updated date in v1.1.0 header, got:\n%s", got)
		}
		// Old date must be gone from the v1.1.0 header.
		if strings.Contains(got, "## v1.1.0 - 2026-02-01") {
			t.Errorf("expected old date to be replaced, got:\n%s", got)
		}
		// New content unique to notesV110NewDate must be present.
		if !strings.Contains(got, "feat: support custom date via --date flag") {
			t.Errorf("expected new feature note in output, got:\n%s", got)
		}
		if !strings.Contains(got, "chore: bump zerolog to v1.35.0") {
			t.Errorf("expected dependency update note in output, got:\n%s", got)
		}
		if !strings.Contains(got, "#### 🧩 Dependency Updates") {
			t.Errorf("expected dependency updates section heading in output, got:\n%s", got)
		}
		// Content only in the old notes must be gone.
		if strings.Contains(got, "fix: preserve trailing newline on changelog write") {
			t.Errorf("expected removed patch note to be gone, got:\n%s", got)
		}
		// v1.0.0 must be untouched.
		if !strings.Contains(got, "## v1.0.0 - 2026-01-01") {
			t.Errorf("expected v1.0.0 header preserved, got:\n%s", got)
		}
		if !strings.Contains(got, "fix: nil pointer on empty tag input") {
			t.Errorf("expected v1.0.0 notes preserved, got:\n%s", got)
		}

	})

	t.Run("basic -- nothing in file", func(t *testing.T) {
		// With auto-insert, a missing version header is no longer an error —
		// it's inserted at the top. Use a completely empty file to force an error.
		dir := writeChangelog(t, "")

		err := Run(Options{
			Tag:   "v9.9.9",
			Notes: notesV100,
			Date:  "2026-01-01",
			Path:  dir,
		})
		if err != nil {
			t.Fatalf("unexpected error on empty file: %v", err)
		}

		got := readChangelog(t, dir)
		if !strings.Contains(got, "## v9.9.9 - 2026-01-01") {
			t.Errorf("expected v9.9.9 header after insert, got:\n%s", got)
		}
	})
}
