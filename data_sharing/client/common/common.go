package common

import (
	"os"
	"sync"
	"sync/atomic"
)

type IPAddrInfo struct {
	LocalPort  string
	PublicPort string
	PublicIP   string
}

type NodeInfo struct {
	IPAddr IPAddrInfo
	ID     string
}

type SocketErr int

const (
	OK SocketErr = iota + 1
	ERR_TIMEOUT
	ERR_NETWORK
	ERR_JSON
	ERR_OTHER
)

type EventType int

const (
	EVENT_TYPE_OPEN_FILE_ERR = 0
	EVENT_TYPE_RECV_TIMEOUT  = 1
)

type P2PHandler struct {
	CtxMutex sync.Mutex
	State    P2PFileTransferState
	// a user-select file path
	DstFilePath string
	// a handler of user-select file path
	DstFile *os.File
	// require a file from a remote i.e. last copy file
	SrcFile *os.File

	//the peer to specify the transmission
	AppointIpAddr string
	// record the last copy "file" info
	CopyFilePath atomic.Value
	// record the last copy "image" info
	CopyImgHeader ImgHeader
	CopyImgData   []byte
	// record the last copy "common" info
	CopyDataSize FileSize
	//file name
	FileName string
	// record the last updater from a remote
	SourceID string
	SourceIP string
}

type P2PFileTransferStateEnum string

const (
	FILE_DROP_INIT          P2PFileTransferStateEnum = "FILE_DROP_INIT" // FileDropCmd only
	SRC_INIT                P2PFileTransferStateEnum = "SRC_INIT"
	DEST_INIT               P2PFileTransferStateEnum = "DEST_INIT"
	PROCESSING_TRAN_ING     P2PFileTransferStateEnum = "PROCESSING_TRAN_ING"
	PROCESSING_TRAN_ING_ACK P2PFileTransferStateEnum = "PROCESSING_TRAN_ING_ACK"
	UNINIT                  P2PFileTransferStateEnum = "UNINIT"
)

type P2PFileTransferState struct {
	State P2PFileTransferStateEnum
}

type P2PFileTransferInbandEnum string

const (
	FILE_DROP_INITIATE   P2PFileTransferInbandEnum = "FILE_DROP_INITIATE" // FileDropCmd only
	FILE_REQ             P2PFileTransferInbandEnum = "FILE_REQ"
	FILE_REQ_ACK         P2PFileTransferInbandEnum = "FILE_REQ_ACK"
	FILE_REQ_DONE        P2PFileTransferInbandEnum = "FILE_REQ_DONE"
	FILE_TRAN_CANCEL_SRC P2PFileTransferInbandEnum = "FILE_TRAN_CANCEL_SRC"
	FILE_TRAN_CANCEL_DST P2PFileTransferInbandEnum = "FILE_TRAN_CANCEL_DST"
	TEXT_TRAN            P2PFileTransferInbandEnum = "TEXT_TRAN"
)

type FileDropCmd string

const (
	FILE_DROP_REQUEST FileDropCmd = "FILE_DROP_REQUEST"
	FILE_DROP_ACCEPT  FileDropCmd = "FILE_DROP_ACCEPT"
	FILE_DROP_CANCEL  FileDropCmd = "FILE_DROP_CANCEL"
)

type ClipboardFmtType string

const (
	TEXT  ClipboardFmtType = "TEXT"
	IMAGE ClipboardFmtType = "IMAGE"
	FILE  ClipboardFmtType = "FILE"
)

type P2PMessage struct {
	SourceID      string
	InbandCmd     P2PFileTransferInbandEnum
	Hash          string
	PacketID      int
	FmtType       ClipboardFmtType
	Buf           []byte
	CopySize      FileSize
	CopyImgHeader ImgHeader // TODO: consider combine CopyImgHeader with Buf
}

type ClipBoardData struct {
	SourceID string
	Hash     string
	FmtType  ClipboardFmtType
	Content  []byte
	CopySize FileSize
}

type FileSize struct {
	SizeHigh uint32
	SizeLow  uint32
}

type FileInfo struct {
	FileSize_ FileSize
	FilePath  string
}

type ImgHeader struct {
	Width       int32
	Height      int32
	Planes      uint16
	BitCount    uint16
	Compression uint32
}

type ConnectMessage struct {
	Tag           string
	ObservedAddrs string
}

type SyncMessage struct {
	Tag string
}

type RegMessage struct {
	HOST  string
	GUEST string
}

type RegResponseMessage struct {
	GUEST_LIST            []string
	GUEST_PUBLIC_TCP_IP   string
	GUEST_PUBLIC_TCP_PORT string
}

const (
	P2P_EVENT_SERVER_CONNEDTED    = 0
	P2P_EVENT_SERVER_CONNECT_FAIL = 1
	P2P_EVENT_CLIENT_CONNEDTED    = 2
	P2P_EVENT_CLIENT_CONNECT_FAIL = 3
)
