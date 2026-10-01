// card-impact indexes Go behavior code and reports conservative review targets.
package main

import (
	"encoding/json"
	"eraofarcane/internal/cardimpact"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", "game", "Go engine source directory")
	card := flag.String("card", "", "card number to inspect")
	file := flag.String("file", "", "changed source file, relative to -dir")
	flag.Parse()
	index, err := cardimpact.Build(*dir)
	if err != nil {
		fail(err)
	}
	queryFile := ""
	if *file != "" {
		queryFile = *file
		if !filepath.IsAbs(queryFile) {
			queryFile = filepath.Join(*dir, queryFile)
		}
	}
	report, err := index.Analyze(*card, queryFile)
	if err != nil {
		fail(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		fail(err)
	}
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
