/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

// Gedcom represents a complete GEDCOM file. It is the top-level container
// returned by [Decoder.Decode] and accepted by [Encoder.Encode].
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

// SubmissionRecord contains information about a batch submission of genealogical data.
type SubmissionRecord struct {
	Xref string
}

// Trailer marks the end of a GEDCOM file. It contains no data.
type Trailer struct{}

// FamilyRecord represents a family unit in GEDCOM. It links individuals
// together as husband, wife, and children, and contains family events
// such as marriage, divorce, and census records.
type FamilyRecord struct {
	Xref              string                 // Unique cross-reference identifier for this family (e.g., "@F1@")
	Husband           *IndividualRecord      // Pointer to the husband/father individual
	Wife              *IndividualRecord      // Pointer to the wife/mother individual
	Child             []*IndividualRecord    // Pointers to child individuals
	Event             []*EventRecord         // Family events (marriage, divorce, census, etc.)
	NumberOfChildren  string                 // Total number of children, may differ from len(Child)
	UserReference     []*UserReferenceRecord // User-provided reference numbers
	AutomatedRecordId string                 // Unique record ID assigned by the source system
	Change            ChangeRecord           // Record of when this record was last modified
	Note              []*NoteRecord          // Notes attached to this family
	Citation          []*CitationRecord      // Source citations for this family
	Media             []*MediaRecord         // Media objects (photos, documents) for this family
	UserDefined       []UserDefinedTag       // User-defined tags (prefixed with underscore)
}

// IndividualRecord represents a person in GEDCOM. It contains the person's
// names, sex, life events, attributes, and links to their families.
type IndividualRecord struct {
	Xref                      string                 // Unique cross-reference identifier (e.g., "@I1@")
	Name                      []*NameRecord          // Names (may have multiple for maiden names, aliases)
	Sex                       string                 // Sex: "M" for male, "F" for female, "U" for unknown
	Event                     []*EventRecord         // Life events (birth, death, burial, etc.)
	Attribute                 []*EventRecord         // Attributes (occupation, residence, education, etc.)
	Parents                   []*FamilyLinkRecord    // Links to families where this person is a child
	Family                    []*FamilyLinkRecord    // Links to families where this person is a spouse
	Submitter                 []*SubmitterRecord     // Submitters of this record
	Association               []*AssociationRecord   // Associations with other individuals
	PermanentRecordFileNumber string                 // Permanent record file number
	AncestralFileNumber       string                 // Ancestral file number
	UserReference             []*UserReferenceRecord // User-provided reference numbers
	AutomatedRecordId         string                 // Unique record ID assigned by the source system
	Change                    ChangeRecord           // Record of when this record was last modified
	Note                      []*NoteRecord          // Notes attached to this individual
	Citation                  []*CitationRecord      // Source citations for this individual
	Media                     []*MediaRecord         // Media objects (photos, documents)
	UserDefined               []UserDefinedTag       // User-defined tags (prefixed with underscore)
}

// MediaRecord represents a multimedia object such as a photo, document,
// or audio/video recording. It can be referenced by individuals, families,
// events, and sources.
type MediaRecord struct {
	Xref              string                 // Unique cross-reference identifier (e.g., "@M1@")
	File              []*FileRecord          // File references for this media object
	Title             string                 // Title or description of the media
	UserReference     []*UserReferenceRecord // User-provided reference numbers
	AutomatedRecordId string                 // Unique record ID assigned by the source system
	Change            ChangeRecord           // Record of when this record was last modified
	Note              []*NoteRecord          // Notes attached to this media
	Citation          []*CitationRecord      // Source citations
	UserDefined       []UserDefinedTag       // User-defined tags
}

// FileRecord contains information about a multimedia file.
type FileRecord struct {
	Name        string           // File path or URL
	Format      string           // File format (e.g., "jpeg", "gif", "pdf")
	FormatType  string           // Media type (e.g., "photo", "document")
	Title       string           // Title or caption for this file
	UserDefined []UserDefinedTag // User-defined tags
}

// UserReferenceRecord contains a user-defined reference number for a record.
type UserReferenceRecord struct {
	Number string // The reference number
	Type   string // The type of reference
}

// ChangeRecord indicates when a record was last modified.
type ChangeRecord struct {
	Date string        // Date of the change in GEDCOM date format
	Time string        // Time of the change
	Note []*NoteRecord // Notes about the change
}

// RepositoryRecord represents a repository where source documents are held,
// such as a library, archive, or private collection.
type RepositoryRecord struct {
	Xref              string                 // Unique cross-reference identifier (e.g., "@R1@")
	Name              string                 // Name of the repository
	Address           AddressRecord          // Address of the repository
	Note              []*NoteRecord          // Notes about the repository
	UserReference     []*UserReferenceRecord // User-provided reference numbers
	AutomatedRecordId string                 // Unique record ID assigned by the source system
	Change            ChangeRecord           // Record of when this record was last modified
	UserDefined       []UserDefinedTag       // User-defined tags
}

// SourceRecord represents a source of genealogical information, such as
// a book, document, website, or oral interview.
type SourceRecord struct {
	Xref              string                  // Unique cross-reference identifier (e.g., "@S1@")
	Title             string                  // Title of the source
	Data              *SourceDataRecord       // Data recorded from the source
	Originator        string                  // Author or creator of the source
	FiledBy           string                  // Person who filed the source
	PublicationFacts  string                  // Publication information
	Text              string                  // Verbatim text from the source
	Repository        *SourceRepositoryRecord // Repository where the source is held
	UserReference     []*UserReferenceRecord  // User-provided reference numbers
	AutomatedRecordId string                  // Unique record ID assigned by the source system
	Change            ChangeRecord            // Record of when this record was last modified
	Note              []*NoteRecord           // Notes about the source
	Media             []*MediaRecord          // Media objects (photos of documents, etc.)
	UserDefined       []UserDefinedTag        // User-defined tags
}

// SourceDataRecord contains data recorded from a source.
type SourceDataRecord struct {
	Event []*SourceEventRecord // Events covered by the source
}

// SourceEventRecord describes an event type covered by a source.
type SourceEventRecord struct {
	Kind  string // Type of event (e.g., "BIRT", "MARR", "DEAT")
	Date  string // Date range covered by the source
	Place string // Place jurisdiction covered by the source
}

// SourceRepositoryRecord links a source to its repository.
type SourceRepositoryRecord struct {
	Repository *RepositoryRecord         // The repository holding the source
	Note       []*NoteRecord             // Notes about the source at this repository
	CallNumber []*SourceCallNumberRecord // Call numbers for locating the source
}

// SourceCallNumberRecord contains a call number for a source in a repository.
type SourceCallNumberRecord struct {
	CallNumber string // The call number or shelf location
	MediaType  string // Type of media (e.g., "book", "microfilm")
}

// CitationRecord represents a citation to a source for a specific claim.
// It links to a [SourceRecord] and provides details about where in the
// source the information was found.
type CitationRecord struct {
	Source      *SourceRecord    // The source being cited
	Page        string           // Page number or location within the source
	Data        DataRecord       // Data extracted from the source
	Quay        string           // Quality assessment (0-3, with 3 being direct evidence)
	Media       []*MediaRecord   // Media objects (e.g., photo of the source page)
	Note        []*NoteRecord    // Notes about this citation
	UserDefined []UserDefinedTag // User-defined tags
}

// SubmitterRecord contains information about the person or organization
// that submitted the genealogical data.
type SubmitterRecord struct {
	Xref                  string         // Unique cross-reference identifier (e.g., "@SUBM1@")
	Name                  string         // Name of the submitter
	Address               *AddressRecord // Address of the submitter
	Media                 []*MediaRecord // Media objects (e.g., photo of submitter)
	Language              []string       // Languages used by the submitter
	SubmitterRecordFileID string         // Submitter record file identifier
	AutomatedRecordId     string         // Unique record ID assigned by the source system
	Note                  []*NoteRecord  // Notes from the submitter
	Change                *ChangeRecord  // Record of when this record was last modified
}

// NameRecord represents a name for an individual. An individual may have
// multiple names (e.g., maiden name, married name, aliases).
// The Name field contains the full name in GEDCOM format: "Given Name /Surname/ Suffix".
// Use [SplitPersonalName] to parse the Name field into components.
type NameRecord struct {
	Name                   string               // Full name in GEDCOM format (e.g., "John /Smith/ Jr.")
	Type                   string               // Name type (e.g., "birth", "married", "aka")
	NamePiecePrefix        string               // Name prefix (e.g., "Dr.", "Rev.")
	NamePieceGiven         string               // Given name(s)
	NamePieceNick          string               // Nickname
	NamePieceSurnamePrefix string               // Surname prefix (e.g., "van", "de")
	NamePieceSurname       string               // Surname
	NamePieceSuffix        string               // Name suffix (e.g., "Jr.", "III")
	Phonetic               []*VariantNameRecord // Phonetic variants of the name
	Romanized              []*VariantNameRecord // Romanized variants of the name
	Citation               []*CitationRecord    // Source citations for this name
	Note                   []*NoteRecord        // Notes about this name
	UserDefined            []UserDefinedTag     // User-defined tags
}

// VariantNameRecord represents a phonetic or romanized variant of a name.
type VariantNameRecord struct {
	Name                   string            // Full variant name
	Type                   string            // Variant type (e.g., "kana", "hangul" for phonetic)
	NamePiecePrefix        string            // Name prefix
	NamePieceGiven         string            // Given name(s)
	NamePieceNick          string            // Nickname
	NamePieceSurnamePrefix string            // Surname prefix
	NamePieceSurname       string            // Surname
	NamePieceSuffix        string            // Name suffix
	Citation               []*CitationRecord // Source citations
	Note                   []*NoteRecord     // Notes
}

// DataRecord contains data extracted from a source citation.
type DataRecord struct {
	Date        string           // Date the data was recorded
	Text        []string         // Verbatim text from the source
	UserDefined []UserDefinedTag // User-defined tags
}

// EventRecord represents an event or attribute in a person's or family's life.
// Common event tags include BIRT (birth), DEAT (death), MARR (marriage),
// BURI (burial), CHR (christening), DIV (divorce).
// Common attribute tags include OCCU (occupation), RESI (residence),
// EDUC (education), RELI (religion).
type EventRecord struct {
	Tag                  string            // Event type tag (e.g., "BIRT", "DEAT", "MARR")
	Value                string            // Event value, often "Y" to indicate event occurred
	Type                 string            // Detailed event type for generic EVEN tags
	Date                 string            // Date in GEDCOM date format
	Place                PlaceRecord       // Location where the event occurred
	Address              AddressRecord     // Address associated with the event
	Age                  string            // Age of the individual at the time of the event
	ResponsibleAgency    string            // Agency responsible for the record
	ReligiousAffiliation string            // Religious affiliation associated with event
	Cause                string            // Cause (e.g., cause of death)
	RestrictionNotice    string            // Privacy restriction (GEDCOM 5.5.1)
	ChildInFamily        *FamilyRecord     // Link to parent family for birth events
	AdoptedByParent      string            // For adoption: "HUSB", "WIFE", or "BOTH"
	Citation             []*CitationRecord // Source citations for this event
	Media                []*MediaRecord    // Media objects (e.g., photos, certificates)
	Note                 []*NoteRecord     // Notes about this event
	UserDefined          []UserDefinedTag  // User-defined tags
}

// NoteRecord contains a note or comment attached to a record.
type NoteRecord struct {
	Note     string            // The note text
	Citation []*CitationRecord // Source citations for the note
}

// PlaceRecord represents a geographic location. The Name field typically
// contains a comma-separated jurisdiction hierarchy (e.g., "City, County, State, Country").
type PlaceRecord struct {
	Name      string                    // Place name (jurisdiction hierarchy)
	Phonetic  []*VariantPlaceNameRecord // Phonetic variants of the place name
	Romanized []*VariantPlaceNameRecord // Romanized variants of the place name
	Latitude  string                    // Latitude in GEDCOM format (e.g., "N50.9333")
	Longitude string                    // Longitude in GEDCOM format (e.g., "W1.8")
	Citation  []*CitationRecord         // Source citations
	Note      []*NoteRecord             // Notes about the place
}

// VariantPlaceNameRecord represents a phonetic or romanized variant of a place name.
type VariantPlaceNameRecord struct {
	Name string // The variant place name
	Type string // Variant type (e.g., "kana", "hangul" for phonetic)
}

// FamilyLinkRecord links an individual to a family. It is used both for
// linking children to their parents (via [IndividualRecord.Parents]) and
// for linking spouses to their families (via [IndividualRecord.Family]).
type FamilyLinkRecord struct {
	Family *FamilyRecord // Pointer to the linked family
	Type   string        // Relationship type (e.g., "birth", "adopted", "foster")
	Note   []*NoteRecord // Notes about this family link
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

// AssociationRecord links an individual to another person with a defined
// relationship, such as godparent, witness, or friend.
type AssociationRecord struct {
	Xref     string            // Cross-reference to the associated individual
	Relation string            // Relationship type (e.g., "godparent", "witness")
	Citation []*CitationRecord // Source citations
	Note     []*NoteRecord     // Notes about this association
}
