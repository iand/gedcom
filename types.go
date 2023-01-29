/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

type Gedcom struct {
	Header      *Header
	Family      []*FamilyRecord
	Individual  []*IndividualRecord
	Media       []*MediaRecord
	Repository  []*RepositoryRecord
	Source      []*SourceRecord
	Submitter   []*SubmitterRecord
	Trailer     *Trailer
	UserDefined []UserDefinedTag
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
	UserDefined     []UserDefinedTag
}

type SubmissionRecord struct {
	Xref string
}

type Trailer struct{}

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
	Xref                      string
	Name                      []*NameRecord
	Sex                       string
	Event                     []*EventRecord
	Attribute                 []*EventRecord
	Parents                   []*FamilyLinkRecord
	Family                    []*FamilyLinkRecord
	Submitter                 []*SubmitterRecord
	Association               []*AssociationRecord
	PermanentRecordFileNumber string
	AncestralFileNumber       string
	UserReference             []*UserReferenceRecord
	AutomatedRecordId         string
	Change                    ChangeRecord
	Note                      []*NoteRecord
	Citation                  []*CitationRecord
	Media                     []*MediaRecord
	UserDefined               []UserDefinedTag
}

type MediaRecord struct {
	Xref              string
	File              []*FileRecord
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	Note              []*NoteRecord
	Citation          []*CitationRecord
	UserDefined       []UserDefinedTag
}

type FileRecord struct {
	Name        string
	Format      string
	FormatType  string
	Title       string
	UserDefined []UserDefinedTag
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
	Xref              string
	Name              string
	Address           AddressRecord
	Note              []*NoteRecord
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	UserDefined       []UserDefinedTag
}

type SourceRecord struct {
	Xref              string
	Title             string
	Data              *SourceDataRecord
	Originator        string
	FiledBy           string
	PublicationFacts  string
	Text              string
	Repository        *SourceRepositoryRecord
	UserReference     []*UserReferenceRecord
	AutomatedRecordId string
	Change            ChangeRecord
	Note              []*NoteRecord
	Media             []*MediaRecord
	UserDefined       []UserDefinedTag
}

type SourceDataRecord struct {
	Event []*SourceEventRecord
}

type SourceEventRecord struct {
	Kind  string
	Date  string
	Place string
}

type SourceRepositoryRecord struct {
	Repository *RepositoryRecord
	Note       []*NoteRecord
	CallNumber []*SourceCallNumberRecord
}

type SourceCallNumberRecord struct {
	CallNumber string
	MediaType  string
}

type CitationRecord struct {
	Source      *SourceRecord
	Page        string
	Data        DataRecord
	Quay        string
	Media       []*MediaRecord
	Note        []*NoteRecord
	UserDefined []UserDefinedTag
}

type SubmitterRecord struct {
	Xref string
}

type NameRecord struct {
	Name                   string
	NamePiecePrefix        string
	NamePieceGiven         string
	NamePieceNick          string
	NamePieceSurnamePrefix string
	NamePieceSurname       string
	NamePieceSuffix        string
	Citation               []*CitationRecord
	Note                   []*NoteRecord
}

type DataRecord struct {
	Date        string
	Text        []string
	UserDefined []UserDefinedTag
}

type EventRecord struct {
	Tag                  string
	Value                string
	Type                 string
	Date                 string
	Place                PlaceRecord
	Address              AddressRecord
	Age                  string
	ResponsibleAgency    string
	ReligiousAffiliation string
	Cause                string
	RestrictionNotice    string        // 5.5.1
	ChildInFamily        *FamilyRecord // link to parent family for birth events
	AdoptedByParent      string        // for adoption event, one of HUSB,WIFE,BOTH
	Citation             []*CitationRecord
	Media                []*MediaRecord
	Note                 []*NoteRecord
	UserDefined          []UserDefinedTag
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

// See https://www.tamurajones.net/GEDCOMADDR.xhtml for very informative analysis of the ADDR structure
type AddressRecord struct {
	Address []*AddressDetail
	Phone   []string
	Email   []string // 5.5.1
	Fax     []string // 5.5.1
	WWW     []string // 5.5.1
}

type AddressDetail struct {
	Full       string // The full address as found in free-form fields which may be optionally broken down using following structured fields
	Line1      string
	Line2      string
	Line3      string // 5.5.1
	City       string
	State      string
	PostalCode string
	Country    string
}

// A UserDefinedTag is a tag that is not defined in the GEDCOM specification but is included by the publisher of the
// data. In GEDCOM user defined tags must be prefixed with an underscore. This is preserved in the Tag field.
type UserDefinedTag struct {
	Tag         string
	Value       string
	Xref        string
	Level       int
	UserDefined []UserDefinedTag
}

type AssociationRecord struct {
	Xref     string
	Relation string
	Citation []*CitationRecord
	Note     []*NoteRecord
}
