package fonttablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
	//"github.com/OpenDiablo2/HellSpawner/hscommon"
	//"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	inputIntW = 25
)

/*const (
	removeItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png"
	upItemButtonPath     = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png"
	downItemButtonPath   = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png"
)*/

type textures struct {
	up, down, del *giu.Texture
}

type FontTableWidget struct {
	fontTable *d2font.Font
	id        string
	textures
}

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

	/*textureLoader.CreateTextureFromFileAsync(removeItemButtonPath, func(texture *giu.Texture) {
		result.textures.del = texture
	})

	textureLoader.CreateTextureFromFileAsync(upItemButtonPath, func(texture *giu.Texture) {
		result.textures.up = texture
	})

	textureLoader.CreateTextureFromFileAsync(downItemButtonPath, func(texture *giu.Texture) {
		fmt.Println(texture)
		result.textures.down = texture
	})*/

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
		giu.Label("Actions"),
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
			/*hsutil.MakeImageButton(
				"##"+p.id+"delete"+string(r),
				15, 15,
				p.textures.del,
				func() {},
			),*/
			giu.ImageButton(p.textures.del).Size(15, 15).OnClick(func() { p.deleteRow(r) }),
			giu.ImageButton(p.textures.up).Size(15, 15).OnClick(func() {}),
			giu.ImageButton(p.textures.down).Size(15, 15).OnClick(func() {}),
		),
		giu.Label(fmt.Sprintf("%d", p.fontTable.Glyphs[r].FrameIndex())),
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
