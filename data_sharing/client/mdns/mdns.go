package mdns

import (
	"context"
	"fmt"
	"log"
	rtkConnection "rtk-cross-share/connection"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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
	if err := ser.Start( /*rtkUtils.GetNetInterfaces()*/ ); err != nil {
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
				// if other end peer id greater than us, don't connect to it, just wait for it to connect us
				log.Println("Found peer:", peer, " id is greater than us, wait for it to connect to us")
				continue
			}
			if !rtkUtils.IsInPeerList(peer.ID.String(), rtkGlobal.MdnsPeerList) {
				rtkGlobal.MdnsPeerList = append(rtkGlobal.MdnsPeerList, peer)
			}

			fmt.Println("Found peer:", peer, ", connecting...")

			if err := node.Connect(ctx, peer); err != nil {
				fmt.Println("Connection failed:", err)
				continue
			}
			// open a stream, this stream will be handled by handleStream other end
			stream, err := node.NewStream(ctx, peer.ID, protocol.ID(rtkGlobal.ProtocolDirectID))
			if err != nil {
				fmt.Println("Stream open failed", err)
			}

			rtkConnection.ExecuteDirectConnect(ctx, stream)
			log.Println("BuildMdnsTalker Connected to:", peer, " success!")
		}
	}()

}
