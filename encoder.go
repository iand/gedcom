/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// An Encoder encodes and writes GEDCOM objects to an input stream.
type Encoder struct {
	w *bufio.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	bw := bufio.NewWriter(w)
	return &Encoder{
		w: bw,
	}
}

func (e *Encoder) Encode(g *Gedcom) error {
	if err := e.header(g.Header); err != nil {
		return fmt.Errorf("header: %w", err)
	}

	return e.flush()
}

func (e *Encoder) flush() error {
	return e.w.Flush()
}

func (e *Encoder) tagFull(level int, tag string, xref string, value string) error {
	if _, err := e.w.WriteString(fmt.Sprintf("%d %s", level, tag)); err != nil {
		return err
	}

	if xref != "" {
		if _, err := e.w.WriteString("@" + xref + "@"); err != nil {
			return err
		}
	}

	if value != "" {
		if _, err := e.w.WriteString(" " + value); err != nil {
			return err
		}
	}

	if _, err := e.w.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) tag(level int, tag string, value string) error {
	return e.tagFull(level, tag, "", value)
}

func (e *Encoder) tagIfValue(level int, tag string, value string) error {
	if value == "" {
		return nil
	}
	return e.tagFull(level, tag, "", value)
}

func (e *Encoder) xref(level int, tag string, xref string) error {
	if _, err := e.w.WriteString(fmt.Sprintf("%d %s @%s@", level, tag, xref)); err != nil {
		return err
	}

	if _, err := e.w.WriteString("\n"); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) text(level int, tag string, value string) error {
	conts := strings.Split(value, "\n")
	if err := e.tag(level, tag, conts[0]); err != nil {
		return err
	}

	for i := 1; i < len(conts); i++ {
		cont := conts[i]
		if len(cont) <= 246 {
			if err := e.tag(level+1, "CONT", cont); err != nil {
				return err
			}
			continue
		}

		if err := e.tag(level+1, "CONT", cont[:246]); err != nil {
			return err
		}
		for len(cont) > 246 {
			cont = cont[246:]
			if err := e.tag(level+1, "CONC", cont[:246]); err != nil {
				return err
			}
		}

	}

	return nil
}

func (e *Encoder) textIfValue(level int, tag string, value string) error {
	if value == "" {
		return nil
	}
	return e.text(level, tag, value)
}

func (e *Encoder) header(h *Header) error {
	if h == nil {
		return nil
	}
	e.tag(0, "HEAD", "")
	e.tagIfValue(1, "CHAR", h.CharacterSet)
	e.tagIfValue(2, "VERS", h.CharacterSetVersion)
	e.sourceSystem(0, h.SourceSystem)
	e.tagIfValue(1, "DEST", h.Destination)
	e.tagIfValue(1, "DATE", h.Date)
	e.tagIfValue(2, "TIME", h.Time)

	if h.Submitter != nil {
		e.xref(1, "SUBM", h.Submitter.Xref)
	}

	if h.Submission != nil {
		e.xref(1, "SUBN", h.Submission.Xref)
	}
	e.tagIfValue(1, "FILE", h.Filename)
	e.tagIfValue(1, "COPR", h.Copyright)

	if h.Version != "" || h.Form != "" {
		e.tag(1, "GEDC", "")
		if h.Version != "" {
			e.tag(2, "VERS", h.Version)
		}
		if h.Form != "" {
			e.tag(2, "FORM", h.Form)
		}
	}
	e.tagIfValue(1, "LANG", h.Language)
	e.textIfValue(1, "NOTE", h.Note)
	e.userDefinedList(1, h.UserDefined)

	return nil
}

func (e *Encoder) sourceSystem(level int, s SystemRecord) error {
	e.tag(level+1, "SOUR", s.Xref)
	e.tagIfValue(level+2, "VERS", s.Version)
	e.tagIfValue(level+2, "NAME", s.ProductName)
	e.tagIfValue(level+2, "CORP", s.BusinessName)

	e.address(level+3, s.Address)

	e.tagIfValue(level+2, "DATA", s.SourceName)
	e.tagIfValue(level+3, "DATE", s.SourceDate)
	e.tagIfValue(level+3, "COPR", s.SourceCopyright)
	e.userDefinedList(1, s.UserDefined)

	return nil
}

func (e *Encoder) userDefinedList(level int, uds []UserDefinedTag) error {
	for _, ud := range uds {
		if err := e.userDefined(level, ud); err != nil {
			return fmt.Errorf("user defined tag list: %w", err)
		}
	}
	return nil
}

func (e *Encoder) userDefined(level int, ud UserDefinedTag) error {
	if err := e.tagFull(level, ud.Tag, ud.Xref, ud.Value); err != nil {
		return fmt.Errorf("address tag: %w", err)
	}

	return e.userDefinedList(level+1, ud.UserDefined)
}

func (e *Encoder) address(level int, ar AddressRecord) error {
	if err := e.addressDetailList(level, ar.Address); err != nil {
		return fmt.Errorf("user defined tag: %w", err)
	}

	for _, v := range ar.Phone {
		if err := e.textIfValue(level, "PHON", v); err != nil {
			return fmt.Errorf("phone number: %w", err)
		}
	}

	for _, v := range ar.Email {
		if err := e.textIfValue(level, "EMAIL", v); err != nil {
			return fmt.Errorf("email: %w", err)
		}
	}

	for _, v := range ar.Fax {
		if err := e.textIfValue(level, "FAX", v); err != nil {
			return fmt.Errorf("fax: %w", err)
		}
	}

	for _, v := range ar.WWW {
		if err := e.textIfValue(level, "WWW", v); err != nil {
			return fmt.Errorf("www: %w", err)
		}
	}

	return nil
}

func (e *Encoder) addressDetailList(level int, ads []*AddressDetail) error {
	for _, ad := range ads {
		if err := e.addressDetail(level, ad); err != nil {
			return fmt.Errorf("address detail list: %w", err)
		}
	}

	return nil
}

func (e *Encoder) addressDetail(level int, ar *AddressDetail) error {
	if ar == nil {
		return nil
	}
	if err := e.text(level, "ADDR", ar.Full); err != nil {
		return fmt.Errorf("address detail: %w", err)
	}

	if err := e.textIfValue(level+1, "ADR1", ar.Line1); err != nil {
		return fmt.Errorf("address line 1: %w", err)
	}

	if err := e.textIfValue(level+1, "ADR2", ar.Line2); err != nil {
		return fmt.Errorf("address line 2: %w", err)
	}

	if err := e.textIfValue(level+1, "ADR3", ar.Line3); err != nil {
		return fmt.Errorf("address line 3: %w", err)
	}

	if err := e.textIfValue(level+1, "CITY", ar.City); err != nil {
		return fmt.Errorf("city: %w", err)
	}

	if err := e.textIfValue(level+1, "STAE", ar.State); err != nil {
		return fmt.Errorf("address state: %w", err)
	}

	if err := e.textIfValue(level+1, "POST", ar.PostalCode); err != nil {
		return fmt.Errorf("postal code: %w", err)
	}

	if err := e.textIfValue(level+1, "CTRY", ar.Country); err != nil {
		return fmt.Errorf("county: %w", err)
	}

	return nil
}
