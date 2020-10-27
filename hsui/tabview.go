package hsui

import (
	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"github.com/hajimehoshi/ebiten"
)

const (
	tabNotFound = -1
)

func CreateTabView(info hsutil.InfoProvider, w, h int) *TabView {
	tv := &TabView{
		info:      info,
		tabs:      make([]*Tab, 0),
		tabBox:    CreateHBox(),
		container: CreateVBox(),
		pager:     CreatePager(0, 0, nil),
		reqWidth:  w,
		reqHeight: h,
		enabled:   true,
		visible:   true,
		expanded:  true,
	}

	tv.container.AddChild(tv.tabBox)
	tv.container.AddChild(tv.pager)

	tv.container.SetExpandChild(true)

	tv.Invalidate()

	return tv
}

type TabView struct {
	info                       hsutil.InfoProvider
	tabs                       []*Tab
	container                  *VBox
	tabBox                     *HBox
	pager                      *Pager
	reqWidth, reqHeight        int
	enabled, visible, expanded bool
}

func (t *TabView) Render(screen *ebiten.Image, x, y, width, height int) {
	_, tabBoxReqH := t.tabBox.GetRequestedSize()

	t.tabBox.Render(screen, x, y, width, tabBoxReqH)
	t.pager.Render(screen, x, y+tabBoxReqH, width, height-tabBoxReqH)
}

func (t *TabView) Update() (dirty bool) {
	dirty = dirty || t.container.Update()
	dirty = dirty || t.pager.Update()
	dirty = dirty || t.tabBox.Update()

	for idx := range t.tabs {
		dirty = dirty || t.tabs[idx].Update()
	}

	if dirty {
		t.Invalidate()
	}

	return dirty
}

func (t *TabView) GetRequestedSize() (width, height int) {
	return t.reqWidth, t.reqHeight
}

func (t *TabView) Invalidate() {
	t.container.Invalidate()
	t.pager.Invalidate()
	t.tabBox.Invalidate()

	t.reqWidth, t.reqHeight = 100, 100

	for idx := range t.tabs {
		t.tabs[idx].Invalidate()
	}
}

func (t *TabView) AddTab(title string, content Widget, closeable bool) {
	tab := newTab(t, title, closeable, content)

	t.tabs = append(t.tabs, tab)

	t.tabBox.AddChild(tab)
	t.pager.AddChild(tab.content)
	t.Invalidate()
}

func (t *TabView) RemoveTab(tab *Tab) {
	idx := t.getTabIndex(tab)

	if idx == tabNotFound {
		return
	}

	t.SelectTabByIndex(idx - 1)

	t.tabs = append(t.tabs[:idx], t.tabs[idx+1:]...)
	t.tabBox = CreateHBox()
	t.pager = CreatePager(0, 0, nil)

	for idx := range t.tabs {
		t.tabBox.AddChild(t.tabs[idx])
		t.pager.AddChild(t.tabs[idx].content)
	}

	t.Invalidate()
}

func (t *TabView) SelectTab(tab *Tab) {
	t.SelectTabByIndex(t.getTabIndex(tab))
}

func (t *TabView) SelectTabByIndex(idx int) {
	if idx == tabNotFound {
		return
	}

	t.pager.SetSelectedChild(idx)
}

func (t *TabView) getTabIndex(tab *Tab) int {
	for idx := range t.tabs {
		if t.tabs[idx] != tab {
			continue
		}

		return idx
	}

	return tabNotFound
}
