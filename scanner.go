/*
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package gedcom

import (
	"fmt"
	"io"
	"strconv"
)

// Line represents a single line from a GEDCOM file after tokenization.
// A GEDCOM line has the format: "level [xref] tag [value]"
// For example: "0 @I1@ INDI" or "1 NAME John /Smith/"
type Line struct {
	Level      int    // Hierarchy level (0 for top-level records)
	Tag        string // GEDCOM tag (e.g., "INDI", "NAME", "BIRT")
	Value      string // Optional value following the tag
	Xref       string // Optional cross-reference identifier (e.g., "I1" from "@I1@")
	LineNumber int    // Line number in the input file (1-indexed)
	Offset     int    // Character offset in the input file
}

// String returns the line in GEDCOM format.
func (l *Line) String() string {
	if l.Xref != "" {
		return fmt.Sprintf("%d @%s@ %s %s", l.Level, l.Xref, l.Tag, l.Value)
	}
	return fmt.Sprintf("%d %s %s", l.Level, l.Tag, l.Value)
}

// Scanner tokenizes GEDCOM input line by line. It is a low-level component
// used by [Decoder]. Most users should use [Decoder] directly instead.
//
// A Scanner reads from an [io.RuneScanner] and breaks the input into
// GEDCOM lines, parsing the level, optional cross-reference, tag, and
// optional value from each line.
type Scanner struct {
	r      io.RuneScanner
	err    error
	state  int
	line   int
	offset int
	level  int
	buf    []rune
	tag    string
	value  string
	xref   string
}

// NewScanner creates a new Scanner that reads from r.
// Use [Scanner.Next] to advance through the input and [Scanner.Line]
// to retrieve the current line after each successful call to Next.
func NewScanner(r io.RuneScanner) *Scanner {
	return &Scanner{
		r:     r,
		state: stateBegin,
		buf:   make([]rune, 0, 4),
	}
}

const (
	stateBegin = iota
	stateLevel
	stateSeekTagOrXref
	stateSeekTag
	stateTag
	stateXref
	stateSeekValue
	stateValue
	stateEnd
	stateError
)

// Next advances the scanner to the next line. It returns false if there are no more lines
// or if an error is encountered. The caller should check the Err method whenever this
// method returns false.
func (s *Scanner) Next() bool {
	s.state = stateBegin
	s.level = 0
	s.buf = s.buf[:0]
	s.xref = ""
	s.tag = ""
	s.value = ""
	s.offset = 0
	s.line++

	for {
		c, n, err := s.r.ReadRune()
		if err != nil {
			if err != io.EOF {
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("read: %w", err),
				}
			}

			if s.state != stateEnd && s.state != stateBegin {
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        io.ErrUnexpectedEOF,
				}
			}

			return false
		}
		s.offset += n

		switch s.state {
		case stateBegin:
			switch {
			case c >= '0' && c <= '9':
				s.buf = append(s.buf, c)
				s.state = stateLevel
			case isSpace(c):
				continue
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("found non-whitespace %q (%#[1]x) before level", c),
				}
				return false
			}
		case stateLevel:
			switch {
			case c >= '0' && c <= '9':
				s.buf = append(s.buf, c)
				continue
			case c == ' ':
				parsedLevel, perr := strconv.ParseInt(string(s.buf), 10, 64)
				if perr != nil {
					s.err = &ScanErr{
						LineNumber: s.line,
						Offset:     s.offset,
						Err:        fmt.Errorf("parse level: %w", perr),
					}
					return false
				}
				s.level = int(parsedLevel)
				s.buf = s.buf[:0]
				s.state = stateSeekTagOrXref
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("level contained non-numerics"),
				}
				return false
			}

		case stateSeekTag:
			switch {
			case isAlphaNumeric(c):
				s.buf = append(s.buf, c)
				s.state = stateTag
			case c == ' ':
				continue
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("tag contained non-alphanumeric (%#x)", c),
				}
				return false
			}
		case stateSeekTagOrXref:
			switch {
			case isAlphaNumeric(c):
				s.buf = append(s.buf, c)
				s.state = stateTag
			case c == '@':
				s.state = stateXref
			case c == ' ':
				continue
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("tag or xref contained non-alphanumeric (%#x)", c),
				}
				return false
			}

		case stateTag:
			switch {
			case isAlphaNumeric(c):
				s.buf = append(s.buf, c)
				continue
			case c == '\n' || c == '\r':
				s.swallowCr(c)
				s.tag = string(s.buf)
				s.buf = s.buf[:0]
				s.state = stateEnd
				return true
			case c == ' ':
				s.tag = string(s.buf)
				s.buf = s.buf[:0]
				s.state = stateSeekValue
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("tag contained non-alphanumeric (%#x)", c),
				}
				return false
			}

		case stateXref:
			switch {
			case isAlphaNumeric(c):
				s.buf = append(s.buf, c)
				continue
			case c == '@':
				continue
			case c == ' ':
				s.xref = string(s.buf)
				s.buf = s.buf[:0]
				s.state = stateSeekTag
			default:
				s.state = stateError
				s.err = &ScanErr{
					LineNumber: s.line,
					Offset:     s.offset,
					Err:        fmt.Errorf("xref contained non-alphanumeric (%#x)", c),
				}
				return false
			}
		case stateSeekValue:
			switch {
			case c == '\n' || c == '\r':
				s.swallowCr(c)
				s.state = stateEnd
				return true
			case c == ' ':
				continue
			default:
				s.buf = append(s.buf, c)
				s.state = stateValue
			}

		case stateValue:
			switch {
			case c == '\n' || c == '\r':
				s.swallowCr(c)

				// Check to see if there is a malformed NOTE that contains an embedded newline
				// For example, Ancestry GEDCOM exports that include source "London, England, Church of England Births and Baptisms, 1813-1917"
				// have the following NOTE tag split over two lines (yet the CONC tag is correctly formatted!)
				//
				//   1 NOTE Board of Guardian Records and Church of England Parish Registers. London Metropolitan Archives, London.
				//   <p>Images produced by permission of the City of London Corporation. The City of London gives n

				if s.tag == "NOTE" {
					next, _, err := s.r.ReadRune()
					s.r.UnreadRune()
					if err == nil {
						if !isNumeric(next) {
							// Looks like it might be a malformed note, so continue parsing
							s.buf = append(s.buf, '\n')
							continue
						}
					}
				}

				s.value = string(s.buf)
				s.buf = s.buf[:0]
				s.state = stateEnd
				return true
			default:
				s.buf = append(s.buf, c)
				continue
			}
		}
	}
}

// swallowCr skips a carriage return if it is followed by a newline
func (s *Scanner) swallowCr(c rune) {
	if c == '\r' {
		next, _, _ := s.r.ReadRune()
		if next == '\n' {
			s.offset++
		} else {
			s.r.UnreadRune()
		}
	}
}

// Line returns the most recent line tokenized by a call to Next.
func (s *Scanner) Line() Line {
	return Line{
		Level:      s.level,
		Tag:        s.tag,
		Value:      s.value,
		Xref:       s.xref,
		LineNumber: s.line,
		Offset:     s.offset,
	}
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if s.err == nil {
		return nil
	}
	return s.err
}

// ScanErr represents a scanning error with location information.
// It wraps the underlying error and includes the line number and
// character offset where the error occurred.
type ScanErr struct {
	Err        error // The underlying error
	LineNumber int   // Line number where the error occurred (1-indexed)
	Offset     int   // Character offset within the line
}

func (e *ScanErr) Error() string {
	return fmt.Sprintf("scan error (line:%d, position:%d): %v", e.LineNumber, e.Offset, e.Err)
}

func (e *ScanErr) Unwrap() error {
	return e.Err
}

// isSpace reports whether c is a whitespace rune. The Unicode BOM (U+FEFF) is
// included so that GEDCOM 7.0 files, which require a leading UTF-8 BOM, are
// handled transparently by the scanner.
func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\uFEFF'
}

func isAlphaNumeric(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isNumeric(c rune) bool {
	return (c >= '0' && c <= '9')
}
