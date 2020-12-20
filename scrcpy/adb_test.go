package scrcpy

import (
	"fmt"
	"testing"
)

func Test_adbDevice(t *testing.T) {
	//adbExec("", "devices")
	//adbRun("", "devices")
	list := adbDevice()
	fmt.Printf("%v", list)
}
