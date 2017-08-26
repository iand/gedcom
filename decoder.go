/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

// Package gedcom provides a functions to parse GEDCOM files.
package gedcom

import (
	"io"
	"strings"
)

// A Decoder reads and decodes GEDCOM objects from an input stream.
type Decoder struct {
	r       io.Reader
	parsers []parser
	refs    map[string]interface{}
}

// NewDecoder returns a new decoder that reads r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next GEDCOM-encoded value from its
// input and stores it in the value pointed to by v.
func (d *Decoder) Decode() (*Gedcom, error) {

	g := &Gedcom{
		Family:     make([]*FamilyRecord, 0),
		Individual: make([]*IndividualRecord, 0),
		Media:      make([]*MediaRecord, 0),
		Repository: make([]*RepositoryRecord, 0),
		Source:     make([]*SourceRecord, 0),
		Submitter:  make([]*SubmitterRecord, 0),
	}

	d.refs = make(map[string]interface{})
	d.parsers = []parser{makeRootParser(d, g)}
	d.scan(g)

	return g, nil
}

func (d *Decoder) scan(g *Gedcom) {
	s := &scanner{}
	buf := make([]byte, 512)

	n, err := d.r.Read(buf)
	if err != nil {
		// TODO
	}

	for n > 0 {
		pos := 0

		for {
			s.reset()
			offset, err := s.nextTag(buf[pos:n])
			pos += offset
			if err != nil {
				if err != io.EOF {
					println(err.Error())
					return
				}
				break
			}

			d.parsers[len(d.parsers)-1](s.level, string(s.tag), string(s.value), string(s.xref))

		}

		// shift unparsed bytes to start of buffer
		rest := copy(buf, buf[pos:])

		// top up buffer
		num, err := d.r.Read(buf[rest:])
		if err != nil {
			break
		}

		n = rest + num - 1

	}

}

type parser func(level int, tag string, value string, xref string) error

func (d *Decoder) pushParser(p parser) {
	d.parsers = append(d.parsers, p)
}

func (d *Decoder) popParser(level int, tag string, value string, xref string) error {
	n := len(d.parsers) - 1
	if n < 1 {
		panic("MASSIVE ERROR") // TODO
	}
	d.parsers = d.parsers[0:n]

	return d.parsers[len(d.parsers)-1](level, tag, value, xref)
}

func (d *Decoder) individual(xref string) *IndividualRecord {
	if xref == "" {
		return &IndividualRecord{}
	}

	ref, found := d.refs[xref].(*IndividualRecord)
	if !found {
		rec := &IndividualRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref

}

func (d *Decoder) family(xref string) *FamilyRecord {
	if xref == "" {
		return &FamilyRecord{}
	}

	ref, found := d.refs[xref].(*FamilyRecord)
	if !found {
		rec := &FamilyRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func (d *Decoder) source(xref string) *SourceRecord {
	if xref == "" {
		return &SourceRecord{}
	}

	ref, found := d.refs[xref].(*SourceRecord)
	if !found {
		rec := &SourceRecord{Xref: xref}
		d.refs[rec.Xref] = rec
		return rec
	}
	return ref
}

func makeRootParser(d *Decoder, g *Gedcom) parser {
	return func(level int, tag string, value string, xref string) error {
		//println(level, tag, value, xref)
		if level == 0 {
			switch tag {
			case "INDI":
				obj := d.individual(xref)
				g.Individual = append(g.Individual, obj)
				d.pushParser(makeIndividualParser(d, obj, level))
			case "SUBM":
				g.Submitter = append(g.Submitter, &SubmitterRecord{})
			case "FAM":
				obj := d.family(xref)
				g.Family = append(g.Family, obj)
				d.pushParser(makeFamilyParser(d, obj, level))
			case "SOUR":
				obj := d.source(xref)
				g.Source = append(g.Source, obj)
				//d.pushParser(makeSourceParser(d, s, level))
			}
		}
		return nil
	}
}

func makeIndividualParser(d *Decoder, i *IndividualRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "NAME":
			n := &NameRecord{Name: value}
			i.Name = append(i.Name, n)
			d.pushParser(makeNameParser(d, n, level))
		case "SEX":
			i.Sex = value
		case "BIRT", "CHR", "DEAT", "BURI", "CREM", "ADOP", "BAPM", "BARM", "BASM", "BLES", "CHRA", "CONF", "FCOM", "ORDN", "NATU", "EMIG", "IMMI", "CENS", "PROB", "WILL", "GRAD", "RETI", "EVEN":
			e := &EventRecord{Tag: tag, Value: value}
			i.Event = append(i.Event, e)
			d.pushParser(makeEventParser(d, e, level))
		case "CAST", "DSCR", "EDUC", "IDNO", "NATI", "NCHI", "NMR", "OCCU", "PROP", "RELI", "RESI", "SSN", "TITL", "FACT":
			e := &EventRecord{Tag: tag, Value: value}
			i.Attribute = append(i.Attribute, e)
			d.pushParser(makeEventParser(d, e, level))
		case "FAMC":
			family := d.family(stripXref(value))
			f := &FamilyLinkRecord{Family: family}
			i.Parents = append(i.Parents, f)
			d.pushParser(makeFamilyLinkParser(d, f, level))
		case "FAMS":
			family := d.family(stripXref(value))
			f := &FamilyLinkRecord{Family: family}
			i.Family = append(i.Family, f)
			d.pushParser(makeFamilyLinkParser(d, f, level))
		}
		return nil
	}
}

func makeNameParser(d *Decoder, n *NameRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			n.Citation = append(n.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			n.Note = append(n.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		}

		return nil
	}
}

func makeSourceParser(d *Decoder, s *SourceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "TITL":
			s.Title = value
			d.pushParser(makeTextParser(d, &s.Title, level))

		case "NOTE":
			r := &NoteRecord{Note: value}
			s.Note = append(s.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		}

		return nil
	}
}

func makeCitationParser(d *Decoder, c *CitationRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "PAGE":
			c.Page = value
		case "QUAY":
			c.Quay = value
		case "NOTE":
			r := &NoteRecord{Note: value}
			c.Note = append(c.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		case "DATA":
			d.pushParser(makeDataParser(d, &c.Data, level))

		}

		return nil
	}
}

func makeNoteParser(d *Decoder, n *NoteRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			n.Note = n.Note + "\n" + value
		case "CONC":
			n.Note = n.Note + value
		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			n.Citation = append(n.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		}

		return nil
	}
}

func makeTextParser(d *Decoder, s *string, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			*s = *s + "\n" + value
		case "CONC":
			*s = *s + value
		}

		return nil
	}
}

func makeDataParser(d *Decoder, r *DataRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "DATE":
			r.Date = value
		case "TEXT":
			r.Text = append(r.Text, value)
			d.pushParser(makeTextParser(d, &r.Text[len(r.Text)-1], level))
		}

		return nil
	}
}

func makeEventParser(d *Decoder, e *EventRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "TYPE":
			e.Type = value
		case "DATE":
			e.Date = value
		case "PLAC":
			e.Place.Name = value
			d.pushParser(makePlaceParser(d, &e.Place, level))
		case "ADDR":
			e.Address.Full = value
			d.pushParser(makeAddressParser(d, &e.Address, level))
		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			e.Citation = append(e.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			e.Note = append(e.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		}

		return nil
	}
}

func makePlaceParser(d *Decoder, p *PlaceRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {

		case "SOUR":
			c := &CitationRecord{Source: d.source(stripXref(value))}
			p.Citation = append(p.Citation, c)
			d.pushParser(makeCitationParser(d, c, level))
		case "NOTE":
			r := &NoteRecord{Note: value}
			p.Note = append(p.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		}

		return nil
	}
}

func makeFamilyLinkParser(d *Decoder, f *FamilyLinkRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "PEDI":
			f.Type = value
		case "NOTE":
			r := &NoteRecord{Note: value}
			f.Note = append(f.Note, r)
			d.pushParser(makeNoteParser(d, r, level))

		}

		return nil
	}
}

func makeFamilyParser(d *Decoder, f *FamilyRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "HUSB":
			f.Husband = d.individual(stripXref(value))
		case "WIFE":
			f.Wife = d.individual(stripXref(value))
		case "CHIL":
			f.Child = append(f.Child, d.individual(stripXref(value)))
		case "ANUL", "CENS", "DIV", "DIVF", "ENGA", "MARR", "MARB", "MARC", "MARL", "MARS", "EVEN":
			e := &EventRecord{Tag: tag, Value: value}
			f.Event = append(f.Event, e)
			d.pushParser(makeEventParser(d, e, level))

		}
		return nil
	}
}

func makeAddressParser(d *Decoder, a *AddressRecord, minLevel int) parser {
	return func(level int, tag string, value string, xref string) error {
		if level <= minLevel {
			return d.popParser(level, tag, value, xref)
		}
		switch tag {
		case "CONT":
			a.Full = a.Full + "\n" + value
		case "ADR1":
			a.Line1 = value
		case "ADR2":
			a.Line2 = value
		case "CITY":
			a.City = value
		case "STAE":
			a.State = value
		case "POST":
			a.PostalCode = value
		case "CTRY":
			a.Country = value
		case "PHON":
			a.Phone = value

		}

		return nil
	}
}

func stripXref(value string) string {
	return strings.Trim(value, "@")
}
