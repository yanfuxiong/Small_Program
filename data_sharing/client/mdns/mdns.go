package mdns

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"log"
	rtkConnection "rtk-cross-share/connection"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start( /*rtkUtils.GetNetInterfaces()*/); err != nil {
		panic(err)
	}
	return n.PeerChan
}

func BuildMdnsListener(node host.Host) {

	node.SetStreamHandler(rtkGlobal.ProtocolDirectID, func(s network.Stream) {
		rtkConnection.MdnsHandleStream(s)
	})
}

func BuildMdnsTalker(ctx context.Context, node host.Host) {
	peerChan := initMDNS(node, rtkPlatform.GetHostID())

	go func() {
		for {
			peer := <-peerChan
			if peer.ID > node.ID() {
				log.Println("Found peer:", peer, " id is greater than us, wait for it to connect to us")
				continue
			}
			fmt.Println("Found peer:", peer, ", connecting...")

			if err := node.Connect(ctx, peer); err != nil {
				fmt.Println("Connection failed:", err)
				continue
			}

			stream, err := node.NewStream(ctx, peer.ID, protocol.ID(rtkGlobal.ProtocolDirectID))
			if err != nil {
				fmt.Println("Stream open failed", err)
				continue
			}

			rtkConnection.ExecuteDirectConnect(ctx, stream)
			log.Println("BuildMdnsTalker Connected to:", peer.ID.String(), " success!")
		}
	}()

}
