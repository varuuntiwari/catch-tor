# tor-detect
tool to check for tor running on the network

# Flags
|   flag    |                       description                     |
|-----------|-------------------------------------------------------|
| -live     | capture packets live from interface                   |
| -silent   | disable the dots that prove it does work              |
| -refresh  | fetch and store the latest list of Tor IPs            |
| -d        | specify interface to capture on if live mode enabled  |
| -offline  | read packets from pcap file                           |
| -file     | specify filename if offline mode enabled              |

# Usage
Simply run the program using `go run catch-tor.go <flags>`, which will analyze every packet from the PacketSource for a match with an IP confirmed as a Tor node, thus establishing proof of Tor connection.
