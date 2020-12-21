package scrcpy

import (
	"fmt"
	"testing"
	"time"
)

func Test_exec(t *testing.T) {
	ip := adbGetIP("")
	if len(ip) == 0 {
		t.Errorf("get ip address failed")
		return
	}
	fmt.Printf("result:\n%s\n", ip)
}

func Test_adbReverse(t *testing.T) {
	err := adbReverse("", "tcpip", 5555)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_adbWireless(t *testing.T) {
	list := adbDevices()
	if len(list) != 1 {
		t.Errorf("not prepared!")
		return
	}
	dvsUSB := list[0]
	if err := adbTCPMod(dvsUSB, 5555); err != nil {
		t.Error(err.Error())
		return
	}
	time.Sleep(time.Duration(2) * time.Second)
	/*cmd := exec.Command(adbCmd, dvsUSB, "tcpip", fmt.Sprintf("%d", 5555))
	cmd.Run()*/
	ip := adbGetIP(dvsUSB)
	if len(ip) == 0 {
		t.Errorf("get ip address failed")
		return
	}
	if err := adbWirelessConnect(dvsUSB, ip, 5555); err != nil {
		t.Error(err.Error())
		return
	}
	list = adbDevices()
	if len(list) != 2 {
		t.Errorf("(2)invalid devices %v", list)
		return
	}
	if err := adbWirelessDisconnect(dvsUSB, ip); err != nil {
		t.Error(err.Error())
		return
	}
	list = adbDevices()
	if len(list) != 1 {
		t.Errorf("(1)invalid devices %v", list)
		return
	}
	if err := adbUSBMod(dvsUSB); err != nil {
		t.Error(err.Error())
		return
	}
	//fmt.Println("IP address: ", ip)
	/*list = adbDevices()
	if len(list) != 2 {
		t.Errorf("invalid devices %v", list)
		return
	}
	if err := adbUSBMod(list[1]); err != nil {
		t.Error(err.Error())
		return
	}
	list = adbDevices()
	if len(list) != 1 {
		t.Errorf("invalid devices %v", list)
		return
	}*/
}
