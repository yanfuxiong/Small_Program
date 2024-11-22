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
)

func MdnsHandleStream(stream network.Stream) {
	netConn := rtkUtils.NewConnFromStream(stream)
	ipAddr := rtkUtils.GetRemoteAddrFromStream(stream)

	ip, port := rtkUtils.ExtractTCPIPandPort(stream.Conn().LocalMultiaddr())
	rtkGlobal.NodeInfo.IPAddr.PublicIP = ip
	rtkGlobal.NodeInfo.IPAddr.PublicPort = port
	fmt.Printf("public ip [%s] port[%s]\n", ip, port)
	fmt.Println("************************************************")
	log.Println("H Connected to ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
	fmt.Println("************************************************")
	rtkPlatform.GoUpdateClientStatus(1, ipAddr, stream.Conn().RemotePeer().String(), ipAddr)
	rtkUtils.InsertMdnsClientList(stream.Conn().RemotePeer().String(), ipAddr)
	rtkPlatform.FoundPeer()

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
		rtkUtils.LostMdnsClientList(stream.Conn().RemotePeer().String())
		rtkPlatform.FoundPeer()
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
		fmt.Println("************************************************")
		rtkPlatform.GoUpdateClientStatus(0, ipAddr, stream.Conn().RemotePeer().String(), ipAddr)
		stream.Close()
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
	log.Println("E Connected to ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
	fmt.Println("************************************************")
	rtkPlatform.GoUpdateClientStatus(1, ipAddr, stream.Conn().RemotePeer().String(), ipAddr)
	rtkUtils.InsertMdnsClientList(stream.Conn().RemotePeer().String(), ipAddr)
	rtkPlatform.FoundPeer()

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
		rtkUtils.LostMdnsClientList(stream.Conn().RemotePeer().String())
		rtkPlatform.FoundPeer()
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer().String(), " IP:", ipAddr)
		fmt.Println("************************************************")
		rtkPlatform.GoUpdateClientStatus(0, ipAddr, stream.Conn().RemotePeer().String(), ipAddr)
		stream.Close()
	}()
}
