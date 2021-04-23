package hsenum

type Locale byte

const (
	LocaleEnglish Locale = iota
	LocaleChinaTraditional
	LocaleKorean
	LocalePolish
)

func (l Locale) String() string {
	lookup := map[Locale]string{
		LocaleEnglish:          "English",
		LocaleChinaTraditional: "China (Traditional)",
		LocaleKorean:           "Korean",
		LocalePolish:           "Polish",
	}

	value, ok := lookup[l]
	if !ok {
		return "Unknown"
	}

	return value
}
