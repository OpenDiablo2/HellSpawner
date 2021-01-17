package abysswrapper

import (
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

var (
	mutex sync.RWMutex = sync.RWMutex{}
)

type AbyssWrapper struct {
	running bool
	output  io.Writer
	cmd     *exec.Cmd
}

func (a *AbyssWrapper) Read(p []byte) (n int, err error) {
	time.Sleep(time.Second * 3)
	bytes := []byte("Hello from HellSpawner! " + time.Now().String() + "\n")
	n = copy(p, bytes)
	err = nil

	return
}

func (a *AbyssWrapper) Write(p []byte) (n int, err error) {
	return a.output.Write(p)
}

func Create() *AbyssWrapper {
	result := &AbyssWrapper{}
	return result
}

func (a *AbyssWrapper) Launch(config *hsconfig.Config, output io.Writer) error {
	mutex.RLock()
	if a.running {
		mutex.RUnlock()

		return nil
	}
	mutex.RUnlock()
	mutex.Lock()

	a.output = output
	a.cmd = exec.Command(config.AbyssEnginePath)
	a.cmd.Stdout = a
	a.cmd.Stderr = a
	a.cmd.Stdin = a

	if err := a.cmd.Start(); err != nil {
		mutex.Unlock()
		return err
	}

	a.running = true

	mutex.Unlock()

	go func() {
		_ = a.cmd.Wait()
		mutex.Lock()
		a.running = false
		mutex.Unlock()
	}()

	return nil
}

func (a *AbyssWrapper) Kill() error {
	mutex.RLock()
	defer mutex.RUnlock()

	if !a.running {
		return nil
	}

	return a.cmd.Process.Kill()
}

func (a *AbyssWrapper) IsRunning() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return a.running
}
