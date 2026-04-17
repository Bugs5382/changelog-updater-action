package process

import (
	"testing"

	"github.com/rs/zerolog"
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

	t.Run("basic -- nothing in file", func(t *testing.T) {})

}
