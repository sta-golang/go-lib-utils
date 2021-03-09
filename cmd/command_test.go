package cmd

import (
	"fmt"
	"testing"
)

func TestExecCmd(t *testing.T) {
	cmd, err := ExecCmdAsync("ping", "www.google.com")
	fmt.Println(err)
	fmt.Println("cmd is running", cmd.IsRunning())
	fmt.Println("cmd is finish", cmd.IsFinish())
	fmt.Println(cmd.status)
	fmt.Println(cmd.OutInfo())
	cmd.Wait()
	fmt.Println(err)
	fmt.Println("cmd is running", cmd.IsRunning())
	fmt.Println("cmd is finish", cmd.IsFinish())
	fmt.Println(cmd.OutInfo())
	fmt.Println(cmd.RunErr)
	cmd.Wait()
	fmt.Println("hello")
}
