package scrcpy

import (
	"regexp"
	"strings"
)

//GetDevices get android devices by adb
func GetDevices() []string {
	var list []string
	ret := adbRun("", "devices")
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
