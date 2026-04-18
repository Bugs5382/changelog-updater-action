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

const baseContent = "# 📝 Changelog Updater Action :: Notes (Unit Testing)\n"

// notesV100 is the initial v1.0.0 release notes in Release Drafter format.
// shiftHeaders() prepends "##" to every heading, so # → ### and ## → ####.
const notesV100 = `# What Changed 👀

## 🐛 Bug Fixes

- fix: nil pointer on empty tag input @Bugs5382 (#1)

# Extra

**Full Changelog**: https://github.com/Bugs5382/changelog-updater-action/compare/v0.0.0...v1.0.0`

// notesV110 is the v1.1.0 release notes — adds a features section.
const notesV110 = `# What Changed 👀

## 🚀 Features

- feat: add --dry flag to skip writing changes @Bugs5382 (#7)

## 🐛 Bug Fixes

- fix: incorrect header shift when notes contain nested headers @Bugs5382 (#9)

# Extra

**Full Changelog**: https://github.com/Bugs5382/changelog-updater-action/compare/v1.0.0...v1.1.0`

// notesV110SameDatePatch is the v1.1.0 notes patched on the same release date.
const notesV110SameDatePatch = `# What Changed 👀

## 🚀 Features

- feat: add --dry flag to skip writing changes @Bugs5382 (#7)

## 🐛 Bug Fixes

- fix: incorrect header shift when notes contain nested headers @Bugs5382 (#9)
- fix: preserve trailing newline on changelog write @Bugs5382 (#11)

# Extra

**Full Changelog**: https://github.com/Bugs5382/changelog-updater-action/compare/v1.0.0...v1.1.0`

// notesV110NewDate is the v1.1.0 notes revised on a later date, adding a
// dependency update section and an additional feature.
const notesV110NewDate = `# What Changed 👀

## 🚀 Features

- feat: add --dry flag to skip writing changes @Bugs5382 (#7)
- feat: support custom date via --date flag @Bugs5382 (#13)

## 🐛 Bug Fixes

- fix: incorrect header shift when notes contain nested headers @Bugs5382 (#9)

## 🧩 Dependency Updates

- chore: bump zerolog to v1.35.0 @Bugs5382 (#15)

# Extra

**Full Changelog**: https://github.com/Bugs5382/changelog-updater-action/compare/v1.0.0...v1.1.0`
