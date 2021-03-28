package hsinput

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ianling/giu"
)

func Test_InputManager_Create(t *testing.T) {
	im := NewInputManager()
	if im == nil {
		t.Fatal("Input manager wasn't set up")
	}

	if im.shortcuts == nil {
		t.Fatal("Wrong input manager created")
	}
}

func Test_inputCombo_create(t *testing.T) {
	key, mod := giu.Key(1), giu.Modifier(1)
	combo := createKeyCombo(key, mod)

	if int(combo.Key) != int(key) || int(combo.Modifier) != int(mod) {
		t.Fatal("wrong key combo created")
	}
}

func Test_InputManager_RegisterShortcut(t *testing.T) {
	im := NewInputManager()
	if im == nil {
		t.Fatal("Error creating input manager")
	}

	key, mod := giu.Key(1), giu.Modifier(2)
	combo := createKeyCombo(key, mod)
	global := true

	cb := func() {}

	im.RegisterShortcut(cb, key, mod, global)

	shortcut, exist := im.shortcuts[combo]

	if !exist {
		t.Fatal("shortcut wasn't registered")
	}

	if shortcut.Global == nil {
		t.Fatal("callback wasn't written")
	}

	if shortcut.Window != nil {
		t.Fatal("callback for window isn't nil")
	}
}

func Test_InputManager_UnregisterWindowShortcuts(t *testing.T) {
	im := NewInputManager()
	if im == nil {
		t.Fatal("Error creating input manager")
	}

	key, mod := giu.Key(1), giu.Modifier(2)
	combo := createKeyCombo(key, mod)
	cb := func() {}

	shortcut := &CallbackGroup{Window: cb}
	im.shortcuts[combo] = shortcut

	im.UnregisterWindowShortcuts()

	if im.shortcuts[combo].Window != nil {
		t.Fatal("Shortcuts wasn't unregistered")
	}
}

func Test_InputManager(t *testing.T) {
	im := NewInputManager()
	if im == nil {
		t.Fatal("Error creating input manager")
	}

	key, mod := giu.Key(1), giu.Modifier(2)
	combo := createKeyCombo(key, mod)
	ok := false
	cb := func() {
		ok = true
	}

	shortcut := &CallbackGroup{Global: cb}
	im.shortcuts[combo] = shortcut

	im.HandleInput(glfw.Key(key), glfw.ModifierKey(mod), glfw.Press)

	if !ok {
		t.Fatal("Callback wasn't handled")
	}
}
