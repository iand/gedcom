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
	w   *bufio.Writer
	err error
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	bw := bufio.NewWriter(w)
	return &Encoder{
		w: bw,
	}
}

func (e *Encoder) Encode(g *Gedcom) error {
	e.header(g.Header)

	for _, r := range g.Individual {
		e.individual(r)
	}

	for _, r := range g.Family {
		e.family(r)
	}

	for _, r := range g.Media {
		e.media(0, r)
	}

	for _, r := range g.Repository {
		e.repository(r)
	}

	for _, r := range g.Source {
		e.source(r)
	}

	for _, r := range g.Submitter {
		e.submitter(0, r)
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

	if len(value) <= 246 {
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
	e.maybeTag(1, "CHAR", h.CharacterSet)
	e.maybeTag(2, "VERS", h.CharacterSetVersion)
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

	if h.Version != "" || h.Form != "" {
		e.tag(1, "GEDC", "")
		if h.Version != "" {
			e.tag(2, "VERS", h.Version)
		}
		if h.Form != "" {
			e.tag(2, "FORM", h.Form)
		}
	}
	e.maybeTag(1, "LANG", h.Language)
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
	if r.Name == "" && len(r.Phonetic) == 0 && len(r.Romanized) == 0 && r.Latitude == "" && r.Longitude == "" && len(r.Note) == 0 && len(r.Citation) == 0 {
		return
	}

	e.tag(level, "PLAC", r.Name)
	for _, sr := range r.Phonetic {
		e.tag(level+1, "FONE", sr.Name)
		e.maybeTag(level+1, "TYPE", sr.Type)
	}

	for _, sr := range r.Romanized {
		e.tag(level+1, "ROMN", sr.Name)
		e.maybeTag(level+2, "TYPE", sr.Type)
	}

	if r.Latitude != "" || r.Longitude != "" {
		e.tag(level+1, "MAP", "")
		e.maybeTag(level+2, "LATI", r.Latitude)
		e.maybeTag(level+2, "LONG", r.Longitude)

	}

	e.noteList(level+1, r.Note)
	e.citationList(level+1, r.Citation)
}

func (e *Encoder) variantPlaceName(level int, r *VariantPlaceNameRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tag(level, r.Name, r.Type)
	e.maybeTag(level+1, "TYPE", r.Type)
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
	if len(r.Association) > 0 {
		// TODO: Association
		// Association               []*AssociationRecord
		e.err = fmt.Errorf("not implemented: Association")
		return
	}

	e.maybeTagWithText(level+1, "RFN", r.PermanentRecordFileNumber)
	e.maybeTagWithText(level+1, "AFN", r.AncestralFileNumber)

	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.noteList(level+1, r.Note)
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
	e.maybeTag(level+1, "NCHI", r.NumberOfChildren)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.noteList(level+1, r.Note)
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
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)

	e.change(level+1, &r.Change)
	e.noteList(level+1, r.Note)
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
	e.noteList(level+1, r.Note)
	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
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

	e.userReferenceList(level+1, r.UserReference)
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.change(level+1, &r.Change)
	e.noteList(level+1, r.Note)
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
	e.maybeTagWithText(level+1, "RIN", r.AutomatedRecordId)
	e.noteList(level+1, r.Note)
	e.change(level+1, r.Change)
}

func (e *Encoder) trailer(r *Trailer) {
	if e.err != nil {
		return
	}
	// nothing to do
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

	if len(r.Phonetic) > 0 {
		// TODO: FONE
		// Phonetic               []*VariantNameRecord
		e.err = fmt.Errorf("not implemented: FONE")
		return
	}

	if len(r.Romanized) > 0 {
		// TODO: ROMN
		// Romanized              []*VariantNameRecord
		e.err = fmt.Errorf("not implemented: ROMN")
		return
	}

	e.citationList(level+1, r.Citation)
	e.noteList(level+1, r.Note)
	e.userDefinedList(level+1, r.UserDefined)

	return
}

func (e *Encoder) change(level int, r *ChangeRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
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
		e.note(level+1, sr)
	}
}

func (e *Encoder) note(level int, r *NoteRecord) {
	if e.err != nil {
		return
	}
	if r == nil {
		return
	}
	e.tagWithText(level, "NOTE", "")
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
	return
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

	e.address(level+1, &r.Address)
	e.place(level+1, &r.Place)
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
	e.citationList(level+1, r.Citation)
	e.mediaRefList(level+1, r.Media)
	e.userDefinedList(level+1, r.UserDefined)
}
