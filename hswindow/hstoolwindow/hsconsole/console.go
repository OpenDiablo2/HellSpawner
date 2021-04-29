// Package hsconsole contains project's console
package hsconsole

import (
	"fmt"
	"os"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
)

const (
	mainWindowW, mainWindowH = 600, 200
	lineW, lineH             = -1, -1
)

// Console represents a console
type Console struct {
	*hstoolwindow.ToolWindow
	outputText string
	fontFixed  imgui.Font
	logFile    *os.File
}

// Create creates a new console
func Create(fontFixed imgui.Font, x, y float32, logFile *os.File) *Console {
	result := &Console{
		fontFixed:  fontFixed,
		ToolWindow: hstoolwindow.New("Console", hsstate.ToolWindowTypeConsole, x, y),
		logFile:    logFile,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result
}

// Build builds a console
func (c *Console) Build() {
	c.IsOpen(&c.Visible).
		Layout(g.Layout{
			g.Custom(func() {
				g.PushFont(c.fontFixed)
			}),
			g.InputTextMultiline("", &c.outputText).
				Size(lineW, lineH).
				Flags(g.InputTextFlags_ReadOnly | g.InputTextFlags_NoUndoRedo),
			g.Custom(func() {
				g.PopFont()
			}),
		})
}

// Write writes input on console
func (c *Console) Write(p []byte) (n int, err error) {
	msg := string(p) // convert message from byte slice into string

	c.outputText = msg + c.outputText // append message

	fmt.Print(msg) // print to terminal

	n, err = c.logFile.Write(p) // print to file
	if err != nil {
		return n, fmt.Errorf("Error writting to log file: %w", err)
	} else if n != len(p) {
		return n, fmt.Errorf("Invalid data writen to log file")
	}

	return len(p), nil
}
