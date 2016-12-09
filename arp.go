package main

import (
	"bytes"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"fmt"
	"strings"
)


func scanButtonArp(manufatures Manufactures) string {
	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var dashMac string

	for _, iface := range ifaces {
		wg.Add(1)
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			defer wg.Done()
			dashMac, err = scan(manufatures, &iface)
		}(iface)
	}
	wg.Wait()
	return dashMac
}

// scan scans an individual interface's local network for machines using ARP requests/replies.
//
// scan loops forever, sending packets out regularly.  It returns an error if
// it's ever unable to write a packet.
func scan(manufatures Manufactures, iface *net.Interface) (string, error) {
	// We just look for IPv4 addresses, so try to find if the interface has one.
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return "", err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return "", errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return "", errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return "", errors.New("mask means network is too large")
	}
	log.Printf("Using network range %v for interface %v", addr, iface.Name)

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return "", err
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			var mac = net.HardwareAddr(arp.SourceHwAddress)
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
				continue
			}
			if m := manufatures.find(mac.String()); m != nil {
				if strings.Contains(strings.ToLower(m.name), "amazon") {
					return string(mac), nil
				}
			}
		}
	}
	return "", fmt.Errorf("can't found")
}
