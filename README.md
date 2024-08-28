# gedcom

Go package to parse GEDCOM files.

[![Check Status](https://github.com/iand/gedcom/actions/workflows/check.yml/badge.svg?branch=master)](https://github.com/iand/gedcom/actions/workflows/check.yml)
[![Test Status](https://github.com/iand/gedcom/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/iand/gedcom/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iand/gedcom)](https://goreportcard.com/report/github.com/iand/gedcom)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/iand/gedcom)

## Purpose

The `gedcom` package provides tools for working with GEDCOM files in Go. GEDCOM (Genealogical Data Communication) is a standard format used for exchanging genealogical data between software applications. This package includes functionality for both parsing existing GEDCOM files and generating new ones.

The package includes a streaming decoder for reading GEDCOM files and an encoder for creating GEDCOM files from Go structs.

## Usage

The package provides a `Decoder` with a single `Decode` method that returns a Gedcom struct. Use the `NewDecoder` method to create a new decoder.

This example shows how to parse a GEDCOM file and list all the individuals. In this example the entire input file is read into memory, but the decoder is streaming so it should be able to deal with very large files: just pass an appropriate Reader.


	package main

	import (
		"bytes"
		"github.com/iand/gedcom"
		"io/ioutil"
	)

    func main() {
		data, _ := ioutil.ReadFile("testdata/kennedy.ged")

		d := gedcom.NewDecoder(bytes.NewReader(data))

		g, _ := d.Decode()

		for _, rec := range g.Individual {
			if len(rec.Name) > 0 {
				println(rec.Name[0].Name)
			}			
		}
	}

The structures produced by the Decoder are in [types.go](types.go) and correspond roughly 1:1 to the structures in the [GEDCOM specification](http://homepages.rootsweb.ancestry.com/~pmcbride/gedcom/55gctoc.htm).

This package does not implement the entire GEDCOM specification, I'm still working on it. It's about 80% complete which is enough for about 99% of GEDCOM files. It has not been extensively tested with non-ASCII character sets nor with pathological cases such as the [GEDCOM 5.5 Torture Test Files](http://www.geditcom.com/gedcom.html).

### Using the Encoder

In addition to decoding GEDCOM files, this package also provides an Encoder for generating GEDCOM files from the structs in [types.go](types.go). You can create an encoder using the `NewEncoder` method, which writes to an `io.Writer`.

To see an example of how to use the encoder, refer to [encoder_example.go](encoder_example.go). This example illustrates how to create individual and family records, populate them with data, and encode them into a valid GEDCOM file.

You can run the example using the following command:

```bash
go run encoder_example.go
```

## Installation

Simply run

Run the following in the directory containing your project's `go.mod` file:

```bash
go get github.com/iand/gedcom@latest
```

Documentation is at [https://pkg.go.dev/github.com/iand/gedcom](https://pkg.go.dev/github.com/iand/gedcom)

## Authors

* [Ian Davis](http://github.com/iand)


## Contributors


## Contributing

* Do submit your changes as a pull request
* Do your best to adhere to the existing coding conventions and idioms.
* Do run `go fmt` on the code before committing 
* Do feel free to add yourself to the [`CREDITS`](CREDITS) file and the
  corresponding Contributors list in the [`README.md`](README.md). 
  Alphabetical order applies.
* Don't touch the [`AUTHORS`](AUTHORS) file. An existing author will add you if 
  your contributions are significant enough.
* Do note that in order for any non-trivial changes to be merged (as a rule
  of thumb, additions larger than about 15 lines of code), an explicit
  Public Domain Dedication needs to be on record from you. Please include
  a copy of the statement found in the [`WAIVER`](WAIVER) file with your pull request

## License

This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.
