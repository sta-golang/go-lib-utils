package cmd

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"sync/atomic"
)

const (
	GB18030 = "GB18030"
)

func startCommand(cmdStr string, args ...string) (*exec.Cmd, *ExecCommand, *bytes.Buffer, error) {
	command := exec.Command(cmdStr, args...)
	execCmd := NewCommand(command)
	buff := bytes.Buffer{}
	command.Stdout = &buff
	err := command.Start()
	if err != nil {
		return nil, nil, nil, err
	}
	atomic.StoreInt32(&execCmd.status, CommandStatusRunning)
	return command, execCmd, &buff, nil
}

func ExecCmd(cmdStr string, args ...string) (*ExecCommand, error) {
	command, execCmd, buff, err := startCommand(cmdStr, args...)
	if err != nil {
		return nil, err
	}

	doExecWait(command, execCmd, buff)

	return execCmd, nil
}

func ExecCmdAsync(cmdStr string, args ...string) (*ExecCommand, error) {
	command, execCmd, buff, err := startCommand(cmdStr, args...)
	if err != nil {
		return nil, err
	}

	execCmd.wg.Add(1)
	go doExecWait(command, execCmd, buff)
	return execCmd, nil
}

func doExecWait(cmd *exec.Cmd, execCmd *ExecCommand, bytesBuff *bytes.Buffer) {
	defer func() {
		execCmd.wg.Done()
	}()
	if err := cmd.Wait(); err != nil {
		if execCmd != nil {
			execCmd.RunErr = err
		}
	}
	if execCmd != nil {
		execCmd.OutMessage, execCmd.RunErr = simplifiedchinese.GB18030.NewDecoder().Bytes(bytesBuff.Bytes())
		atomic.CompareAndSwapInt32(&execCmd.status, CommandStatusRunning, CommandStatusFinish)
	}
}
