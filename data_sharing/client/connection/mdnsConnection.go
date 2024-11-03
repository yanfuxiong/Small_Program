package connection

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkP2P "rtk-cross-share/p2p"
	rtkUtils "rtk-cross-share/utils"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			continue
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			//fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
			log.Printf("%s", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		log.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func MdnsHandleStream(stream network.Stream) {
	remotePeer := peer.AddrInfo{ID: stream.Conn().RemotePeer(), Addrs: []ma.Multiaddr{stream.Conn().RemoteMultiaddr()}}
	// fmt.Println("Got a new stream: ", remotePeer)

	if !rtkUtils.IsInPeerList(remotePeer.ID.String(), rtkGlobal.MdnsPeerList) {
		rtkGlobal.MdnsPeerList = append(rtkGlobal.MdnsPeerList, remotePeer)
	}

	netConn := rtkUtils.NewConnFromStream(stream)
	ipAddr := rtkUtils.GetRemoteAddrFromStream(stream)

	ip, port := rtkUtils.ExtractTCPIPandPort(stream.Conn().LocalMultiaddr())
	rtkGlobal.NodeInfo.IPAddr.PublicIP = ip
	rtkGlobal.NodeInfo.IPAddr.PublicPort = port

	fmt.Println("************************************************")
	log.Println("H Connected to ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
	fmt.Println("************************************************")

	rtkGlobal.CBData[ipAddr] = rtkCommon.ClipBoardData{
		SourceID: "",
		Hash:     "",
		FmtType:  rtkCommon.TEXT,
		Content:  []byte{},
	}

	connCtx, cancel := context.WithCancel(context.Background())
	go rtkP2P.P2PRead(netConn, ipAddr, connCtx, cancel)
	go rtkP2P.P2PWrite(netConn, ipAddr, connCtx)
	go func() {
		<-connCtx.Done()
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
		fmt.Println("************************************************")
	}()

}

func ExecuteDirectConnect(ctx context.Context, stream network.Stream) {
	netConn := rtkUtils.NewConnFromStream(stream)
	ipAddr := rtkUtils.GetRemoteAddrFromStream(stream)

	ip, port := rtkUtils.ExtractTCPIPandPort(stream.Conn().LocalMultiaddr())
	rtkGlobal.NodeInfo.IPAddr.PublicIP = ip
	rtkGlobal.NodeInfo.IPAddr.PublicPort = port

	fmt.Println("************************************************")
	log.Println("E Connected to ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
	fmt.Println("************************************************")

	rtkGlobal.CBData[ipAddr] = rtkCommon.ClipBoardData{
		SourceID: "",
		Hash:     "",
		FmtType:  rtkCommon.TEXT,
		Content:  []byte{},
	}

	connCtx, cancel := context.WithCancel(context.Background())
	go rtkP2P.P2PRead(netConn, ipAddr, connCtx, cancel)
	go rtkP2P.P2PWrite(netConn, ipAddr, connCtx)
	go func() {
		<-connCtx.Done()
		fmt.Println("************************************************")
		log.Println("Lost connection with ID:", stream.Conn().RemotePeer(), " IP:", ipAddr)
		fmt.Println("************************************************")
	}()
}
