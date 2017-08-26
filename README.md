# gedcom

Go package to parse GEDCOM files.

[![Build Status](https://travis-ci.org/iand/gedcom.svg?branch=master)](https://travis-ci.org/iand/gedcom)

## Usage

The package provides a Decoder with a single Decode method that returns a Gedcom struct. Use the NewDecoder method to create a new decoder.

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

This package does not implement the entire GEDCOM specification, I'm still working on it. It's about 80% complete which is enough for about 99% of GEDCOM files. It has not been extensively tested with non-ASCII character sets nor with pathological cases such as the [http://www.geditcom.com/gedcom.html](GEDCOM 5.5 Torture Test Files).

## Installation

Simply run

	go get github.com/iand/gedcom

Documentation is at [http://godoc.org/github.com/iand/gedcom](http://godoc.org/github.com/iand/gedcom)

## Authors

* [Ian Davis](http://github.com/iand) - <http://iandavis.com/>


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
