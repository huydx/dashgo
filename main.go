package main

import (
	"fmt"
	"os"
	//"time"

	//"github.com/google/gopacket"
	//"github.com/google/gopacket/layers"
	//"github.com/google/gopacket/pcap"
	//"strconv"
	"log"
	"bufio"
)

func main() {
	parseDevices()
	//var device string
	//var timeout = 30 * time.Second

	//print all device
	//devices, err := pcap.FindAllDevs()
	//logAndDie(err)
	//
	//for i, device := range devices {
	//	fmt.Printf("Device %d: %s\n", i, device)
	//	d := prompt("Enter your device: ")
	//	di, err := strconv.Atoi(d)
	//	logAndDie(err)
	//	device = devices[di]
	//}
	//handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
}

func logAndDie(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func prompt(prefix string) string {
	fmt.Print(prefix)
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	logAndDie(scanner.Err())

	return line
}
