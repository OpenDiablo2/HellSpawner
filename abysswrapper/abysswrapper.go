package abysswrapper

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

const (
	waitTime = 3
)

// AbyssWrapper represents abyss wrapper
type AbyssWrapper struct {
	running bool
	output  io.Writer
	cmd     *exec.Cmd
	mutex   sync.RWMutex
}

// Create creates new Abyss Wrapper
func Create() *AbyssWrapper {
	result := &AbyssWrapper{}
	result.mutex = sync.RWMutex{}

	return result
}

func (a *AbyssWrapper) Read(p []byte) (n int, err error) {
	time.Sleep(time.Second * waitTime)
	bytes := []byte("Hello from HellSpawner! " + time.Now().String() + "\n")
	n = copy(p, bytes)
	err = nil

	return
}

func (a *AbyssWrapper) Write(p []byte) (n int, err error) {
	n, err = a.output.Write(p)
	if err != nil {
		return n, fmt.Errorf("error writing to output: %w", err)
	}

	return n, nil
}

// Launch launches abyss wrapper
func (a *AbyssWrapper) Launch(config *hsconfig.Config, output io.Writer) error {
	a.mutex.RLock()
	if a.running {
		a.mutex.RUnlock()

		return nil
	}
	a.mutex.RUnlock()
	a.mutex.Lock()

	a.output = output
	a.cmd = exec.Command(config.AbyssEnginePath) // nolint:gosec // is ok
	a.cmd.Stdout = a
	a.cmd.Stderr = a
	a.cmd.Stdin = a

	if err := a.cmd.Start(); err != nil {
		a.mutex.Unlock()
		return fmt.Errorf("error while running AbyssWrapper: %w", err)
	}

	a.running = true

	a.mutex.Unlock()

	go func() {
		_ = a.cmd.Wait()

		a.mutex.Lock()
		a.running = false
		a.mutex.Unlock()
	}()

	return nil
}

// Kill stops abyss wrapper
func (a *AbyssWrapper) Kill() error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	if !a.running {
		return nil
	}

	if err := a.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error closing AbyssWrapper: %w", err)
	}

	return nil
}

// IsRunning returns true, if AbyssWrapper is running
func (a *AbyssWrapper) IsRunning() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.running
}
