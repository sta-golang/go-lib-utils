package cmd

import (
	"github.com/sta-golang/go-lib-utils/str"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
)

const (
	CommandStatusInit = iota
	CommandStatusRunning
	CommandStatusFinish
)

type ExecCommand struct {
	ProcessState *os.ProcessState
	OutMessage   []byte
	status       int32
	RunErr       error
	wg           sync.WaitGroup
}

func NewCommand(cmd *exec.Cmd) *ExecCommand {
	return &ExecCommand{
		ProcessState: cmd.ProcessState,
		OutMessage:   nil,
		RunErr:       nil,
	}
}

func (ec *ExecCommand) Wait() {
	if atomic.LoadInt32(&ec.status) != CommandStatusRunning {
		return
	}
	ec.wg.Wait()
}

func (ec *ExecCommand) Pid() int {
	return ec.ProcessState.Pid()
}

func (ec *ExecCommand) ExitCode() int {
	return ec.ProcessState.ExitCode()
}

func (ec *ExecCommand) IsFinish() bool {
	return atomic.LoadInt32(&ec.status) == CommandStatusFinish
}

func (ec *ExecCommand) IsRunning() bool {
	return atomic.LoadInt32(&ec.status) == CommandStatusRunning
}

func (ec *ExecCommand) OutInfo() string {
	return str.BytesToString(ec.OutMessage)
}
