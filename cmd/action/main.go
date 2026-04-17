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
	"os"

	"github.com/Bugs5382/changelog-updater-action/internal/logging"
	"github.com/Bugs5382/changelog-updater-action/internal/process"
	"github.com/enescakir/emoji"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

var (
	Version = "local"
	Gitsha  = "Unknown"
)

func main() {

	// setup flags
	tag := flag.StringP("tag", "t", "", "The release tag name")
	notes := flag.StringP("notes", "n", "", "The release notes body")
	_ = flag.BoolP("diff", "d", false, "Show the diff (if any) of changes")
	verbose := flag.BoolP("verbose", "v", false, "Enable debug level logging")

	// parse
	flag.Parse()

	// setup logging
	logging.Init(*verbose)

	// start
	log.Info().Msgf("%s Changelog Updater Action by Shane", emoji.Information.String())
	log.Debug().Msgf("%s Version: %s", emoji.Construction.String(), Version)
	log.Debug().Msgf("%s Build SHA: %s", emoji.Construction.String(), Gitsha)

	// process
	err := process.Run(*tag, *notes)
	if err != nil {
		log.Error().Msgf("%s Update failed: %s", emoji.Bomb.String(), err)
		os.Exit(1)
	}

}
