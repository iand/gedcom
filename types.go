/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

type Gedcom struct {
	Header           *Header
	SubmissionRecord *SubmissionRecord
	Family           []*FamilyRecord
	Individual       []*IndividualRecord
	Media            []*MediaRecord
	Repository       []*RepositoryRecord
	Source           []*SourceRecord
	Submitter        []*SubmitterRecord
	Trailer          *Trailer
}

type Header struct {
	SourceSystem SystemRecord
}

type SystemRecord struct {
	Id              string
	Version         string
	ProductName     string
	BusinessName    string
	Address         AddressRecord
	SourceName      string
	SourceDate      string
	SourceCopyright string
}

type SubmissionRecord struct {
	Xref string
}

type Trailer struct {
}

type FamilyRecord struct {
	Xref    string
	Husband *IndividualRecord
	Wife    *IndividualRecord
	Child   []*IndividualRecord
	Event   []*EventRecord
}

type IndividualRecord struct {
	Xref      string
	Name      []*NameRecord
	Sex       string
	Event     []*EventRecord
	Attribute []*EventRecord
	Parents   []*FamilyLinkRecord
	Family    []*FamilyLinkRecord
}

type MediaRecord struct {
}

type RepositoryRecord struct {
}

type SourceRecord struct {
	Xref  string
	Title string
	Media []*MediaRecord
	Note  []*NoteRecord
}

type CitationRecord struct {
	Source *SourceRecord
	Page   string
	Data   DataRecord
	Quay   string
	Media  []*MediaRecord
	Note   []*NoteRecord
}

type SubmitterRecord struct {
}

type NameRecord struct {
	Name     string
	Citation []*CitationRecord
	Note     []*NoteRecord
}

type DataRecord struct {
	Date string
	Text []string
}

type EventRecord struct {
	Tag      string
	Value    string
	Type     string
	Date     string
	Place    PlaceRecord
	Address  AddressRecord
	Age      string
	Agency   string
	Cause    string
	Citation []*CitationRecord
	Media    []*MediaRecord
	Note     []*NoteRecord
}

type NoteRecord struct {
	Note     string
	Citation []*CitationRecord
}

type PlaceRecord struct {
	Name     string
	Citation []*CitationRecord
	Note     []*NoteRecord
}

type FamilyLinkRecord struct {
	Family *FamilyRecord
	Type   string
	Note   []*NoteRecord
}

type AddressRecord struct {
	Full       string
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
	Phone      string
}
