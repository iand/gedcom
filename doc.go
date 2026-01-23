/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

/*
Package gedcom provides functions to parse and produce GEDCOM files.

GEDCOM (Genealogical Data Communication) is a standard format used for
exchanging genealogical data between software applications. This package
includes functionality for both parsing existing GEDCOM files and generating
new ones.

# Decoding GEDCOM Files

The package provides a streaming [Decoder] for reading GEDCOM files. Use
[NewDecoder] to create a decoder that reads from an [io.Reader]:

	data, err := os.ReadFile("family.ged")
	if err != nil {
		log.Fatal(err)
	}

	d := gedcom.NewDecoder(bytes.NewReader(data))
	g, err := d.Decode()
	if err != nil {
		log.Fatal(err)
	}

	for _, ind := range g.Individual {
		if len(ind.Name) > 0 {
			fmt.Println(ind.Name[0].Name)
		}
	}

The decoder is streaming and can handle large files without loading the entire
contents into memory.

# Encoding GEDCOM Files

The package also provides an [Encoder] for generating GEDCOM files. Use
[NewEncoder] to create an encoder that writes to an [io.Writer]:

	g := &gedcom.Gedcom{
		Header: &gedcom.Header{
			SourceSystem: gedcom.SystemRecord{
				Xref:        "MyApp",
				ProductName: "My Application",
			},
			CharacterSet: "UTF-8",
		},
		Individual: []*gedcom.IndividualRecord{
			{
				Xref: "I1",
				Name: []*gedcom.NameRecord{
					{Name: "John /Doe/"},
				},
				Sex: "M",
			},
		},
		Trailer: &gedcom.Trailer{},
	}

	f, err := os.Create("output.ged")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := gedcom.NewEncoder(f)
	if err := enc.Encode(g); err != nil {
		log.Fatal(err)
	}

# Data Model

The [Gedcom] struct is the top-level container returned by the decoder and
accepted by the encoder. It contains slices of records for individuals,
families, sources, and other GEDCOM record types.

[IndividualRecord] represents a person and contains their names, sex, life
events (birth, death, etc.), family links, and citations.

[FamilyRecord] represents a family unit and links to husband, wife, and
children as [IndividualRecord] pointers.

[EventRecord] is a flexible type used for both events (birth, death, marriage)
and attributes (occupation, residence). The Tag field indicates the event type.

[SourceRecord] and [CitationRecord] handle source citations for genealogical
claims.

# Name Parsing

The [SplitPersonalName] helper function parses GEDCOM-formatted names:

	parsed := gedcom.SplitPersonalName("John \"Jack\" /Smith/ Jr.")
	// parsed.Given = "John"
	// parsed.Nickname = "Jack"
	// parsed.Surname = "Smith"
	// parsed.Suffix = "Jr."

# User-Defined Tags

GEDCOM allows custom tags prefixed with an underscore. These are captured in
[UserDefinedTag] slices on most record types, preserving vendor-specific
extensions.

# Specification Coverage

This package implements approximately 80% of the GEDCOM 5.5 specification,
which is sufficient for parsing about 99% of real-world GEDCOM files. It has
not been extensively tested with non-ASCII character sets.
*/
package gedcom
