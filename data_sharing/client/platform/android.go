//go:build android

/*
	typedef struct IMAGE_HEADER {
	    int width;
	    int height;
	    unsigned short planes;
	    unsigned short bitCount;
	    unsigned long compression;
	} IMAGE_HEADER;
*/
package platform

import (
	"context"
	"log"
	"os"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkUtils "rtk-cross-share/utils"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
)

const (
	hostID      = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.HostID"
	nodeID      = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.ID"
	receiveFile = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/"
)

var ImageData []byte
var curInputText string

type Callback interface {
	CallbackMethod(result string)
	CallbackMethodImage(content string)
	LogMessageCallback(msg string)
	EventCallback(event int)
	CallbackMethodFileConfirm(platform, fileName string, fileSize int64)
	CallbackMethodFileDone(name string, fileSize int64)
	CallbackMethodFoundPeer()
	CallbackUpdateProgressBar(size int)
	CallbackMethodFileError(name, err string)
}

var CallbackInstance Callback = nil

func SetCallback(cb Callback) {
	CallbackInstance = cb
}

func PerformTextTask(msg string) {
	log.Printf("PerformTextTask - start - msg: %s", msg)
	if CallbackInstance == nil {
		log.Println("PerformTextTask - failed - callbackInstance is nil")
		return
	}
	CallbackInstance.CallbackMethod(msg)
	log.Println("PerformTextTask - CallbackMethod - done")
}

type CallbackFunc func(rtkCommon.ClipboardResetType)

var CallbackInstanceResetCB CallbackFunc = nil

func SetResetCBCallback(cb CallbackFunc) {
	CallbackInstanceResetCB = cb
}

func WatchClipboardText(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	var lastInputText string // this must be local variable

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(100 * time.Millisecond):
			if len(curInputText) > 0 && !strings.EqualFold(curInputText, lastInputText) {
				log.Println("watchClipboardText - got new message: ", curInputText)
				lastInputText = curInputText

				hash, err := rtkUtils.CreateMD5Hash([]byte(curInputText))
				if err != nil {
					log.Fatalf("Failed to create hash: %v", err)
				}
				resultChan <- rtkCommon.ClipBoardData{
					SourceID: rtkGlobal.NodeInfo.ID,
					Content:  []byte(curInputText),
					FmtType:  rtkCommon.TEXT,
					Hash:     hash.B58String(),
				}
			}
		}
	}

}

func SendMessage(s string) {
	log.Printf("SendMessage:[%s] ", s)
	if s == "" || len(s) == 0 {
		return
	}
	curInputText = s
}

func SetupCallbackSettings() {

}

func GoClipboardPasteFileCallback(content string) {

}

func GoSetupDstPasteFile(desc, fileName, platform string, fileSizeHigh uint32, fileSizeLow uint32) {
	fileSize := int64(fileSizeHigh)<<32 | int64(fileSizeLow)
	log.Printf("(DST) GoSetupDstPasteFile  sourceID:%s fileName:[%s] fileSize:[%d]", desc, fileName, fileSize)
	rtkGlobal.Handler.CopyFileName = fileName
	CallbackInstance.CallbackMethodFileConfirm(platform, fileName, fileSize)
}

func GoSetupFileDrop(desc, fileName, platform string, fileSizeHigh uint32, fileSizeLow uint32) {
	fileSize := int64(fileSizeHigh)<<32 | int64(fileSizeLow)
	log.Printf("(DST) GoSetupFileDrop  source:%s fileName:%s  fileSize:%d", desc, fileName, fileSize)
	rtkGlobal.Handler.CopyFileName = fileName
	CallbackInstance.CallbackMethodFileConfirm(platform, fileName, fileSize)
}

func ReceiveCopyDataDone(fmtType rtkCommon.ClipboardFmtType, fileSize int64) {
	log.Printf("ReceiveCopyDataDone: %s  size:%d", fmtType, fileSize)
	if CallbackInstance == nil {
		log.Println(" CallbackInstance is null !")
		return
	}

	if fmtType == rtkCommon.FILE {
		CallbackInstance.CallbackMethodFileDone(rtkGlobal.Handler.DstFilePath, fileSize)
	} else if fmtType == rtkCommon.IMAGE {
		go func() {
			imageBase64 := rtkUtils.Base64Encode(rtkUtils.BitmapToImage(ImageData, int(rtkGlobal.Handler.CopyImgHeader.Width), int(rtkGlobal.Handler.CopyImgHeader.Height)))
			// log.Printf("len[%d][%d][%d][%+v]", len(ImageData), len(imageBase64), rtkGlobal.Handler.CopyImgHeader.Width, imageBase64)
			CallbackInstance.CallbackMethodImage(imageBase64)
		}()
	}
}

func FoundPeer() {
	log.Println("CallbackMethodFoundPeer")
	CallbackInstance.CallbackMethodFoundPeer()
}

func GoSetupDstPasteImage(desc string, content []byte, imgHeader rtkCommon.ImgHeader, dataSize uint32) {
	log.Printf("GoSetupDstPasteImage from ID %s, len:[%d] CopySize.SizeLow:[%d]\n\n", desc, len(content), dataSize)

	rtkGlobal.Handler.CopyImgHeader.Height = imgHeader.Height
	rtkGlobal.Handler.CopyImgHeader.Width = imgHeader.Width
	rtkGlobal.Handler.CopyImgHeader.BitCount = imgHeader.BitCount
	rtkGlobal.Handler.CopyImgHeader.Planes = imgHeader.Planes
	rtkGlobal.Handler.CopyImgHeader.Compression = imgHeader.Compression
	rtkGlobal.Handler.CtxMutex.Lock()
	defer rtkGlobal.Handler.CtxMutex.Unlock()
	rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
	ImageData = []byte{}

}

func GoDataTransfer(data []byte) {
	ImageData = append(ImageData, data...)
}

func GoUpdateProgressBar(size int) {
	//log.Println("GoUpdateProgressBar size:", size)
	CallbackInstance.CallbackUpdateProgressBar(size)
}

func GoDeinitProgressBar() {

}

func GoUpdateClientStatus(status uint32, ip string, id string, name string) {

}

func GoEventHandle(fmtType rtkCommon.ClipboardFmtType, eventType int) {
	var strErr string
	if eventType == rtkCommon.EVENT_TYPE_OPEN_FILE_ERR {
		strErr = "open src file err!"
	} else if eventType == rtkCommon.EVENT_TYPE_RECV_TIMEOUT {
		strErr = "revive file time out!"
	}

	if fmtType == rtkCommon.FILE {
		CallbackInstance.CallbackMethodFileError(rtkGlobal.Handler.DstFilePath, strErr)
	}

}

func GoSetupDstPasteText(content []byte) {
	log.Printf("GoSetupDstPasteText:%s \n\n", string(content))
	PerformTextTask(string(content))
}

func GenKey() crypto.PrivKey {
	privKeyFile := "/storage/emulated/0/Android/data/com.rtk.myapplication/files/priv.pem"
	return rtkUtils.GenKey(privKeyFile)
}

func IsHost() bool {
	return rtkUtils.FileExists(hostID)
}

func GetHostID() string {
	file, err := os.Open(hostID)
	if err != nil {
		log.Println(err)
		return rtkGlobal.HOST_ID
	}
	defer file.Close()

	data := make([]byte, 1024)
	_, err = file.Read(data)
	if err != nil {
		log.Println(err)
		return rtkGlobal.HOST_ID
	}

	return string(data)
}

func GetIDPath() string {
	return nodeID
}

func GetHostIDPath() string {
	return hostID
}

func GetReceiveFilePath() string {
	return receiveFile
}
