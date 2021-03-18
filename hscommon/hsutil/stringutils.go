package hsutil

import (
	"strings"
	"unicode"
)

// SplitIntoLinesWithMaxWidth splits the given string into lines considering the given maxChars
func SplitIntoLinesWithMaxWidth(fullSentence string, maxChars int) []string {
	lines := make([]string, 0)
	line := ""
	totalLength := 0
	words := strings.Split(fullSentence, " ")

	if len(words[0]) > maxChars {
		// mostly happened within CJK characters (no whitespace)
		return splitCjkIntoChunks(fullSentence, maxChars)
	}

	for _, word := range words {
		totalLength += 1 + len(word)
		if totalLength > maxChars {
			totalLength = len(word)

			lines = append(lines, line)
			line = ""
		} else {
			line += " "
		}

		line += word
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}

	return lines
}

func splitCjkIntoChunks(str string, chars int) []string {
	chunks := make([]string, 0)
	i, count := 0, 0

	for j, ch := range str {
		if ch < unicode.MaxLatin1 {
			count++
		} else {
			// assume we're truncating CJK characters
			count += 2
		}

		if count >= chars {
			chunks = append(chunks, str[i:j])
			i, count = j, 0
		}
	}

	return append(chunks, str[i:])
}
