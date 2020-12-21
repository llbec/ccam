package scrcpy

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func adbDevices() []string {
	var list []string
	ret, err := adbExecOut("", "devices")
	if err != nil {
		if debugOpt.Error() {
			log.Printf("adb devices failed: %s(%s)", err.Error(), ret)
			return list
		}
	}
	//fmt.Println(ret)
	if len(ret) != 0 {
		reg, err := regexp.Compile("([\\S]+)[\\s]+device\\b")
		if err == nil {
			list = reg.FindAllString(ret, -1)
			//fmt.Println(list)
		}
	}
	for i, v := range list {
		list[i] = strings.Fields(strings.TrimSpace(v))[0]
	}
	return list
}

func adbTCPMod(serial string, port int) error {
	return adbExec(serial, "tcpip", fmt.Sprintf("%d", port))
}

func adbGetIP(serial string) string {
	res, err := adbExecOut(serial, "shell", "ip", "-o", "a")
	if err != nil {
		if debugOpt.Error() {
			log.Printf("adb shell ip failed: %s(%s)", err.Error(), res)
			return ""
		}
	}
	if len(res) > 0 {
		reg, err := regexp.Compile("\\swlan0[\\s]+inet\\s([\\d]+.[\\d]+.[\\d]+.[\\d]+)")
		if err == nil {
			list := reg.FindAllString(res, -1)
			if len(list) > 0 {
				ips := strings.Fields(strings.TrimSpace(list[0]))
				if len(ips) == 3 {
					return ips[2]
				}
				if debugOpt.Error() {
					log.Printf("adbGetIP regex invalid format: %d", len(ips))
				}
			} else {
				if debugOpt.Error() {
					log.Printf("adbGetIP regex invalid result: %d", len(list))
				}
			}
		} else {
			if debugOpt.Error() {
				log.Printf("adbGetIP regex failed: %s", err.Error())
			}
		}
	}
	return ""
}

func adbWirelessConnect(serial, ip string, port int) error {
	return adbExec(serial, "connect", fmt.Sprintf("%s:%d", ip, port))
}

func adbWirelessDisconnect(serial, ip string) error {
	return adbExec(serial, "disconnect", ip)
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

/*func adbRun(serial string, params ...string) string {
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
	outinfo := bytes.Buffer{}
	cmd.Stdout = &outinfo
	//stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		if debugOpt.Error() {
			log.Printf("cmd start failed %s", err.Error())
		}
		return ""
	}

	//outinfo, _ := ioutil.ReadAll(stdout)
	//stdout.Close()

	if err = cmd.Wait(); err != nil {
		if debugOpt.Error() {
			log.Printf("cmd wait failed %s", err.Error())
		}
		return ""
	}
	if debugOpt.Debug() {
		log.Printf("PID %v\n", cmd.ProcessState.Pid())
		//log.Printf("ExitCode %v\n", cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
		log.Printf(outinfo.String())
	}
	return outinfo.String()
}*/

func adbExecOut(serial string, params ...string) (string, error) {
	var buf bytes.Buffer
	var berr bytes.Buffer
	cmd, err := adbExecAsync(&buf, &berr, serial, params...)
	if err != nil {
		return "", err
	}
	//fmt.Println("adb execute waiting")
	err = cmd.Wait()
	if err != nil {
		return berr.String(), err
	}
	return buf.String(), nil
}

func adbExec(serial string, params ...string) error {
	var buf bytes.Buffer
	cmd, err := adbExecAsync(&buf, &buf, serial, params...)
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func adbExecAsync(pbuf, perr *bytes.Buffer, serial string, params ...string) (*exec.Cmd, error) {
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
	/*if debugOpt.Debug() {
		cmd.Stderr = os.Stderr
		//cmd.Stdout = os.Stdout
	}*/
	cmd.Stderr = perr
	cmd.Stdout = pbuf

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
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
