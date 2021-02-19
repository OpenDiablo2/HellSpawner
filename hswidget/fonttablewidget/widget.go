package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

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
				state.mode = fontTableWidgetEditRune
			}),
			giu.Label(string(r)),
		),
		giu.InputInt("##"+p.id+"width"+string(r), &width32).Size(inputIntW).OnChange(func() {
			_, h := p.fontTable.Glyphs[r].Size()
			p.fontTable.Glyphs[r].SetSize(int(width32), h)
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
			giu.Button("Save##"+p.id+"editRuneSave").Size(saveCancelW, saveCancelH).OnClick(func() {
				_, exist := p.fontTable.Glyphs[state.editRune.editedRune]
				if !exist {
					p.fontTable.Glyphs[state.editRune.editedRune] = p.fontTable.Glyphs[state.editRune.startRune]
					delete(p.fontTable.Glyphs, state.editRune.startRune)
				} else {
					dialog.Message("only one rune of one type is possible").Error()
				}

				state.mode = fontTableWidgetViewer
			}),
			giu.Button("Cancel##"+p.id+"editRuneSave").Size(saveCancelW, saveCancelH).OnClick(func() {
				state.mode = fontTableWidgetViewer
			}),
		),
	}
}

func (p *widget) makeAddItemLayout() giu.Layout {
	state := p.getState()

	firstFreeIndex := 0

	usedIndexes := make([]int, 0)

	for _, i := range p.fontTable.Glyphs {
		usedIndexes = append(usedIndexes, i.FrameIndex())
	}

	sort.Ints(usedIndexes)

	for index, used := range usedIndexes {
		if index != used {
			firstFreeIndex = index

			break
		}
	}

	r := string(state.addItem.newRune.editedRune)

	return giu.Layout{
		giu.Line(
			giu.Label(fmt.Sprintf("Frame index: %d", firstFreeIndex)),
		),
		giu.Line(
			giu.Label("Rune: "),
			giu.InputText("##"+p.id+"addItemRune", &r).Size(inputIntW).OnChange(func() {
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
			// this allows as to click save, only, when rune doesn't exist in map
			giu.Custom(func() {
				cancel := giu.Button("Cancel##"+p.id+"addItemCancel").Size(saveCancelW, saveCancelH).OnClick(func() {
					state.mode = fontTableWidgetViewer
				})

				_, exist := p.fontTable.Glyphs[state.addItem.newRune.editedRune]
				if !exist {
					giu.Line(
						giu.Button("Save##"+p.id+"addItemSave").Size(saveCancelW, saveCancelH).OnClick(func() {
							newGlyph := d2fontglyph.Create(
								firstFreeIndex,
								int(state.addItem.width),
								int(state.addItem.height),
							)

							p.fontTable.Glyphs[state.addItem.newRune.editedRune] = newGlyph
							state.mode = fontTableWidgetViewer
						}),
						cancel,
					).Build()
				} else {
					cancel.Build()
				}
			}),
		),
	}
}
