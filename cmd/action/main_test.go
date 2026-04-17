package main

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
	"testing"

	"github.com/Bugs5382/changelog-updater-action/internal/process"
	"github.com/rs/zerolog"
)

func TestRun(t *testing.T) {

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	tests := []struct {
		name    string
		version string
		notes   string
		wantErr bool
	}{
		{
			name:    "valid input",
			version: "v1.0.0",
			notes:   "Fixed a bug in the logger",
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			notes:   "Some notes here",
			wantErr: true,
		},
		{
			name:    "notes too short",
			version: "v1.1.0",
			notes:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := process.Run(tt.version, tt.notes)

			// Check if we got an error when we expected one (or vice versa)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
