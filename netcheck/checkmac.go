// Package netcheck provides functions to verify the interface and MAC address of the device.
package netcheck

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/google/gopacket/pcap"
)

// VerifyInterface verifies whether the interface exists and is available.
func VerifyInterface(dev string) (string, bool) {
	if dev == "" {
		return "", false
	}

	switch runtime.GOOS {
	case "windows":
		devs, _ := net.Interfaces()
		for _, i := range devs {
			if i.Name == dev {
				if name := getDeviceNameWindows(i.HardwareAddr.String()); name != "" {
					return name, true
				}
			}
		}

	case "linux":
		devs, _ := pcap.FindAllDevs()
		for _, i := range devs {
			if i.Name == dev {
				fmt.Printf("[+] %v interface found\n", dev)
				return dev, true
			}
		}
		fmt.Fprintln(os.Stderr, "[-] network interface not found")
		os.Exit(1)
	default:
		fmt.Fprintln(os.Stderr, "[-] unsupported OS")
		os.Exit(1)
	}
	return "", false
}

// getDeviceNameWindows returns the actual device name from the friendly name, such as wlan0 or eth0.
// It is used as a workaround for the issue with Windows using long and arduous names for devices.
func getDeviceNameWindows(mac string) string {
	ip := getIPfromMac(mac)
	devs, _ := pcap.FindAllDevs()

	for _, i := range devs {
		for _, addr := range i.Addresses {
			if addr.IP.String() == ip {
				return i.Name
			}
		}
	}
	return ""
}

// getIPfromMac returns the IP address of the device with the specified MAC address. It is used in
// combination with getDeviceNameWindows to locate the device name using the IP of the device.
func getIPfromMac(mac string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		addresses, err := iface.Addrs()
		if err != nil {
			return ""
		}

		for _, addr := range addresses {
			switch v := addr.(type) {
			case *net.IPNet:
				hwAddr := iface.HardwareAddr.String()
				if hwAddr == mac {
					return v.IP.String()
				}
			}
		}
	}
	return ""
}
