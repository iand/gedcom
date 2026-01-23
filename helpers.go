package gedcom

import (
	"strings"
)

// ParsedName contains the components of a personal name after parsing
// with [SplitPersonalName].
type ParsedName struct {
	Full     string // Reconstructed full name without GEDCOM delimiters
	Given    string // Given name(s) / first name(s)
	Surname  string // Surname / family name / last name
	Suffix   string // Name suffix (e.g., "Jr.", "III", "PhD")
	Nickname string // Nickname, if present in quotes
}

// SplitPersonalName parses a GEDCOM-formatted personal name into its components.
// GEDCOM names use slashes to delimit the surname: "Given Names /Surname/ Suffix".
//
// Examples:
//
//	SplitPersonalName("John /Smith/")
//	// Returns: Given="John", Surname="Smith"
//
//	SplitPersonalName("John \"Jack\" /Smith/ Jr.")
//	// Returns: Given="John", Nickname="Jack", Surname="Smith", Suffix="Jr."
//
//	SplitPersonalName("Mary Jane /van der Berg/")
//	// Returns: Given="Mary Jane", Surname="van der Berg"
//
// The function also handles alternative surnames separated by slashes within
// the surname delimiters (e.g., "/Smith/Smyth/" becomes Surname="Smith/Smyth").
func SplitPersonalName(name string) ParsedName {
	name = strings.TrimSpace(name)

	parts := strings.Split(name, "/")
	if len(parts) == 1 {
		return ParsedName{
			Full:  name,
			Given: name,
		}
	}

	// Find a part that was delimited by slashes with no whitespace after the leading slash or before the following slash
	// That part is treated as the surname, anything before that part is treated as the given name, anything after is assumed to
	// be a suffix.
	for i := 1; i < len(parts); i++ {
		p := parts[i]
		if len(p) == 0 || p[0] == ' ' || p[len(p)-1] == ' ' {
			continue
		}

		pn := ParsedName{
			Given:   strings.TrimSpace(strings.Join(parts[:i], "/")),
			Surname: parts[i],
			Suffix:  strings.TrimSpace(strings.Join(parts[i+1:], "/")),
		}

		// Check for a nickname, which is usually in quotes after the given name
		if strings.HasSuffix(pn.Given, `"`) {
			pos := strings.Index(pn.Given, `"`)
			if pos != -1 && pos < len(pn.Given)-2 {
				pn.Nickname = strings.TrimSpace(pn.Given[pos+1 : len(pn.Given)-1])
				pn.Given = strings.TrimSpace(pn.Given[:pos])
			}
		}

		// See if there is a following part that could be part of the surname.
		// Some surnames may have alternatives: smith/smyth
		for j := i + 1; j < len(parts); j++ {
			p := parts[j]
			if len(p) == 0 || p[0] == ' ' || p[len(p)-1] == ' ' {
				// This part can't be part of the surname
				break
			}
			// Append this part to the surname and recalculate the suffix
			pn.Surname += "/" + p
			pn.Suffix = strings.TrimSpace(strings.Join(parts[j+1:], "/"))
		}

		pn.Full = pn.Given
		if len(pn.Surname) > 0 {
			if len(pn.Full) > 0 {
				pn.Full += " "
			}
			pn.Full += pn.Surname
		}
		if len(pn.Suffix) > 0 {
			if len(pn.Full) > 0 {
				pn.Full += " "
			}
			pn.Full += pn.Suffix
		}

		return pn

	}

	// Could not find a surname
	return ParsedName{
		Full:  strings.TrimRight(name, "/ "),
		Given: strings.TrimRight(name, "/ "),
	}
}
