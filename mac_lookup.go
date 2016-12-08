package main

import (
	"regexp"
	"os"
	"bufio"
	"fmt"
)

var hexPat = regexp.MustCompile(`(.*) \(hex\) (.*)`)
type MacManufacture struct {
	macprefix string
	name string
}

func parseDevices() []MacManufacture {
	ret := make([]MacManufacture, 0)

	defaultDevicesFile := "./devices.txt"
	f, err := os.Open(defaultDevicesFile)
	logAndDie(err)
	defer f.Close()

	scanner := bufio.NewScanner(f) //read file line by line
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if hexPat.MatchString(scanner.Text()) {
			fmt.Println(line)
		}
	}
	return ret
}
