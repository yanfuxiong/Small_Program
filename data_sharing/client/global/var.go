package global

import (
	rtkCommon "rtk-cross-share/common"
	"sync"
	"time"
)

var NodeInfo = rtkCommon.NodeInfo{
	IPAddr: rtkCommon.IPAddrInfo{
		LocalPort:  "",
		PublicPort: "",
		PublicIP:   "",
	},
	ID: "",
}

var Handler = rtkCommon.P2PHandler{
	State: rtkCommon.P2PFileTransferState{
		State: rtkCommon.UNINIT,
	},
	DstFile:       nil,
	SrcFile:       nil,
	DstFilePath:   "",
	SourceID:      "",
	SourceIP:      "",
	AppointIpAddr: "",
	IsFileDropMap: make(map[string]bool),
	CopyFileName:  "",
	CopyImgHeader: rtkCommon.ImgHeader{
		Width:       0,
		Height:      0,
		Planes:      0,
		BitCount:    0,
		Compression: 0,
	},
}

var (
	//RelayServerID     = "QmT4ZCzr1Jhnk2B81fgSsuu9t2qnexo8oP5b1m5eUcSxrg"
	//RelayServerIPInfo = "/ip4/180.218.164.96/tcp/8878/p2p/"
	RelayServerID     = "QmYsgB9x85LDSn8aon1kaaFnJP5LbZNv3apJVi67gnr8gv"
	RelayServerIPInfo = "/ip4/192.168.153.40/tcp/7999/p2p/"
	GuestList         []string
	MdnsClientList    []rtkCommon.ClientInfo
	MdnsListRWMutex   = sync.RWMutex{}

	CBData sync.Map
	RTT    map[string]time.Duration = make(map[string]time.Duration)
)
