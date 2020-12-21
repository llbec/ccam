package scrcpy

import (
	"fmt"
	"os/exec"
	"testing"
)

func Test_exec(t *testing.T) {
	adbCmdOnce.Do(getAdbCommand)
	cmd := exec.Command(adbCmd, "devices")
	ret, err := cmd.Output()
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(string(ret))
}

func Test_adbReverse(t *testing.T) {
	err := adbReverse("", "tcpip", 5555)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_adbTcpMod(t *testing.T) {
	list := GetDevices()
	if len(list) != 1 {
		t.Errorf("not prepared!")
		return
	}
	dvsUSB := list[0]
	if err := adbTCPMod(dvsUSB, 5555); err != nil {
		t.Error(err.Error())
		return
	}
	list = GetDevices()
	if len(list) != 2 {
		t.Errorf("invalid devices %v", list)
		return
	}
	if err := adbUSBMod(list[1]); err != nil {
		t.Error(err.Error())
		return
	}
	list = GetDevices()
	if len(list) != 1 {
		t.Errorf("invalid devices %v", list)
		return
	}
}
