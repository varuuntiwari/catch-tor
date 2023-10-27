# catch-tor
Check for tor running on the network. The program analyzes every packet from the PacketSource for a match with an IP confirmed as a Tor node, thus establishing proof of Tor connection. The list of Tor nodes is taken from `https://github.com/SecOps-Institute/Tor-IP-Addresses` and can be refreshed using the flag -refresh.

# Flags
|   flag    |                       description                     |
|-----------|-------------------------------------------------------|
| -live     | capture packets live from interface                   |
| -refresh  | fetch and store the latest list of Tor IPs            |
| -d        | specify interface to capture on if live mode enabled  |
| -offline  | read packets from pcap file                           |
| -file     | specify filename if offline mode enabled              |

# Usage
Run the program using `go run catch-tor.go <flags>`.
- -live always requires an interface name specified with -d.
- -offline always requires a file specified with -file.

## Examples
- `.\catch-tor -live -d wlan0` will look for live Tor connections running on the interface named `wlan0`.
- `.\catch-tor -offline -file test.pcap` will look for Tor connections through the file `test.pcap`.