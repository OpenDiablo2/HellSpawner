package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font/d2fontglyph"
)

const (
	inputIntW                = 30
	delW, delH               = 50, 30
	upW, upH                 = 30, 30
	downW, downH             = 40, 30
	addW, addH               = 400, 30
	editRuneW, editRuneH     = 50, 30
	saveCancelW, saveCancelH = 80, 30
)

type textures struct {
	up, down, del *giu.Texture
}

// widget represents a font table widget
type widget struct {
	fontTable *d2font.Font
	id        string
	textures
}

// Create creates a new FontTable widget
func Create(
	up, down, del *giu.Texture,
	id string, fontTable *d2font.Font,
) giu.Widget {
	result := &widget{
		fontTable: fontTable,
		id:        id,
		textures: textures{
			up, down, del,
		},
	}

	return result
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case fontTableWidgetViewer:
		p.makeTableLayout().Build()
	case fontTableWidgetEditRune:
		p.makeEditRuneLayout().Build()
	case fontTableWidgetAddItem:
		p.makeAddItemLayout().Build()
	}
}

func (p *widget) makeTableLayout() giu.Layout {
	state := p.getState()

	rows := make(giu.Rows, 0)

	rows = append(rows, giu.Row(
		giu.Label("Delete"),
		giu.Label("Index"),
		giu.Label("Character"),
		giu.Label("Width (px)"),
		giu.Label("Height (px)"),
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
			giu.Separator(),
			giu.Button("add##"+p.id+"addItem").Size(addW, addH).OnClick(func() {
				state.mode = fontTableWidgetAddItem
			}),
		}),
	}
}

func (p *widget) makeGlyphLayout(r rune) *giu.RowWidget {
	state := p.getState()

	if p.fontTable.Glyphs[r] == nil {
		return &giu.RowWidget{}
	}

	w := p.fontTable.Glyphs[r].Width()
	width32 := int32(w)

	h := p.fontTable.Glyphs[r].Height()
	height32 := int32(h)

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
		giu.Line(
			giu.Button("edit##"+p.id+"editRune"+string(r)).Size(editRuneW, editRuneH).OnClick(func() {
				state.editRune.startRune = r
				state.editRune.editedRune = r
				state.mode = fontTableWidgetEditRune
			}),
			giu.Label(string(r)),
		),
		giu.InputInt("##"+p.id+"width"+string(r), &width32).Size(inputIntW).OnChange(func() {
			h := p.fontTable.Glyphs[r].Height()
			p.fontTable.Glyphs[r].SetSize(int(width32), h)
		}),
		giu.InputInt("##"+p.id+"height"+string(r), &height32).Size(inputIntW).OnChange(func() {
			w := p.fontTable.Glyphs[r].Width()
			p.fontTable.Glyphs[r].SetSize(w, int(height32))
		}),
	)

	return row
}

func (p *widget) deleteRow(r rune) {
	delete(p.fontTable.Glyphs, r)
}

func (p *widget) itemUp(r rune) {
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
func (p *widget) itemDown(r rune) {
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

func (p *widget) makeEditRuneLayout() giu.Layout {
	state := p.getState()

	r := string(state.editRune.editedRune)

	return giu.Layout{
		giu.Label("Edit rune:"),
		giu.Line(
			giu.Label("Rune: "),
			giu.InputText("##"+p.id+"editRuneRune", &r).Size(inputIntW).OnChange(func() {
				state.editRune.editedRune = int32(r[0])
			}),
		),
		giu.Line(
			giu.Label("Int: "),
			giu.InputInt("##"+p.id+"editRuneInt", &state.editRune.editedRune).Size(inputIntW),
		),
		giu.Separator(),
		giu.Line(
			p.makeSaveCancelLine(func() {
				p.fontTable.Glyphs[state.editRune.editedRune] = p.fontTable.Glyphs[state.editRune.startRune]
				p.deleteRow(state.editRune.startRune)

				state.mode = fontTableWidgetViewer
			}, state.editRune.editedRune),
		),
	}
}

// adds new item on the first free position (frame index)
func (p *widget) makeAddItemLayout() giu.Layout {
	state := p.getState()

	// first free index determinates a first frame index, when
	// we can place our new item
	firstFreeIndex := -1

	// frame indexes, which are already taken
	usedIndexes := make([]int, 0)

	for _, i := range p.fontTable.Glyphs {
		usedIndexes = append(usedIndexes, i.FrameIndex())
	}

	sort.Ints(usedIndexes)

	for index, used := range usedIndexes {
		// simple condition:
		// if index != used, it means that for example
		// used indexes are [0, 2, 3], so
		// index in list is 1, but usedIndexes[1] is 2, so
		// frame 1 is free
		if index != used {
			firstFreeIndex = index

			break
		}
	}

	// if no free indexes found, then set to next index
	if firstFreeIndex == -1 {
		firstFreeIndex = len(usedIndexes)
	}

	r := string(state.addItem.newRune.editedRune)

	return giu.Layout{
		giu.Line(
			giu.Label(fmt.Sprintf("Frame index: %d", firstFreeIndex)),
		),
		giu.Line(
			giu.Label("Rune: "),
			// if user put here more then one letter,
			// second and further letters will be skipped
			giu.InputText("##"+p.id+"addItemRune", &r).Size(inputIntW).OnChange(func() {
				if r == "" {
					state.addItem.newRune.editedRune = 0

					return
				}

				state.addItem.newRune.editedRune = int32(r[0])
			}),
		),
		giu.Line(
			giu.Label("Int: "),
			giu.InputInt("##"+p.id+"addItemRuneInt", &state.addItem.newRune.editedRune).Size(inputIntW),
		),
		giu.Line(
			giu.Label("Width: "),
			giu.InputInt("##"+p.id+"addItemWidth", &state.addItem.width).Size(inputIntW),
		),
		giu.Line(
			giu.Label("Height: "),
			giu.InputInt("##"+p.id+"addItemHeight", &state.addItem.height).Size(inputIntW),
		),
		giu.Separator(),
		giu.Line(
			p.makeSaveCancelLine(func() {
				p.addItem(firstFreeIndex)
			}, state.addItem.newRune.editedRune),
		),
	}
}

func (p *widget) addItem(idx int) {
	state := p.getState()

	newGlyph := d2fontglyph.Create(
		idx,
		int(state.addItem.width),
		int(state.addItem.height),
	)

	p.fontTable.Glyphs[state.addItem.newRune.editedRune] = newGlyph
	state.mode = fontTableWidgetViewer
}

// this giant custom function allows us to
// check if letter entered by user already exists in map
// end depending on it, build save and cancel buttons
// or cancel only
func (p *widget) makeSaveCancelLine(saveCB func(), r rune) giu.Layout {
	state := p.getState()
	return giu.Layout{
		giu.Custom(func() {
			cancel := giu.Button("Cancel##"+p.id+"addItemCancel").Size(saveCancelW, saveCancelH).OnClick(func() {
				state.mode = fontTableWidgetViewer
			})

			_, exist := p.fontTable.Glyphs[r]
			if exist {
				cancel.Build()

				return
			}

			giu.Line(
				giu.Button("Save##"+p.id+"addItemSave").Size(saveCancelW, saveCancelH).OnClick(func() {
					saveCB()
				}),
				cancel,
			).Build()
		}),
	}
}
