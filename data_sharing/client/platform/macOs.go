//go:build darwin

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

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.design/x/clipboard"
)

type Callback interface {
	CallbackMethod(result string)
	CallbackMethodImage(content []byte)
	LogMessageCallback(msg string)
	EventCallback(event int)
	CallbackMethodFileConfirm(platform string, fileSize int64)
}

var CallbackInstance Callback = nil

func SetCallback(cb Callback) {
	CallbackInstance = cb
}

type CallbackFunc func(rtkCommon.ClipboardResetType)

var CallbackInstanceResetCB CallbackFunc = nil

func SetResetCBCallback(cb CallbackFunc) {
	CallbackInstanceResetCB = cb
}

func WatchClipboardText(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	changed := clipboard.Watch(ctx, clipboard.FmtText)
	for {
		select {
		case <-ctx.Done():
			return
		case content := <-changed:
			if CallbackInstanceResetCB != nil {
				CallbackInstanceResetCB(rtkCommon.CLIPBOARD_RESET_TYPE_TEXT)
			}
			hash, err := rtkUtils.CreateMD5Hash(content)
			if err != nil {
				log.Fatalf("Failed to create hash: %v", err)
			}

			resultChan <- rtkCommon.ClipBoardData{
				SourceID: rtkGlobal.NodeInfo.ID,
				Content:  content,
				FmtType:  rtkCommon.TEXT,
				Hash:     hash.B58String(),
			}
		}
	}
}

func SetupCallbackSettings() {

}

func GoClipboardPasteFileCallback(content string) {

}

func GoSetupDstPasteFile(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {

}

func GoSetupFileDrop(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {

}

func GoSetupDstPasteImage(desc string, content []byte, imgHeader rtkCommon.ImgHeader, dataSize uint32) {

}

func GoDataTransfer(data []byte) {

}

func GoUpdateProgressBar(size int) {

}

func GoDeinitProgressBar() {

}

func GoUpdateClientStatus(status uint32, ip string, id string, name string) {

}

func GoEventHandle(fmtType rtkCommon.ClipboardFmtType, eventType int) {

}

func GoSetupDstPasteText(content []byte) {
	clipboard.Write(clipboard.FmtText, content)
}

func ReceiveFileConfirm(fileSize int64) {

}

func ReceiveFileDone(fileType rtkCommon.ClipboardFmtType, fileSize int64) {

}

func FoundPeer() {

}

func GenKey() crypto.PrivKey {
	privKeyFile := ".priv.pem"
	return rtkUtils.GenKey(privKeyFile)
}

func GetIDPath() string {
	return ".ID"
}

func GetHostIDPath() string {
	return ".HostID"
}

func IsHost() bool {
	return rtkUtils.FileExists(".HostID")
}

func GetHostID() string {
	file, err := os.Open(".HostID")
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
