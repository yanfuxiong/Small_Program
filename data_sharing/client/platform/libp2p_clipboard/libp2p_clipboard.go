package libp2p_clipboard

import (
	"log"
	rtkCmd "rtk-cross-share/cmd"
	rtkCommon "rtk-cross-share/common"
	rtkFileDrop "rtk-cross-share/filedrop"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"
	"strings"
)

type Callback interface {
	rtkPlatform.Callback
}

func MainInit(cb Callback, serverId string, serverIpInfo, listentHost string, listentPort int) {
	rtkCmd.MainInit(cb, serverId, serverIpInfo, listentHost, listentPort)
}

func SendMessage(s string) {
	rtkPlatform.SendMessage(s)
}

func GetClientList() string {

	log.Printf("GetClientList  len:%d", len(rtkGlobal.MdnsPeerList))
	var clientList string
	for _, val := range rtkGlobal.MdnsPeerList {
		ip, port := rtkUtils.ExtractTCPIPandPort(val.Addrs[0])
		clientList = clientList + rtkUtils.ConcatIP(ip, port) + "#"
	}

	return strings.Trim(clientList, "#")
}

func SendImage(content string, ipAddr string) {
	if content == "" || len(content) == 0 {
		return
	}
	data := rtkUtils.Base64Decode(content)
	if data == nil {
		return
	}

	w, h, size := rtkUtils.GetByteImageInfo(data)

	log.Printf("android SendImage:[%d][%d][%d][%s]", len(content), len(data), size, ipAddr)

	rtkGlobal.Handler.CopyImgHeader.Width = int32(w)
	rtkGlobal.Handler.CopyImgHeader.Height = int32(h)
	rtkGlobal.Handler.CopyImgHeader.Compression = 0
	rtkGlobal.Handler.CopyImgHeader.Planes = 1
	rtkGlobal.Handler.CopyImgHeader.BitCount = uint16((size * 8) / (w * h))
	rtkGlobal.Handler.CopyImgData = rtkUtils.ImageToBitmap(data)
	rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
	rtkGlobal.Handler.CopyDataSize.SizeLow = uint32(size)
	rtkGlobal.Handler.AppointIpAddr = ipAddr
}

func SendAddrsFromJava(addrsList string) {
	parts := strings.Split(addrsList, "#")
	rtkUtils.GetAddrsFromJava(parts)
}

func SendNetInterfaces(name, mac string) {
	log.Printf("SendNetInterfaces [%s][%s]", name, mac)
	rtkUtils.SetNetInterfaces(name, mac)
}

func SendCopyFile(filePath, ipAddr string, fileSizeHigh, fileSizeLow int) {
	rtkGlobal.Handler.AppointIpAddr = ipAddr
	var fileInfo = rtkCommon.FileInfo{
		FileSize_: rtkCommon.FileSize{
			SizeHigh: uint32(fileSizeHigh),
			SizeLow:  uint32(fileSizeLow),
		},
		FilePath: filePath,
	}
	rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_REQUEST, fileInfo)
	log.Println("andriod Send file content:", rtkGlobal.Handler.CopyFilePath.Load().(string), "fileSize high:", fileSizeHigh, "low:", fileSizeLow)
}

// 接收方确认是否接收文档
func IfClipboardPasteFile(isReceive bool) {
	log.Println("andriod IfClipboardPasteFile:  ", isReceive)
	if isReceive {
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
		fileFullName := rtkPlatform.GetReceiveFilePath()
		if rtkGlobal.Handler.FileName != "" {
			fileFullName = fileFullName + rtkGlobal.Handler.FileName
		} else {
			fileFullName = fileFullName + "recevie.file"
		}
		rtkGlobal.Handler.DstFilePath = fileFullName
		log.Printf("DstFilePath1:[%s] ", rtkGlobal.Handler.DstFilePath)
	}
}
