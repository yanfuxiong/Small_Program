package libp2p_clipboard

import (
	"log"
	rtkCmd "rtk-cross-share/cmd"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"
	"strings"
)

type Callback interface {
	rtkPlatform.Callback
}

func MainInit(cb Callback, serverId string, serverIpInfo string, listentPort int) {
	rtkCmd.MainInit(cb, serverId, serverIpInfo, listentPort)
}

func SendMessage(s string) {
	rtkPlatform.SendMessage(s)
}

func GetClientList() string {
	//clientList := "QmVzXN1gJcz5kturdM1CY8KCyiinRbDWWmucyL9PBcvhEr#uEiDbhL6yOiMzyjXaGnHgBDayDsgyIHwXfE6Lrvb0JhkTuA#QmVw7ZAfgAqWn1sv25oMHtD3VgMdeeue8Bk9CdPnb7cd15"
	log.Printf("GetClientList  len:%d", len(rtkGlobal.MdnsPeerList))
	var clientList string
	for _, val := range rtkGlobal.MdnsPeerList {
		ip, port := rtkUtils.ExtractTCPIPandPort(val.Addrs[0])
		clientList = clientList + rtkUtils.ConcatIP(ip, port) + "#"
	}
	return clientList
}

func SendImage(content []byte, ipAddr string) {
	nLen := len(content)
	log.Printf("xyf-android SendImage:[%d][%+v]", nLen, content)
	rtkGlobal.Handler.CopyImgHeader.Height = 100
	rtkGlobal.Handler.CopyImgHeader.Width = 100
	rtkGlobal.Handler.CopyImgData = content
	rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
	rtkGlobal.Handler.CopyDataSize.SizeLow = uint32(len(content))
	rtkGlobal.Handler.IpAddr = ipAddr
}

func ClipboardImagePaste(path string) {
	log.Printf("ClipboardPaste %s", path)
	rtkGlobal.Handler.CtxMutex.Lock()
	defer rtkGlobal.Handler.CtxMutex.Unlock()
	rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
}

func SendAddrsFromJava(addrsList string) {
	parts := strings.Split(addrsList, "#")
	rtkUtils.GetAddrsFromJava(parts)
}

func SendNetInterfaces(name, mac string) {
	log.Printf("SendNetInterfaces [%s][%s]", name, mac)
	rtkUtils.SetNetInterfaces(name, mac)
}

func SendCopyFile(filePath, ipAddr string, fileSizeHigh, fileSizeLow uint32) {
	rtkGlobal.Handler.CopyFilePath.Store(filePath)
	rtkGlobal.Handler.IpAddr = ipAddr
	rtkGlobal.Handler.CopyDataSize.SizeHigh = fileSizeHigh
	rtkGlobal.Handler.CopyDataSize.SizeLow = fileSizeLow
	log.Println("andriod Clipboard file content:", rtkGlobal.Handler.CopyFilePath, "fileSize high:", fileSizeHigh, "low:", fileSizeLow)
}

// 接收方确认是否接收文档
func IfClipboardPasteFile(isReceive bool) {
	if isReceive {
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
		rtkGlobal.Handler.DstFilePath = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/recevie.file"
	}
	log.Println("andriod IfClipboardPasteFile:  ", isReceive)
}
