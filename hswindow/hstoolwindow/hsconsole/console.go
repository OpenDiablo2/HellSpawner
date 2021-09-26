// Package hsconsole provides a graphical console for logging output while the app is running.
package hsconsole

import (
	"fmt"
	"os"

	g "github.com/AllenDang/giu"

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
	fontFixed  *g.FontInfo
	logFile    *os.File
}

// Create creates a new console
func Create(fontFixed *g.FontInfo, x, y float32, logFile *os.File) *Console {
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
			g.Style().SetFont(c.fontFixed).To(
				g.InputTextMultiline(&c.outputText).
					Size(lineW, lineH).
					Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoUndoRedo),
			),
		})
}

// Write writes input on console, stdout and (if exists) to the log file
func (c *Console) Write(p []byte) (n int, err error) {
	msg := string(p) // convert message from byte slice into string

	c.outputText = msg + c.outputText // append message

	fmt.Print(msg) // print to terminal

	if c.logFile != nil {
		n, err = c.logFile.Write(p) // print to file
		if err != nil {
			return n, fmt.Errorf("error writing to log file: %w", err)
		} else if n != len(p) {
			return n, fmt.Errorf("invalid data written to log file")
		}
	}

	return len(p), nil
}
