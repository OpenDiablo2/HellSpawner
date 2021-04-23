package hsenum

// PolishSpecialCharacters are the characters which should be added
// to imgui's glyph range. They aren't by default
const PolishSpecialCharacters = "ĄąĘęŁłŃńÓóŚśŹźŻż"

// Locale represents an app locale
type Locale byte

// this is a list of locales supported by Diablo II
const (
	LocaleEnglish Locale = iota
	LocaleGerman
	LocaleFrench
	LocaleKorean
	LocaleChinaTraditional
	LocaleSpanish
	LocaleItalien
	LocalePolish
)

func (l Locale) String() string {
	lookup := map[Locale]string{
		LocaleEnglish:          "English",
		LocaleGerman:           "German",
		LocaleFrench:           "French",
		LocaleKorean:           "Korean",
		LocaleChinaTraditional: "Chinese (Traditional)",
		LocaleSpanish:          "Spanish",
		LocaleItalien:          "Italien",
		LocalePolish:           "Polish",
	}

	value, ok := lookup[l]
	if !ok {
		return "Unknown"
	}

	return value
}
