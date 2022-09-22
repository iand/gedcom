package gedcom

import (
	"strings"
)

type ParsedName struct {
	Full    string
	Given   string
	Surname string
	Suffix  string
}

// SplitPersonalName parses a name in the format "First Name /Surname/ suffix" into its components.
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
