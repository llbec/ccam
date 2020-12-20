package scrcpy

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

func adbDevice() []string {
	var list []string
	ret := adbRun("", "devices")
	if len(ret) != 0 {
	}
	return list
}

func adbPush(serial, local, remote string) error {
	return adbExec(serial, "push", local, remote)
}

func adbInstall(serial, local string) error {
	return adbExec(serial, "install", "-r", local)
}

func adbRemovePath(serial, path string) error {
	return adbExec(serial, "shell", "rm", "-rf", path)
}

func adbReverse(serial, sockName string, localPort int) error {
	return adbExec(serial, "reverse",
		fmt.Sprintf("localabstract:%s", sockName),
		fmt.Sprintf("tcp:%d", localPort))
}

func adbReverseRemove(serial, sockName string) error {
	return adbExec(serial, "reverse", "--remove",
		fmt.Sprintf("localabstract:%s", sockName))
}

func adbForward(serial string, localPort int, sockName string) error {
	return adbExec(serial, "forward",
		fmt.Sprintf("tcp:%d", localPort),
		fmt.Sprintf("localabstract:%s", sockName))
}

func adbForwardRemove(serial string, localPort int) error {
	return adbExec(serial, "forward", "--remove",
		fmt.Sprintf("tcp:%d", localPort))
}

func adbRun(serial string, params ...string) string {
	args := make([]string, 0, 8)
	if len(serial) > 0 {
		args = append(args, "-s", serial)
	}
	args = append(args, params...)

	adbCmdOnce.Do(getAdbCommand)
	if debugOpt.Debug() {
		log.Printf("执行 %s %s\n", adbCmd, strings.Join(args, " "))
	}
	outinfo := bytes.Buffer{}
	cmd := exec.Command(adbCmd, args...)
	cmd.Stdout = &outinfo
	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err = cmd.Wait(); err != nil {
		if debugOpt.Error() {
			log.Printf(err.Error())
		}
	}
	if debugOpt.Debug() {
		log.Printf("PID %v\n", cmd.ProcessState.Pid())
		log.Printf("ExitCode %v\n", cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
		log.Printf(outinfo.String())
	}
	return outinfo.String()
}

func adbExec(serial string, params ...string) error {
	cmd, ret, err := adbExecAsync(serial, params...)
	if err != nil {
		return err
	}
	err = cmd.Wait()
	fmt.Println(ret.String())
	return err
}

func adbExecAsync(serial string, params ...string) (*exec.Cmd, bytes.Buffer, error) {
	args := make([]string, 0, 8)
	if len(serial) > 0 {
		args = append(args, "-s", serial)
	}
	args = append(args, params...)

	adbCmdOnce.Do(getAdbCommand)
	if debugOpt.Debug() {
		log.Printf("adbExecAsync %s %s\n", adbCmd, strings.Join(args, " "))
	}
	outinfo := bytes.Buffer{}
	cmd := exec.Command(adbCmd, args...)
	cmd.Stdout = &outinfo
	if debugOpt.Debug() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Start(); err != nil {
		return nil, outinfo, err
	}
	return cmd, outinfo, nil
}

var adbCmd = "adb"
var adbCmdOnce sync.Once

func getAdbCommand() {
	adbEnv := os.Getenv("ADB")
	if len(adbEnv) > 0 {
		adbCmd = adbEnv
	} else {
		switch runtime.GOOS {
		case "windows":
			adbCmd = "../3rd/adb/win/adb.exe"
		case "linux":
			adbCmd = "../3rd/adb/linux/adb"
		default:
			adbCmd = "../3rd/adb/mac/adb"
		}
	}
}
