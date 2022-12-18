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

// Flag variables
var (
	live    bool
	silent  bool
	refresh bool
	offline bool
	dev     string
	file    string
)

// Package variables
var (
	h    *pcap.Handle
	err  error
	scan bool = false
)

// Default variables
var (
	snaplen     uint32 = 1024
	promiscuous bool   = true
)

func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	version := pcap.Version()
	if version != "" {
		fmt.Printf("[+] %s found\n", version)
	}

	flag.BoolVar(&live, "live", false, "capture packets live from interface")
	flag.BoolVar(&silent, "silent", false, "disable the dots that prove it does work")
	flag.BoolVar(&refresh, "refresh", false, "fetch and store the latest list of Tor IPs")
	flag.StringVar(&dev, "d", "", "specify interface to capture on if live mode enabled")
	flag.BoolVar(&offline, "offline", false, "read packets from pcap file")
	flag.StringVar(&file, "file", "", "specify filename if offline mode enabled")
	flag.Parse()

	// Temporary flag checks
	if !live && !refresh && !offline {
		fmt.Fprintln(os.Stderr, "[-] plis enable live flag as it will not work any other way for now")
		os.Exit(0)
	}

	// Check both offline and live modes are not enabled, if one of them is enabled scan continues
	if live && offline {
		fmt.Fprintln(os.Stderr, "[-] cannot read from file and capture live packets simultaneously")
	} else if live || offline {
		scan = true
	}

	// Refresh list if specified
	if refresh {
		size := torips.RefreshList()
		fmt.Printf("[+] written %v bytes, list of IPs refreshed\n", size)
	}

	if !scan {
		os.Exit(0)
	}

	// Verify interface for live capture
	if live {
		devs, _ := pcap.FindAllDevs()
		flag := false
		for _, i := range devs {
			if i.Name == dev {
				fmt.Printf("[+] %v interface found\n", dev)
				scan = true
				flag = true
				break
			}
		}
		if !flag {
			fmt.Fprintln(os.Stderr, "[-] network interface not found")
			os.Exit(1)
		}
	}

	if live {
		h, err = pcap.OpenLive(dev, int32(snaplen), promiscuous, (-1 * time.Second))
	} else if offline && file != "" {
		h, err = pcap.OpenOffline(file)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Println("Exiting..")
		os.Exit(2)
	}
	defer h.Close()

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
