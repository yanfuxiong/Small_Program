//go:build android

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

// TODO: consider to replace int with long long type
func MainInit(cb Callback, serverId, serverIpInfo, listentHost string, listentPort int) {
	rtkCmd.MainInit(cb, serverId, serverIpInfo, listentHost, listentPort)
}

func SetMainCallback(cb Callback) {
	rtkPlatform.SetCallback(cb)
}

func SendMessage(s string) {
	rtkPlatform.SendMessage(s)
}

func GetClientList() string {
	clientList := rtkUtils.GetClientList()
	log.Printf("GetClientList :[%s]", clientList)
	return clientList
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
	if size == 0 {
		log.Println("GetByteImageInfo err!")
		return
	}
	log.Printf("SendImage:[%d][%d][%d][%s]", len(content), len(data), size, ipAddr)

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

// TODO: consider to replace int with long long type
func SendCopyFile(filePath, ipAddr string, fileSizeHigh, fileSizeLow int) {
	if filePath == "" || len(filePath) == 0 || fileSizeLow == 0 {
		log.Printf("filePath:[%s] or fileSizeLow:[%d] is null", filePath, fileSizeLow)
		return
	}
	rtkGlobal.Handler.AppointIpAddr = ipAddr
	var fileInfo = rtkCommon.FileInfo{
		FileSize_: rtkCommon.FileSize{
			SizeHigh: uint32(fileSizeHigh),
			SizeLow:  uint32(fileSizeLow),
		},
		FilePath: filePath,
	}
	rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_REQUEST, fileInfo)
	log.Println("(SRC)Send file:", rtkGlobal.Handler.CopyFilePath.Load().(string), "fileSize high:", fileSizeHigh, "low:", fileSizeLow)
}

func IfClipboardPasteFile(isReceive bool) {
	if isReceive {
		fileFullName := rtkPlatform.GetReceiveFilePath()
		if rtkGlobal.Handler.CopyFileName != "" {
			fileFullName += rtkGlobal.Handler.CopyFileName
		} else {
			fileFullName += "recevie.file"
		}

		rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_ACCEPT, fileFullName)
		log.Printf("(DST) FilePath:[%s] confirm receipt", rtkGlobal.Handler.DstFilePath)
	} else {
		log.Printf("(DST) Filename:[%s] reject.", rtkGlobal.Handler.CopyFileName)
	}
}
