# gedcom

Go package to parse and produce GEDCOM files.

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

### Using the Encoder

In addition to decoding GEDCOM files, this package also provides an Encoder for generating GEDCOM files from the structs in [types.go](types.go). You can create an encoder using the `NewEncoder` method, which writes to an `io.Writer`.

By default the encoder produces GEDCOM 5.5 output. Use `SetVersion` to produce GEDCOM 7.0 output:

	enc := gedcom.NewEncoder(w)
	enc.SetVersion(gedcom.Gedcom70)

To see a full example of how to use the encoder, refer to [encoder_example.go](encoder_example.go).

## Compatibility

### GEDCOM 5.5/5.5.1

The decoder supports the core GEDCOM 5.5/5.5.1 record types: HEAD, INDI, FAM, SOUR, REPO, OBJE, SUBM, SUBN, and TRLR. Within these records, the commonly used tags are handled, covering the vast majority of real-world GEDCOM files. Tags that are not recognized by the decoder are preserved as `UserDefinedTag` values, so no data is silently discarded. The encoder can round-trip all decoded data, with the exception of `Submitter` references on individual records.

### GEDCOM 7.0

The decoder reads GEDCOM 7.0 files transparently, including files with a UTF-8 BOM. The following GEDCOM 7.0 tags are supported in addition to the 5.5/5.5.1 tags:

| Tag | Description | Context |
|-----|-------------|---------|
| `SNOTE` | Shared note record and references | Top-level record; references in individuals, families, sources, repositories, media, submitters, events, names |
| `SCHMA` | Schema for extension tag definitions | Header |
| `TAG` | Extension tag-to-URI mapping | Schema |
| `EXID` | External identifier | Individuals, families, sources, repositories, media, submitters |
| `UID` | Unique identifier | Individuals, families, sources, repositories, media, submitters |
| `CREA` | Record creation date | Individuals, families, sources, repositories, media, submitters, shared notes |
| `SDATE` | Sort date for ordering events | Events |
| `NO` | Non-event (assertion that an event did not occur) | Individuals, families |
| `CROP` | Image crop region | Media file |
| `TRAN` | Translation of text | Shared notes, notes, names, places |
| `MIME` | MIME type of text content | Shared notes, notes |
| `LANG` | BCP 47 language tag | Shared notes, notes, places |
| `RESN` | Restriction notice | Individuals, families, sources, names |
| `ROLE` | Role in an association (with `PHRASE` sub-tag) | Associations |
| `INIL` | Initiatory (LDS) | Individual events |
| `BAPL` | Baptism (LDS) | Individual events |
| `CONL` | Confirmation (LDS) | Individual events |
| `ENDL` | Endowment (LDS) | Individual events |
| `SLGC` | Sealing to parents (LDS) | Individual events |
| `SLGS` | Sealing to spouse (LDS) | Family events |
| `FACT` | Fact | Family events |
| `AGE` | Age at event | Events |

The encoder can produce GEDCOM 7.0 output via `SetVersion(Gedcom70)`. In this mode it writes a UTF-8 BOM, omits the `CHAR` and `GEDC.FORM` tags, sets `GEDC.VERS` to `7.0`, and writes long text lines without `CONC` splitting.

### Ancestry

[Ancestry](https://www.ancestry.com/) is a major producer of GEDCOM files. The decoder includes workarounds for known issues in Ancestry exports:

- **Malformed NOTE values**: Ancestry sometimes produces NOTE tags with embedded newlines in the value (i.e. line breaks that are not preceded by a CONT tag). The scanner detects and recovers from this by treating the continuation as part of the NOTE value.
- **Non-standard PUBL sub-tags**: Ancestry exports may include DATE and PLAC tags nested under the PUBL (publication facts) tag, which is not part of the GEDCOM specification. The decoder incorporates these values into the publication facts string rather than discarding them.

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
