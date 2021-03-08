package stringtablewidget

import (
	"sort"
	"strconv"
	"strings"
)

func (p *widget) formatKey(s *string) {
	*s = strings.ReplaceAll(*s, " ", "_")
}

func (p *widget) updateValueText() {
	state := p.getState()

	str, found := p.dict[state.key]
	if found {
		state.value = str
	} else {
		state.value = ""
	}
}

func (p *widget) calculateFirstFreeNoName() (firstFreeNoName int) {
	state := p.getState()

	ints := make([]int, 0)

	for _, key := range state.keys {
		if key[0] == '#' {
			idx, err := strconv.Atoi(key[1:])
			if err != nil {
				continue
			}

			ints = append(ints, idx)
		}
	}

	sort.Ints(ints)

	for n, i := range ints {
		if n != i {
			firstFreeNoName = n
			break
		}
	}

	return
}
