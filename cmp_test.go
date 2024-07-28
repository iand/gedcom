package gedcom

import "github.com/google/go-cmp/cmp"

// familyXrefComparer is a Comparer that compares FamilyLinkRecords only by Family xref
var familyXrefComparer = cmp.Comparer(func(a, b *FamilyLinkRecord) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return a == nil
	}

	if a.Family == nil {
		return b.Family == nil
	}

	if b.Family == nil {
		return a.Family == nil
	}

	return a.Family.Xref == b.Family.Xref
})

// individualXrefComparer is a Comparer that compares IndividualRecords only by xref
var individualXrefComparer = cmp.Comparer(func(a, b *IndividualRecord) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return a == nil
	}

	return a.Xref == b.Xref
})

// sourceXrefComparer is a Comparer that compares CitationRecords only by source xref
var sourceXrefComparer = cmp.Comparer(func(a, b *CitationRecord) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return a == nil
	}

	if a.Source == nil {
		return b.Source == nil
	}

	if b.Source == nil {
		return a.Source == nil
	}

	return a.Source.Xref == b.Source.Xref
})

// eventIgnoreComparer is a Comparer that ignores event comparisons
var eventIgnoreComparer = cmp.Comparer(func(a, b []*EventRecord) bool {
	return true
})

// mediaFileNameCompare is a Comparer that compares MediaRecord only by first file name
var mediaFileNameCompare = cmp.Comparer(func(a, b *MediaRecord) bool {
	if a == nil {
		return b == nil
	}

	if b == nil {
		return a == nil
	}

	if len(a.File) == 0 {
		return len(b.File) == 0
	}

	if len(b.File) == 0 {
		return len(a.File) == 0
	}

	return a.File[0].Name == b.File[0].Name
})
