package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"golang.design/x/clipboard"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	rtkBuildConfig "rtk-cross-share/buildConfig"
	rtkCommon "rtk-cross-share/common"
	rtkDebug "rtk-cross-share/debug"
	rtkGlobal "rtk-cross-share/global"
	rtkMdns "rtk-cross-share/mdns"
	rtkPlatform "rtk-cross-share/platform"
	rtkRelay "rtk-cross-share/relay"
	rtkUtils "rtk-cross-share/utils"
	"time"
)

func SetupGlobalHandler() {
	rtkGlobal.Handler.CopyFilePath.Store("")
	rtkGlobal.Handler.CopyImgHeader = rtkCommon.ImgHeader{
		Width:       0,
		Height:      0,
		Planes:      0,
		BitCount:    0,
		Compression: 0,
	}
	rtkGlobal.Handler.CopyDataSize = rtkCommon.FileSize{
		SizeHigh: 0,
		SizeLow:  0,
	}
}

func SetupSettings() {
	rtkPlatform.SetupCallbackSettings()
	SetupGlobalHandler()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	flag.Parse()
}

func SetupLogFileSetting() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "p2p.log",
		MaxSize:    256,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
}

func listen_addrs(port int) []string {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/%d",
		"/ip4/0.0.0.0/udp/%d/quic",
		"/ip6/::/tcp/%d",
		"/ip6/::/udp/%d/quic",
	}

	for i, a := range addrs {
		addrs[i] = fmt.Sprintf(a, port)
	}

	return addrs
}

func SetupNode() host.Host {
	priv := rtkPlatform.GenKey()

	if rtkMdns.MdnsCfg.ListenPort <= 0 {
		log.Fatal("listen port is not set")
	}

	node, err := libp2p.New(
		libp2p.ListenAddrStrings(listen_addrs(rtkMdns.MdnsCfg.ListenPort)...),
		libp2p.NATPortMap(),
		libp2p.Identity(priv),
		libp2p.ForceReachabilityPrivate(),
		libp2p.ResourceManager(&network.NullResourceManager{}),
		libp2p.EnableHolePunching(),
		libp2p.EnableRelay(),
	)
	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return nil
	}

	node.Network().Listen(node.Addrs()[0]) //mdns包此处需要指定listen地址，否则报错

	log.Println("Self ID: ", node.ID().String())
	log.Println("Self node Addr: ", node.Addrs())
	log.Println("========================\n\n")

	for _, p := range node.Peerstore().Peers() {
		node.Peerstore().ClearAddrs(p)
	}

	if rtkPlatform.IsHost() {
		rtkUtils.WriteNodeID(node.ID().String(), rtkPlatform.GetHostIDPath())
	}

	rtkUtils.WriteNodeID(node.ID().String(), rtkPlatform.GetIDPath())

	rtkGlobal.NodeInfo.IPAddr.LocalPort = rtkUtils.GetLocalPort(node.Addrs())
	rtkGlobal.NodeInfo.ID = node.ID().String()

	return node
}
func Run() {
	SetupSettings()
	//SetupLogFileSetting()

	ctx := context.Background()
	node := SetupNode()

	rtkRelay.BuildListener(ctx, node)

	rtkMdns.BuildMdnsListener(node)
	rtkMdns.BuildMdnsTalker(ctx, node)

	<-time.After(3 * time.Second)                                                              // wait mdns discovery all peers
	rtkGlobal.GuestList = rtkUtils.DistinctMdnsID(rtkGlobal.GuestList, rtkGlobal.MdnsPeerList) //delete mdns found peer  in global list

	if len(rtkGlobal.GuestList) == 0 {
		log.Println("Wait for node")
		go rtkDebug.DebugCmdLine()
		select {}
	}

	for _, target := range rtkGlobal.GuestList {
		rtkRelay.BuildTalker(ctx, node, target)
	}

	go rtkDebug.DebugCmdLine()
	select {}
}

func MainInit(cb rtkPlatform.Callback, serverId string, serverIpInfo string, listentPort int) {
	log.Println("========================")
	log.Println("Version: ", rtkBuildConfig.Version)
	log.Println("Build Date: ", rtkBuildConfig.BuildDate)
	log.Printf("========================\n\n")

	rtkPlatform.CallbackInstance = cb

	rtkMdns.MdnsCfg = &rtkMdns.Config{}
	if len(serverId) > 0 && len(serverIpInfo) > 0 {
		rtkGlobal.RelayServerID = serverId
		rtkGlobal.RelayServerIPInfo = serverIpInfo
		rtkMdns.MdnsCfg.ListenPort = listentPort

		log.Printf("set relayServerID: %s", serverId)
		log.Printf("set relayServerIPInfo: %s", serverIpInfo)
		log.Printf("set mdns listen port: %d\n", listentPort)
	}

	SetupSettings()
	//SetupLogFileSetting()

	ctx := context.Background()
	node := SetupNode()

	fmt.Println("Self ID: ", node.ID())
	fmt.Printf("========================\n\n")

	rtkRelay.BuildListener(ctx, node)

	rtkMdns.BuildMdnsListener(node)
	rtkMdns.BuildMdnsTalker(ctx, node)

	<-time.After(3 * time.Second) // wait mdns discovery all peers
	if len(rtkGlobal.GuestList) == 0 {
		log.Println("Wait for node")
		go rtkDebug.DebugCmdLine()
		select {}
	}

	rtkGlobal.GuestList = rtkUtils.DistinctMdnsID(rtkGlobal.GuestList, rtkGlobal.MdnsPeerList) //delete mdns fund peer
	for _, targetId := range rtkGlobal.GuestList {
		rtkRelay.BuildTalker(ctx, node, targetId)
	}

	go rtkDebug.DebugCmdLine()
	select {}
}

func SendFileTransCmd(cmd rtkCommon.FileTransferCmd, data interface{}) {
	switch cmd {
	case rtkCommon.FILE_TRANS_REQUEST:
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		rtkGlobal.Handler.State.State = rtkCommon.FILE_TRANS_INIT
	case rtkCommon.FILE_TRANS_ACCEPT:
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
	case rtkCommon.FILE_TRANS_CANCEL:
		// Do nothing
	case rtkCommon.FILE_TRANS_SETUP_PATH:
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		s, ok := data.(string)
		if ok {
			rtkGlobal.Handler.DstFilePath = s
		}
	}
}
