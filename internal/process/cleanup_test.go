package process

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
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestRun_CleanupExclusivity(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "cleanup with tag",
			opts: Options{Cleanup: true, Tag: "v1.0.0"},
		},
		{
			name: "cleanup with notes",
			opts: Options{Cleanup: true, Notes: "some notes"},
		},
		{
			name: "cleanup with date",
			opts: Options{Cleanup: true, Date: "2026-01-01"},
		},
		{
			name: "cleanup with all three",
			opts: Options{Cleanup: true, Tag: "v1.0.0", Notes: "n", Date: "2026-01-01"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.opts)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, ErrCleanupExclusive) {
				t.Errorf("expected ErrCleanupExclusive, got %v", err)
			}
		})
	}
}

func TestRun_CleanupAloneIsAllowed(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	// Minimal well-formed changelog with versions out of order — cleanup
	// should not return an exclusivity error and should reorder them.
	initial := "# Title\n\n## v1.0.0 - 2026-01-01\n\n- first\n\n## v2.0.0 - 2026-02-01\n\n- second\n"
	dir := writeChangelog(t, initial)

	err := Run(Options{Cleanup: true, Path: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readChangelog(t, dir)
	v2 := strings.Index(got, "## v2.0.0")
	v1 := strings.Index(got, "## v1.0.0")
	if v2 < 0 || v1 < 0 {
		t.Fatalf("expected both versions to be present, got:\n%s", got)
	}
	if v2 >= v1 {
		t.Errorf("expected v2.0.0 before v1.0.0, got:\n%s", got)
	}
}

func TestParseSemver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		in    string
		ok    bool
		major uint64
		minor uint64
		patch uint64
		pre   []string
	}{
		{name: "plain", in: "1.2.3", ok: true, major: 1, minor: 2, patch: 3},
		{name: "v-prefix", in: "v1.2.3", ok: true, major: 1, minor: 2, patch: 3},
		{name: "capital V", in: "V1.2.3", ok: true, major: 1, minor: 2, patch: 3},
		{name: "missing patch", in: "v1.2", ok: true, major: 1, minor: 2, patch: 0},
		{name: "missing minor+patch", in: "v1", ok: true, major: 1, minor: 0, patch: 0},
		{name: "pre-release single", in: "v1.0.0-beta", ok: true, major: 1, pre: []string{"beta"}},
		{name: "pre-release dotted", in: "v1.0.0-rc.10", ok: true, major: 1, pre: []string{"rc", "10"}},
		{name: "build metadata ignored", in: "1.0.0+build.5", ok: true, major: 1},
		{name: "pre-release + build", in: "1.0.0-rc.1+build.5", ok: true, major: 1, pre: []string{"rc", "1"}},

		{name: "empty", in: "", ok: false},
		{name: "non-numeric core", in: "vX.Y.Z", ok: false},
		{name: "too many parts", in: "1.2.3.4", ok: false},
		{name: "unreleased-label", in: "[Unreleased]", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseSemver(tt.in)
			if ok != tt.ok {
				t.Fatalf("parseSemver(%q) ok = %v, want %v", tt.in, ok, tt.ok)
			}
			if !tt.ok {
				return
			}
			if got.major != tt.major || got.minor != tt.minor || got.patch != tt.patch {
				t.Errorf("parseSemver(%q) core = %d.%d.%d, want %d.%d.%d",
					tt.in, got.major, got.minor, got.patch, tt.major, tt.minor, tt.patch)
			}
			if !equalStrings(got.preRelease, tt.pre) {
				t.Errorf("parseSemver(%q) pre = %v, want %v", tt.in, got.preRelease, tt.pre)
			}
		})
	}
}

func TestCompareSemver(t *testing.T) {
	t.Parallel()

	// Each case: a and b are raw strings; want is the expected sign of
	// compareSemver(a, b) (-1, 0, or 1).
	tests := []struct {
		name string
		a, b string
		want int
	}{
		// Core comparison.
		{name: "major greater", a: "v2.0.0", b: "v1.9.9", want: 1},
		{name: "minor greater", a: "v1.2.0", b: "v1.1.99", want: 1},
		{name: "patch greater", a: "v1.0.2", b: "v1.0.1", want: 1},
		{name: "equal", a: "v1.2.3", b: "1.2.3", want: 0},
		{name: "padded missing patch equal", a: "v1.2", b: "v1.2.0", want: 0},

		// Stable vs. pre-release (spec: stable > pre-release).
		{name: "stable > rc", a: "v1.0.0", b: "v1.0.0-rc.1", want: 1},
		{name: "stable > beta", a: "v1.0.0", b: "v1.0.0-beta", want: 1},

		// Pre-release ordering rules.
		// Numeric identifiers compare numerically (so rc.10 > rc.2).
		{name: "numeric pre identifier ordering", a: "v1.0.0-rc.10", b: "v1.0.0-rc.2", want: 1},
		// Numeric identifiers have lower precedence than alphanumerics.
		{name: "alphanumeric > numeric pre field", a: "v1.0.0-rc.alpha", b: "v1.0.0-rc.1", want: 1},
		// Alphanumeric compares lexically (so rc > beta > alpha).
		{name: "rc > beta", a: "v1.0.0-rc.1", b: "v1.0.0-beta.1", want: 1},
		{name: "beta > alpha", a: "v1.0.0-beta", b: "v1.0.0-alpha", want: 1},
		// Larger field-count wins when all prior fields equal.
		{name: "more pre fields wins", a: "v1.0.0-rc.1.0", b: "v1.0.0-rc.1", want: 1},

		// Build metadata MUST be ignored for precedence.
		{name: "build metadata ignored", a: "v1.0.0+build.5", b: "v1.0.0+build.1", want: 0},
		{name: "build metadata ignored with pre", a: "v1.0.0-rc.1+x", b: "v1.0.0-rc.1+y", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av, ok := parseSemver(tt.a)
			if !ok {
				t.Fatalf("parseSemver(%q) failed", tt.a)
			}
			bv, ok := parseSemver(tt.b)
			if !ok {
				t.Fatalf("parseSemver(%q) failed", tt.b)
			}

			got := sign(compareSemver(av, bv))
			if got != tt.want {
				t.Errorf("compareSemver(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}

			// Anti-symmetry: compareSemver(b, a) must be the negation.
			rev := sign(compareSemver(bv, av))
			if rev != -tt.want {
				t.Errorf("compareSemver(%q, %q) = %d, want %d (anti-symmetry)", tt.b, tt.a, rev, -tt.want)
			}
		})
	}
}

// TestCompareSemver_Sort locks in the full descending order of a mixed
//
//	set, so we catch regressions in the ordering transitively.
func TestCompareSemver_Sort(t *testing.T) {
	t.Parallel()

	in := []string{
		"v1.0.0",
		"v2.0.0-rc.1",
		"v2.0.0",
		"v1.0.1",
		"v2.0.0-beta",
		"v2.0.0-rc.10",
		"v2.0.0-rc.2",
		"v1.0.0-alpha",
	}

	want := []string{
		"v2.0.0",       // stable beats every 2.0.0 pre-release
		"v2.0.0-rc.10", // numeric identifier: 10 > 2
		"v2.0.0-rc.2",
		"v2.0.0-rc.1",
		"v2.0.0-beta", // lexical: rc > beta
		"v1.0.1",
		"v1.0.0",
		"v1.0.0-alpha",
	}

	parsed := make([]semver, len(in))
	for i, s := range in {
		v, ok := parseSemver(s)
		if !ok {
			t.Fatalf("parseSemver(%q) failed", s)
		}
		parsed[i] = v
	}

	sort.SliceStable(parsed, func(i, j int) bool {
		return compareSemver(parsed[i], parsed[j]) > 0
	})

	got := make([]string, len(parsed))
	for i, v := range parsed {
		got[i] = v.raw
	}
	if !equalStrings(got, want) {
		t.Errorf("sort order mismatch:\n got:  %v\n want: %v", got, want)
	}
}

func TestParseHeaderVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
		wantOK bool
		wantV  string // raw version token we expect to be extracted
	}{
		{name: "bare", header: "## v1.2.3", wantOK: true, wantV: "v1.2.3"},
		{name: "with date suffix", header: "## v1.2.3 - 2026-01-01", wantOK: true, wantV: "v1.2.3"},
		{name: "bracketed", header: "## [1.2.3] - 2026-01-01", wantOK: true, wantV: "1.2.3"},
		{name: "pre-release with date", header: "## v2.0.0-rc.1 - 2026-02-02", wantOK: true, wantV: "v2.0.0-rc.1"},
		{name: "unreleased", header: "## [Unreleased]", wantOK: false},
		{name: "arbitrary heading", header: "## Some arbitrary section", wantOK: false},
		{name: "empty after prefix", header: "## ", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := parseHeaderVersion(tt.header)
			if ok != tt.wantOK {
				t.Fatalf("parseHeaderVersion(%q) ok = %v, want %v", tt.header, ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if v.raw != tt.wantV {
				t.Errorf("parseHeaderVersion(%q) raw = %q, want %q", tt.header, v.raw, tt.wantV)
			}
		})
	}
}

// TestRunCleanup_EndToEnd verifies that cleanup mode reorders version
// sections descending, leaves non-version "## " sections at the top in
// their original relative order, and preserves the file preamble.
func TestRunCleanup_EndToEnd(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	initial := strings.Join([]string{
		"# 📝 Changelog",
		"",
		"Some preamble paragraph that must survive cleanup.",
		"",
		"## [Unreleased]",
		"",
		"- in-progress item",
		"",
		"## v1.0.0 - 2026-01-01",
		"",
		"- first release",
		"",
		"## v2.0.0 - 2026-03-01",
		"",
		"- second release",
		"",
		"## v2.0.0-rc.1 - 2026-02-15",
		"",
		"- rc build",
		"",
		"## v1.5.0 - 2026-02-01",
		"",
		"- minor bump",
		"",
	}, "\n")

	dir := writeChangelog(t, initial)

	if err := Run(Options{Cleanup: true, Path: dir}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readChangelog(t, dir)

	// Preamble preserved.
	if !strings.Contains(got, "# 📝 Changelog") {
		t.Errorf("expected top-level title preserved, got:\n%s", got)
	}
	if !strings.Contains(got, "Some preamble paragraph that must survive cleanup.") {
		t.Errorf("expected preamble paragraph preserved, got:\n%s", got)
	}

	// Non-version section stays above all versioned ones.
	unreleased := strings.Index(got, "## [Unreleased]")
	if unreleased < 0 {
		t.Fatalf("expected [Unreleased] preserved, got:\n%s", got)
	}

	// Descending order: v2.0.0 > v2.0.0-rc.1 > v1.5.0 > v1.0.0, all below [Unreleased].
	expectOrder := []string{
		"## [Unreleased]",
		"## v2.0.0 - 2026-03-01",
		"## v2.0.0-rc.1 - 2026-02-15",
		"## v1.5.0 - 2026-02-01",
		"## v1.0.0 - 2026-01-01",
	}
	prev := -1
	for _, marker := range expectOrder {
		idx := strings.Index(got, marker)
		if idx < 0 {
			t.Fatalf("expected %q in output, got:\n%s", marker, got)
		}
		if idx <= prev {
			t.Errorf("section %q out of order (idx=%d, prev=%d) in:\n%s", marker, idx, prev, got)
		}
		prev = idx
	}
}

// --- tiny local helpers kept here so cleanup tests don't leak symbols ---

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
