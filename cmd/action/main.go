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
	"github.com/Bugs5382/changelog-updater-action/cmd/action/init/logging"
	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
)

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	logging.Init()

	log.Info().Msgf("%s Changelog Updater Action by Shane", logging.Emj(emoji.Information))
	log.Debug().Msgf("%s Version: %s", logging.Emj(emoji.Construction), Version)
	log.Debug().Msgf("%s Build SHA: %s", logging.Emj(emoji.Construction), Gitsha)
}
