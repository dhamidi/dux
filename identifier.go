package dux

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// IdentifierStyle described a way of encoding an identifier as a string.
type IdentifierStyle interface {
	Parse(string) []string
	Original([]string) string
	Upper([]string) string
	Lower([]string) string
	Title([]string) string
}

// SeparatedIdentifier describes identifiers using a separator to
// distinguish between words.
type SeparatedIdentifier struct {
	Separator string
}

// Parse parses the identifier by splitting it according to the separator
func (s *SeparatedIdentifier) Parse(identifier string) []string {
	return strings.Split(identifier, s.Separator)
}

// Original renders the identifier as a string using the original casing
func (s *SeparatedIdentifier) Original(constituents []string) string {
	return strings.Join(constituents, s.Separator)
}

// Upper renders the identifier by converting all consituents to upper case before joining them.
func (s *SeparatedIdentifier) Upper(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		fmt.Fprintf(out, "%s", strings.ToUpper(c))
		if i > 0 && i < len(constituents)-1 {
			fmt.Fprintf(out, "%s", s.Separator)
		}
	}

	return out.String()
}

// Lower renders the identifier by converting all constituents to lower case before joining them.
func (s *SeparatedIdentifier) Lower(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		fmt.Fprintf(out, "%s", strings.ToLower(c))
		if i > 0 && i < len(constituents)-1 {
			fmt.Fprintf(out, "%s", s.Separator)
		}
	}

	return out.String()
}

// Title is an alias for Original.
func (s *SeparatedIdentifier) Title(constituents []string) string { return s.Original(constituents) }

// CasedIdentifier described identifier that distinguish constituents using letter casing.
type CasedIdentifier struct{}

// Parse parses the identifier by splitting on upper-case letters
func (s *CasedIdentifier) Parse(identifier string) []string {
	constituent := []rune{}
	result := []string{}
	for _, r := range identifier {
		if unicode.IsUpper(r) {
			result = append(result, string(constituent))
			constituent = []rune{r}
		} else {
			constituent = append(constituent, r)
		}
	}
	if len(constituent) > 0 {
		result = append(result, string(constituent))
	}

	return result
}

// Original renders the identifier preserving original casing.
func (s *CasedIdentifier) Original(constituents []string) string {
	return strings.Join(constituents, "")
}

// Upper renders the identifier by converting the first letter of each constituent to upper case.
func (s *CasedIdentifier) Upper(constituents []string) string {
	out := bytes.NewBufferString("")
	for _, c := range constituents {
		if len(c) == 0 {
			continue
		}
		part := []rune(c)
		part[0] = unicode.ToUpper(part[0])
		fmt.Fprintf(out, "%s", string(part))
	}
	return out.String()
}

// Lower renders the identifier by converting the first consituent to lower case and the rest to title case.
func (s *CasedIdentifier) Lower(constituents []string) string {
	out := bytes.NewBufferString("")
	for i, c := range constituents {
		if len(c) == 0 {
			continue
		}
		part := []rune(c)
		if i == 0 {
			part[0] = unicode.ToLower(part[0])
		} else {
			part[0] = unicode.ToUpper(part[0])
		}
		fmt.Fprintf(out, "%s", string(part))
	}
	return out.String()
}

// Title is an alias for upper.
func (s *CasedIdentifier) Title(constituents []string) string {
	return s.Upper(constituents)
}

// Identifier reprents a programming language identifier that can be
// expressed in various casing styles.
type Identifier struct {
	Constituents []string
	Style        IdentifierStyle
}

// String renders the identifier in the style that was detected during creation of the identifier.
func (i *Identifier) String() string {
	if i.Style == nil {
		return ""
	}
	return i.Style.Original(i.Constituents)
}

// Set implements flag.Value by parsing the identifier
func (i *Identifier) Set(s string) error {
	newIdentifier := ParseIdentifier(s)
	*i = *newIdentifier
	return nil
}

// Get implements flag.Value by returning the identifier itself.
func (i *Identifier) Get() interface{} {
	return i
}

var (
	// SnakeCaseStyle is an identifier that separates words using underscores.
	SnakeCaseStyle = &SeparatedIdentifier{Separator: "_"}

	// LispCaseStyle is an identifier that separates words using hyphens.
	LispCaseStyle = &SeparatedIdentifier{Separator: "-"}

	// CamelCasedStyle is an identifier that separates words using letter casing.
	CamelCasedStyle = new(CasedIdentifier)
)

// ParseIdentifier analyzes a string as an identifier.
func ParseIdentifier(identifier string) *Identifier {
	result := &Identifier{
		Constituents: []string{},
		Style:        CamelCasedStyle,
	}

runes:
	for _, r := range identifier {
		switch r {
		case '-':
			result.Style = LispCaseStyle
			break runes
		case '_':
			result.Style = SnakeCaseStyle
			break runes
		}
	}

	result.Constituents = result.Style.Parse(identifier)
	return result
}

// Upper returns the identifier in upper case
func (i *Identifier) Upper() string {
	return i.Style.Upper(i.Constituents)
}

// Lower returns the identifier in title case
func (i *Identifier) Lower() string {
	return i.Style.Lower(i.Constituents)
}

// Title returns the identifier in title case
func (i *Identifier) Title() string {
	return i.Style.Title(i.Constituents)
}

// ToSnake converts the identifier into snake case
func (i *Identifier) ToSnake() *Identifier {
	return &Identifier{
		Constituents: i.Constituents,
		Style:        SnakeCaseStyle,
	}
}

// ToSnake converts the identifier into lisp case
func (i *Identifier) ToLisp() *Identifier {
	return &Identifier{
		Constituents: i.Constituents,
		Style:        LispCaseStyle,
	}
}
