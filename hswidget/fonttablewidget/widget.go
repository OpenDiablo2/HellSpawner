package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
)

const (
	inputIntW = 25
)

type FontTableWidget struct {
	fontTable *d2font.Font
	id        string
}

func Create(id string, fontTable *d2font.Font) *FontTableWidget {
	result := &FontTableWidget{
		fontTable: fontTable,
		id:        id,
	}

	return result
}

func (p *FontTableWidget) Build() {
	giu.Layout(giu.Layout{
		p.makeTableLayout(),
	}).Build()
}

func (p *FontTableWidget) makeTableLayout() giu.Layout {
	rows := make(giu.Rows, 0)

	rows = append(rows, giu.Row(
		giu.Label("Index"),
		giu.Label("Character"),
		giu.Label("Width (px)"),
	))

	// we need to get keys from map[rune]*d2font.fontGlyph
	// and then sort them
	chars := make([]rune, len(p.fontTable.Glyphs))

	// reading runes from map
	idx := 0

	for r := range p.fontTable.Glyphs {
		chars[idx] = r
		idx++
	}

	// sorting runes
	sort.Slice(chars, func(i, j int) bool {
		return p.fontTable.Glyphs[chars[i]].FrameIndex() < p.fontTable.Glyphs[chars[j]].FrameIndex()
	})

	for _, idx := range chars {
		rows = append(rows, p.makeGlyphLayout(idx))
	}

	return giu.Layout{
		giu.Child("").Border(false).Layout(giu.Layout{
			giu.FastTable("").Border(true).Rows(rows),
		}),
	}
}

func (p *FontTableWidget) makeGlyphLayout(r rune) *giu.RowWidget {
	if p.fontTable.Glyphs[r] == nil {
		return &giu.RowWidget{}
	}

	w, _ := p.fontTable.Glyphs[r].Size()
	width32 := int32(w)
	row := giu.Row(
		giu.Label(fmt.Sprintf("%d", p.fontTable.Glyphs[r].FrameIndex())),
		giu.Label(string(r)),
		giu.InputInt("##"+p.id+"width"+string(r), &width32).Size(inputIntW).OnChange(func() {
			_, h := p.fontTable.Glyphs[r].Size()
			p.fontTable.Glyphs[r].SetSize(int(width32), h)
		}),
	)

	return row
}
