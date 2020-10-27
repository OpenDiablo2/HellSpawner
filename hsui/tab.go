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
		tabSelect: CreateHBox(),
		tabClose:  CreateHBox(),
		enabled:   true,
		visible:   true,
		dirty:     true,
	}

	tab.container.SetChildSpacing(-2)

	tab.tabSelect.AddChild(CreateButton(tab.info, title, func() { tab.Select() }))
	tab.tabSelect.SetExpandChild(true)

	tab.closeButton = CreateButton(tab.info, closeButtonSymbol, func() { tab.Close() })
	tab.tabClose.SetExpandChild(false)

	tab.container.AddChild(tab.tabSelect)
	tab.container.AddChild(tab.tabClose)

	tab.container.SetExpandChild(false)

	tab.SetCloseable(closeable)
	tab.Invalidate()

	return tab
}

type Tab struct {
	info    hsutil.InfoProvider
	viewer  *TabView
	content Widget
	title   string

	container *HBox // expands, contains
	tabSelect *HBox // expands, contains a button
	tabClose  *HBox // doesn't expand, contains a button

	closeButton *Button

	reqWidth, reqHeight int

	enabled,
	visible,
	expanded,
	dirty,
	closeable bool
}

func (t *Tab) Render(screen *ebiten.Image, x, y, width, height int) {
	t.container.Render(screen, x, y, width, height)
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

	if t.closeable {
		hbox := CreateHBox()
		hbox.AddChild(t.closeButton)
		t.tabClose = hbox
	} else {
		hbox := CreateHBox()
		t.tabClose = hbox
	}

	t.container.children[1] = t.tabClose
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
}
