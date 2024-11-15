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
	"fmt"
	"log"
	"net"
	"os"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkUtils "rtk-cross-share/utils"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
)

const (
	hostID      = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.HostID"
	nodeID      = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.ID"
	receiveFile = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/"
)

type ClipBoardTextData struct {
	LastInputText string
	CurInputText  string
}

var CBTextData = ClipBoardTextData{
	LastInputText: "",
	CurInputText:  "",
}

var ImageData []byte

type Callback interface {
	CallbackMethod(result string)
	CallbackMethodImage(content string)
	LogMessageCallback(msg string)
	EventCallback(event int)
	CallbackMethodFileConfirm(platform string, fileSize int64)
	CallbackMethodFileDone(name string, fileSize int64)
	CallbackMethodFoundPeer()
}

var CallbackInstance Callback = nil

func SetCallback(cb Callback) {
	CallbackInstance = cb
}

func PerformTextTask(msg string) {
	r := fmt.Sprintf("PerformTextTask - start - msg: %s", msg)
	log.Println(r)
	if CallbackInstance != nil {
		go func() {
			CallbackInstance.CallbackMethod(msg)
			log.Println("PerformTextTask - CallbackMethod - done")
		}()
	} else {
		log.Println("PerformTextTask - failed - callbackInstance is nil")
	}
}

type CallbackFunc func(rtkCommon.ClipboardFmtType)

var CallbackInstanceResetCB CallbackFunc = nil

func SetResetCBCallback(cb CallbackFunc) {
	CallbackInstanceResetCB = cb
}

func WatchClipboardText(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	for {
		if len(CBTextData.CurInputText) > 0 && CBTextData.CurInputText != CBTextData.LastInputText {
			log.Println("android watchClipboardText - got new message: ", CBTextData.CurInputText)
			CBTextData.LastInputText = CBTextData.CurInputText
			hash, err := rtkUtils.CreateMD5Hash([]byte(CBTextData.LastInputText))
			if err != nil {
				log.Fatalf("Failed to create hash: %v", err)
			}
			resultChan <- rtkCommon.ClipBoardData{
				SourceID: rtkGlobal.NodeInfo.ID,
				Content:  []byte(CBTextData.LastInputText),
				FmtType:  rtkCommon.TEXT,
				Hash:     hash.B58String(),
			}
		} else {
			time.Sleep(time.Second)
			//log.Println("watchClipboardText - no new message")
		}
	}
}

func WatchAndroidClipboardText(s net.Conn) {
	log.Println("watchAndroidClipboardText")
	for {
		if len(CBTextData.CurInputText) > 0 && CBTextData.CurInputText != CBTextData.LastInputText {
			CBTextData.LastInputText = CBTextData.CurInputText
			if _, err := s.Write([]byte(CBTextData.LastInputText)); err != nil {
				log.Println("Cannot write")
			}
			fmt.Println(">>>")
		}
		time.Sleep(time.Second)
	}
}

func SendMessage(s string) {
	log.Printf("android SendMessage:[%s] ", s)
	if s == "" || len(s) == 0 {
		return
	}
	CBTextData.CurInputText = s
}

func SetupCallbackSettings() {

}

func GoClipboardPasteFileCallback(content string) {

}

func GoSetupDstPasteFile(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {
	log.Printf("GoSetupDstPasteFile  sourceID：%s fileName:[%s] [%d][%d]", desc, fileName, fileSizeHigh, fileSizeLow)
}

func GoSetupFileDrop(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {
	fileSize := int64(fileSizeHigh)<<32 | int64(fileSizeLow)
	log.Printf("GoSetupFileDrop  source:%s name:%s  fileSize:%d", desc, fileName, fileSize)
	rtkGlobal.Handler.FileName = fileName
	CallbackInstance.CallbackMethodFileConfirm("android", fileSize)
}

func ReceiveCopyDataDone(fileType rtkCommon.ClipboardFmtType, fileSize int64) {
	log.Printf("ReceiveFileDone: %s  size:%d", fileType, fileSize)
	if CallbackInstance == nil {
		log.Println(" CallbackInstance is null !")
		return
	}

	if fileType == rtkCommon.FILE {
		CallbackInstance.CallbackMethodFileDone(rtkGlobal.Handler.DstFilePath, fileSize)
	} else if fileType == rtkCommon.IMAGE {
		go func() {

			imageBase64 := rtkUtils.Base64Encode(rtkUtils.BitmapToImage(ImageData, int(rtkGlobal.Handler.CopyImgHeader.Width), int(rtkGlobal.Handler.CopyImgHeader.Height)))

			//log.Printf("len[%d][%d][%d][%+v]", len(ImageData), len(imageBase64), rtkGlobal.Handler.CopyImgHeader.Width, imageBase64)
			CallbackInstance.CallbackMethodImage(imageBase64)
		}()
	}
}

func FoundPeer() {
	log.Println("android CallbackMethodFoundPeer")
	CallbackInstance.CallbackMethodFoundPeer()
}

func GoSetupDstPasteImage(desc string, content []byte, imgHeader rtkCommon.ImgHeader, dataSize uint32) {
	log.Printf("android GoSetupDstPasteImage from ID %s, len:[%d] CopySize.SizeLow:[%d]\n\n", desc, len(content), dataSize)

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
	log.Println("GoDataTransfer len:", len(data))
	ImageData = append(ImageData, data...)

}

func GoEventHandle(eventType int) {

}

func GoSetupDstPasteText(content []byte) {
	log.Printf("android GoSetupDstPasteText:%s \n\n", string(content))
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
