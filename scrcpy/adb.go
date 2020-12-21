package scrcpy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func adbTCPMod(serial string, port int) error {
	return adbExec(serial, "tcpip", fmt.Sprintf("%d", port))
}

func adbUSBMod(serial string) error {
	return adbExec(serial, "usb")
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

	cmd := exec.Command(adbCmd, args...)
	//outinfo := bytes.Buffer{}
	//cmd.Stdout = &outinfo
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		if debugOpt.Error() {
			log.Printf(err.Error())
		}
		return ""
	}

	outinfo, _ := ioutil.ReadAll(stdout)
	stdout.Close()

	if err = cmd.Wait(); err != nil {
		if debugOpt.Error() {
			log.Printf(err.Error())
		}
		return ""
	}
	if debugOpt.Debug() {
		log.Printf("PID %v\n", cmd.ProcessState.Pid())
		//log.Printf("ExitCode %v\n", cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
		//log.Printf(outinfo.String())
	}
	return string(outinfo)
}

func adbRun1(serial string, params ...string) string {
	args := make([]string, 0, 8)
	if len(serial) > 0 {
		args = append(args, "-s", serial)
	}
	args = append(args, params...)

	adbCmdOnce.Do(getAdbCommand)
	if debugOpt.Debug() {
		log.Printf("执行 %s %s\n", adbCmd, strings.Join(args, " "))
	}

	cmd := exec.Command(adbCmd, args...)
	ret, err := cmd.Output()
	if err != nil {
		if debugOpt.Error() {
			log.Printf(err.Error())
			return ""
		}
	}
	return string(ret)
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
