package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
)

const (
	inputIntW    = 25
	delW, delH   = 50, 30
	upW, upH     = 30, 30
	downW, downH = 40, 30
)

type textures struct {
	up, down, del *giu.Texture
}

// FontTableWidget represents a font table widget
type FontTableWidget struct {
	fontTable *d2font.Font
	id        string
	textures
}

// Create creates a new FontTable widget
func Create(
	up, down, del *giu.Texture,
	id string, fontTable *d2font.Font,
) *FontTableWidget {
	result := &FontTableWidget{
		fontTable: fontTable,
		id:        id,
		textures: textures{
			up, down, del,
		},
	}

	return result
}

// Build builds a widget
func (p *FontTableWidget) Build() {
	p.makeTableLayout().Build()
}

func (p *FontTableWidget) makeTableLayout() giu.Layout {
	rows := make(giu.Rows, 0)

	rows = append(rows, giu.Row(
		giu.Label("Delete"),
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
		giu.Child("##" + p.id + "tableArea").Border(false).Layout(giu.Layout{
			giu.FastTable("##" + p.id + "table").Border(true).Rows(rows),
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
		giu.Line(
			// in fact, it should be ImageButton, but this shit doesn't work
			// (imgui.PushID returns panic) :-/
			giu.Button("del##"+p.id+"deleteFrame"+string(r)).Size(delW, delH).OnClick(func() {
				p.deleteRow(r)
			}),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("%d", p.fontTable.Glyphs[r].FrameIndex())),
			giu.Button("up##"+p.id+"upItem"+string(r)).Size(upW, upH).OnClick(func() {
				p.itemUp(r)
			}),
			giu.Button("down##"+p.id+"upItem"+string(r)).Size(downW, downH).OnClick(func() {
				p.itemDown(r)
			}),
		),
		giu.Label(string(r)),
		giu.InputInt("##"+p.id+"width"+string(r), &width32).Size(inputIntW).OnChange(func() {
			_, h := p.fontTable.Glyphs[r].Size()
			p.fontTable.Glyphs[r].SetSize(int(width32), h)
		}),
	)

	return row
}

func (p *FontTableWidget) deleteRow(r rune) {
	delete(p.fontTable.Glyphs, r)
}

func (p *FontTableWidget) itemUp(r rune) {
	// currentFrame is frame index of 'r'
	currentFrame := p.fontTable.Glyphs[r].FrameIndex()

	// checks if above current index (r) is another one
	for cr, i := range p.fontTable.Glyphs {
		if i.FrameIndex() == currentFrame-1 {
			// if above current index ('r') is another one,
			// this above row gets down
			p.fontTable.Glyphs[cr].SetFrameIndex(
				p.fontTable.Glyphs[cr].FrameIndex() + 1,
			)

			break
		}
	}

	// current row's frame count gets up
	p.fontTable.Glyphs[r].SetFrameIndex(
		p.fontTable.Glyphs[r].FrameIndex() - 1,
	)
}

// itemDown does the sam as itemUp
func (p *FontTableWidget) itemDown(r rune) {
	currentFrame := p.fontTable.Glyphs[r].FrameIndex()

	for cr, i := range p.fontTable.Glyphs {
		if i.FrameIndex() == currentFrame+1 {
			p.fontTable.Glyphs[cr].SetFrameIndex(
				p.fontTable.Glyphs[cr].FrameIndex() - 1,
			)

			break
		}
	}

	p.fontTable.Glyphs[r].SetFrameIndex(
		p.fontTable.Glyphs[r].FrameIndex() + 1,
	)
}
