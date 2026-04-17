package process

import (
	"errors"

	"github.com/enescakir/emoji"
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

// Run Lets go!
func Run(tag any, notes string) error {
	// tag check
	if tag == "" {
		log.Error().Msgf("%s Missing required --tag flag", emoji.Bomb.String())
		return errors.New("missing required --tag flag")
	}

	// notes check
	if len(notes) <= 0 {
		log.Error().Msgf("%s Notes are too short", emoji.Bomb.String())
		return errors.New("notes are too short")
	}

	return nil
}
