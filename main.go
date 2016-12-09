package main

import (
	"fmt"
	"os"
	"time"

	"strconv"
	"log"
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os/exec"
)

var (
	device string
	timeout = 30 * time.Second
	promiscuous = true
	snapshotLen = int32(2 << 16)
)

//34:d2:70:e3:ad:65
func main() {
	fmt.Println("reading manufactures list...")
	ms := manufactures()
	fmt.Println("reading done")



	devices, err := pcap.FindAllDevs()
	logAndDie(err)

	for i, device := range devices {
		fmt.Printf("Device %d: %s\n", i, device.Name)
	}

	d := prompt("Enter your device: ")
	di, err := strconv.Atoi(d)

	logAndDie(err)
	device = devices[di].Name

	var buttonMac string
	if len(os.Args) > 1 {
		buttonMac = os.Args[1]
	}
	if  buttonMac == "" {
		fmt.Println("start scan arp to find dash button")
		buttonMac = scanButtonArp(ms)
		fmt.Printf("found button with mac %s \n", buttonMac)
	}

	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	logAndDie(err)

	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var filter string = fmt.Sprintf("ether host %s", buttonMac)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for range packetSource.Packets() {
		err := exec.Command("./execute.sh").Run()
		fmt.Println(err)
	}

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
