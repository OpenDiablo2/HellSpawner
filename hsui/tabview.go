package hsui

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"github.com/hajimehoshi/ebiten"
)

const (
	tabNotFound       = -1
	defaultTabPadding = 2
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

	tv.padding = defaultTabPadding

	tv.tabBox.SetChildSpacing(2)
	tv.tabBox.SetAlignment(hscommon.HAlignLeft)

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
	padding                    int
	reqWidth, reqHeight        int
	enabled, visible, expanded bool
}

func (t *TabView) Render(screen *ebiten.Image, x, y, width, height int) {
	_, tabBoxReqH := t.tabBox.GetRequestedSize()

	p := t.padding

	t.tabBox.Render(screen, x+p, y+p, width-p*2, tabBoxReqH-p*2)
	t.pager.Render(screen, x+p, y+tabBoxReqH+p, width-p*2, height-tabBoxReqH-p*2)
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

	if len(t.tabs) == 1 {
		t.SelectTab(tab)
	}

	t.Invalidate()
}

func (t *TabView) RemoveTab(tab *Tab) {
	removedIdx := t.getTabIndex(tab)

	if removedIdx == tabNotFound {
		return
	}

	// we cant remove children from hbox, so just make a new one
	t.tabs = append(t.tabs[:removedIdx], t.tabs[removedIdx+1:]...)
	t.tabBox = CreateHBox()
	t.pager = CreatePager(0, 0, nil)

	oneTabSelected := false

	for tabIdx := range t.tabs {
		t.tabBox.AddChild(t.tabs[tabIdx])
		t.pager.AddChild(t.tabs[tabIdx].content)

		if t.tabs[tabIdx].selected {
			t.SelectTabByIndex(tabIdx)

			oneTabSelected = true
		}
	}

	if !oneTabSelected {
		if removedIdx-1 < 0 {
			removedIdx = 1
		}

		t.SelectTabByIndex(removedIdx - 1)
	}
}

func (t *TabView) SelectTab(tab *Tab) {
	t.SelectTabByIndex(t.getTabIndex(tab))
}

func (t *TabView) SelectTabByIndex(selectedIdx int) {
	if selectedIdx == tabNotFound {
		return
	}

	t.pager.SetSelectedChild(selectedIdx)

	for tabIdx := range t.tabs {
		t.tabs[tabIdx].selected = tabIdx == selectedIdx
	}

	t.Invalidate()
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

func (t *TabView) SetTabPadding(p int) {
	t.padding = p
	t.Invalidate()
}
