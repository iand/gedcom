/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/
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

// NewDecoder returns a new decoder that reads from r.
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
		num, err := d.r.Read(buf[rest:len(buf)])
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

func makeRootParser(d *Decoder, g *Gedcom) parser {
	return func(level int, tag string, value string, xref string) error {
		//println(level, tag, value, xref)
		if level == 0 {
			switch tag {
			case "INDI":
				var obj *IndividualRecord
				if xref != "" {
					ref, ok := d.refs[xref].(*IndividualRecord)
					if !ok {
						obj = &IndividualRecord{Xref: stripXref(xref)}
						d.refs[obj.Xref] = obj
					} else {
						obj = ref
					}
				}

				g.Individual = append(g.Individual, obj)
				d.pushParser(makeIndividualParser(d, obj, level))

			case "SUBM":
				g.Submitter = append(g.Submitter, &SubmitterRecord{})
			case "FAM":
				var obj *FamilyRecord
				if xref != "" {
					ref, ok := d.refs[xref].(*FamilyRecord)
					if !ok {
						obj = &FamilyRecord{Xref: stripXref(xref)}
						d.refs[obj.Xref] = obj
					} else {
						obj = ref
					}
				}

				g.Family = append(g.Family, obj)
				d.pushParser(makeFamilyParser(d, obj, level))
			case "SOUR":
				s := &SourceRecord{
					Xref: stripXref(xref),
				}
				g.Source = append(g.Source, s)
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
			xref := stripXref(value)

			family, ok := d.refs[xref].(*FamilyRecord)
			if !ok {
				family = &FamilyRecord{Xref: xref}
				d.refs[family.Xref] = family
			}

			f := &FamilyLinkRecord{Family: family}
			i.Parents = append(i.Parents, f)
			d.pushParser(makeFamilyLinkParser(d, f, level))
		case "FAMS":
			xref := stripXref(value)
			family, ok := d.refs[xref].(*FamilyRecord)
			if !ok {
				family = &FamilyRecord{Xref: xref}
				d.refs[family.Xref] = family
			}

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
			s := &SourceRecord{}
			n.Source = append(n.Source, s)
			d.pushParser(makeSourceParser(d, s, level))
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
		case "PAGE":
			s.Page = value
		case "QUAY":
			s.Quay = value
		case "NOTE":
			r := &NoteRecord{Note: value}
			s.Note = append(s.Note, r)
			d.pushParser(makeNoteParser(d, r, level))
		case "DATA":
			d.pushParser(makeDataParser(d, &s.Data, level))

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
			s := &SourceRecord{}
			n.Source = append(n.Source, s)
			d.pushParser(makeSourceParser(d, s, level))
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
			s := &SourceRecord{}
			e.Source = append(e.Source, s)
			d.pushParser(makeSourceParser(d, s, level))
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
			s := &SourceRecord{}
			p.Source = append(p.Source, s)
			d.pushParser(makeSourceParser(d, s, level))
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
			xref = stripXref(value)
			i, ok := d.refs[xref].(*IndividualRecord)
			if !ok {
				i = &IndividualRecord{Xref: xref}
				d.refs[i.Xref] = i
			}
			f.Husband = i
		case "WIFE":
			xref = stripXref(value)
			i, ok := d.refs[xref].(*IndividualRecord)
			if !ok {
				i = &IndividualRecord{Xref: xref}
				d.refs[i.Xref] = i
			}
			f.Wife = i
		case "CHIL":
			xref = stripXref(value)
			i, ok := d.refs[xref].(*IndividualRecord)
			if !ok {
				i = &IndividualRecord{Xref: stripXref(xref)}
				d.refs[i.Xref] = i
			}
			f.Child = append(f.Child, i)

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
