/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

type Gedcom struct {
	Header     *Header
	Family     []*FamilyRecord
	Individual []*IndividualRecord
	Media      []*MediaRecord
	Repository []*RepositoryRecord
	Source     []*SourceRecord
	Submitter  []*SubmitterRecord
	Trailer    *Trailer
}

// A Header contains information about the GEDCOM file.
type Header struct {
	SourceSystem        SystemRecord
	Destination         string
	Date                string
	Time                string
	Submitter           *SubmitterRecord
	Submission          *SubmissionRecord
	Filename            string
	Copyright           string
	Version             string
	Form                string
	CharacterSet        string
	CharacterSetVersion string
	Language            string
	Place               PlaceRecord
	Note                string
	UserDefined         []UserDefinedTag
}

// A SystemRecord contains information about the system that produced the GEDCOM.
type SystemRecord struct {
	Xref            string
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
	Xref              string
	Husband           *IndividualRecord
	Wife              *IndividualRecord
	Child             []*IndividualRecord
	Event             []*EventRecord
	NumberOfChildren  string
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	Note              []*NoteRecord
	Citation          []*CitationRecord
	Media             []*MediaRecord
	UserDefined       []UserDefinedTag
}

type IndividualRecord struct {
	Xref              string
	Name              []*NameRecord
	Sex               string
	Event             []*EventRecord
	Attribute         []*EventRecord
	Parents           []*FamilyLinkRecord
	Family            []*FamilyLinkRecord
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	Note              []*NoteRecord
	Citation          []*CitationRecord
	Media             []*MediaRecord
	UserDefined       []UserDefinedTag
}

type MediaRecord struct {
	Xref              string
	File              []*FileRecord
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	Note              []*NoteRecord
	Citation          []*CitationRecord
}

type FileRecord struct {
	Name       string
	Format     string
	FormatType string
	Title      string
}

type UserReferenceRecord struct {
	Number string
	Type   string
}

type ChangeRecord struct {
	Date string
	Time string
	Note []*NoteRecord
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
	Xref string
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
	Phone      []string
	Email      []string
	Fax        []string
	WWW        []string
}

// A UserDefinedTag is a tag that is not defined in the GEDCOM specification but is included by the publisher of the
// data. In GEDCOM user defined tags must be prefixed with an underscore. This is preserved in the Tag field.
type UserDefinedTag struct {
	Tag   string
	Value string
	Xref  string
	Level int
}
