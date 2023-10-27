package torips

import (
	"io"
	"net/http"
	"os"
)

// Set the GitHub repository URL and file path
var (
	re = "https://raw.githubusercontent.com/SecOps-Institute/Tor-IP-Addresses/master/tor-nodes.lst"
)

// RefreshList fetches the latest list of Tor IP addresses from the GitHub repository.
func RefreshList() (size int64) {
	// Create the file
	file, err := os.Create(filename)
	size = 0
	if err == nil {
		defer file.Close()
		resp, err := http.Get(re)
		if err == nil {
			b, err := io.Copy(file, resp.Body)
			if err == nil {
				size = b
			}
			defer resp.Body.Close()
		}
	}
	return
}