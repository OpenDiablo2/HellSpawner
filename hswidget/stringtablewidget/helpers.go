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

	str, found := p.dict[state.Key]
	if found {
		state.Value = str
	} else {
		state.Value = ""
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

func (p *widget) generateTableKeys() (keys []string) {
	state := p.getState()

	switch {
	case state.NumOnly:
		for _, key := range state.keys {
			if key[0] == '#' {
				keys = append(keys, key)
			} else {
				// labels are sorted, so no-name (starting from # are on top)
				break
			}
		}
	case state.Search != "":
		for _, key := range state.keys {
			s := strings.ToLower(state.Search)
			k := strings.ToLower(key)
			v := strings.ToLower(p.dict[key])

			switch {
			case strings.Contains(k, s),
				strings.Contains(v, s):
				keys = append(keys, key)
			}
		}
	default:
		keys = state.keys
	}

	return
}
