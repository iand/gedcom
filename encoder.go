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

// GedcomVersion specifies the GEDCOM version for encoding output.
type GedcomVersion int

const (
	Gedcom55 GedcomVersion = iota // GEDCOM 5.5/5.5.1 (default)
	Gedcom70                      // GEDCOM 7.0
)

// Encoder writes GEDCOM-encoded data to an output stream.
// Use [NewEncoder] to create an Encoder and [Encoder.Encode] to write
// a [Gedcom] structure.
//
// The encoder handles GEDCOM line length limits automatically, using
// CONT (continuation) and CONC (concatenation) tags to split long text.
type Encoder struct {
	w       *bufio.Writer
	err     error
	version GedcomVersion
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	bw := bufio.NewWriter(w)
	return &Encoder{
		w: bw,
	}
}

// SetVersion sets the GEDCOM version for the encoder output.
// By default the encoder produces GEDCOM 5.5 output. Use [Gedcom70]
// to produce GEDCOM 7.0 output.
func (e *Encoder) SetVersion(v GedcomVersion) {
	e.version = v
}

// Encode writes the GEDCOM-encoded representation of g to the encoder's output stream.
// It writes the header, all records (individuals, families, sources, etc.), and trailer.
func (e *Encoder) Encode(g *Gedcom) error {
	if e.version == Gedcom70 {
		if _, err := e.w.WriteString("\xEF\xBB\xBF"); err != nil {
			return fmt.Errorf("write BOM: %w", err)
		}
	}
	e.header(g.Header)

	for _, r := range g.Individual {
		e.individual(r)
	}

	for _, r := range g.Family {
		e.family(r)
	}

	e.mediaList(0, g.Media)

	for _, r := range g.Repository {
		e.repository(r)
	}

	for _, r := range g.Source {
		e.source(r)
	}

	for _, r := range g.Submitter {
		e.submitter(0, r)
	}

	for _, r := range g.SharedNote {
		e.sharedNoteRecord(r)
	}

	e.userDefinedList(0, g.UserDefined)
	e.trailer(g.Trailer)

	return e.flush()
}

func (e *Encoder) flush() error {
	if e.err != nil {
		return e.err
	}
	return e.w.Flush()
}

func (e *Encoder) tagWithID(level int, tag string, id string) {
	if e.err != nil {
		return
	}
	if id == "" {
		e.err = fmt.Errorf("tag %s missing id", tag)
		return
	}
	if _, err := e.w.WriteString(fmt.Sprintf("%d @%s@ %s", level, id, tag)); err != nil {
		e.err = fmt.Errorf("write tag with id %s @%s@: %w", tag, id, err)
		return
	}

	if _, err := e.w.WriteString("\n"); err != nil {
		e.err = fmt.Errorf("write tag %s: %w", tag, err)
		return
	}
}

func (e *Encoder) tag(level int, tag string, value string) {
	if e.err != nil {
		return
	}

	if _, err := e.w.WriteString(fmt.Sprintf("%d %s", level, tag)); err != nil {
		e.err = fmt.Errorf("write tag %s: %w", tag, err)
		return
	}

	if value != "" {
		if _, err := e.w.WriteString(" " + value); err != nil {
			e.err = fmt.Errorf("write tag %s: %w", tag, err)
			return
		}
	}
	if _, err := e.w.WriteString("\n"); err != nil {
		e.err = fmt.Errorf("write tag %s: %w", tag, err)
		return
	}
}

// maybeTag writes a tag with a level if the value is not empty
func (e *Encoder) maybeTag(level int, tag string, value string) {
	if e.err != nil {
		return
	}
	if value == "" {
		return
	}
	e.tag(level, tag, value)
}

// tagWithPointer writes a tag with a pointer reference
func (e *Encoder) tagWithPointer(level int, tag string, xref string) {
	if e.err != nil {
		return
	}
	if _, err := e.w.WriteString(fmt.Sprintf("%d %s @%s@\n", level, tag, xref)); err != nil {
		e.err = fmt.Errorf("write tag with pointer %s @%s@: %w", tag, xref, err)
		return
	}
}

// tagWithOptionalPointer writes a tag with a pointer reference if it is non empty
func (e *Encoder) tagWithOptionalPointer(level int, tag string, xref string) {
	if e.err != nil {
		return
	}
	if xref != "" {
		e.tagWithPointer(level, tag, xref)
	} else {
		e.tag(level, tag, "")
	}
}

// tagWithText writes a tag with text, handling continuations
func (e *Encoder) tagWithText(level int, tag string, value string) {
	if e.err != nil {
		return
	}

	conts := strings.Split(value, "\n")
	e.textOneLine(level, tag, conts[0])

	for i := 1; i < len(conts); i++ {
		e.textOneLine(level+1, "CONT", conts[i])
	}
}

func (e *Encoder) textOneLine(level int, tag string, value string) {
	if e.err != nil {
		return
	}

	if e.version == Gedcom70 || len(value) <= 246 {
		e.tag(level, tag, value)
		return
	}
	e.tag(level, tag, value[:246])

	for len(value) > 246 {
		value = value[246:]
		if len(value) <= 246 {
			e.tag(level+1, "CONC", value)
			return
		}

		e.tag(level+1, "CONC", value[:246])
	}
}

// maybeTagWithText writes a tag with text only if the text is not empty
func (e *Encoder) maybeTagWithText(level int, tag string, value string) {
	if e.err != nil {
		return
	}
	if value == "" {
		return
	}
	e.tagWithText(level, tag, value)
}

func (e *Encoder) header(h *Header) {
	if e.err != nil {
		return
	}
	if h == nil {
		return
	}
	e.tag(0, "HEAD", "")
	if e.version != Gedcom70 {
		e.maybeTag(1, "CHAR", h.CharacterSet)
		e.maybeTag(2, "VERS", h.CharacterSetVersion)
	}
	e.sourceSystem(0, h.SourceSystem)
	e.maybeTag(1, "DEST", h.Destination)
	e.maybeTag(1, "DATE", h.Date)
	e.maybeTag(2, "TIME", h.Time)

	if h.Submitter != nil {
		e.tagWithPointer(1, "SUBM", h.Submitter.Xref)
	}

	if h.Submission != nil {
		e.tagWithPointer(1, "SUBN", h.Submission.Xref)
	}
	e.maybeTag(1, "FILE", h.Filename)
	e.maybeTag(1, "COPR", h.Copyright)

	if e.version == Gedcom70 {
		e.tag(1, "GEDC", "")
		e.tag(2, "VERS", "7.0")
	} else if h.Version != "" || h.Form != "" {
		e.tag(1, "GEDC", "")
		if h.Version != "" {
			e.tag(2, "VERS", h.Version)
		}
		if h.Form != "" {
			e.tag(2, "FORM", h.Form)
		}
	}
	e.maybeTag(1, "LANG", h.Language)
	e.schema(1, h.Schema)
	e.maybeTagWithText(1, "NOTE", h.Note)
	e.userDefinedList(1, h.UserDefined)
}

func (e *Encoder) sourceSystem(level int, s SystemRecord) {
	if e.err != nil {
		return
	}
	e.tag(level+1, "SOUR", s.Xref)
	e.maybeTag(level+2, "VERS", s.Version)
	e.maybeTag(level+2, "NAME", s.ProductName)
	e.maybeTag(level+2, "CORP", s.BusinessName)

	e.address(level+3, &s.Address)

	e.maybeTag(level+2, "DATA", s.SourceName)
	e.maybeTag(level+3, "DATE", s.SourceDate)
	e.maybeTag(level+3, "COPR", s.SourceCopyright)
	e.userDefinedList(1, s.UserDefined)
}

func (e *Encoder) userDefinedList(level int, uds []UserDefinedTag) {
	if e.err != nil {
		return
	}
	for _, ud := range uds {
		e.userDefined(level, ud)
	}
}

func (e *Encoder) userDefined(level int, r UserDefinedTag) {
	if e.err != nil {
		return
	}
	if r.Xref != "" {
		e.tagWithPointer(level, r.Tag, r.Xref)
	} else {
		e.tag(level, r.Tag, r.Value)
	}
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) address(level int, r *AddressRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.addressDetailList(level, r.Address)

	for _, v := range r.Phone {
		e.maybeTagWithText(level, "PHON", v)
	}

	for _, v := range r.Email {
		e.maybeTagWithText(level, "EMAIL", v)
	}

	for _, v := range r.Fax {
		e.maybeTagWithText(level, "FAX", v)
	}

	for _, v := range r.WWW {
		e.maybeTagWithText(level, "WWW", v)
	}
}

func (e *Encoder) addressDetailList(level int, rs []*AddressDetail) {
	if e.err != nil {
		return
	}
	for _, ad := range rs {
		e.addressDetail(level, ad)
	}
}

func (e *Encoder) addressDetail(level int, r *AddressDetail) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithText(level, "ADDR", r.Full)
	e.maybeTagWithText(level+1, "ADR1", r.Line1)
	e.maybeTagWithText(level+1, "ADR2", r.Line2)
	e.maybeTagWithText(level+1, "ADR3", r.Line3)
	e.maybeTagWithText(level+1, "CITY", r.City)
	e.maybeTagWithText(level+1, "STAE", r.State)
	e.maybeTagWithText(level+1, "POST", r.PostalCode)
	e.maybeTagWithText(level+1, "CTRY", r.Country)
}

func (e *Encoder) place(level int, r *PlaceRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Name == "" && len(r.Phonetic) == 0 && len(r.Romanized) == 0 && len(r.Translation) == 0 && r.Language == "" && r.Latitude == "" && r.Longitude == "" && len(r.Note) == 0 && len(r.Citation) == 0 {
		return
	}

	e.tag(level, "PLAC", r.Name)
	e.maybeTag(level+1, "LANG", r.Language)
	for _, sr := range r.Phonetic {
		e.tag(level+1, "FONE", sr.Name)
		e.maybeTag(level+1, "TYPE", sr.Type)
	}

	for _, sr := range r.Romanized {
		e.tag(level+1, "ROMN", sr.Name)
		e.maybeTag(level+2, "TYPE", sr.Type)
	}

	e.translationList(level+1, r.Translation)

	if r.Latitude != "" || r.Longitude != "" {
		e.tag(level+1, "MAP", "")
		e.maybeTag(level+2, "LATI", r.Latitude)
		e.maybeTag(level+2, "LONG", r.Longitude)

	}

	e.noteList(level+1, r.Note)
	e.citationList(level+1, r.Citation)
}

func (e *Encoder) individual(r *IndividualRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}

	level := 0
	e.tagWithID(level, "INDI", r.Xref)
	for _, v := range r.Name {
		e.name(level+1, v)
	}
	e.maybeTagWithText(level+1, "SEX", r.Sex)

	e.eventList(level+1, r.Event)
	e.eventList(level+1, r.Attribute)
	e.nonEventList(level+1, r.NonEvent)

	for _, sr := range r.Parents {
		e.familyLink(level+1, "FAMC", sr)
	}
	for _, sr := range r.Family {
		e.familyLink(level+1, "FAMS", sr)
	}

	if len(r.Submitter) > 0 {
		// Submitter                 []*SubmitterRecord
		e.err = fmt.Errorf("not implemented: Submitter")
		return
	}
	for _, a := range r.Association {
		e.association(level+1, a)
	}

	e.maybeTagWithText(level+1, "RFN", r.PermanentRecordFileNumber)
	e.maybeTagWithText(level+1, "AFN", r.AncestralFileNumber)
	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.maybeTag(level+1, "RESN", r.RestrictionNotice)

	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.creation(level+1, &r.Creation)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.citationList(level+1, r.Citation)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) family(r *FamilyRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}

	level := 0
	e.tagWithID(level, "FAM", r.Xref)
	e.individualRef(level+1, "HUSB", r.Husband)
	e.individualRef(level+1, "WIFE", r.Wife)
	for _, sr := range r.Child {
		e.individualRef(level+1, "CHIL", sr)
	}
	e.eventList(level+1, r.Event)
	e.nonEventList(level+1, r.NonEvent)
	e.maybeTag(level+1, "NCHI", r.NumberOfChildren)
	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.maybeTag(level+1, "RESN", r.RestrictionNotice)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.creation(level+1, &r.Creation)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.citationList(level+1, r.Citation)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) mediaList(level int, rs []*MediaRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.media(level, r)
	}
}

func (e *Encoder) media(level int, r *MediaRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if level == 0 {
		e.tagWithID(level, "OBJE", r.Xref)
	} else {
		e.tagWithOptionalPointer(level, "OBJE", r.Xref)
	}

	for _, sr := range r.File {
		e.file(level+1, sr)
	}
	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)

	e.change(level+1, &r.Change)
	e.creation(level+1, &r.Creation)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.citationList(level+1, r.Citation)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) repository(r *RepositoryRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}

	level := 0
	e.tagWithID(level, "REPO", r.Xref)
	e.maybeTag(level+1, "NAME", r.Name)
	e.address(level+1, &r.Address)
	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.creation(level+1, &r.Creation)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) source(r *SourceRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}

	level := 0
	e.tagWithID(level, "SOUR", r.Xref)
	e.maybeTagWithText(level+1, "TITL", r.Title)
	if r.Data != nil {
		e.tag(level+1, "DATA", "")
		for _, sr := range r.Data.Event {
			e.tag(level+2, "EVEN", sr.Kind)
			e.maybeTag(level+3, "DATE", sr.Date)
			e.maybeTag(level+3, "PLAC", sr.Place)
		}

	}

	e.maybeTagWithText(level+1, "AUTH", r.Originator)
	e.maybeTagWithText(level+1, "ABBR", r.FiledBy)
	e.maybeTagWithText(level+1, "PUBL", r.PublicationFacts)
	e.maybeTagWithText(level+1, "TEXT", r.Text)

	if r.Repository != nil && r.Repository.Repository != nil && r.Repository.Repository.Xref != "" {
		e.tagWithPointer(level+1, "REPO", r.Repository.Repository.Xref)
		e.noteList(level+2, r.Repository.Note)
		for _, sr := range r.Repository.CallNumber {
			e.tag(level+2, "CALN", sr.CallNumber)
			e.maybeTag(level+3, "MEDI", sr.MediaType)
		}
	}

	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.maybeTag(level+1, "RESN", r.RestrictionNotice)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.creation(level+1, &r.Creation)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) submitter(level int, r *SubmitterRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithID(level, "SUBM", r.Xref)
	e.maybeTagWithText(level+1, "NAME", r.Name)
	e.address(level+1, r.Address)
	e.mediaRefList(level+1, r.Media)

	for _, l := range r.Language {
		e.maybeTagWithText(level+1, "LANG", l)
	}
	e.maybeTagWithText(level+1, "RFN", r.SubmitterRecordFileID)
	e.maybeTag(level+1, "UID", r.UID)
	e.externalIDList(level+1, r.ExternalID)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.change(level+1, r.Change)
	e.creation(level+1, &r.Creation)
}

func (e *Encoder) trailer(r *Trailer) {
	if e.err != nil {
		return
	}
	e.tag(0, "TRLR", "")
}

func (e *Encoder) name(level int, r *NameRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.maybeTagWithText(level, "NAME", r.Name)
	e.maybeTagWithText(level+1, "TYPE", r.Type)
	e.maybeTagWithText(level+1, "NPFX", r.NamePiecePrefix)
	e.maybeTagWithText(level+1, "GIVN", r.NamePieceGiven)
	e.maybeTagWithText(level+1, "NICK", r.NamePieceNick)
	e.maybeTagWithText(level+1, "SPFX", r.NamePieceSurnamePrefix)
	e.maybeTagWithText(level+1, "SURN", r.NamePieceSurname)
	e.maybeTagWithText(level+1, "NSFX", r.NamePieceSuffix)
	e.maybeTag(level+1, "RESN", r.RestrictionNotice)

	for _, vn := range r.Phonetic {
		e.variantName(level+1, "FONE", vn)
	}
	for _, vn := range r.Romanized {
		e.variantName(level+1, "ROMN", vn)
	}

	e.translationList(level+1, r.Translation)
	e.citationList(level+1, r.Citation)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) change(level int, r *ChangeRecord) {
	if e.err != nil {
		return
	}
	if r == nil || (r.Date == "" && r.Time == "" && len(r.Note) == 0) {
		return
	}
	e.tagWithText(level, "CHAN", "")
	e.maybeTagWithText(level+1, "DATE", r.Date)
	e.maybeTagWithText(level+2, "TIME", r.Time)

	e.noteList(level+1, r.Note)
}

func (e *Encoder) noteList(level int, rs []*NoteRecord) {
	if e.err != nil {
		return
	}
	for _, sr := range rs {
		e.note(level, sr)
	}
}

func (e *Encoder) note(level int, r *NoteRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithText(level, "NOTE", r.Note)
	e.maybeTag(level+1, "MIME", r.Mime)
	e.maybeTag(level+1, "LANG", r.Language)
	e.translationList(level+1, r.Translation)
	e.citationList(level+1, r.Citation)
}

func (e *Encoder) citationList(level int, rs []*CitationRecord) {
	if e.err != nil {
		return
	}
	for _, sr := range rs {
		e.citation(level, sr)
	}
}

func (e *Encoder) citation(level int, r *CitationRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Source == nil {
		e.err = fmt.Errorf("source missing")
		return
	}
	if r.Source.Xref == "" {
		e.tag(level, "SOUR", "")
	} else {
		e.tagWithPointer(level, "SOUR", r.Source.Xref)
	}
	e.maybeTagWithText(level+1, "PAGE", r.Page)
	e.maybeTagWithText(level+1, "QUAY", r.Quay)

	if r.Data.Date != "" || len(r.Data.Text) != 0 || len(r.Data.UserDefined) != 0 {
		e.data(level+1, &r.Data)
	}

	e.noteList(level+1, r.Note)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) data(level int, r *DataRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}

	e.tag(level, "DATA", "")
	e.maybeTag(level+1, "DATE", r.Date)
	for _, sr := range r.Text {
		e.maybeTagWithText(level+1, "TEXT", sr)
	}
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) familyLink(level int, tag string, r *FamilyLinkRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Family == nil {
		e.err = fmt.Errorf("family missing")
		return
	}
	if r.Family.Xref == "" {
		e.err = fmt.Errorf("family missing xref")
		return
	}
	e.tagWithPointer(level, tag, r.Family.Xref)
	e.maybeTagWithText(level+1, "PEDI", r.Type)
	e.noteList(level+1, r.Note)
}

func (e *Encoder) file(level int, r *FileRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.maybeTagWithText(level, "FILE", r.Name)
	e.maybeTagWithText(level+1, "FORM", r.Format)
	e.maybeTagWithText(level+2, "TYPE", r.FormatType)
	e.maybeTagWithText(level+1, "TITL", r.Title)
	e.crop(level+1, r.Crop)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) userReferenceList(level int, rs []*UserReferenceRecord) {
	if e.err != nil {
		return
	}
	for _, sr := range rs {
		e.userReference(level+1, sr)
	}
}

func (e *Encoder) userReference(level int, r *UserReferenceRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.maybeTagWithText(level, "REFN", r.Number)
	e.maybeTagWithText(level+1, "TYPE", r.Type)
}

func (e *Encoder) individualRef(level int, tag string, r *IndividualRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Xref == "" {
		e.err = fmt.Errorf("individual missing xref for %s", tag)
		return
	}
	e.tagWithPointer(level, tag, r.Xref)
}

func (e *Encoder) familyRef(level int, tag string, r *FamilyRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Xref == "" {
		e.err = fmt.Errorf("family missing xref")
		return
	}
	e.tagWithPointer(level, tag, r.Xref)
}

func (e *Encoder) mediaRefList(level int, rs []*MediaRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.mediaRef(level, r)
	}
}

func (e *Encoder) mediaRef(level int, r *MediaRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Xref != "" {
		e.tagWithPointer(level, "OBJE", r.Xref)
		return
	}

	// inline media
	e.tag(level, "OBJE", "")
	for _, sr := range r.File {
		e.file(level+1, sr)
	}
	e.maybeTagWithText(level+1, "TITL", r.Title)
}

func (e *Encoder) eventList(level int, rs []*EventRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.event(level, r)
	}
}

func (e *Encoder) event(level int, r *EventRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		e.err = fmt.Errorf("event not specified")
		return
	}
	e.tag(level, r.Tag, r.Value)
	e.maybeTagWithText(level+1, "TYPE", r.Type)
	e.maybeTagWithText(level+1, "DATE", r.Date)
	e.maybeTag(level+1, "SDATE", r.SortDate)

	e.address(level+1, &r.Address)
	e.place(level+1, &r.Place)
	e.maybeTag(level+1, "AGE", r.Age)
	e.maybeTag(level+1, "AGNC", r.ResponsibleAgency)
	e.maybeTag(level+1, "RELI", r.ReligiousAffiliation)
	e.maybeTag(level+1, "CAUS", r.Cause)
	e.maybeTag(level+1, "RESN", r.RestrictionNotice)

	if r.ChildInFamily != nil {
		e.familyRef(level+1, "FAMC", r.ChildInFamily)
		if r.Tag == "ADOP" {
			e.maybeTag(level+2, "ADOP", r.AdoptedByParent)
		}
	}
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.citationList(level+1, r.Citation)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}

func (e *Encoder) nonEventList(level int, rs []*EventRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.nonEvent(level, r)
	}
}

func (e *Encoder) nonEvent(level int, r *EventRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tag(level, "NO", r.Value)
	e.maybeTagWithText(level+1, "DATE", r.Date)
	e.noteList(level+1, r.Note)
	e.sharedNoteRefList(level+1, r.SharedNote)
	e.citationList(level+1, r.Citation)
}

func (e *Encoder) association(level int, r *AssociationRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithOptionalPointer(level, "ASSO", r.Xref)
	e.maybeTag(level+1, "RELA", r.Relation)
	if r.Role != "" {
		e.tag(level+1, "ROLE", r.Role)
		e.maybeTag(level+2, "PHRASE", r.Phrase)
	}
	e.citationList(level+1, r.Citation)
	e.noteList(level+1, r.Note)
}

func (e *Encoder) sharedNoteRecord(r *SharedNoteRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithIDAndText(0, "SNOTE", r.Xref, r.Note)
	e.maybeTag(1, "MIME", r.Mime)
	e.maybeTag(1, "LANG", r.Language)
	e.translationList(1, r.Translation)
	e.citationList(1, r.Citation)
	e.userReferenceList(1, r.UserReference)
	e.maybeTagWithText(1, "RIN", r.AutomatedRecordId)
	e.change(1, &r.Change)
	e.creation(1, &r.Creation)
	e.userDefinedList(1, r.UserDefined)
}

func (e *Encoder) tagWithIDAndText(level int, tag string, id string, value string) {
	if e.err != nil {
		return
	}
	if id == "" {
		e.err = fmt.Errorf("tag %s missing id", tag)
		return
	}

	conts := strings.Split(value, "\n")
	first := conts[0]
	if first != "" {
		first = " " + first
	}
	if _, err := e.w.WriteString(fmt.Sprintf("%d @%s@ %s%s\n", level, id, tag, first)); err != nil {
		e.err = fmt.Errorf("write tag with id %s @%s@: %w", tag, id, err)
		return
	}

	for i := 1; i < len(conts); i++ {
		e.textOneLine(level+1, "CONT", conts[i])
	}
}

func (e *Encoder) sharedNoteRefList(level int, rs []*SharedNoteRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.sharedNoteRef(level, r)
	}
}

func (e *Encoder) sharedNoteRef(level int, r *SharedNoteRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	if r.Xref != "" {
		e.tagWithPointer(level, "SNOTE", r.Xref)
	}
}

func (e *Encoder) translationList(level int, rs []*TranslationRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.translation(level, r)
	}
}

func (e *Encoder) translation(level int, r *TranslationRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithText(level, "TRAN", r.Value)
	e.maybeTag(level+1, "LANG", r.Language)
	e.maybeTag(level+1, "MIME", r.Mime)
}

func (e *Encoder) externalIDList(level int, rs []*ExternalIDRecord) {
	if e.err != nil {
		return
	}
	for _, r := range rs {
		e.externalID(level, r)
	}
}

func (e *Encoder) externalID(level int, r *ExternalIDRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tag(level, "EXID", r.ID)
	e.maybeTag(level+1, "TYPE", r.Type)
}

func (e *Encoder) creation(level int, r *CreationRecord) {
	if e.err != nil {
		return
	}
	if r == nil || (r.Date == "" && r.Time == "") {
		return
	}
	e.tag(level, "CREA", "")
	e.maybeTag(level+1, "DATE", r.Date)
	e.maybeTag(level+2, "TIME", r.Time)
}

func (e *Encoder) schema(level int, r *SchemaRecord) {
	if e.err != nil {
		return
	}
	if r == nil || len(r.Tag) == 0 {
		return
	}
	e.tag(level, "SCHMA", "")
	for _, st := range r.Tag {
		if st == nil {
			continue
		}
		value := st.Tag
		if st.URI != "" {
			value = st.Tag + " " + st.URI
		}
		e.tag(level+1, "TAG", value)
	}
}

func (e *Encoder) crop(level int, r *CropRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tag(level, "CROP", "")
	e.maybeTag(level+1, "TOP", r.Top)
	e.maybeTag(level+1, "LEFT", r.Left)
	e.maybeTag(level+1, "WIDTH", r.Width)
	e.maybeTag(level+1, "HEIGHT", r.Height)
}

func (e *Encoder) variantName(level int, tag string, r *VariantNameRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tag(level, tag, r.Name)
	e.maybeTag(level+1, "TYPE", r.Type)
	e.maybeTagWithText(level+1, "NPFX", r.NamePiecePrefix)
	e.maybeTagWithText(level+1, "GIVN", r.NamePieceGiven)
	e.maybeTagWithText(level+1, "NICK", r.NamePieceNick)
	e.maybeTagWithText(level+1, "SPFX", r.NamePieceSurnamePrefix)
	e.maybeTagWithText(level+1, "SURN", r.NamePieceSurname)
	e.maybeTagWithText(level+1, "NSFX", r.NamePieceSuffix)
	e.citationList(level+1, r.Citation)
	e.noteList(level+1, r.Note)
}
