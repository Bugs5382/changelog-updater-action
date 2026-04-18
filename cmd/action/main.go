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
	"time"

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
	diff := flag.Bool("diff", false, "Show the diff (if any) of changes")
	dry := flag.Bool("dry", false, "Dry run, make no changes")
	path := flag.StringP("path", "p", ".", "Directory relative to root containing CHANGELOG.md")
	date := flag.String("date", time.Now().Format("2006-01-02"), "Release date")
	verbose := flag.BoolP("verbose", "v", false, "Enable debug level logging")

	// parse
	flag.Parse()

	// setup logging
	logging.Init(*verbose)

	opts := process.Options{
		Tag:    *tag,
		Notes:  *notes,
		Diff:   *diff,
		DryRun: *dry,
		Path:   *path,
		Date:   *date,
	}

	// start
	log.Info().Msgf("%s Changelog Updater Action by Shane", emoji.Information.String())
	log.Debug().Msgf("%s Version: %s", emoji.Construction.String(), Version)
	log.Debug().Msgf("%s Build SHA: %s", emoji.Construction.String(), Gitsha)

	// process
	err := process.Run(opts)
	if err != nil {
		log.Error().Msgf("%s Update failed: %s", emoji.Bomb.String(), err)
		os.Exit(1)
	}

	log.Info().Msgf("%s Changelog updated successfully", emoji.Rocket.String())

}
