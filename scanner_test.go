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
		s := newScanner(bytes.NewReader(ex.input))
		if !s.next() {
			if s.err != nil {
				t.Fatalf(`nextTag for "%s" returned error "%v", expected no error`, ex.input, s.err)
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

		if s.next() {
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
		s := newScanner(bytes.NewReader(ex))
		if s.next() {
			t.Fatalf(`got tag for %q, wanted no tag`, ex)
		}
		if !errors.Is(s.err, io.ErrUnexpectedEOF) {
			t.Errorf("got error %v, wanted %v", s.err, io.ErrUnexpectedEOF)
		}

	}
}

func TestWindowsLineEndings(t *testing.T) {
	input := []byte("0 HEAD\r\n1 CHAR UTF-8\r\n1 GEDC\r\n")

	s := newScanner(bytes.NewReader(input))
	if !s.next() {
		t.Fatalf("missing first tag, err=%v", s.err)
	}

	if s.level != 0 {
		t.Errorf("first tag, got level %d, wanted %d", s.level, 0)
	}

	if !s.next() {
		t.Fatalf("missing second tag, err=%v", s.err)
	}

	if s.level != 1 {
		t.Errorf("second tag, got level %d, wanted %d", s.level, 1)
	}

	if !s.next() {
		t.Fatalf("missing third tag, err=%v", s.err)
	}

	if s.level != 1 {
		t.Errorf("third tag, got level %d, wanted %d", s.level, 1)
	}

	if s.next() {
		t.Errorf("got an unexpected tag")
	}

	if s.err != nil {
		t.Errorf("got an unexpected error: %v", s.err)
	}
}
