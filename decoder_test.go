/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"bytes"
	"io/ioutil"
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
	event1 := &EventRecord{
		Tag:  "BIRT",
		Date: "31 DEC 1997",
		Place: PlaceRecord{
			Name: "The place",
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

	if diff := cmp.Diff(i1.Event[0], event1, sourceOpt); diff != "" {
		t.Errorf("Individual 0, event 0 mismatch (-want +got):\n%s", diff)
	}

	if len(i1.Attribute) != 15 {
		t.Fatalf(`Individual 0 had %d attributes, expected 15`, len(i1.Attribute))
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

	if diff := cmp.Diff(i1.Attribute[0], att1, sourceOpt); diff != "" {
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

	if diff := cmp.Diff(g.Submitter, submitters); diff != "" {
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

	families := []*FamilyRecord{
		{
			Xref:    "FAMILY1",
			Husband: &IndividualRecord{Xref: "PERSON1"},
			Wife:    &IndividualRecord{Xref: "PERSON2"},
			Child: []*IndividualRecord{
				{Xref: "PERSON3"},
				{Xref: "PERSON4"},
			},
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

	if diff := cmp.Diff(g.Family, families, indOpt, eventOpt); diff != "" {
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
			Xref:  "SOURCE1",
			Title: "A bmp picture",
			Note: []*NoteRecord{
				{
					Note: "A note about whatever\nNote continued here. The word TEST should not be broken!",
				},
				{
					Note: "A note\nNote continued here. The word TEST should not be broken!",
				},
				{
					Note: "A note about the family\nNote continued here. The word TEST should not be broken!",
				},
				{
					Note: "A note\nNote continued here. The word TEST should not be broken!",
				},
			},
		},
	}

	if diff := cmp.Diff(g.Source, sources); diff != "" {
		t.Errorf("source mismatch (-want +got):\n%s", diff)
	}
}
