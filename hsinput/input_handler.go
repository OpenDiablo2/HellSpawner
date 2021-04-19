// Package hsinput handles keyboard inputs (e.g. for shortcuts) on a per-window basis.
// Shortcuts are stored in a map. Every combination of keys can be assigned a function that is
// executed when the key combo is detected.
// The list of functions assigned to a key combo operates like a stack.
// When a key combo is detected, the function that was most recently pushed into the stack is executed.
// A global shortcut, if one has been defined for a certain key combo,
// would be at the very bottom of a key combo's list, since those are registered on startup.
// When a window gains focus, its keyboard shortcuts are pushed to the stack.
// When a window loses focus, its keyboard shortcuts are popped from the stack.
// This ensures that only the window that has focus and the global shortcuts are active at any given time.
package hsinput

import (
	"github.com/AllenDang/giu"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	// ModNone is used when you want to specify a key combo without any modifier keys
	ModNone = 0
)

// InputManager represents registered game shortcuts
type InputManager struct {
	shortcuts map[KeyCombo]*CallbackGroup
}

// NewInputManager creates a new input manager
func NewInputManager() *InputManager {
	result := &InputManager{
		shortcuts: make(map[KeyCombo]*CallbackGroup),
	}

	return result
}

// InputCallbackFunc is the function signature for functions that are called on input events
type InputCallbackFunc func()

// CallbackGroup defines the Global and Window-specific callbacks associated with a KeyCombo
type CallbackGroup struct {
	Global InputCallbackFunc
	Window InputCallbackFunc
}

// KeyCombo defines a Key and a Modifier (e.g. Ctrl+A)
type KeyCombo struct {
	Key      glfw.Key
	Modifier glfw.ModifierKey
}

func createKeyCombo(key giu.Key, modifier glfw.ModifierKey) KeyCombo {
	return KeyCombo{
		Key:      glfw.Key(key),
		Modifier: modifier,
	}
}

// RegisterShortcut registers a new shortcut
func (im *InputManager) RegisterShortcut(callbackFunc InputCallbackFunc, key giu.Key, modifier glfw.ModifierKey, isGlobal bool) {
	combo := createKeyCombo(key, modifier)

	shortcut, alreadyRegistered := im.shortcuts[combo]
	if !alreadyRegistered {
		shortcut = &CallbackGroup{}
	}

	if isGlobal {
		shortcut.Global = callbackFunc
	} else {
		shortcut.Window = callbackFunc
	}

	im.shortcuts[combo] = shortcut
}

// UnregisterWindowShortcuts removes registered window's shortcuts
func (im *InputManager) UnregisterWindowShortcuts() {
	for _, callbackFuncs := range im.shortcuts {
		callbackFuncs.Window = nil
	}
}

// HandleInput handles input shortcut
func (im *InputManager) HandleInput(key glfw.Key, mods glfw.ModifierKey, action glfw.Action) {
	for keyCombo, callbackFuncs := range im.shortcuts {
		if key == keyCombo.Key && mods == keyCombo.Modifier && action == glfw.Press {
			if callbackFuncs.Window != nil {
				callbackFuncs.Window()
			} else if callbackFuncs.Global != nil {
				callbackFuncs.Global()
			}

			return
		}
	}
}
