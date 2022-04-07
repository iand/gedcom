/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
package gedcom

import (
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
	s := &scanner{}
	for _, ex := range examples {
		s.reset()
		offset, err := s.nextTag(ex.input)
		if err != nil {
			t.Fatalf(`nextTag for "%s" returned error "%v", expected no error`, ex.input, err)
		}

		if offset == 0 {
			t.Fatalf(`nextTag for "%s" did not find tag, expected it to find`, ex.input)
		}

		if s.level != ex.level {
			t.Errorf(`nextTag for "%s" returned level %d, expected %d`, ex.input, s.level, ex.level)
		}

		if string(s.tag) != ex.tag {
			t.Errorf(`nextTag for "%s" returned tag "%s", expected "%s"`, ex.input, s.tag, ex.tag)
		}

		if string(s.value) != ex.value {
			t.Errorf(`nextTag for "%s" returned value "%s", expected "%s"`, ex.input, s.value, ex.value)
		}

		if string(s.xref) != ex.xref {
			t.Errorf(`nextTag for "%s" returned xref "%s", expected "%s"`, ex.input, s.xref, ex.xref)
		}

	}
}

var examplesNot = [][]byte{
	[]byte("1 SEX F"),
	[]byte(" 1 SEX F "),
}

func TestNextTagNotFound(t *testing.T) {
	s := &scanner{}
	for _, ex := range examplesNot {
		s.reset()
		_, err := s.nextTag(ex)

		if err != io.EOF {
			t.Fatalf(`nextTag for "%s" returned unexpected error "%v", expected io.EOF`, ex, err)
		}

	}
}
