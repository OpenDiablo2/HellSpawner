package hsui

import (
	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"github.com/hajimehoshi/ebiten"
)

const closeButtonSymbol = "X"

func newTab(tabViewer *TabView, title string, closeable bool, content Widget) *Tab {
	tab := &Tab{
		info:      tabViewer.info,
		viewer:    tabViewer,
		title:     title,
		content:   content,
		container: CreateHBox(),
		enabled:   true,
		visible:   true,
		dirty:     true,
	}

	tab.container.SetChildSpacing(-2)

	tab.selectButton = CreateButton(tab.info, title, func() { tab.Select() })
	tab.closeButton = CreateButton(tab.info, closeButtonSymbol, func() { tab.Close() })

	tab.container.AddChild(tab.selectButton)
	tab.container.AddChild(tab.closeButton)

	tab.SetCloseable(closeable)

	return tab
}

type Tab struct {
	info    hsutil.InfoProvider
	viewer  *TabView
	content Widget
	title   string

	container    *HBox // expands, contains
	selectButton *Button
	closeButton  *Button

	reqWidth, reqHeight int

	enabled,
	visible,
	expanded,
	dirty,
	closeable,
	selected bool
}

func (t *Tab) Render(screen *ebiten.Image, x, y, width, height int) {
	/*
		tab looks like this:
		+------------------+---+
		|  <label>         | X |
		+------------------+---+
	*/
	closeWidth, closeHeight := height, height
	selectWidth, selectHeight := width-closeWidth, height
	closeX, closeY := x+selectWidth, y
	// label x,y is incoming x,y

	t.selectButton.Render(screen, x, y, selectWidth, selectHeight)
	t.closeButton.Render(screen, closeX, closeY, closeWidth, closeHeight)

	if !t.selected {
		return
	}

	rgba := t.info.GetAppConfig().Colors.TabSelected
	hsutil.DrawColoredRect(screen, x, y, width, height, rgba[0], rgba[1], rgba[2], rgba[3])
}

func (t *Tab) Update() bool {
	t.dirty = false

	if t.content == nil {
		return t.dirty
	}

	t.dirty = t.dirty || t.container.Update()
	t.dirty = t.dirty || t.content.Update()

	if !t.dirty {
		return t.dirty
	}

	t.Invalidate()
	t.dirty = true

	return t.dirty
}

func (t *Tab) GetRequestedSize() (int, int) {
	return t.reqWidth, t.reqHeight
}

func (t *Tab) Invalidate() {
	t.reqWidth, t.reqHeight = t.container.GetRequestedSize()

	if t.content != nil {
		t.content.Invalidate()
	}
}

func (t *Tab) SetCloseable(b bool) {
	if t.closeable == b {
		return
	}

	t.closeable = b
	t.Invalidate()
}

func (t *Tab) Close() {
	if !t.closeable || t.viewer == nil {
		return
	}

	t.enabled = false
	t.visible = false
	t.viewer.RemoveTab(t)
}

func (t *Tab) Select() {
	t.viewer.SelectTab(t)
	t.selected = true
}
