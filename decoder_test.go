/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data []byte

func init() {
	var err error
	data, err = ioutil.ReadFile("testdata/allged.ged")
	if err != nil {
		panic(err)
	}
}

func TestStructuresAreInitialized(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("Result of decoding gedcom gave error %v, expected no error", err)
	}

	if g == nil {
		t.Fatalf("Result of decoding gedcom was nil, expected valid object")
	}
	if g.Individual == nil {
		t.Fatalf("Individual list was nil, expected valid slice")
	}

	if g.Family == nil {
		t.Fatalf("Family list was nil, expected valid slice")
	}

	if g.Media == nil {
		t.Fatalf("Media list was nil, expected valid slice")
	}

	if g.Repository == nil {
		t.Fatalf("Repository list was nil, expected valid slice")
	}

	if g.Source == nil {
		t.Fatalf("Source list was nil, expected valid slice")
	}

	if g.Submitter == nil {
		t.Fatalf("Submitter list was nil, expected valid slice")
	}
}

func TestIndividual(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a comparison option that ignores events
	eventOpt := cmp.Comparer(func(a, b []*EventRecord) bool {
		return true
	})

	// Create a comparison option that compares just names
	nameOpt := cmp.Comparer(func(a, b *NameRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		return a.Name == b.Name
	})

	// Create a comparison option that compares families by xref
	familyOpt := cmp.Comparer(func(a, b *FamilyLinkRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		if a.Family == nil {
			return b.Family == nil
		}

		if b.Family == nil {
			return a.Family == nil
		}

		return a.Family.Xref == b.Family.Xref
	})

	// Create a comparison option that compares citations by source xref only
	sourceOpt := cmp.Comparer(func(a, b *CitationRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		if a.Source == nil {
			return b.Source == nil
		}

		if b.Source == nil {
			return a.Source == nil
		}

		return a.Source.Xref == b.Source.Xref
	})

	// Create a comparison option that compares media files by name only
	fileOpt := cmp.Comparer(func(a, b *MediaRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		if len(a.File) == 0 {
			return len(b.File) == 0
		}

		if len(b.File) == 0 {
			return len(a.File) == 0
		}

		return a.File[0].Name == b.File[0].Name
	})

	individuals := []*IndividualRecord{
		{
			Xref: "PERSON1",
			Sex:  "M",
			Name: []*NameRecord{
				{
					Name: "given name /surname/jr.",
				},
				{
					Name: "another name /surname/",
				},
			},
			Family: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY1"},
				},
				{
					Family: &FamilyRecord{Xref: "FAMILY2"},
				},
			},
			Parents: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "PARENTS"},
				},
				{
					Family: &FamilyRecord{Xref: "ADOPTIVE_PARENTS"},
				},
			},
			Citation: []*CitationRecord{
				{
					Source: &SourceRecord{
						Xref: "SOURCE1",
					},
				},
			},

			Change: ChangeRecord{
				Date: "1 APR 1998",
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
			Note: []*NoteRecord{
				{
					Note: "A note about the inidvidual\nNote continued here. The word TEST should not be broken!",
				},
			},
			Media: []*MediaRecord{
				{
					File: []*FileRecord{
						{
							Name: `\\network\drive\path\file name.gif`,
						},
					},
				},
			}, UserDefined: []UserDefinedTag{
				{Tag: "_MYOWNTAG", Value: "This is a non-standard tag. Not recommended but allowed", Level: 1},
			},
		},
		{
			Xref: "PERSON2",
			Name: []*NameRecord{
				{
					Name: "/Wife/",
				},
			},
			Sex: "F",
			Family: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY1"},
				},
			},
		},
		{
			Xref: "PERSON3",
			Name: []*NameRecord{
				{
					Name: "/Child 1/",
				},
			},
			Parents: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY1"},
				},
			},
		},
		{
			Xref: "PERSON4",
			Name: []*NameRecord{
				{
					Name: "/Child 2/",
				},
			},
			Parents: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY1"},
				},
			},
		},
		{
			Xref: "PERSON5",
			Sex:  "M",
			Name: []*NameRecord{
				{
					Name: "/Father/",
				},
			},
			Family: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "PARENTS"},
				},
			},
		},
		{
			Xref: "PERSON6",
			Name: []*NameRecord{
				{
					Name: "/Adoptive mother/",
				},
			},
			Sex: "F",
			Family: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "ADOPTIVE_PARENTS"},
				},
			},
		},
		{
			Xref: "PERSON7",
			Name: []*NameRecord{
				{
					Name: "/Child 3/",
				},
			},
			Parents: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY2"},
				},
			},
		},
		{
			Xref: "PERSON8",
			Name: []*NameRecord{
				{
					Name: "/2nd Wife/",
				},
			},
			Sex: "F",
			Family: []*FamilyLinkRecord{
				{
					Family: &FamilyRecord{Xref: "FAMILY2"},
				},
			},
		},
	}

	if diff := cmp.Diff(individuals, g.Individual, eventOpt, familyOpt, nameOpt, sourceOpt, fileOpt); diff != "" {
		t.Errorf("submitter mismatch (-want +got):\n%s", diff)
	}
}

func TestIndividualDetail(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Individual) != 8 {
		t.Fatalf("Individual list length was %d, expected 8", len(g.Individual))
	}

	i1 := g.Individual[0]

	if i1.Xref != "PERSON1" {
		t.Errorf(`Individual 0 xref was "%s", expected @PERSON1@`, i1.Xref)
	}

	if i1.Sex != "M" {
		t.Errorf(`Individual 0 sex "%s" names, expected "M"`, i1.Sex)
	}

	if len(i1.Name) != 2 {
		t.Fatalf(`Individual 0 had %d names, expected 2`, len(i1.Name))
	}

	// Create a comparison option that compares sources by xref only
	sourceOpt := cmp.Comparer(func(a, b *SourceRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		return a.Xref == b.Xref
	})

	name1 := &NameRecord{
		Name: "given name /surname/jr.",
		Citation: []*CitationRecord{
			{
				Source: &SourceRecord{
					Xref: "SOURCE1",
				},

				Page: "42",
				Quay: "0",
				Data: DataRecord{
					Date: "BEF 1 JAN 1900",
					Text: []string{
						"a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			{
				Note: "Personal Name note\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	if diff := cmp.Diff(i1.Name[0], name1, sourceOpt); diff != "" {
		t.Errorf("Individual 0, name 0 mismatch (-want +got):\n%s", diff)
	}

	if len(i1.Event) != 24 {
		t.Fatalf(`Individual 0 had %d events, expected 24`, len(i1.Event))
	}

	// Create a comparison option that compares families by xref
	familyOpt := cmp.Comparer(func(a, b *FamilyRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		return a.Xref == b.Xref
	})

	event1 := &EventRecord{
		Tag:  "BIRT",
		Date: "31 DEC 1997",
		Place: PlaceRecord{
			Name: "The place",
		},

		ChildInFamily: &FamilyRecord{
			Xref: "PARENTS",
		},
		Citation: []*CitationRecord{
			{
				Source: &SourceRecord{
					Xref: "SOURCE1",
				},
				Page: "42",
				Quay: "2",
				Data: DataRecord{
					Date: "31 DEC 1900",
					Text: []string{
						"a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			{
				Note: "BIRTH event note (the event of entering into life)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	if diff := cmp.Diff(i1.Event[0], event1, sourceOpt, familyOpt); diff != "" {
		t.Errorf("Individual 0, event 0 mismatch (-want +got):\n%s", diff)
	}

	if len(i1.Attribute) != 14 {
		t.Fatalf(`Individual 0 had %d attributes, expected 14`, len(i1.Attribute))
	}

	att1 := &EventRecord{
		Tag:   "CAST",
		Value: "Cast name",
		Date:  "31 DEC 1997",
		Place: PlaceRecord{
			Name: "The place",
		},
		Citation: []*CitationRecord{
			{
				Source: &SourceRecord{
					Xref: "SOURCE1",
				},
				Page: "42",
				Quay: "3",
				Data: DataRecord{
					Date: "31 DEC 1900",
					Text: []string{
						"a sample text\nSample text continued here. The word TEST should not be broken!",
					},
				},
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
		},
		Note: []*NoteRecord{
			{
				Note: "CASTE event note (the name of an individual's rank or status in society, based   on racial or religious differences, or differences in wealth, inherited   rank, profession, occupation, etc)\nNote continued here. The word TEST should not be broken!",
			},
		},
	}

	if diff := cmp.Diff(att1, i1.Attribute[0], sourceOpt); diff != "" {
		t.Errorf("Individual 0, attribute 0 mismatch (-want +got):\n%s", diff)
	}

	if len(i1.Parents) != 2 {
		t.Fatalf(`Individual 0 had %d parent families, expected 2`, len(i1.Parents))
	}
}

func TestSubmitter(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	submitters := []*SubmitterRecord{{}}

	if diff := cmp.Diff(submitters, g.Submitter); diff != "" {
		t.Errorf("submitter mismatch (-want +got):\n%s", diff)
	}
}

func TestFamily(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a comparison option that compares individuals by xref only
	indOpt := cmp.Comparer(func(a, b *IndividualRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		return a.Xref == b.Xref
	})

	// Create a comparison option that compares events by tag and date
	eventOpt := cmp.Comparer(func(a, b *EventRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		return a.Tag == b.Tag && a.Date == b.Date
	})

	// Create a comparison option that compares citations by source xref
	citeOpt := cmp.Comparer(func(a, b *CitationRecord) bool {
		if a == nil {
			return b == nil
		}

		if b == nil {
			return a == nil
		}

		if a.Source == nil {
			return b.Source == nil
		}

		if b.Source == nil {
			return a.Source == nil
		}

		return a.Source.Xref == b.Source.Xref
	})

	families := []*FamilyRecord{
		{
			Xref:    "FAMILY1",
			Husband: &IndividualRecord{Xref: "PERSON1"},
			Wife:    &IndividualRecord{Xref: "PERSON2"},
			Child: []*IndividualRecord{
				{Xref: "PERSON3"},
				{Xref: "PERSON4"},
			},
			NumberOfChildren: "42",
			Event: []*EventRecord{
				{Tag: "ANUL", Date: "31 DEC 1997"},
				{Tag: "CENS", Date: "31 DEC 1997"},
				{Tag: "DIV", Date: "31 DEC 1997"},
				{Tag: "DIVF", Date: "31 DEC 1997"},
				{Tag: "ENGA", Date: "31 DEC 1997"},
				{Tag: "MARR", Date: "31 DEC 1997"},
				{Tag: "MARB", Date: "31 DEC 1997"},
				{Tag: "MARC", Date: "31 DEC 1997"},
				{Tag: "MARL", Date: "31 DEC 1997"},
				{Tag: "MARS", Date: "31 DEC 1997"},
				{Tag: "EVEN", Date: "31 DEC 1997"},
			},
			Citation: []*CitationRecord{
				{
					Source: &SourceRecord{Xref: "SOURCE1"},
				},
			},
			Change: ChangeRecord{
				Date: "1 APR 1998",
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
			Note: []*NoteRecord{
				{
					Note: "A note about the family\nNote continued here. The word TEST should not be broken!",
				},
			},
			Media: []*MediaRecord{
				{
					File: []*FileRecord{
						{
							Name:   `\\network\drive\path\file name.bmp`,
							Format: "bmp",
							Title:  "A bmp picture",
						},
					},
					Note: []*NoteRecord{
						{
							Note: "A note\nNote continued here. The word TEST should not be broken!",
						},
					},
				},
			},
			UserDefined: []UserDefinedTag{
				{Tag: "_MYOWNTAG", Value: "This is a non-standard tag. Not recommended but allowed", Level: 1},
			},
		},
		{
			Xref:    "PARENTS",
			Husband: &IndividualRecord{Xref: "PERSON5"},
			Child: []*IndividualRecord{
				{Xref: "PERSON1"},
			},
		},
		{
			Xref: "ADOPTIVE_PARENTS",
			Wife: &IndividualRecord{Xref: "PERSON6"},
			Child: []*IndividualRecord{
				{Xref: "PERSON1"},
			},
		},
		{
			Xref:    "FAMILY2",
			Husband: &IndividualRecord{Xref: "PERSON1"},
			Wife:    &IndividualRecord{Xref: "PERSON8"},
			Child: []*IndividualRecord{
				{Xref: "PERSON7"},
			},
		},
	}

	if diff := cmp.Diff(families, g.Family, indOpt, eventOpt, citeOpt); diff != "" {
		t.Errorf("family mismatch (-want +got):\n%s", diff)
	}
}

func TestSource(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sources := []*SourceRecord{
		{
			Xref: "SOURCE1",
			Data: &SourceDataRecord{
				Event: []*SourceEventRecord{
					{
						Kind:  "BIRT, CHR",
						Date:  "FROM 1 JAN 1980 TO 1 FEB 1982",
						Place: "Place",
					},

					{
						Kind:  "DEAT",
						Date:  "FROM 1 JAN 1980 TO 1 FEB 1982",
						Place: "Another place",
					},
				},
			},
			Title:            "Title of source\nTitle continued here. The word TEST should not be broken!",
			Originator:       "Author of source\nAuthor continued here. The word TEST should not be broken!",
			FiledBy:          "Short title",
			PublicationFacts: "Source publication facts\nPublication facts continued here. The word TEST should not be broken!",
			Text:             "Citation from source\nCitation continued here. The word TEST should not be broken!",
			Change: ChangeRecord{
				Date: "1 APR 1998",
				Note: []*NoteRecord{
					{
						Note: "A note\nNote continued here. The word TEST should not be broken!",
					},
				},
			},
			Note: []*NoteRecord{
				{
					Note: "A note about the family\nNote continued here. The word TEST should not be broken!",
				},
			},
			Media: []*MediaRecord{
				{
					File: []*FileRecord{
						{
							Name:   `\\network\drive\path\file name.bmp`,
							Format: "bmp",
							Title:  "A bmp picture",
						},
					},
					Note: []*NoteRecord{
						{
							Note: "A note\nNote continued here. The word TEST should not be broken!",
						},
					},
				},
			},
			UserDefined: []UserDefinedTag{
				{Tag: "_MYOWNTAG", Value: "This is a non-standard tag. Not recommended but allowed", Level: 1},
			},
		},
	}

	if diff := cmp.Diff(sources, g.Source); diff != "" {
		t.Errorf("source mismatch (-want +got):\n%s", diff)
	}
}

func TestHeader(t *testing.T) {
	d := NewDecoder(bytes.NewReader(data))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	header := &Header{
		SourceSystem: SystemRecord{
			Xref:         "APPROVED_SOURCE_NAME",
			Version:      "Version number of source-program",
			ProductName:  "Name of source-program",
			BusinessName: "Corporation name",
			Address: AddressRecord{
				Full:       "Corporation address line 1\nCorporation address line 2\nCorporation address line 3\nCorporation address line 4",
				Line1:      "Corporation address line 1",
				Line2:      "Corporation address line 2",
				City:       "Corporation address city",
				State:      "Corporation address state",
				PostalCode: "Corporation address ZIP code",
				Country:    "Corporation address country",
				Phone: []string{
					"Corporation phone number 1",
					"Corporation phone number 2",
					"Corporation phone number 3 (last one!)",
				},
			},
			SourceName:      "Name of source data",
			SourceDate:      "1 JAN 1998",
			SourceCopyright: "Copyright of source data",
		},
		Destination:         "Destination of transmission",
		Date:                "1 JAN 1998",
		Time:                "13:57:24.80",
		Submitter:           &SubmitterRecord{Xref: "SUBMITTER"},
		Submission:          &SubmissionRecord{Xref: "SUBMISSION"},
		Filename:            "ALLGED.GED",
		Copyright:           "(C) 1997-2000 by H. Eichmann. You can use and distribute this file freely as long as you do not charge for it",
		Version:             "5.5",
		Form:                "LINEAGE-LINKED",
		CharacterSet:        "ASCII",
		CharacterSetVersion: "Version number of ASCII (whatever it means) ",
		Language:            "language",
		Note: "A general note about this file:" + "\n" +
			"It demonstrates most of the data which can be submitted using GEDCOM5.5. It shows the relatives of PERSON1:" + "\n" +
			"His 2 wifes (PERSON2, PERSON8), his parents (father: PERSON5, mother not given), " + "\n" +
			"adoptive parents (mother: PERSON6, father not given) and his 3 children (PERSON3, PERSON4 and PERSON7)." + "\n" +
			"In PERSON1, FAMILY1, SUBMITTER, SUBMISSION and SOURCE1 as many datafields as possible are used." + "\n" +
			"All other individuals/families contain no data. Note, that many data tags can appear more than once" + "\n" +
			"(in this transmission this is demonstrated with tags: NAME, OCCU, PLACE and NOTE. Seek the word 'another'." + "\n" +
			"The data transmitted here do not make sence. Just the HEAD.DATE tag contains the date of the creation" + "\n" +
			"of this file and will change in future Versions!" + "\n" +
			"This file is created by H. Eichmann: h.eichmann@@gmx.de. Feel free to copy and use it for any " + "\n" +
			"non-commercial purpose. For the creation the GEDCOM standard Release 5.5 (2 JAN 1996) has been used." + "\n" +
			"Copyright: The church of Jesus Christ of latter-day saints, gedcom@@gedcom.org" + "\n" +
			"Download it (the GEDCOM 5.5 specs) from: ftp.gedcom.com/pub/genealogy/gedcom." + "\n" +
			"Some Specials: This line is very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long but not too long (255 caharcters is the limit). " + "\n" +
			"This @@ (commercial at) character may only appear ONCE!" + "\n" +
			"Note continued here. The word TEST should not be broken!",
		UserDefined: []UserDefinedTag{
			{Tag: "_MYOWNTAG", Value: "This is a non-standard tag. Not recommended but allowed", Level: 1},
		},
	}

	if diff := cmp.Diff(header, g.Header); diff != "" {
		t.Errorf("header mismatch (-want +got):\n%s", diff)
	}
}

func TestIndividualAlia(t *testing.T) {
	aliaData := []byte(`
0 @PERSON1@ INDI
1 SEX F
1 NAME Margaret /Smith/
1 ALIA Peggy
`)

	d := NewDecoder(bytes.NewReader(aliaData))

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Individual) == 0 {
		t.Fatalf("no individual was decoded")
	}

	individual := &IndividualRecord{
		Xref: "PERSON1",
		Name: []*NameRecord{
			{Name: "Margaret /Smith/"},
			{Name: "Peggy"}, // alias becomes alternate name
		},
		Sex: "F",
	}

	if diff := cmp.Diff(individual, g.Individual[0]); diff != "" {
		t.Errorf("individual mismatch (-want +got):\n%s", diff)
	}
}

func TestFixupAncestryBadNode(t *testing.T) {
	f, err := os.Open("testdata/badnote.ged")
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	d := NewDecoder(f)

	g, err := d.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Source) == 0 {
		t.Fatalf("no source was decoded")
	}

	want := &SourceRecord{
		Xref:             "S507927087",
		Title:            "London, England, Church of England Births and Baptisms, 1813-1917",
		Originator:       "Ancestry.com",
		PublicationFacts: "Ancestry.com Operations, Inc.",
		Note: []*NoteRecord{
			{
				Note: "Board of Guardian Records and Church of England Parish Registers. London Metropolitan Archives, London.\n<p>Images produced by permission of the City of London Corporation. The City of London gives no warranty as to the accuracy, completeness or fitness for the purpose of the information provided. Images may be used only for purposes of research, private study or education. Applications for any other use should be made to London Metropolitan Archives, 40 Northampton Road, London EC1R 0HB. Email -   ask.lma@@cityoflondon.gov.uk. Infringement of the above condition may result in legal action.</p>",
			},
		},
		UserDefined: []UserDefinedTag{
			{Tag: "_APID", Value: "1,1558::0", Level: 1},
		},
	}

	if diff := cmp.Diff(want, g.Source[0]); diff != "" {
		t.Errorf("source mismatch (-want +got):\n%s", diff)
	}
}
