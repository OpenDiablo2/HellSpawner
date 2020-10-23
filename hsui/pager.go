package hsui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

func CreatePager(w, h int, children []Widget) *Pager {
	p := &Pager{
		children:  children,
		reqWidth:  w,
		reqHeight: h,
		dirty:     true,
		visible:   true,
		enabled:   true,
	}

	return p
}

type Pager struct {
	children                []Widget
	selectedChild           int
	reqWidth, reqHeight     int
	dirty, visible, enabled bool
}

const (
	noChildren = -1
)

func (p *Pager) Render(screen *ebiten.Image, x, y, width, height int) {
	if child, err := p.GetSelectedChild(); err == nil {
		child.Render(screen, x, y, p.reqWidth, p.reqHeight)
	}
}

func (p *Pager) Update() (dirty bool) {
	if p.dirty {
		p.Invalidate()
	}

	dirty = false

	child := p.children[p.selectedChild]
	if child == nil {
		return
	}

	childDirty := child.Update()

	if childDirty {
		dirty = true
	}

	if dirty {
		p.dirty = true
	}

	return dirty
}

func (p *Pager) GetRequestedSize() (int, int) {
	return p.reqWidth, p.reqHeight
}

func (p *Pager) Invalidate() {
	for idx := range p.children {
		p.children[idx].Invalidate()
	}
}

func (p *Pager) SetSelectedChild(idx int) {
	if idx < 0 || idx >= len(p.children) {
		idx = noChildren
	}

	if p.selectedChild == idx {
		return
	}

	p.Invalidate()

	p.selectedChild = idx
	fmt.Println(idx)
}

func (p *Pager) GetSelectedChild() (Widget, error) {
	if p.selectedChild < 0 || p.selectedChild >= len(p.children) {
		return nil, fmt.Errorf("hsui.Pager: no child at index %d", p.selectedChild)
	}

	return p.children[p.selectedChild], nil
}

func (p *Pager) AddChild(child Widget) {
	p.children = append(p.children, child)
}

func (p *Pager) RemoveChild(idx int) {
	_, err := p.GetSelectedChild()
	if err != nil {
		return
	}

	p.children = append(p.children[:idx], p.children[idx:]...)

	if p.selectedChild >= idx {
		p.SetSelectedChild(p.selectedChild - 1)
	}
}

func (p *Pager) SelectNextChild() {
	p.SetSelectedChild((p.selectedChild + 1) % len(p.children))
}

func (p *Pager) SelectPreviousChild() {
	idx := p.selectedChild - 1

	for idx < 0 {
		idx += len(p.children)
	}

	idx %= len(p.children)

	p.SetSelectedChild(idx)
}
