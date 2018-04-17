package cmd

import (
        "os/exec"
        "time"
        "fmt"
        "syscall"
        "bytes"
)

func RunCombineCommandWithTimeOut(cmds []string, timeout time.Duration) (err error, returnCode int, stderrout []byte) {
        if len(cmds) <= 0 {
                err = fmt.Errorf("Invalid Command %v", cmds)
                return
        }
        cmd := exec.Command(cmds[0], cmds[1:]...)
        done := make(chan error)
        go func() {
                stderrout, err = cmd.CombinedOutput()
                done <- err
        }()

        select {
        case <-time.After(timeout):
        // timeout
                if errr := cmd.Process.Kill(); errr != nil {
                        err = fmt.Errorf("failed to kill: %s, error: %s", cmd.Path, err)
                        return
                }
                go func() {
                        <-done // allow goroutine to exit
                }()
                err = fmt.Errorf("TimeOut %s", cmd.Path)
                return
        case err = <-done:
                if exiterr, ok := err.(*exec.ExitError); ok {
                        if status, ok := exiterr.ProcessState.Sys().(syscall.WaitStatus); ok {
                                returnCode = status.ExitStatus()

                        }
                }
        }
        return
}
func RunCommandWithTimeOut(cmds []string, timeout time.Duration) (err error, returnCode int, stdout []byte, stderr []byte) {
        if len(cmds) <= 0 {
                err = fmt.Errorf("Invalid Command %v", cmds)
                return
        }
        cmd := exec.Command(cmds[0], cmds[1:]...)
        cmdStdOut := &bytes.Buffer{}
        cmdStdErr := &bytes.Buffer{}
        cmd.Stdout = cmdStdOut
        cmd.Stderr = cmdStdErr

        cmd.Start()
        done := make(chan error)
        go func() {
                done <- cmd.Wait()
        }()

        select {
        case <-time.After(timeout):
        // timeout
                if errr := cmd.Process.Kill(); errr != nil {
                        err = fmt.Errorf("failed to kill: %s, error: %s", cmd.Path, err)
                        return
                }
                go func() {
                        <-done // allow goroutine to exit
                }()
                err = fmt.Errorf("TimeOut %s", cmd.Path)
                return
        case err = <-done:
                stdout = cmdStdOut.Bytes()
                stderr = cmdStdErr.Bytes()
                if exiterr, ok := err.(*exec.ExitError); ok {
                        if status, ok := exiterr.ProcessState.Sys().(syscall.WaitStatus); ok {
                                returnCode = status.ExitStatus()

                        }
                }
        }
        return
}
func RunTwoCommandWithTimeOut(cmds1 []string, cmds2 []string, timeout time.Duration) (err error, returnCode int, stdout []byte, stderr []byte) {
        if len(cmds1) <= 0 || len(cmds2) <= 0 {
                err = fmt.Errorf("Invalid Command %v, %v", cmds1, cmds2)
                return
        }

        cmdStdOut := &bytes.Buffer{}
        cmdStdErr := &bytes.Buffer{}
        cmd1 := exec.Command(cmds1[0], cmds1[1:]...)
        cmd2 := exec.Command(cmds2[0], cmds2[1:]...)
        cmd2.Stdin, _ = cmd1.StdoutPipe()
        cmd2.Stdout = cmdStdOut
        cmd2.Stderr = cmdStdErr

        cmd2.Start()
        cmd1.Start()
        done := make(chan error)
        go func() {
                cmd1.Wait()
                done <- cmd2.Wait()
        }()

        select {
        case <-time.After(timeout):
        // timeout
                if errr := cmd2.Process.Kill(); errr != nil {
                        err = fmt.Errorf("failed to kill: %s, error: %s", cmd2.Path, err)
                        return
                }
                go func() {
                        <-done // allow goroutine to exit
                }()
                err = fmt.Errorf("TimeOut %s", cmd2.Path)
                return
        case err = <-done:
                stdout = cmdStdOut.Bytes()
                stderr = cmdStdErr.Bytes()
                if exiterr, ok := err.(*exec.ExitError); ok {
                        if status, ok := exiterr.ProcessState.Sys().(syscall.WaitStatus); ok {
                                returnCode = status.ExitStatus()

                        }
                }
        }
        return
}

func RunMultiCommandWithTimeOut(timeout time.Duration, cmdStrs ...[]string) (err error, returnCode int, stdout []byte, stderr []byte) {
        if len(cmdStrs) < 1 {
                err = fmt.Errorf("Invalid Command %v", cmdStrs)
                return
        }
        cmds := make([]*exec.Cmd, 0)
        for _, cmd := range cmdStrs {
                cmds = append(cmds, exec.Command(cmd[0], cmd[1:]...))
        }
        cmdStdOut := &bytes.Buffer{}
        cmdStdErr := &bytes.Buffer{}

        cmdLen := len(cmds)

        cmds[cmdLen-1].Stdout = cmdStdOut
        cmds[cmdLen-1].Stderr = cmdStdErr

        for i := 1; i < cmdLen; i++ {
                cmds[i].Stdin, _ = cmds[i - 1].StdoutPipe()
        }
        // start and wait
        done := make(chan error)
        for i := cmdLen - 1; i >= 0; i-- {
                cmds[i].Start()
        }
        go func() {
                for i := 0; i < cmdLen - 1; i++ {
                        cmds[i].Wait()
                }
                done <- cmds[cmdLen-1].Wait()
        }()

        select {
        case <-time.After(timeout):
        // timeout
                errStr := ""
                for i := 0; i < cmdLen; i++ {
                        errr := cmds[i].Process.Kill()
                        if errr != nil {
                                errStr += fmt.Sprintf("failed to kill: %s, error: %s\n", cmds[i].Path, err)
                        }
                }
                if errStr != "" {
                        err = fmt.Errorf("%s", errStr)
                        return
                }
                go func() {
                        <-done // allow goroutine to exit
                }()
                err = fmt.Errorf("TimeOut %s", cmds[cmdLen-1].Path)
                return
        case err = <-done:
                stdout = cmdStdOut.Bytes()
                stderr = cmdStdErr.Bytes()
                if exiterr, ok := err.(*exec.ExitError); ok {
                        if status, ok := exiterr.ProcessState.Sys().(syscall.WaitStatus); ok {
                                returnCode = status.ExitStatus()

                        }
                }
        }
        return
}