package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/varuuntiwari/catch-tor/torips"
)

var (
	live   bool
	local  bool
	silent bool
)

var (
	snaplen     uint32 = 1024
	promiscuous bool   = true
	dev         string = "wlo1"
)

func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	flag.BoolVar(&live, "live", false, "capture packets live from interface")
	flag.BoolVar(&local, "local", false, "detect local traffic for tor connections")
	flag.BoolVar(&silent, "silent", false, "disable the dots that prove it does work")

	flag.Parse()

	if !live || !local {
		fmt.Println("[-] plis enable both live and local flags as it will not work any other way for now")
		os.Exit(0)
	}

	version := pcap.Version()
	if version != "" {
		fmt.Printf("[+] %s found\n", version)
	}

	h, err := pcap.OpenLive(dev, int32(snaplen), promiscuous, (-1 * time.Second))
	if err != nil {
		fmt.Println("[-] run as superuser to capture packets")
		fmt.Println("Exiting..")
		os.Exit(2)
	}
	defer h.Close()

	// size := torips.RefreshList()
	// fmt.Printf("written %v bytes, list of IPs refreshed\n", size)

	fmt.Println("[+] scanning packets now")
	count := 0
	source := gopacket.NewPacketSource(h, h.LinkType())
	for p := range source.Packets() {
		ipLayer := p.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if !silent {
				if count == 50 {
					count = 0
					fmt.Println()
				} else {
					count++
					fmt.Print(".")
				}
			}
			if torips.IPinList(ip.DstIP) {
				fmt.Printf("\n[+] tor connection found from %v, connecting to %v\n", ip.SrcIP, ip.DstIP)
			}
		}
	}
}
