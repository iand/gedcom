package gedcom

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncodeHeader(t *testing.T) {
	testCases := []struct {
		name   string
		header *Header
		want   []string
	}{
		{
			name: "allged",
			header: &Header{
				SourceSystem: SystemRecord{
					Xref:         "APPROVED_SOURCE_NAME",
					Version:      "Version number of source-program",
					ProductName:  "Name of source-program",
					BusinessName: "Corporation name",
					Address: AddressRecord{
						Address: []*AddressDetail{
							{
								Full:       "Corporation address line 1\nCorporation address line 2\nCorporation address line 3\nCorporation address line 4",
								Line1:      "Corporation address line 1",
								Line2:      "Corporation address line 2",
								City:       "Corporation address city",
								State:      "Corporation address state",
								PostalCode: "Corporation address ZIP code",
								Country:    "Corporation address country",
							},
						},
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
			},

			want: []string{
				"0 HEAD",
				"1 CHAR ASCII",
				"2 VERS Version number of ASCII (whatever it means) ",
				"1 SOUR APPROVED_SOURCE_NAME",
				"2 VERS Version number of source-program",
				"2 NAME Name of source-program",
				"2 CORP Corporation name",
				"3 ADDR Corporation address line 1",
				"4 CONT Corporation address line 2",
				"4 CONT Corporation address line 3",
				"4 CONT Corporation address line 4",
				"4 ADR1 Corporation address line 1",
				"4 ADR2 Corporation address line 2",
				"4 CITY Corporation address city",
				"4 STAE Corporation address state",
				"4 POST Corporation address ZIP code",
				"4 CTRY Corporation address country",
				"3 PHON Corporation phone number 1",
				"3 PHON Corporation phone number 2",
				"3 PHON Corporation phone number 3 (last one!)",
				"2 DATA Name of source data",
				"3 DATE 1 JAN 1998",
				"3 COPR Copyright of source data",
				"1 DEST Destination of transmission",
				"1 DATE 1 JAN 1998",
				"2 TIME 13:57:24.80",
				"1 SUBM @SUBMITTER@",
				"1 SUBN @SUBMISSION@",
				"1 FILE ALLGED.GED",
				"1 COPR (C) 1997-2000 by H. Eichmann. You can use and distribute this file freely as long as you do not charge for it",
				"1 GEDC",
				"2 VERS 5.5",
				"2 FORM LINEAGE-LINKED",
				"1 LANG language",
				"1 NOTE A general note about this file:",
				"2 CONT It demonstrates most of the data which can be submitted using GEDCOM5.5. It shows the relatives of PERSON1:",
				"2 CONT His 2 wifes (PERSON2, PERSON8), his parents (father: PERSON5, mother not given), ",
				"2 CONT adoptive parents (mother: PERSON6, father not given) and his 3 children (PERSON3, PERSON4 and PERSON7).",
				"2 CONT In PERSON1, FAMILY1, SUBMITTER, SUBMISSION and SOURCE1 as many datafields as possible are used.",
				"2 CONT All other individuals/families contain no data. Note, that many data tags can appear more than once",
				"2 CONT (in this transmission this is demonstrated with tags: NAME, OCCU, PLACE and NOTE. Seek the word 'another'.",
				"2 CONT The data transmitted here do not make sence. Just the HEAD.DATE tag contains the date of the creation",
				"2 CONT of this file and will change in future Versions!",
				"2 CONT This file is created by H. Eichmann: h.eichmann@@gmx.de. Feel free to copy and use it for any ",
				"2 CONT non-commercial purpose. For the creation the GEDCOM standard Release 5.5 (2 JAN 1996) has been used.",
				"2 CONT Copyright: The church of Jesus Christ of latter-day saints, gedcom@@gedcom.org",
				"2 CONT Download it (the GEDCOM 5.5 specs) from: ftp.gedcom.com/pub/genealogy/gedcom.",
				"2 CONT Some Specials: This line is very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very very long but not too long (255 caharcters is the limit). ",
				"2 CONT This @@ (commercial at) character may only appear ONCE!",
				"2 CONT Note continued here. The word TEST should not be broken!",
				"1 _MYOWNTAG This is a non-standard tag. Not recommended but allowed",
				// "0 @SUBMITTER@ SUBM",
				// "1 NAME /Submitter-Name/",
				// "1 ADDR Submitter address line 1",
				// "2 CONT Submitter address line 2",
				// "2 CONT Submitter address line 3",
				// "2 CONT Submitter address line 4",
				// "2 ADR1 Submitter address line 1",
				// "2 ADR2 Submitter address line 2",
				// "2 CITY Submitter address city",
				// "2 STAE Submitter address state",
				// "2 POST Submitter address ZIP code",
				// "2 CTRY Submitter address country",
				// "1 PHON Submitter phone number 1",
				// "1 PHON Submitter phone number 2",
				// "1 PHON Submitter phone number 3 (last one!)",
				// "1 LANG English",
				// "1 CHAN ",
				// "2 DATE 19 JUN 2000",
				// "3 TIME 12:34:56.789",
				// "2 NOTE A note",
				// "3 CONT Note continued here. The word TE",
				// "3 CONC ST should not be broken!",
				// "1 _MYOWNTAG This is a non-standard tag. Not recommended but allowed",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := &Gedcom{
				Header: tc.header,
			}

			buf := new(bytes.Buffer)
			enc := NewEncoder(buf)

			err := enc.Encode(g)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if diff := cmp.Diff(tc.want, lines); diff != "" {
				t.Errorf("header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEncodeText(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "multiline",
			text: "line 1\nline 2\nline 3",
			want: []string{
				"1 NOTE line 1",
				"2 CONT line 2",
				"2 CONT line 3",
			},
		},
		{
			name: "max length line",
			text: strings.Repeat("0123456789", 24) + "012345", // 246 characters
			want: []string{
				"1 NOTE " + strings.Repeat("0123456789", 24) + "012345",
			},
		},
		{
			name: "long line",
			text: strings.Repeat("0123456789", 24) + "0123456789",
			want: []string{
				"1 NOTE " + strings.Repeat("0123456789", 24) + "012345",
				"2 CONC 6789",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			enc := NewEncoder(buf)

			err := enc.text(1, "NOTE", tc.text)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			err = enc.flush()
			if err != nil {
				t.Fatalf("unexpected error during flush: %v", err)
			}

			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if diff := cmp.Diff(tc.want, lines); diff != "" {
				t.Errorf("header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
