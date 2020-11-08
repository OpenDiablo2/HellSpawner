package hsui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	noChildren = -1
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

func (p *Pager) GetChildren() []Widget {
	return p.children
}

func (p *Pager) GetChild(idx int) Widget {
	if idx < 0 || idx >= len(p.children) {
		return nil
	}

	return p.children[idx]
}

func (p *Pager) Render(screen *ebiten.Image, x, y, width, height int) {
	if width < p.reqWidth {
		width = p.reqWidth
	}

	if height < p.reqHeight {
		height = p.reqHeight
	}

	if child, err := p.GetSelectedChild(); err == nil {
		child.Render(screen, x, y, width, height)
	}
}

func (p *Pager) Update() (dirty bool) {
	dirty = false

	if len(p.children) < 1 {
		return false
	}

	child := p.children[p.selectedChild]
	if child == nil {
		return
	}

	childDirty := child.Update()

	if childDirty {
		dirty = true
	}

	if dirty {
		p.Invalidate()
		p.dirty = true
	}

	return dirty
}

func (p *Pager) GetRequestedSize() (int, int) {
	w := 0
	h := 0

	for idx := range p.children {
		cw, ch := p.children[idx].GetRequestedSize()
		w += cw
		h += ch
	}

	return w, h
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

	p.selectedChild = idx
	p.Invalidate()
}

func (p *Pager) GetSelectedChild() (Widget, error) {
	if p.selectedChild < 0 || p.selectedChild >= len(p.children) {
		return nil, fmt.Errorf("hsui.Pager: no child at index %d", p.selectedChild)
	}

	return p.GetChild(p.selectedChild), nil
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
	if len(p.children) < 1 {
		return
	}

	p.SetSelectedChild((p.selectedChild + 1) % len(p.children))
}

func (p *Pager) SelectPreviousChild() {
	if len(p.children) < 1 {
		return
	}

	idx := p.selectedChild - 1

	for idx < 0 {
		idx += len(p.children)
	}

	idx %= len(p.children)

	p.SetSelectedChild(idx)
}
