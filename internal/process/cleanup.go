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

/*
ISC License

Copyright (c) 2026 Shane

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

// runCleanup re-orders every version section inside CHANGELOG.md so that
// releases appear in descending semver order. Non-version "## " sections
// (anything whose first token isn't parseable as a semver) are preserved
// at the top, in their original relative order — this keeps blocks like
// "## [Unreleased]" exactly where authors expect them.
func runCleanup(opts Options) error {
	targetFile := filepath.Join(opts.Path, "CHANGELOG.md")
	log.Debug().Msgf("%s Cleanup target file: %s", emoji.Construction.String(), targetFile)

	content, err := os.ReadFile(targetFile)
	if err != nil {
		log.Error().Msgf("%s could not read file %s: %s", emoji.Bomb.String(), targetFile, err)
		return fmt.Errorf("could not read file %s: %w", targetFile, err)
	}

	original := string(content)
	lines := strings.Split(original, "\n")

	// Split the file into: preamble (everything before the first "## "),
	// followed by a list of sections each beginning with a "## " heading.
	var preamble []string
	type section struct {
		header string
		body   []string
		ver    semver
		isVer  bool
	}
	var sections []section

	i := 0
	for ; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			break
		}
		preamble = append(preamble, lines[i])
	}
	for i < len(lines) {
		if !strings.HasPrefix(lines[i], "## ") {
			// Should not happen on a well-formed split, but guard anyway.
			i++
			continue
		}
		header := lines[i]
		i++
		var body []string
		for i < len(lines) && !strings.HasPrefix(lines[i], "## ") {
			body = append(body, lines[i])
			i++
		}
		ver, ok := parseHeaderVersion(header)
		sections = append(sections, section{
			header: header,
			body:   body,
			ver:    ver,
			isVer:  ok,
		})
	}

	// Partition: non-version sections stay at the top in their original
	// relative order; version sections are sorted descending.
	var nonVersion []section
	var versioned []section
	for _, s := range sections {
		if s.isVer {
			versioned = append(versioned, s)
		} else {
			nonVersion = append(nonVersion, s)
			log.Debug().Msgf("%s Cleanup: preserving non-version section %q", emoji.TestTube.String(), s.header)
		}
	}

	sort.SliceStable(versioned, func(a, b int) bool {
		// Descending: larger version first.
		return compareSemver(versioned[a].ver, versioned[b].ver) > 0
	})

	// Reassemble.
	out := make([]string, 0, len(lines))
	out = append(out, preamble...)
	for _, s := range nonVersion {
		out = append(out, s.header)
		out = append(out, s.body...)
	}
	for _, s := range versioned {
		out = append(out, s.header)
		out = append(out, s.body...)
	}

	output := normalizeSpacing(strings.Join(out, "\n"))

	if opts.Diff {
		logDiff(original, output)
	}

	if opts.DryRun {
		log.Info().Msgf("%s Dry run enabled; %s left unchanged", emoji.TestTube.String(), targetFile)
		return nil
	}

	if err := os.WriteFile(targetFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}
	return nil
}

// parseHeaderVersion extracts the first whitespace-delimited token after
// "## " and tries to parse it as a semver. It tolerates a leading "v",
// surrounding brackets (e.g. "## [1.2.3]"), and trailing " - 2026-01-01"
// style date suffixes that this tool itself writes.
func parseHeaderVersion(header string) (semver, bool) {
	text := strings.TrimSpace(strings.TrimPrefix(header, "## "))
	if text == "" {
		return semver{}, false
	}
	// Take the first whitespace-separated token.
	token := strings.Fields(text)[0]
	// Strip common wrappers users put around versions.
	token = strings.Trim(token, "[]()")
	return parseSemver(token)
}

// semver is a minimal SemVer 2.0.0 representation sufficient for ordering
// pre-releases. Build metadata (after '+') is parsed but ignored for
// comparison, per the spec.
type semver struct {
	raw        string
	major      uint64
	minor      uint64
	patch      uint64
	preRelease []string // split on '.'; empty slice == stable release
}

// parseSemver parses strings like "v1.2.3", "1.2.3-rc.1", "v2.0.0-beta",
// "1.2.3+build.5", "1.2". A missing patch (or minor) defaults to 0 so
// that truncated tags still sort sensibly.
func parseSemver(s string) (semver, bool) {
	if s == "" {
		return semver{}, false
	}
	v := semver{raw: s}
	core := strings.TrimPrefix(s, "v")
	core = strings.TrimPrefix(core, "V")

	// Split off build metadata (ignored for ordering).
	if idx := strings.Index(core, "+"); idx >= 0 {
		core = core[:idx]
	}
	// Split off pre-release.
	var pre string
	if idx := strings.Index(core, "-"); idx >= 0 {
		pre = core[idx+1:]
		core = core[:idx]
	}
	parts := strings.Split(core, ".")
	if len(parts) == 0 || len(parts) > 3 {
		return semver{}, false
	}
	nums := make([]uint64, 3)
	for i, p := range parts {
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			return semver{}, false
		}
		nums[i] = n
	}
	v.major, v.minor, v.patch = nums[0], nums[1], nums[2]
	if pre != "" {
		v.preRelease = strings.Split(pre, ".")
	}
	return v, true
}

// compareSemver returns -1, 0, or 1 per SemVer 2.0.0 precedence rules.
// Stable > pre-release; within pre-releases numeric identifiers compare
// numerically and have lower precedence than alphanumeric identifiers;
// a longer set of fields wins when all prior fields are equal.
func compareSemver(a, b semver) int {
	if c := cmpUint(a.major, b.major); c != 0 {
		return c
	}
	if c := cmpUint(a.minor, b.minor); c != 0 {
		return c
	}
	if c := cmpUint(a.patch, b.patch); c != 0 {
		return c
	}
	// Stable release > any pre-release.
	switch {
	case len(a.preRelease) == 0 && len(b.preRelease) == 0:
		return 0
	case len(a.preRelease) == 0:
		return 1
	case len(b.preRelease) == 0:
		return -1
	}
	n := len(a.preRelease)
	if len(b.preRelease) < n {
		n = len(b.preRelease)
	}
	for i := 0; i < n; i++ {
		if c := cmpPreReleaseField(a.preRelease[i], b.preRelease[i]); c != 0 {
			return c
		}
	}
	return cmpInt(len(a.preRelease), len(b.preRelease))
}

func cmpPreReleaseField(a, b string) int {
	aNum, aIsNum := parseUintStrict(a)
	bNum, bIsNum := parseUintStrict(b)
	switch {
	case aIsNum && bIsNum:
		return cmpUint(aNum, bNum)
	case aIsNum: // numeric identifiers have lower precedence
		return -1
	case bIsNum:
		return 1
	default:
		return strings.Compare(a, b)
	}
}

func parseUintStrict(s string) (uint64, bool) {
	if s == "" {
		return 0, false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	n, err := strconv.ParseUint(s, 10, 64)
	return n, err == nil
}

func cmpUint(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}
