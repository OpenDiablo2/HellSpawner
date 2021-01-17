package abysswrapper

import (
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

var (
	mutex sync.Mutex = sync.Mutex{}
)

type AbyssWrapper struct {
	running bool
	cmd     *exec.Cmd
}

func Create() *AbyssWrapper {
	result := &AbyssWrapper{}
	runtime.SetFinalizer(result, dispose)
	return result
}

func dispose(a *AbyssWrapper) {
	mutex.Lock()
	defer mutex.Unlock()

	if !a.running {
		return
	}
	_ = a.cmd.Process.Kill()
}

func (a *AbyssWrapper) Launch(config *hsconfig.Config) error {
	mutex.Lock()
	if a.running {
		mutex.Unlock()

		return nil
	}

	a.cmd = exec.Command(config.AbyssEnginePath)

	a.cmd.Stdout = os.Stdout
	a.cmd.Stderr = os.Stderr

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
	mutex.Lock()
	defer mutex.Unlock()

	if !a.running {
		return nil
	}

	return a.cmd.Process.Kill()
}

func (a *AbyssWrapper) IsRunning() bool {
	return a.running
}
