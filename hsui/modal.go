package hsui

import "github.com/hajimehoshi/ebiten/v2"

func CreateModal() *Modal {
	return &Modal{children: make([]Widget, 0)}
}

type Modal struct {
	children []Widget
}

func (m *Modal) Render(screen *ebiten.Image, x, y, width, height int) {
	for idx := range m.children {
		m.children[idx].Render(screen, x, y, width, height)
	}
}

func (m *Modal) Update() (dirty bool) {
	if currentChild := m.getLastChild(false); currentChild != nil {
		return currentChild.Update()
	}

	return false
}

func (m *Modal) GetRequestedSize() (int, int) {
	if currentChild := m.getLastChild(false); currentChild != nil {
		return currentChild.GetRequestedSize()
	}

	return 0, 0
}

func (m *Modal) Invalidate() {
	for idx := range m.children {
		m.children[idx].Invalidate()
	}
}

func (m *Modal) getLastChild(pop bool) Widget {
	numChildren := len(m.children)

	if numChildren < 1 {
		return nil
	}

	lastIdx := numChildren - 1
	lastChild := m.children[lastIdx]

	if pop {
		m.children = append(m.children[:lastIdx], m.children[lastIdx:]...)
	}

	return lastChild
}

func (m *Modal) Push(child Widget) {
	m.children = append(m.children, child)
}

func (m *Modal) Pop() Widget {
	return m.getLastChild(true)
}
