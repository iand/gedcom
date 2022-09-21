/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type example struct {
	input []byte
	level int
	tag   string
	value string
	xref  string
}

var examples = []example{
	{[]byte("1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte(" 1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte("  \r\n\t 1 SEX F\n"), 1, `SEX`, `F`, ""},
	{[]byte("  \r\n\t 1     SEX      F\n"), 1, `SEX`, `F`, ""},
	{[]byte("1 SEX F\r"), 1, `SEX`, `F`, ""},
	{[]byte("1 SEX F \r"), 1, `SEX`, `F `, ""},
	{[]byte("0 HEAD\r"), 0, `HEAD`, ``, ""},
	{[]byte("0 @OTHER@ SUBM\n"), 0, `SUBM`, ``, "OTHER"},
	{[]byte("1 PUBL Corp, Inc.\n"), 1, `PUBL`, `Corp, Inc.`, ""},
	{[]byte("1 NOTE <i>markup</i>. plain\n"), 1, `NOTE`, `<i>markup</i>. plain`, ""},
}

func TestNextTagFound(t *testing.T) {
	for _, ex := range examples {
		s := NewScanner(bytes.NewReader(ex.input))
		if !s.Next() {
			if s.Err() != nil {
				t.Fatalf(`nextTag for "%s" returned error "%v", expected no error`, ex.input, s.Err())
			}
		}
		if s.level != ex.level {
			t.Errorf(`nextTag for "%s" returned level %d, expected %d`, ex.input, s.level, ex.level)
		}

		if s.tag != ex.tag {
			t.Errorf(`nextTag for "%s" returned tag "%s", expected "%s"`, ex.input, s.tag, ex.tag)
		}

		if s.value != ex.value {
			t.Errorf(`nextTag for "%s" returned value "%s", expected "%s"`, ex.input, s.value, ex.value)
		}

		if s.xref != ex.xref {
			t.Errorf(`nextTag for "%s" returned xref "%s", expected "%s"`, ex.input, s.xref, ex.xref)
		}

		if s.Next() {
			t.Errorf(`got another tag for %q, wanted no tag`, ex)
		}

	}
}

var examplesNot = [][]byte{
	// These are not terminated by a newline
	[]byte("1 SEX F"),
	[]byte(" 1 SEX F "),
}

func TestNextTagNotFound(t *testing.T) {
	for _, ex := range examplesNot {
		s := NewScanner(bytes.NewReader(ex))
		if s.Next() {
			t.Fatalf(`got tag for %q, wanted no tag`, ex)
		}
		if !errors.Is(s.Err(), io.ErrUnexpectedEOF) {
			t.Errorf("got error %v, wanted %v", s.err, io.ErrUnexpectedEOF)
		}

	}
}

func TestLineEndings(t *testing.T) {
	testCases := []struct {
		platform string
		input    []byte
	}{
		{
			platform: "unix",
			input:    []byte("0 HEAD\n1 CHAR UTF-8\n1 GEDC\n1 NOTE first line\n2 CONT second line\n"),
		},
		{
			platform: "windows",
			input:    []byte("0 HEAD\r\n1 CHAR UTF-8\r\n1 GEDC\r\n1 NOTE first line\r\n2 CONT second line\r\n"),
		},
		{
			platform: "macclassic",
			input:    []byte("0 HEAD\r1 CHAR UTF-8\r\n1 GEDC\r1 NOTE first line\r2 CONT second line\r"),
		},
	}

	want := []Line{
		{Level: 0, Tag: "HEAD"},
		{Level: 1, Tag: "CHAR", Value: "UTF-8"},
		{Level: 1, Tag: "GEDC"},
		{Level: 1, Tag: "NOTE", Value: "first line"},
		{Level: 2, Tag: "CONT", Value: "second line"},
	}

	for _, tc := range testCases {
		t.Run(tc.platform, func(t *testing.T) {
			s := NewScanner(bytes.NewReader(tc.input))

			for i := range want {
				if !s.Next() {
					t.Fatalf("missing line %d, err=%v", i+1, s.Err())
				}

				l := s.Line()
				if l.Level != want[i].Level {
					t.Errorf("line %d got level %d, wanted %d", i+1, l.Level, want[i].Level)
				}

				if l.Tag != want[i].Tag {
					t.Errorf("line %d got tag %s, wanted %s", i+1, l.Tag, want[i].Tag)
				}

				if l.Value != want[i].Value {
					t.Errorf("line %d got value %q wanted %q", i+1, l.Value, want[i].Value)
				}

			}

			if s.Err() != nil {
				t.Errorf("got an unexpected error: %v", s.Err())
			}

			if s.Next() {
				t.Errorf("got an unexpected tag")
			}
		})
	}
}
