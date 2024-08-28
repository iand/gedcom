//go:build ignore

// run this using go run encoder_example.go
package main

import (
	"log"
	"os"

	"github.com/iand/gedcom"
)

func main() {
	g := new(gedcom.Gedcom)

	sub := &gedcom.SubmitterRecord{
		Xref: "SUBM",
		Name: "John Doe",
	}

	g.Submitter = append(g.Submitter, sub)

	g.Header = &gedcom.Header{
		SourceSystem: gedcom.SystemRecord{
			Xref:       "gedcom",
			SourceName: "github.com/iand/gedcom",
		},
		Submitter:    sub,
		CharacterSet: "UTF-8",
		Language:     "English",
		Version:      "5.5.1",
		Form:         "LINEAGE-LINKED",
	}

	// Define the family record
	family := &gedcom.FamilyRecord{
		Xref: "F1",
	}
	// Add the family to the GEDCOM
	g.Family = append(g.Family, family)

	// Define the father individual record
	father := &gedcom.IndividualRecord{
		Xref: "I1",
		Name: []*gedcom.NameRecord{
			{
				Name: "John /Doe/",
			},
		},
		Sex: "M",
		Event: []*gedcom.EventRecord{
			{
				Tag:  "BIRT",
				Date: "1 JAN 1950",
				Place: gedcom.PlaceRecord{
					Name: "London, England",
				},
			},
		},
	}
	// Add the father to the GEDCOM
	g.Individual = append(g.Individual, father)

	// Add the father to the family
	family.Husband = father

	// Define the mother individual record
	mother := &gedcom.IndividualRecord{
		Xref: "I2",
		Name: []*gedcom.NameRecord{
			{
				Name: "Jane /Smith/",
			},
		},
		Sex: "F",
		Event: []*gedcom.EventRecord{
			{
				Tag:  "BIRT",
				Date: "5 MAY 1952",
				Place: gedcom.PlaceRecord{
					Name: "Manchester, England",
				},
			},
		},
	}
	// Add the mother to the GEDCOM
	g.Individual = append(g.Individual, mother)

	// Add the mother to the family
	family.Wife = mother

	// Create individual record for a child
	child := &gedcom.IndividualRecord{
		Xref: "I3",
		Name: []*gedcom.NameRecord{
			{
				Name: "Michael /Doe/",
			},
		},
		Sex: "M",
		Event: []*gedcom.EventRecord{
			{
				Tag:  "BIRT",
				Date: "15 JUL 1980",
				Place: gedcom.PlaceRecord{
					Name: "London, England",
				},
			},
		},
		// Link child to the family
		Parents: []*gedcom.FamilyLinkRecord{
			{
				Family: family,
				Type:   "CHIL",
			},
		},
	}
	// Add the child to the GEDCOM
	g.Individual = append(g.Individual, child)

	// Add the child to the family
	family.Child = append(family.Child, child)

	// Add an event to the family
	family.Event = append(family.Event, &gedcom.EventRecord{
		Tag:  "MARR",
		Date: "10 JUN 1975",
		Place: gedcom.PlaceRecord{
			Name: "London, England",
		},
	},
	)

	g.Trailer = &gedcom.Trailer{}

	enc := gedcom.NewEncoder(os.Stdout)
	if err := enc.Encode(g); err != nil {
		log.Fatalf("failed to encode gedcom: %v", err)
		return
	}
}
