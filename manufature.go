package main

import (
	"regexp"
	"os"
	"bufio"
	"strings"
)

var hexPat = regexp.MustCompile(`^([[:alnum:]]{2})-([[:alnum:]]{2})-([[:alnum:]]{2})\s+\(hex\)\s+(.*)`)
type Manufacture struct {
	m0 string //mac hex 1
	m1 string //mac hex 2
	m2 string //mac hex 3
	name string
}

type Manufactures []Manufacture
func (manufactures Manufactures) find(addr string) *Manufacture {
	addr = strings.ToUpper(addr)
	for _, m := range manufactures {
		mac := strings.Split(addr, ":")
		if mac[0] == m.m0 && mac[1] == m.m1 && mac[2] == m.m2 {
			return &m
		}
	}
	return nil
}

func manufactures() Manufactures {
	ret := make([]Manufacture, 0)

	defaultDevicesFile := "./devices.txt"
	f, err := os.Open(defaultDevicesFile)
	logAndDie(err)
	defer f.Close()

	scanner := bufio.NewScanner(f) //read file line by line
	for scanner.Scan() {
		line := scanner.Text()
		matched := hexPat.FindStringSubmatch(line)
		if len(matched) > 0 {
			ret = append(ret, Manufacture{matched[1], matched[2], matched[3], matched[4]})
		}
	}
	return ret
}
