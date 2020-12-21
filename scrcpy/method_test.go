package scrcpy_test

import (
	"ccam/scrcpy"
	"fmt"
	"testing"
)

func Test_adbDevice(t *testing.T) {
	list := scrcpy.GetDevices()
	fmt.Printf("%v\n", list)
	if len(list) == 0 {
		t.Errorf("no devices!")
	}
}
