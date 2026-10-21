//go:build ignore

// run this using go run unhandledtags.go <file.ged>
//
// It decodes the GEDCOM file and logs every tag the decoder does not handle to
// standard output, one line per tag occurrence.
package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/iand/gedcom"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run unhandledtags.go <file.ged>")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "open gedcom file: %v\n", err)
		os.Exit(1)
	}

	d := gedcom.NewDecoder(bytes.NewReader(data))
	d.LogUnhandledTags(os.Stdout)

	if _, err := d.Decode(); err != nil {
		fmt.Fprintf(os.Stderr, "decode gedcom: %v\n", err)
		os.Exit(1)
	}
}
