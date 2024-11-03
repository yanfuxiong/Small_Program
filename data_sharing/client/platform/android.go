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
	hostID = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.HostID"
	nodeID = "/storage/emulated/0/Android/data/com.rtk.myapplication/files/ID.ID"
)

type ClipBoardTextData struct {
	LastInputText string
	CurInputText  string
}

var CBTextData = ClipBoardTextData{
	LastInputText: "",
	CurInputText:  "",
}

type Callback interface {
	CallbackMethod(result string)
	CallbackMethodImage(content string)
	LogMessageCallback(msg string)
	EventCallback(event int)
	CallbackMethodFileConfirm(platform string, fileSize int64)
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
	log.Println("android SendMessage: ", s)
	CBTextData.CurInputText = s
}

func SendFileTransCmdCallback(cmd rtkCommon.FileTransferCmd) {
	switch cmd {
	case rtkCommon.FILE_TRANS_REQUEST:
		CallbackInstance.CallbackMethod("FILE_TRANS_REQUEST")
	}
}

func SetupCallbackSettings() {

}

func GoClipboardPasteFileCallback(content string) {

}

func GoSetupDstPasteFile(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {

}

func ReceiveFileConfirm(fileSize int64) {

	log.Println("ReceiveFileConfirm:", fileSize)
	CallbackInstance.CallbackMethodFileConfirm("android", fileSize)
}

func GoSetupDstPasteImage(desc string, content []byte, imgHeader rtkCommon.ImgHeader, dataSize uint32) {
	log.Printf("android GoSetupDstPasteImage from ID %s, len:[%d] CopySize.SizeLow:[%d]\n\n", desc, len(content), dataSize /*, string(content)*/)
	//clipboard.Write(clipboard.FmtImage, content)

	/*if CallbackInstance != nil {
		go func() {
			CallbackInstance.CallbackMethodImage(string(content))
			log.Println("android GoSetupDstPasteImage - CallbackMethod - done")
		}()
	} else {
		log.Println("android GoSetupDstPasteImage - failed - callbackInstance is nil")
	}*/
}

func GoDataTransfer(data []byte) {
	log.Printf("andriod Get Data transfer:%d", len(data))
	rtkGlobal.AndriodDataTransfer = append(rtkGlobal.AndriodDataTransfer, data...)

	if rtkGlobal.Handler.DstFile != nil {
		n, err := rtkGlobal.Handler.DstFile.Write(data)
		if err != nil {
			log.Printf(".Handler.DstFile.Write err:%+v", err)
		}
		log.Printf("andriod Get Data DstFile.Write:%d size success!", n)
	} else {
		log.Printf(".Handler.DstFile is not open! ")
	}

}

func GoEventHandle(eventType int) {

}

func GoSetupDstPasteText(content []byte) {
	log.Printf("android GoSetupDstPasteText:%s \n\n", string(content))
	PerformTextTask(string(content))

	//clipboard.Write(clipboard.FmtText, content)  测试好像不行
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
