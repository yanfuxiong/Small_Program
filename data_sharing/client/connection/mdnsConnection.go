package connection

import (
	"context"
	"fmt"
	"log"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkP2P "rtk-cross-share/p2p"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func MdnsHandleStream(stream network.Stream) {
	remotePeer := peer.AddrInfo{ID: stream.Conn().RemotePeer(), Addrs: []ma.Multiaddr{stream.Conn().RemoteMultiaddr()}}

	if !rtkUtils.IsInPeerList(remotePeer.ID.String(), rtkGlobal.MdnsPeerList) {
		rtkGlobal.MdnsPeerList = append(rtkGlobal.MdnsPeerList, remotePeer)
		rtkPlatform.FoundPeer()
	}

	netConn := rtkUtils.NewConnFromStream(stream)
	ipAddr := rtkUtils.GetRemoteAddrFromStream(stream)

	ip, port := rtkUtils.ExtractTCPIPandPort(stream.Conn().LocalMultiaddr())
	rtkGlobal.NodeInfo.IPAddr.PublicIP = ip
	rtkGlobal.NodeInfo.IPAddr.PublicPort = port
	fmt.Println("************************************************")
	log.Println("H Connected to ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
	fmt.Println("************************************************")

	rtkGlobal.CBData.Store(ipAddr, rtkCommon.ClipBoardData{
		SourceID: "",
		Hash:     "",
		FmtType:  rtkCommon.TEXT,
		Content:  []byte{},
	})

	connCtx, cancel := context.WithCancel(context.Background())
	go rtkP2P.P2PRead(netConn, ipAddr, connCtx, cancel)
	go rtkP2P.P2PWrite(netConn, ipAddr, connCtx)
	go func() {
		<-connCtx.Done()
		rtkGlobal.MdnsPeerList = rtkUtils.LostsMdnsPeerList(remotePeer.ID.String(), rtkGlobal.MdnsPeerList)
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
		fmt.Println("************************************************")
	}()

}

func ExecuteDirectConnect(ctx context.Context, stream network.Stream) {
	netConn := rtkUtils.NewConnFromStream(stream)
	ipAddr := rtkUtils.GetRemoteAddrFromStream(stream)

	ip, port := rtkUtils.ExtractTCPIPandPort(stream.Conn().LocalMultiaddr())
	rtkGlobal.NodeInfo.IPAddr.PublicIP = ip
	rtkGlobal.NodeInfo.IPAddr.PublicPort = port
	fmt.Printf("public ip [%s] port[%s]\n", ip, port)
	fmt.Println("************************************************")
	log.Println("E Connected to ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
	fmt.Println("************************************************")

	rtkGlobal.CBData.Store(ipAddr, rtkCommon.ClipBoardData{
		SourceID: "",
		Hash:     "",
		FmtType:  rtkCommon.TEXT,
		Content:  []byte{},
	})

	connCtx, cancel := context.WithCancel(context.Background())
	go rtkP2P.P2PRead(netConn, ipAddr, connCtx, cancel)
	go rtkP2P.P2PWrite(netConn, ipAddr, connCtx)
	go func() {
		<-connCtx.Done()
		rtkGlobal.MdnsPeerList = rtkUtils.LostsMdnsPeerList(stream.Conn().RemotePeer().String(), rtkGlobal.MdnsPeerList)
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
		fmt.Println("************************************************")
	}()
}
