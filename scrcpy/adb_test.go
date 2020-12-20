package scrcpy

import "testing"

func Test_adbExec(t *testing.T) {
	//adbExec("", "devices")
	adbRun("", "devices")
}
