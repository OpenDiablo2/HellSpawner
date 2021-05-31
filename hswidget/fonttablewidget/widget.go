package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font/d2fontglyph"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	inputIntW                = 30
	delSize                  = 20
	addW, addH               = 400, 30
	editRuneW, editRuneH     = 50, 30
	saveCancelW, saveCancelH = 80, 30
)

type widget struct {
	fontTable     *d2font.Font
	id            string
	textureLoader hscommon.TextureLoader
}

// Create creates a new FontTable widget
func Create(
	state []byte,
	tl hscommon.TextureLoader,
	id string, fontTable *d2font.Font,
) giu.Widget {
	result := &widget{
		fontTable:     fontTable,
		id:            id,
		textureLoader: tl,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
	}

	return result
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case modeViewer:
		p.makeTableLayout().Build()
	case modeEditRune:
		p.makeEditRuneLayout().Build()
	case modeAddItem:
		p.makeAddItemLayout().Build()
	}
}

func (p *widget) makeTableLayout() giu.Layout {
	state := p.getState()

	rows := make([]*giu.TableRowWidget, 0)

	rows = append(rows, giu.TableRow(
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
		giu.Button("Add new glyph...##"+p.id+"addItem").Size(addW, addH).OnClick(func() {
			state.mode = modeAddItem
		}),
		giu.Separator(),
		giu.Child("##" + p.id + "tableArea").Border(false).Layout(giu.Layout{
			giu.Table("##" + p.id + "table").FastMode(true).Rows(rows...),
		}),
	}
}

func (p *widget) makeGlyphLayout(r rune) *giu.TableRowWidget {
	state := p.getState()

	if p.fontTable.Glyphs[r] == nil {
		return &giu.TableRowWidget{}
	}

	w := p.fontTable.Glyphs[r].Width()
	width32 := int32(w)

	h := p.fontTable.Glyphs[r].Height()
	height32 := int32(h)

	row := giu.TableRow(
		hswidget.MakeImageButton("##"+p.id+"deleteFrame"+string(r),
			delSize, delSize,
			state.deleteButtonTexture,
			func() { p.deleteRow(r) },
		),
		giu.Row(
			giu.Label(fmt.Sprintf("%d", p.fontTable.Glyphs[r].FrameIndex())),
			giu.ArrowButton("##"+p.id+"upItem"+string(r), giu.DirectionUp).OnClick(func() {
				p.itemUp(r)
			}),
			giu.ArrowButton("##"+p.id+"downItem"+string(r), giu.DirectionDown).OnClick(func() {
				p.itemDown(r)
			}),
		),
		giu.Row(
			giu.Button("edit##"+p.id+"editRune"+string(r)).Size(editRuneW, editRuneH).OnClick(func() {
				state.editRuneState.runeBefore = r
				state.editRuneState.editedRune = r
				state.mode = modeEditRune
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

	r := string(state.editRuneState.editedRune)

	return giu.Layout{
		giu.Label("Edit rune:"),
		giu.Row(
			giu.Label("Rune: "),
			giu.InputText("##"+p.id+"editRuneRune", &r).Size(inputIntW).OnChange(func() {
				if len(r) > 0 {
					state.editRuneState.editedRune = int32(r[0])
				}
			}),
		),
		giu.Row(
			giu.Label("Int: "),
			giu.InputInt("##"+p.id+"editRuneInt", &state.editRuneState.editedRune).Size(inputIntW),
		),
		giu.Separator(),
		giu.Row(
			p.makeSaveCancelRow(func() {
				p.fontTable.Glyphs[state.editRuneState.editedRune] = p.fontTable.Glyphs[state.editRuneState.runeBefore]
				p.deleteRow(state.editRuneState.runeBefore)

				state.mode = modeViewer
			}, state.editRuneState.editedRune),
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

	r := string(state.addItemState.newRune)

	return giu.Layout{
		giu.Row(
			giu.Label(fmt.Sprintf("Frame index: %d", firstFreeIndex)),
		),
		giu.Row(
			giu.Label("Rune: "),
			// if user put here more then one letter,
			// second and further letters will be skipped
			giu.InputText("##"+p.id+"addItemRune", &r).Size(inputIntW).OnChange(func() {
				if r == "" {
					state.addItemState.newRune = 0

					return
				}

				state.addItemState.newRune = int32(r[0])
			}),
		),
		giu.Row(
			giu.Label("Int: "),
			giu.InputInt("##"+p.id+"addItemRuneInt", &state.addItemState.newRune).Size(inputIntW),
		),
		giu.Row(
			giu.Label("Width: "),
			giu.InputInt("##"+p.id+"addItemWidth", &state.addItemState.width).Size(inputIntW),
		),
		giu.Row(
			giu.Label("Height: "),
			giu.InputInt("##"+p.id+"addItemHeight", &state.addItemState.height).Size(inputIntW),
		),
		giu.Separator(),
		giu.Row(
			p.makeSaveCancelRow(func() {
				p.addItem(firstFreeIndex)
			}, state.addItemState.newRune),
		),
	}
}

func (p *widget) addItem(idx int) {
	state := p.getState()

	newGlyph := d2fontglyph.Create(
		idx,
		int(state.addItemState.width),
		int(state.addItemState.height),
	)

	p.fontTable.Glyphs[state.addItemState.newRune] = newGlyph
	state.mode = modeViewer
}

// makeSaveCancelRow creates  line of action buttons for an editor
// if given rune already exists in glyph's table, save button isn't
// created
func (p *widget) makeSaveCancelRow(saveCB func(), r rune) giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Custom(func() {
			cancel := giu.Button("Cancel##"+p.id+"addItemCancel").Size(saveCancelW, saveCancelH).OnClick(func() {
				state.mode = modeViewer
			})

			_, exist := p.fontTable.Glyphs[r]
			if exist {
				cancel.Build()

				return
			}

			giu.Row(
				giu.Button("Save##"+p.id+"addItemSave").Size(saveCancelW, saveCancelH).OnClick(func() {
					saveCB()
				}),
				cancel,
			).Build()
		}),
	}
}
