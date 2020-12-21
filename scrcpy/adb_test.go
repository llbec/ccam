package scrcpy

import (
	"fmt"
	"testing"
)

func Test_adbDevice(t *testing.T) {
	list := adbDevice()
	fmt.Printf("%v", list)
	if len(list) == 0 {
		t.Errorf("no devices!")
	}
}

func Test_adbReverse(t *testing.T) {
	err := adbReverse("", "tcpip", 5555)
	if err != nil {
		t.Errorf(err.Error())
	}
}
