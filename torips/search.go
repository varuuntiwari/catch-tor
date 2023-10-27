package torips

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var filename = "tor-nodes.lst"

func IPinList(ip net.IP) (found bool) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		found = (net.ParseIP(scanner.Text()).Equal(ip))
		if found {
			fmt.Println("Found Tor IP address")
			return
		}
	}
	return
}
