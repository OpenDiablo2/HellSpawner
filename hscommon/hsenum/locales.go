package hsenum

// PolishSpecialCharacters are the characters which should be added
// to imgui's glyph range. They aren't by default
const PolishSpecialCharacters = "ĄąĘęŁłŃńÓóŚśŹźŻż"

// Locale represents an app locale
//go:generate stringer -linecomment -type Locale -output locales_string.go
type Locale byte

// this is a list of locales supported by Diablo II
const (
	LocaleEnglish            Locale = iota // English
	LocaleGerman                           // German
	LocaleFrench                           // French
	LocaleKorean                           // Korean
	LocaleChineseTraditional               // Chinese (Traditional)
	LocaleSpanish                          // Spanish
	LocaleItalien                          // Italien
	LocalePolish                           // Polish
)
