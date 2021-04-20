package hsenum

type Locale byte

const (
	LocaleEnglish Locale = iota
	LocaleChinaTraditional
	LocaleKorean
)

func (l Locale) String() string {
	lookup := map[Locale]string{
		LocaleEnglish:          "English",
		LocaleChinaTraditional: "China (Traditional)",
		LocaleKorean:           "Korean",
	}

	value, ok := lookup[l]
	if !ok {
		return "Unknown"
	}

	return value
}
