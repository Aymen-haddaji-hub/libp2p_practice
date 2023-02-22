/* module to start a libp2p node with a simple protocol */

package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"time"
	libp2p "github.com/libp2p/go-libp2p"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func Ping(addr string) {
	// Parse the multiaddress of the peer we want to ping
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a libp2p node with the ping protocol enabled
	node, err := libp2p.New(libp2p.Ping(false))
	if err != nil {
		log.Fatal(err)
	}

	// Add the peer's multiaddress to the node's peerstore
	pi, err := peerstore.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}
	node.Peerstore().AddAddrs(pi.ID, pi.Addrs, time.Hour)

	// Ping the peer
	ctx := context.Background()
	resultChan := ping.Ping(ctx, node, pi.ID)

	// Wait for the result to be available on the channel
	select {
	case res := <-resultChan:
		// Print the ping result
		fmt.Printf("Ping Result: RTT=%s\n", res.RTT)
	case <-time.After(time.Second):
		log.Fatal("Ping timed out")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./ping_node <node address>")
		return
	}

	addr := os.Args[1]

	Ping(addr)
}
